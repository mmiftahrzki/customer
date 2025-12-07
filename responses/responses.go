package responses

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type GetSingleResponse[T any] struct {
	Data T `json:"data"`
}

type GetMultipleResponse[T any] struct {
	Data []T    `json:"data"`
	Prev string `json:"__prev,omitempty"`
	Next string `json:"__next,omitempty"`
}

type errorResponse struct {
	Message string `json:"message" example:"internal server error"`
}

func WithJson(w http.ResponseWriter, code int, data any) {
	buffer := bytes.NewBuffer(nil)
	json_encoder := json.NewEncoder(buffer)
	err := json_encoder.Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(buffer.Bytes())
}

func Error(w http.ResponseWriter, code int, errorMessage string) {
	WithJson(w, code, errorResponse{Message: errorMessage})
}
