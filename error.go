package couchdb2_goclient

import "fmt"

type ErrorResponse struct {
	ErrorS string `json:"error, omitempty"`
	Reason string `json:"reason, omitempty"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("Error: %s\nReason: %s", e.ErrorS, e.Reason)
}