package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Code      string      `json:"code,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIResponse{Success: true, Data: data})
}

func Error(w http.ResponseWriter, status int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIResponse{Success: false, Error: message, Code: code})
}

func ValidationError(w http.ResponseWriter, errors map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	_ = json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error:   "Revise os campos informados.",
		Code:    "VALIDATION_ERROR",
		Errors:  errors,
	})
}
