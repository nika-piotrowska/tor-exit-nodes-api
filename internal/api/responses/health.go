package responses

import "github.com/nika-piotrowska/tor-exit-nodes-api/internal/generated/oas"

const (
	healthType    = "health"
	readinessType = "readiness"
	serviceID     = "service"
)

// HealthOK returns a JSON:API document representing a healthy service state.
func HealthOK() *oas.HealthDocument {
	return &oas.HealthDocument{
		Data: oas.HealthResource{
			Type: healthType,
			ID:   serviceID,
			Attributes: oas.HealthAttributes{
				Status: "ok",
			},
		},
	}
}

// ReadyOK returns a JSON:API document indicating the service is ready.
func ReadyOK() *oas.ReadinessDocument {
	return &oas.ReadinessDocument{
		Data: oas.ReadinessResource{
			Type: readinessType,
			ID:   serviceID,
			Attributes: oas.ReadinessAttributes{
				Status: "ready",
				Redis:  "ok",
			},
		},
	}
}
