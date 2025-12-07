package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mux http.Handler

type expectedResponse[T any] struct {
	statusCode int
	data       T
}

type testScenario[expectedType any] struct {
	expected expectedResponse[expectedType]
}

type testScenarioWithInput[inputType any, expectedType any] struct {
	input    inputType
	expected expectedResponse[expectedType]
}

func excpectedStr(expected, got any) string {
	return fmt.Sprintf("Expected: %v but got: %v instead.", expected, got)
}

func ParseToJSON[T any](response http.Response) (T, error) {
	defer response.Body.Close()

	var expected T
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return expected, err
	}

	bytes_reader := bytes.NewReader(body)
	json_decoder := json.NewDecoder(bytes_reader)

	err = json_decoder.Decode(&expected)
	if err != nil {
		return expected, err
	}

	return expected, nil
}

func NewExpectedResponse[T any](statusCode int, data T) expectedResponse[T] {
	return expectedResponse[T]{
		statusCode,
		data,
	}
}

func NewTestScenario[T any](statusCode int, expectedData T) testScenario[T] {
	return testScenario[T]{
		NewExpectedResponse(statusCode, expectedData),
	}
}

func NewTestScenarioWithInput[T any, U any](input T, statusCode int, expectedData U) testScenarioWithInput[T, U] {
	return testScenarioWithInput[T, U]{
		input,
		NewExpectedResponse(statusCode, expectedData),
	}
}

func init() {
	mux = New().Handler
}

func TestAuth(t *testing.T) {
	payload := AuthCreateModel{Email: "shirohige65@rocketmail.com"}

	testScenarios := []testScenarioWithInput[AuthCreateModel, string]{
		NewTestScenarioWithInput(payload, http.StatusOK, ""),
	}

	for _, testScenario := range testScenarios {
		byteBuffer := bytes.NewBuffer(nil)
		jsonEncoder := json.NewEncoder(byteBuffer)
		err := jsonEncoder.Encode(&testScenario.input)
		if !assert.Nil(t, err, excpectedStr(nil, err)) {
			return
		}

		expected := testScenario.expected
		req := httptest.NewRequest(http.MethodPost, "/api/auth", byteBuffer)
		req.Header.Add("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mux.ServeHTTP(recorder, req)

		actualResponse := recorder.Result()
		actual := actualResponse

		if !assert.Equal(t, expected.statusCode, actual.StatusCode, excpectedStr(expected.statusCode, actual.StatusCode)) {
			return
		}

		actualResponseBody, err := ParseToJSON[AuthReadModel](*actual)
		if !assert.Nil(t, err, excpectedStr(nil, err)) {
			return
		}

		assert.NotEqual(t, expected.data, actualResponseBody.Token, excpectedStr(expected.data, actualResponseBody.Token))
	}
}
