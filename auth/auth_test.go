package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func excpectedStr(expected, got any) string {
	return fmt.Sprintf("Expected: %v but got: %v instead.", expected, got)
}

func TestAuth(t *testing.T) {
	mux := newMux(NewHandler(NewService()))

	// go test ./auth/ -v -run "TestAuth/request auth token"
	t.Run("request auth token", func(t *testing.T) {
		payload := AuthCreateModel{Email: "shirohige65@rocketmail.com"}

		byte_buffer := bytes.NewBuffer(nil)
		json_encoder := json.NewEncoder(byte_buffer)
		err := json_encoder.Encode(&payload)

		if assert.Nil(t, err, excpectedStr(nil, err)) {
			req := httptest.NewRequest(http.MethodPost, "http://localhost:1312/api/auth/", byte_buffer)
			req.Header.Add("Content-Type", "application/json")
			res := AuthReadModel{}
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, req)

			result := recorder.Result()
			defer result.Body.Close()

			assert.Equal(t, http.StatusOK, result.StatusCode, excpectedStr(http.StatusOK, result.StatusCode))

			json_decoder := json.NewDecoder(result.Body)
			err := json_decoder.Decode(&res)
			if assert.Nil(t, err, excpectedStr(nil, err)) {
				// assert.Equal(t, expected, res, excpectedStr(expected, res))
			}
		}
	})
}
