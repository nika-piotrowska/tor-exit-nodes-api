#!/usr/bin/env bash
set -euo pipefail

VERSION_FILE="VERSION"
CHANGELOG_FILE="CHANGELOG.md"
PUSH=false

die() { echo "Error: $*" >&2; exit 1; }

usage() {
  cat <<EOF
Usage:
  $0 <version> [--push]

Examples:
  $0 v0.1.0
  $0 0.1.0 --push

Notes:
- Creates/updates VERSION and CHANGELOG.md
- Creates an annotated git tag: vX.Y.Z
- If this is NOT the first tag, requires a previous tag to exist.
- No fallback changelog generation without a previous tag.
EOF
}

normalize_tag() {
  local v="$1"
  v="${v#v}"
  [[ "$v" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || die "Invalid version '$1' (expected X.Y.Z or vX.Y.Z)"
  echo "v$v"
}

ensure_clean() {
  git diff --quiet && git diff --cached --quiet || die "Working tree is not clean. Commit/stash changes first."
}

latest_tag() {
  git describe --tags --abbrev=0 2>/dev/null || true
}

main() {
  [[ $# -ge 1 ]] || { usage; exit 1; }

  local input_version="$1"
  shift || true

  while [[ $# -gt 0 ]]; do
    case "$1" in
      --push) PUSH=true; shift ;;
      -h|--help) usage; exit 0 ;;
      *) die "Unknown argument: $1" ;;
    esac
  done

  git rev-parse --is-inside-work-tree >/dev/null 2>&1 || die "Not a git repository."
  ensure_clean

  local new_tag
  new_tag="$(normalize_tag "$input_version")"

  # Prevent duplicate tags
  if git rev-parse -q --verify "refs/tags/$new_tag" >/dev/null; then
    die "Tag $new_tag already exists."
  fi

  local prev_tag
  prev_tag="$(latest_tag)"

  # If VERSION file exists, require prev tag (no fallback).
  # If no previous tag exists, allow only for the first release.
  if [[ -f "$VERSION_FILE" && -z "$prev_tag" ]]; then
    die "VERSION exists but no previous tag found. Create the first tag manually (or delete VERSION)."
  fi

  # Generate changelog section:
  # - If prev tag exists: list commits since prev tag
  # - If no prev tag: do NOT fallback to full history; use a fixed note
  local changes=""
  if [[ -n "$prev_tag" ]]; then
    changes="$(git log --pretty=format:"- %s" "${prev_tag}..HEAD")"
  else
    changes="- Initial release"
  fi

  # safeguard: avoid empty body/changelog
  [[ -n "$changes" ]] || changes="- No notable changes"

  # Write VERSION (store without leading v)
  echo "${new_tag#v}" > "$VERSION_FILE"

  # Prepend CHANGELOG.md
  local date
  date="$(date +%Y-%m-%d)"
  {
    echo "## ${new_tag} (${date})"
    echo
    echo "$changes"
    echo
    if [[ -f "$CHANGELOG_FILE" ]]; then
      cat "$CHANGELOG_FILE"
    fi
  } > .changelog.tmp
  mv .changelog.tmp "$CHANGELOG_FILE"

  git add "$VERSION_FILE" "$CHANGELOG_FILE"
  git commit -m "release: ${new_tag}" -m "${changes}"
  git tag -a "$new_tag" -m "Release ${new_tag}"

  if [[ "$PUSH" == "true" ]]; then
    branch="$(git rev-parse --abbrev-ref HEAD)"
    git push origin "$branch"
    git push origin "$new_tag"
  fi

  echo "Created release ${new_tag}"
}

main "$@"