package customer

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/mmiftahrzki/customer/config"
	"github.com/mmiftahrzki/customer/database"
	"github.com/mmiftahrzki/customer/logger"
	"github.com/mmiftahrzki/customer/responses"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB
var mux http.Handler
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
	var err error = nil
	logger := logger.GetLogger()
	cfg_db := config.DatabaseConfig{
		Host:          "localhost",
		Port:          3306,
		User:          "root",
		Password:      "toor",
		Name:          "portfolio",
		MaxConnection: 10,
	}

	db, err = database.New(cfg_db)
	if err != nil {
		logger.Fatalf("Database Error: %v\n", err)
	}

	mux = New(db).Handler
	baseURL = "http://localhost:1312/api/customer"
}

func TestCustomerHandler(t *testing.T) {
	// go test ./customer/ -v -run "TestCustomer/get single by id"
	t.Run("get single by id", func(t *testing.T) {
		testScenarios := []testScenarioWithInput[int, string]{
			NewTestScenarioWithInput(1, http.StatusOK, "MARY.SMITH@sakilacustomer.org"),
			NewTestScenarioWithInput(7, http.StatusOK, "MARIA.MILLER@sakilacustomer.org"),
			NewTestScenarioWithInput(13, http.StatusOK, "KAREN.JACKSON@sakilacustomer.org"),
			NewTestScenarioWithInput(20, http.StatusOK, "SHARON.ROBINSON@sakilacustomer.org"),
			NewTestScenarioWithInput(26, http.StatusOK, "JESSICA.HALL@sakilacustomer.org"),
		}

		for _, testScenario := range testScenarios {
			id := testScenario.input
			expected := testScenario.expected
			url := fmt.Sprintf("%s/%d", baseURL, id)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, req)

			actualResponse := recorder.Result()
			actual := actualResponse

			if !assert.Equal(t, expected.statusCode, actual.StatusCode, excpectedStr(expected.statusCode, actual.StatusCode)) {
				return
			}

			actualResponseBody, err := ParseToJSON[responses.GetSingleResponse[readModel]](*actual)
			if !assert.Nil(t, err, excpectedStr(nil, err)) {
				return
			}

			actualData := actualResponseBody.Data
			actualEmail := actualData.Email

			assert.EqualValues(t, expected.data, actualEmail, excpectedStr(expected.data, actualEmail))
		}
	})

	// go test ./customer/ -v -run "TestCustomer/get first customer set"
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
			url := "/api/customer/"
			req := httptest.NewRequest(http.MethodGet, url, nil)
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, req)

			actualResponse := recorder.Result()
			actual := actualResponse

			if !assert.Equal(t, expected.statusCode, actual.StatusCode, excpectedStr(expected.statusCode, actual.StatusCode)) {
				return
			}

			actualResponseBody, err := ParseToJSON[responses.GetMultipleResponse[readModel]](*actual)
			if !assert.Nil(t, err, excpectedStr(nil, err)) {
				return
			}

			actualData := actualResponseBody.Data

			for _, expected := range expected.data {
				for k, v := range expected {
					assert.EqualValues(t, v, actualData[k].Email)
				}
			}
		}
	})

	// go test ./customer/ -v -run "TestCustomer/get previous customer set before a customer with some ids"
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

			actualResponse := recorder.Result()
			actual := actualResponse

			if !assert.Equal(t, expected.statusCode, actual.StatusCode, excpectedStr(expected.statusCode, actual.StatusCode)) {
				return
			}

			actualResponseBody, err := ParseToJSON[responses.GetMultipleResponse[readModel]](*actual)
			if !assert.Nil(t, err, excpectedStr(nil, err)) {
				return
			}

			actualData := actualResponseBody.Data
			for _, expected := range expected.data {
				for k, v := range expected {
					assert.EqualValues(t, v, actualData[k].Email)
				}
			}
		}
	})

	// go test ./customer/ -v -run "TestCustomer/get next customer set after a customer with some ids"
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

			if !assert.Equal(t, expected.statusCode, actual.StatusCode, excpectedStr(expected.statusCode, actual.StatusCode)) {
				return
			}

			actualResponseBody, err := ParseToJSON[responses.GetMultipleResponse[readModel]](*actual)
			if !assert.Nil(t, err, excpectedStr(nil, err)) {
				return
			}

			actualData := actualResponseBody.Data
			for _, expected := range expected.data {
				for k, v := range expected {
					assert.EqualValues(t, v, fmt.Sprintf("%d %s", actualData[k].Id, actualData[k].Email))
				}
			}
		}
	})
}

func TestCustomerProtectedHandler(t *testing.T) {
	// go test ./customer/ -v -run "TestCustomer/create new single customer"
	t.Run("create new single customer", func(t *testing.T) {
		new_customer := createModel{
			FirstName: "Muhammad Miftah",
			LastName:  "Rizki",
			Email:     "muhammadmiftahrizki@gmail.com",
		}
		payload := bytes.NewBuffer(nil)
		json_encoder := json.NewEncoder(payload)
		err := json_encoder.Encode(&new_customer)
		if assert.Nil(t, err, excpectedStr(nil, err)) {
			req := httptest.NewRequest(http.MethodPost, "/api/customer/", payload)
			req.Header.Add("Content-Length", strconv.Itoa(payload.Len()))
			recoder := httptest.NewRecorder()

			mux.ServeHTTP(recoder, req)

			result := recoder.Result()
			defer result.Body.Close()

			assert.Equal(t, http.StatusCreated, result.StatusCode, excpectedStr(http.StatusCreated, result.StatusCode))
		}

		new_customer_str := `{"first_name": "Muhammad Miftah","last_name": "Rizki","email": "muhammadmiftahrizki@gmail.com"}`
		payload2 := bytes.NewBuffer(nil)
		len_payload2, err := payload2.WriteString(new_customer_str)
		if assert.Nil(t, err, excpectedStr(err, nil)) {
			req := httptest.NewRequest(http.MethodPost, "/api/customer/", payload2)
			req.Header.Add("Content-Length", strconv.Itoa(len_payload2))
			res := responses.GetSingleResponse[readModel]{}
			recoder := httptest.NewRecorder()

			mux.ServeHTTP(recoder, req)

			result := recoder.Result()
			defer result.Body.Close()
			assert.Equal(t, http.StatusCreated, result.StatusCode, excpectedStr(http.StatusCreated, result.StatusCode))

			json_decoder := json.NewDecoder(result.Body)

			err := json_decoder.Decode(&res)
			if assert.Nil(t, err, excpectedStr(nil, err)) {
			}
		}
	})

	// go test ./customer/ -v -run "TestCustomer/edit single customer"
	t.Run("edit single customer", func(t *testing.T) {
		// 13 KAREN.JACKSON@sakilacustomer.org

		first_name := "KAREN EDITED"
		last_name := "JACKSON EDITED"
		email := "KAREN.JACKSON.EDITED@sakilacustomer.org"
		customer := customerUpdateModel{
			FirstName: &first_name,
			LastName:  &last_name,
			Email:     &email,
		}
		payload := bytes.NewBuffer(nil)
		json_encoder := json.NewEncoder(payload)
		err := json_encoder.Encode(customer)
		if assert.Nil(t, err, excpectedStr(nil, err)) {
			req := httptest.NewRequest(http.MethodPut, "/api/customer/13", payload)
			req.Header.Add("Content-Length", strconv.Itoa(payload.Len()))
			res := responses.GetSingleResponse[readModel]{}
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, req)

			result := recorder.Result()
			defer result.Body.Close()
			assert.Equal(t, http.StatusOK, result.StatusCode, excpectedStr(http.StatusOK, result.StatusCode))

			json_decoder := json.NewDecoder(result.Body)

			err := json_decoder.Decode(&res)
			if assert.Nil(t, err, excpectedStr(nil, err)) {
			}
		}

		first_name = "KAREN"
		last_name = "JACKSON"
		email = "KAREN.JACKSON@sakilacustomer.org"
		customer = customerUpdateModel{
			FirstName: &first_name,
			LastName:  &last_name,
			Email:     &email,
		}
		payload = bytes.NewBuffer(nil)
		json_encoder = json.NewEncoder(payload)
		err = json_encoder.Encode(customer)
		if assert.Nil(t, err, excpectedStr(nil, err)) {
			req := httptest.NewRequest(http.MethodPut, "/api/customer/13", payload)
			req.Header.Add("Content-Length", strconv.Itoa(payload.Len()))
			res := responses.GetSingleResponse[readModel]{}
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, req)

			result := recorder.Result()
			defer result.Body.Close()
			assert.Equal(t, http.StatusOK, result.StatusCode, excpectedStr(http.StatusOK, result.StatusCode))

			json_decoder := json.NewDecoder(result.Body)

			err := json_decoder.Decode(&res)
			if assert.Nil(t, err, excpectedStr(nil, err)) {
			}
		}
	})

	// go test ./customer/ -v -run "TestCustomer/delete single customer by its id"
	t.Run("delete single customer by its id", func(t *testing.T) {
		testScenarios := []testScenarioWithInput[int, string]{
			NewTestScenarioWithInput(1, http.StatusNoContent, http.StatusText(http.StatusNoContent)),
		}

		for _, testScenario := range testScenarios {
			id := testScenario.input
			expected := testScenario.expected

			url := fmt.Sprintf("/api/customer/%d", id)
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ""))
			recoder := httptest.NewRecorder()

			mux.ServeHTTP(recoder, req)

			actualResponse := recoder.Result()
			actual := actualResponse

			if !assert.Equal(t, expected.statusCode, actual.StatusCode, excpectedStr(expected.statusCode, actual.StatusCode)) {
				return
			}

			actualResponseBody, err := ParseToJSON[responses.GetSingleResponse[readModel]](*actual)
			if !assert.Nil(t, err, excpectedStr(nil, err)) {
				return
			}

			assert.Nil(t, actualResponseBody, excpectedStr(nil, actualResponseBody))
		}
	})
}
