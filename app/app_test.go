package app

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/mmiftahrzki/customer/auth"
	"github.com/mmiftahrzki/customer/config"
	"github.com/mmiftahrzki/customer/customer"
	"github.com/mmiftahrzki/customer/database"
	"github.com/mmiftahrzki/customer/logger"
	"github.com/mmiftahrzki/customer/responses"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB
var mux http.Handler
var testServer *httptest.Server
var baseURL string

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

func ParseToJSON[T any](response *http.Response) (T, error) {
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

func login(t *testing.T, payloadStr string) string {
	t.Helper()

	payload := bytes.NewBuffer(nil)
	writtenPayload, err := payload.WriteString(payloadStr)
	if !assert.Nil(t, err, excpectedStr(nil, err)) {
		return ""
	}

	req := httptest.NewRequest(http.MethodPost, "/api/auth", payload)
	req.Header.Add("Content-Length", strconv.Itoa(writtenPayload))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	recoder := httptest.NewRecorder()

	mux.ServeHTTP(recoder, req)

	response := recoder.Result()

	result, err := ParseToJSON[auth.ModelRead](response)

	return result.Token
}

func init() {
	var err error = nil
	logger := logger.GetLogger()

	cfgApp := config.AppConfig{
		AdminEmail: "muhammadmiftahrizki@gmail.com",
		Port:       1312,
	}
	cfgDb := config.DatabaseConfig{
		Host:          "localhost",
		Port:          3306,
		User:          "root",
		Password:      "SuperSecurePassword123!@#",
		Name:          "portfolio",
		MaxConnection: 10,
	}

	db, err = database.New(cfgDb)
	if err != nil {
		logger.Fatalf("Database Error: %v\n", err)
	}

	mux = newMux(db, cfgApp)
	testServer = httptest.NewServer(mux)
	testServer.Config.Addr = fmt.Sprintf(":%d", cfgApp.Port)
}

func TestCustomerHandler(t *testing.T) {
	baseURL = "/api/customer"

	// go test ./app/ -v -run "TestCustomerHandler/get single by id"
	t.Run("get single by id", func(t *testing.T) {
		customerId := rand.Intn(600)
		endpoint := fmt.Sprintf("%s/%d", baseURL, customerId)
		req := httptest.NewRequest(http.MethodGet, endpoint, nil)
		req.Header.Add("Accept", "application/json")

		recorder := httptest.NewRecorder()

		mux.ServeHTTP(recorder, req)

		response := recorder.Result()

		assert.Equal(t, http.StatusOK, response.StatusCode, excpectedStr(http.StatusOK, response.StatusCode))

		actualResponseBody, err := ParseToJSON[responses.GetSingleResponse[customer.ModelRead]](response)
		if assert.Nil(t, err, excpectedStr(nil, err)) {
			assert.NotEqual(t, "", actualResponseBody.Data.FullName, excpectedStr("", actualResponseBody.Data.FullName))
		}
	})

	// go test ./app/ -v -run "TestCustomerHandler/get first customer set"
	t.Run("get first customer set", func(t *testing.T) {
		testScenarios := []testScenario[[]map[int]string]{
			NewTestScenario(http.StatusOK, []map[int]string{
				{0: "MARY.SMITH@sakilacustomer.org"},
				{6: "MARIA.MILLER@sakilacustomer.org"},
				{12: "KAREN.JACKSON@sakilacustomer.org"},
				{18: "SHARON.ROBINSON@sakilacustomer.org"},
				{24: "JESSICA.HALL@sakilacustomer.org"},
			}),
		}

		for _, testScenario := range testScenarios {
			expected := testScenario.expected
			req := httptest.NewRequest(http.MethodGet, baseURL, nil)
			req.Header.Add("Accept", "application/json")
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, req)

			actual := recorder.Result()
			defer actual.Body.Close()

			if assert.Equal(t, expected.statusCode, actual.StatusCode, excpectedStr(expected.statusCode, actual.StatusCode)) {
				actualResponseBody, err := ParseToJSON[responses.GetMultipleResponse[customer.ModelRead]](actual)

				if assert.Nil(t, err, excpectedStr(nil, err)) {
					actualData := actualResponseBody.Data
					assert.True(t, len(actualData) > 0, excpectedStr(true, len(actualData) > 0))

					for _, expected := range expected.data {
						for k, v := range expected {
							assert.EqualValues(t, v, actualData[k].Email)
						}
					}
				}
			}
		}
	})

	// go test ./app/ -v -run "TestCustomerHandler/get previous customer set before a customer with some ids"
	t.Run("get previous customer set before a customer with some ids", func(t *testing.T) {
		testScenarios := []testScenarioWithInput[int, []map[int]string]{
			NewTestScenarioWithInput(51, http.StatusOK, []map[int]string{
				{0: "JESSICA.HALL@sakilacustomer.org"},
				{6: "AMY.LOPEZ@sakilacustomer.org"},
				{12: "MARTHA.GONZALEZ@sakilacustomer.org"},
				{18: "MARIE.TURNER@sakilacustomer.org"},
				{24: "DIANE.COLLINS@sakilacustomer.org"},
			}),
			NewTestScenarioWithInput(77, http.StatusOK, []map[int]string{
				{0: "ALICE.STEWART@sakilacustomer.org"},
				{6: "EVELYN.MORGAN@sakilacustomer.org"},
				{12: "ASHLEY.RICHARDSON@sakilacustomer.org"},
				{18: "CHRISTINA.RAMIREZ@sakilacustomer.org"},
				{24: "IRENE.PRICE@sakilacustomer.org"},
			}),
			NewTestScenarioWithInput(102, http.StatusOK, []map[int]string{
				{0: "JANE.BENNETT@sakilacustomer.org"},
				{6: "LOUISE.JENKINS@sakilacustomer.org"},
				{12: "JULIA.FLORES@sakilacustomer.org"},
				{18: "PAULA.BRYANT@sakilacustomer.org"},
				{24: "PEGGY.MYERS@sakilacustomer.org"},
			}),
			NewTestScenarioWithInput(128, http.StatusOK, []map[int]string{
				{0: "CRYSTAL.FORD@sakilacustomer.org"},
				{6: "TRACY.COLE@sakilacustomer.org"},
				{12: "GRACE.ELLIS@sakilacustomer.org"},
				{18: "SYLVIA.ORTIZ@sakilacustomer.org"},
				{24: "ELAINE.STEVENS@sakilacustomer.org"},
			}),
		}

		for _, testScenario := range testScenarios {
			expected := testScenario.expected
			url := fmt.Sprintf("/api/customer/%d/prev", testScenario.input)
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, req)

			actual := recorder.Result()

			if assert.Equal(t, expected.statusCode, actual.StatusCode, excpectedStr(expected.statusCode, actual.StatusCode)) {
				actualResponseBody, err := ParseToJSON[responses.GetMultipleResponse[customer.ModelRead]](actual)

				if assert.Nil(t, err, excpectedStr(nil, err)) {
					actualData := actualResponseBody.Data
					assert.True(t, len(actualData) > 0, excpectedStr(true, len(actualData) > 0))

					for _, expected := range expected.data {
						for k, v := range expected {
							assert.EqualValues(t, v, actualData[k].Email)
						}
					}
				}
			}
		}
	})

	// go test ./app/ -v -run "TestCustomerHandler/get next customer set after a customer with some ids"

	t.Run("get next customer set after a customer with some ids", func(t *testing.T) {
		testScenarios := []testScenarioWithInput[int, []map[int]string]{
			NewTestScenarioWithInput(26, http.StatusOK, []map[int]string{
				{0: "27 SHIRLEY.ALLEN@sakilacustomer.org"},
				{6: "33 ANNA.HILL@sakilacustomer.org"},
				{12: "39 DEBRA.NELSON@sakilacustomer.org"},
				{18: "45 JANET.PHILLIPS@sakilacustomer.org"},
				{24: "51 ALICE.STEWART@sakilacustomer.org"},
			}),

			NewTestScenarioWithInput(51, http.StatusOK, []map[int]string{
				{0: "52 JULIE.SANCHEZ@sakilacustomer.org"},
				{6: "58 JEAN.BELL@sakilacustomer.org"},
				{12: "65 ROSE.HOWARD@sakilacustomer.org"},
				{18: "71 KATHY.JAMES@sakilacustomer.org"},
				{24: "77 JANE.BENNETT@sakilacustomer.org"},
			}),

			NewTestScenarioWithInput(77, http.StatusOK, []map[int]string{
				{0: "78 LORI.WOOD@sakilacustomer.org"},
				{6: "84 SARA.PERRY@sakilacustomer.org"},
				{12: "90 RUBY.WASHINGTON@sakilacustomer.org"},
				{18: "96 DIANA.ALEXANDER@sakilacustomer.org"},
				{24: "102 CRYSTAL.FORD@sakilacustomer.org"},
			}),

			NewTestScenarioWithInput(102, http.StatusOK, []map[int]string{
				{0: "103 GLADYS.HAMILTON@sakilacustomer.org"},
				{6: "109 EDNA.WEST@sakilacustomer.org"},
				{12: "115 WENDY.HARRISON@sakilacustomer.org"},
				{18: "121 JOSEPHINE.GOMEZ@sakilacustomer.org"},
				{24: "128 MARJORIE.TUCKER@sakilacustomer.org"},
			}),
		}

		for _, testScenario := range testScenarios {
			expected := testScenario.expected
			url := fmt.Sprintf("/api/customer/%d/next", testScenario.input)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, req)

			actualResponse := recorder.Result()
			actual := actualResponse

			if assert.Equal(t, expected.statusCode, actual.StatusCode, excpectedStr(expected.statusCode, actual.StatusCode)) {
				actualResponseBody, err := ParseToJSON[responses.GetMultipleResponse[customer.ModelRead]](actual)
				if assert.Nil(t, err, excpectedStr(nil, err)) {
					actualData := actualResponseBody.Data

					assert.True(t, len(actualData) > 0, excpectedStr(true, len(actualData) > 0))

					for _, expected := range expected.data {
						for k, v := range expected {
							assert.EqualValues(t, v, fmt.Sprintf("%d %s", actualData[k].Id, actualData[k].Email))
						}
					}
				}
			}
		}
	})
}

func TestCustomerProtectedHandler(t *testing.T) {
	baseURL = "/api/customer"

	// go test ./app/ -v -run "TestCustomerHandler/create new single customer"
	t.Run("create new single customer by user admin role", func(t *testing.T) {
		loginPayload := `{"email": "muhammadmiftahrizki@gmail.com"}`
		baseURL = "/api/customer"
		token := login(t, loginPayload)
		if token == "" {
			t.Fatalf("token is empty")
		}
		auth := "Bearer " + token
		newCustomer := `{"first_name": "Muhammad Miftah","last_name": "Rizki","email": "muhammadmiftahrizki@gmail.com"}`
		payload := bytes.NewBuffer(nil)
		payloadLen, err := payload.WriteString(newCustomer)
		if !assert.Nil(t, err, excpectedStr(err, nil)) {
			return
		}

		recorder := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodPost, "/api/customer/", payload)
		req.Header.Add("Content-Length", strconv.Itoa(payloadLen))
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
		req.Header.Add("Authorization", auth)

		mux.ServeHTTP(recorder, req)

		result := recorder.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusCreated, result.StatusCode, excpectedStr(http.StatusCreated, result.StatusCode))
	})

	// go test ./app/ -v -run "TestCustomerHandler/edit single customer"
	t.Run("edit single customer", func(t *testing.T) {
		// 13 KAREN.JACKSON@sakilacustomer.org

		loginPayload := `{"email": "shirohige65@rocketmail.com"}`
		token := login(t, loginPayload)
		if token == "" {
			t.Fatal("token is empty")
		}
		auth := "Bearer " + token
		endpoint := fmt.Sprintf("%s/%d", baseURL, 13)
		updateCustomer :=
			`{
				"first_name": "KAREN EDITED",
				"last_name": "JACKSON EDITED",
				"email": "KAREN.JACKSON.EDITED@sakilacustomer.org",
				"address": {
					"address": "270 Amroha Parkway",
					"district": "Osmaniye",
					"postal_code": "29610",
					"city_id": 123,
					"country_id": 12
				}
			}`
		payload := bytes.NewBuffer(nil)
		payloadLen, err := payload.WriteString(updateCustomer)
		if assert.Nil(t, err, excpectedStr(nil, err)) {
			if assert.Nil(t, err, excpectedStr(nil, err)) {
				req := httptest.NewRequest(http.MethodPut, endpoint, payload)
				req.Header.Add("Content-Length", strconv.Itoa(payloadLen))
				req.Header.Add("Content-Type", "application/json; charset=utf-8")
				req.Header.Add("Authorization", auth)

				recorder := httptest.NewRecorder()

				mux.ServeHTTP(recorder, req)

				result := recorder.Result()

				if assert.Equal(t, http.StatusOK, result.StatusCode, excpectedStr(http.StatusOK, result.StatusCode)) {
					updateCustomer :=
						`{
							"first_name": "KAREN",
							"last_name": "JACKSON",
							"email": "KAREN.JACKSON@sakilacustomer.org",
							"address": {
								"address": "270 Amroha Parkway",
								"district": "Osmaniye",
								"postal_code": "29610",
								"city_id": 384,
								"country_id": 97
							}
						}`
					payload := bytes.NewBuffer(nil)
					payloadLen, err := payload.WriteString(updateCustomer)
					if assert.Nil(t, err, excpectedStr(nil, err)) {
						if assert.Nil(t, err, excpectedStr(nil, err)) {
							req := httptest.NewRequest(http.MethodPut, endpoint, payload)
							req.Header.Add("Content-Length", strconv.Itoa(payloadLen))
							req.Header.Add("Content-Type", "application/json; charset=utf-8")
							req.Header.Add("Authorization", auth)

							recorder := httptest.NewRecorder()

							mux.ServeHTTP(recorder, req)

							result := recorder.Result()

							assert.Equal(t, http.StatusOK, result.StatusCode, excpectedStr(http.StatusOK, result.StatusCode))
						}
					}
				}
			}
		}
	})

	// go test ./app/ -v -run "TestCustomerProtectedHandler/delete single customer by its id by user user role and should return 403 http status code"
	t.Run("delete single customer by its id by user user role and should return 403 http status code", func(t *testing.T) {
		toBeDeletedCustomerId := rand.Intn(600)
		endpoint := fmt.Sprintf("%s/%d", baseURL, toBeDeletedCustomerId)
		loginPayload := `{"email": "shirohige65@rocketmail.com"}`
		token := login(t, loginPayload)
		if token == "" {
			t.Fatalf("token is empty")
		}
		auth := "Bearer " + token

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, endpoint, nil)
		req.Header.Add("Authorization", auth)

		mux.ServeHTTP(recorder, req)

		result := recorder.Result()

		assert.Equal(t, http.StatusForbidden, result.StatusCode, excpectedStr(http.StatusForbidden, result.StatusCode))
	})

	// go test ./app/ -v -run "TestCustomerProtectedHandler/delete single customer by its id by user admin role and should return 204 http status code"
	t.Run("delete single customer by its id by user admin role and should return 204 http status code", func(t *testing.T) {
		toBeDeletedCustomerId := rand.Intn(600)
		endpoint := fmt.Sprintf("%s/%d", baseURL, toBeDeletedCustomerId)
		loginPayload := `{"email": "muhammadmiftahrizki@gmail.com"}`
		token := login(t, loginPayload)
		if token == "" {
			t.Fatalf("token is empty")
		}
		auth := "Bearer " + token

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, endpoint, nil)
		req.Header.Add("Authorization", auth)

		mux.ServeHTTP(recorder, req)

		result := recorder.Result()

		assert.Equal(t, http.StatusNoContent, result.StatusCode, excpectedStr(http.StatusNoContent, result.StatusCode))
	})
}
