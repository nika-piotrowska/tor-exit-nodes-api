// Package responses provides helper builders for API responses used by handlers.
package responses

import "github.com/nika-piotrowska/tor-exit-nodes-api/internal/generated/oas"

// RedisUnavailable builds an ErrorDocument describing an unavailable Redis dependency.
func RedisUnavailable() oas.ErrorDocument {
	return oas.ErrorDocument{
		Errors: []oas.ErrorObject{
			{
				Status: "503",
				Title:  "Service Unavailable",
				Detail: oas.NewOptString("Redis is unavailable"),
			},
		},
	}
}
