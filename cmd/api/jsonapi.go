package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const jsonAPIContentType = "application/vnd.api+json"

type jsonAPIData struct {
	Type       string         `json:"type"`
	ID         string         `json:"id"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

type jsonAPISuccess struct {
	Data jsonAPIData `json:"data"`
}

type jsonAPIError struct {
	Status string `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
}

type jsonAPIErrors struct {
	Errors []jsonAPIError `json:"errors"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", jsonAPIContentType)
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	_ = enc.Encode(v)
}

func writeJSONAPISuccess(w http.ResponseWriter, status int, dataType, id string, attrs map[string]any) {
	writeJSON(w, status, jsonAPISuccess{
		Data: jsonAPIData{
			Type:       dataType,
			ID:         id,
			Attributes: attrs,
		},
	})
}

func writeJSONAPIError(w http.ResponseWriter, status int, title, detail string) {
	writeJSON(w, status, jsonAPIErrors{
		Errors: []jsonAPIError{
			{
				Status: strconv.Itoa(status),
				Title:  title,
				Detail: detail,
			},
		},
	})
}
