package httpx

import (
	"encoding/json"
	"net/http"
)

type Code string

const (
	ErrorCodeInvalidId        Code = "invalid_id"
	ErrorCodeInternalError    Code = "internal_error"
	ErrorCodeNotFound         Code = "not_found"
	ErrorCodeMalformedJson    Code = "malformed_json"
	ErrorCodeValidationFailed Code = "validation_failed"
	ErrorCodeUnauthenticated  Code = "unauthenticated"
	ErrorCodeForbidden        Code = "forbidden"
	ErrorCodeConflict         Code = "conflict"
	ErrorCodeRateLimited      Code = "rate_limited"
)

type errorEnvelope struct {
	Error errorPayload `json:"error"`
}

type errorPayload struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func Error(w http.ResponseWriter, status int, message string, code Code) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorEnvelope{Error: errorPayload{Message: message, Code: code}})
}

func ValidationError(w http.ResponseWriter, status int, message string, code Code, field string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorEnvelope{Error: errorPayload{Message: message, Code: code, Field: field}})
}
