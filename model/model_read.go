package model

import (
	"encoding/json"
	"net/http"
)

type readModel struct {
	Message string         `json:"message"`
	Data    map[string]any `json:"data,omitempty" `
}

func NewReadModel() *readModel {
	return &readModel{
		Message: "Terjadi kesalahan di sisi penyedia layanan.",
		Data:    map[string]any{},
	}
}

func (m *readModel) Send(w http.ResponseWriter) {
	w.Write([]byte("sukses"))
}

func (m *readModel) ToJson() []byte {
	json_encoded_model, _ := json.Marshal(m)

	return json_encoded_model
}
