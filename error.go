package couchdb2_go

import "fmt"

type ErrorResponse struct {
	ErrorS string `json:"error, omitempty"`
	Reason string `json:"reason, omitempty"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("Error: %s. Reason: %s", e.ErrorS, e.Reason)
}