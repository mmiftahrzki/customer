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
	"strings"
	"testing"

	"github.com/mmiftahrzki/customer/auth"
	"github.com/mmiftahrzki/customer/config"
	"github.com/mmiftahrzki/customer/database"
	"github.com/mmiftahrzki/customer/logger"
	"github.com/mmiftahrzki/customer/responses"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB

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
}

func excpectedStr(expected, got any) string {
	return fmt.Sprintf("Expected: %v but got: %v instead.", expected, got)
}

func TestCustomer(t *testing.T) {
	customerHandler := NewHandler(db)

	// go test ./customer/ -v -run "TestCustomer/get single customer by its id"
	t.Run("get single customer by its id", func(t *testing.T) {
		expecteds := []string{
			"1 MARY.SMITH@sakilacustomer.org",
			"7 MARIA.MILLER@sakilacustomer.org",
			"13 KAREN.JACKSON@sakilacustomer.org",
			"20 SHARON.ROBINSON@sakilacustomer.org",
			"26 JESSICA.HALL@sakilacustomer.org",
		}

		for _, expected := range expecteds {
			id := strings.Split(expected, " ")[0]
			url := fmt.Sprintf("http://localhost:1312/api/customer/%s", id)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			res := responses.GetSingleResponse[customerReadModel]{}

			recorder := httptest.NewRecorder()

			customerHandler.GetSingleById(recorder, req)

			result := recorder.Result()
			defer result.Body.Close()

			assert.Equal(t, http.StatusOK, result.StatusCode, excpectedStr(http.StatusOK, result.StatusCode))

			body, err := io.ReadAll(result.Body)

			if assert.Nil(t, err, excpectedStr(nil, err)) {
				bytes_reader := bytes.NewReader(body)
				json_decoder := json.NewDecoder(bytes_reader)
				err := json_decoder.Decode(&res)

				if assert.Nil(t, err, excpectedStr(nil, err)) {
					actual := fmt.Sprintf("%d %s", res.Data.Id, res.Data.Email)
					assert.EqualValues(t, expected, actual, excpectedStr(expected, actual))
				}
			}
		}
	})

	// go test ./customer/ -v -run "TestCustomer/get first customer set"
	t.Run("get first customer set", func(t *testing.T) {
		expecteds := []string{
			"1 MARY.SMITH@sakilacustomer.org",
			"7 MARIA.MILLER@sakilacustomer.org",
			"13 KAREN.JACKSON@sakilacustomer.org",
			"20 SHARON.ROBINSON@sakilacustomer.org",
			"26 JESSICA.HALL@sakilacustomer.org",
		}
		url := "/api/customer/"
		req := httptest.NewRequest(http.MethodGet, url, nil)
		res := responses.GetMultipleResponse[customerReadModel]{}
		recorder := httptest.NewRecorder()

		customerHandler.GetMultiple(recorder, req)

		result := recorder.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusOK, result.StatusCode, excpectedStr(http.StatusOK, result.StatusCode))

		body_bytes, err := io.ReadAll(result.Body)

		if assert.Nil(t, err, excpectedStr(nil, err)) {
			bytes_reader := bytes.NewReader(body_bytes)
			json_decoder := json.NewDecoder(bytes_reader)
			err := json_decoder.Decode(&res)
			if assert.Nil(t, err, excpectedStr(nil, err)) {
				assert.EqualValues(t, expecteds[0], fmt.Sprintf("%d %s", res.Data[0].Id, res.Data[0].Email))
				assert.EqualValues(t, expecteds[1], fmt.Sprintf("%d %s", res.Data[6].Id, res.Data[6].Email))
				assert.EqualValues(t, expecteds[2], fmt.Sprintf("%d %s", res.Data[12].Id, res.Data[12].Email))
				assert.EqualValues(t, expecteds[3], fmt.Sprintf("%d %s", res.Data[18].Id, res.Data[18].Email))
				assert.EqualValues(t, expecteds[4], fmt.Sprintf("%d %s", res.Data[24].Id, res.Data[24].Email))
			}
		}
	})

	// go test ./customer/ -v -run "TestCustomer/get previous customer set before a customer with some ids"
	t.Run("get previous customer set before a customer with some ids", func(t *testing.T) {
		expecteds := map[int][]string{
			51: {
				"JESSICA.HALL@sakilacustomer.org",
				"AMY.LOPEZ@sakilacustomer.org",
				"MARTHA.GONZALEZ@sakilacustomer.org",
				"MARIE.TURNER@sakilacustomer.org",
				"DIANE.COLLINS@sakilacustomer.org",
			},
			77: {
				"ALICE.STEWART@sakilacustomer.org",
				"EVELYN.MORGAN@sakilacustomer.org",
				"ASHLEY.RICHARDSON@sakilacustomer.org",
				"CHRISTINA.RAMIREZ@sakilacustomer.org",
				"IRENE.PRICE@sakilacustomer.org",
			},
			102: {
				"JANE.BENNETT@sakilacustomer.org",
				"LOUISE.JENKINS@sakilacustomer.org",
				"JULIA.FLORES@sakilacustomer.org",
				"PAULA.BRYANT@sakilacustomer.org",
				"PEGGY.MYERS@sakilacustomer.org",
			},
			128: {
				"CRYSTAL.FORD@sakilacustomer.org",
				"TRACY.COLE@sakilacustomer.org",
				"GRACE.ELLIS@sakilacustomer.org",
				"SYLVIA.ORTIZ@sakilacustomer.org",
				"ELAINE.STEVENS@sakilacustomer.org",
			},
		}

		for k, v := range expecteds {
			url := fmt.Sprintf("/api/customer/%d/prev", k)
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			response := responses.GetMultipleResponse[customerReadModel]{}
			recorder := httptest.NewRecorder()

			customerHandler.GetMultiplePrev(recorder, req)
			assert.Equal(t, http.StatusOK, recorder.Code, excpectedStr(http.StatusOK, recorder.Code))

			body_bytes := recorder.Body.Bytes()
			bytes_reader := bytes.NewReader(body_bytes)
			json_decoder := json.NewDecoder(bytes_reader)
			err := json_decoder.Decode(&response)
			if assert.NoError(t, err, excpectedStr(nil, err)) {
				actual_customers := response.Data

				assert.EqualValues(t, v[0], actual_customers[0].Email, excpectedStr(v[0], actual_customers[0].Email))
				assert.EqualValues(t, v[1], actual_customers[6].Email, excpectedStr(v[1], actual_customers[6].Email))
				assert.EqualValues(t, v[2], actual_customers[12].Email, excpectedStr(v[2], actual_customers[12].Email))
				assert.EqualValues(t, v[3], actual_customers[18].Email, excpectedStr(v[3], actual_customers[18].Email))
				assert.EqualValues(t, v[4], actual_customers[24].Email, excpectedStr(v[4], actual_customers[24].Email))
			}
		}
	})

	// go test ./customer/ -v -run "TestCustomer/get next customer set after a customer with some ids"
	t.Run("get next customer set after a customer with some ids", func(t *testing.T) {
		expecteds := map[int][]string{
			26: {
				"27 SHIRLEY.ALLEN@sakilacustomer.org",
				"33 ANNA.HILL@sakilacustomer.org",
				"39 DEBRA.NELSON@sakilacustomer.org",
				"45 JANET.PHILLIPS@sakilacustomer.org",
				"51 ALICE.STEWART@sakilacustomer.org",
			},
			51: {
				"52 JULIE.SANCHEZ@sakilacustomer.org",
				"58 JEAN.BELL@sakilacustomer.org",
				"65 ROSE.HOWARD@sakilacustomer.org",
				"71 KATHY.JAMES@sakilacustomer.org",
				"77 JANE.BENNETT@sakilacustomer.org",
			},
			77: {
				"78 LORI.WOOD@sakilacustomer.org",
				"84 SARA.PERRY@sakilacustomer.org",
				"90 RUBY.WASHINGTON@sakilacustomer.org",
				"96 DIANA.ALEXANDER@sakilacustomer.org",
				"102 CRYSTAL.FORD@sakilacustomer.org",
			},
			102: {
				"103 GLADYS.HAMILTON@sakilacustomer.org",
				"109 EDNA.WEST@sakilacustomer.org",
				"115 WENDY.HARRISON@sakilacustomer.org",
				"121 JOSEPHINE.GOMEZ@sakilacustomer.org",
				"128 MARJORIE.TUCKER@sakilacustomer.org",
			},
		}

		for k, v := range expecteds {
			url := fmt.Sprintf("/api/customer/%d/next", k)
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			res := responses.GetMultipleResponse[customerReadModel]{}
			recorder := httptest.NewRecorder()

			customerHandler.GetMultipleNext(recorder, req)
			assert.Equal(t, http.StatusOK, recorder.Code, excpectedStr(http.StatusOK, recorder.Code))

			body_bytes := recorder.Body.Bytes()
			bytes_reader := bytes.NewReader(body_bytes)
			json_decoder := json.NewDecoder(bytes_reader)
			err := json_decoder.Decode(&res)
			if assert.NoError(t, err, excpectedStr(nil, err)) {
				actual_customers := res.Data

				ex := fmt.Sprintf("%d %s", actual_customers[0].Id, actual_customers[0].Email)
				ex6 := fmt.Sprintf("%d %s", actual_customers[6].Id, actual_customers[6].Email)
				ex12 := fmt.Sprintf("%d %s", actual_customers[12].Id, actual_customers[12].Email)
				ex18 := fmt.Sprintf("%d %s", actual_customers[18].Id, actual_customers[18].Email)
				ex24 := fmt.Sprintf("%d %s", actual_customers[24].Id, actual_customers[24].Email)

				assert.EqualValues(t, v[0], ex, excpectedStr(v[0], ex))
				assert.EqualValues(t, v[1], ex6, excpectedStr(v[1], ex6))
				assert.EqualValues(t, v[2], ex12, excpectedStr(v[2], ex12))
				assert.EqualValues(t, v[3], ex18, excpectedStr(v[3], ex18))
				assert.EqualValues(t, v[4], ex24, excpectedStr(v[4], ex24))
			}
		}
	})

	// go test ./customer/ -v -run "TestCustomer/create new single customer"
	t.Run("create new single customer", func(t *testing.T) {
		new_customer := customerCreateModel{
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

			customerHandler.PostSingle(recoder, req)

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
			res := responses.GetSingleResponse[customerReadModel]{}
			recoder := httptest.NewRecorder()

			customerHandler.PostSingle(recoder, req)

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
			res := responses.GetSingleResponse[customerReadModel]{}
			recorder := httptest.NewRecorder()

			customerHandler.PutSingleById(recorder, req)

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
			res := responses.GetSingleResponse[customerReadModel]{}
			recorder := httptest.NewRecorder()

			customerHandler.PutSingleById(recorder, req)

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
		req := httptest.NewRequest(http.MethodDelete, "/api/customer/13", nil)
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InNoaXJvaGlnZTY1QHJvY2tldG1haWwuY29tIiwiZXhwIjoxNzYyNjExNzQwfQ.4vM_xjU-_o8csybUPIJLmqVCCwO5DDcDCh-aClY417Y")
		res := responses.GetSingleResponse[customerReadModel]{}
		recoder := httptest.NewRecorder()

		authMiddleware := auth.NewMiddleware(auth.NewService())
		deleteSingleById := authMiddleware.VerifyJWT(customerHandler.DeleteSingleById)
		deleteSingleById(recoder, req)

		result := recoder.Result()
		assert.Equal(t, http.StatusNoContent, result.StatusCode, excpectedStr(http.StatusNoContent, result.StatusCode))

		result_body := result.Body
		defer result.Body.Close()

		json_decoder := json.NewDecoder(result_body)
		err := json_decoder.Decode(&res)
		if assert.Nil(t, err, excpectedStr(nil, err)) {
			assert.Nil(t, result_body, excpectedStr(nil, result_body))
		}
	})
}
