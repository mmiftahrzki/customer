package customer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/mmiftahrzki/customer/customer/address"
	"github.com/mmiftahrzki/customer/logger"
	"github.com/mmiftahrzki/customer/response"
	"github.com/sirupsen/logrus"
)

type handler struct {
	repo *repo
	log  *logrus.Entry
}

func newHandler(repo *repo) *handler {
	return &handler{
		repo: repo,
		log:  logger.GetLogger().WithField("component", "customer_handler"),
	}
}

func (h *handler) PostSingle(w http.ResponseWriter, r *http.Request) {
	res := response.New()

	// new_customer, err := validation.ExtractCustomerFromContext(r.Context())

	payload := CustomerCreateModel{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		logger.Logger.Error(err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = h.repo.InsertSingle(r.Context(), payload)
	if err != nil {
		h.log.Error(err)

		mysql_error, ok := err.(*mysql.MySQLError)
		if ok {
			if mysql_error.Number == 1062 {
				res.Message = fmt.Sprintf("customer dengan email: %s sudah ada", payload.Email)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				w.Write(res.ToJson())

				return
			}
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))

		return
	}

	res.Message = "berhasil membuat customer baru"

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res.ToJson()))
}

func (h *handler) GetMultiple(w http.ResponseWriter, r *http.Request) {
	if fmt.Sprintf("%s %s", r.Method, r.RequestURI) != r.Pattern {
		http.NotFound(w, r)

		return
	}

	response := response.New()

	customers, err := h.repo.SelectAll(r.Context())
	if err != nil {
		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(response.ToJson()))

		return
	}

	if len(customers) == LIMIT+1 {
		response.Data["__next"] = fmt.Sprintf("/api/customer/%d/next", customers[LIMIT-1].Id)

		customers = customers[:LIMIT]
	}

	response.Data["customers"] = customers
	response.Message = "berhasil mendapatkan data customer"

	w.Header().Set("Content-Type", "application/json")
	w.Write(response.ToJson())
	h.log.Info("customers data retrieved successfully")
}

func (h *handler) GetSingleById(w http.ResponseWriter, r *http.Request) {
	res := response.New()

	id, err := parseIdStringToUint(r.PathValue("id"))
	if err != nil {
		res.Message = "id tidak valid"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	customer, err := h.repo.SelectSingleById(r.Context(), id)
	if err != nil {
		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))

		return
	}

	empty_customer := Customer{}
	if customer == empty_customer {
		res.Message = fmt.Sprintf("customer dengan id: %d tidak ditemukan", id)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(res.ToJson())

		return
	}

	res.Message = "berhasil mendapatkan data customer"
	res.Data["customer"] = customer

	w.Header().Set("Content-Type", "application/json")
	w.Write(res.ToJson())
}

func (h *handler) GetMultipleNext(w http.ResponseWriter, r *http.Request) {
	res := response.New()

	id, err := parseIdStringToUint(r.PathValue("id"))
	if err != nil {
		h.log.Error(err)

		res.Message = "id tidak valid"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	customer, err := h.repo.SelectSingleById(r.Context(), id)
	if err != nil {
		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))

		return
	}

	if reflect.ValueOf(customer).IsZero() {
		res.Message = http.StatusText(http.StatusNotFound)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(res.ToJson())

		return
	}

	customers, err := h.repo.SelectAllNext(r.Context(), customer)
	if err != nil {
		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))

		return
	}

	if len(customers) == LIMIT+1 {
		res.Data["__prev"] = fmt.Sprintf("/api/customer/%d/prev", customers[0].Id)
		res.Data["__next"] = fmt.Sprintf("/api/customer/%d/next", customers[LIMIT-1].Id)

		customers = customers[:LIMIT]
	}

	res.Data["customers"] = customers
	res.Message = "berhasil mendapatkan data customer"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res.ToJson())
}

func (h *handler) GetMultiplePrev(w http.ResponseWriter, r *http.Request) {
	res := response.New()

	id, err := parseIdStringToUint(r.PathValue("id"))
	if err != nil {
		h.log.Error(err)

		res.Message = "id tidak valid"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	customer, err := h.repo.SelectSingleById(r.Context(), id)
	if err != nil {
		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))

		return
	}

	if reflect.ValueOf(customer).IsZero() {
		res.Message = http.StatusText(http.StatusNotFound)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(res.ToJson())

		return
	}

	customers, err := h.repo.SelectAllPrev(r.Context(), customer)
	if err != nil {
		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))

		return
	}

	if len(customers) == LIMIT+1 {
		res.Data["__prev"] = fmt.Sprintf("/api/customer/%d/prev", customers[1].Id)

		customers = customers[1 : LIMIT+1]
	}

	if len(customers) > 0 {
		res.Data["__next"] = fmt.Sprintf("/api/customer/%d/next", customers[len(customers)-1].Id)
	}

	res.Data["customers"] = customers
	res.Message = "berhasil mendapatkan data customer"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res.ToJson())
}

func (h *handler) PutSingleById(w http.ResponseWriter, r *http.Request) {
	var err error
	res := response.New()

	id, err := parseIdStringToUint(r.PathValue("id"))
	if err != nil {
		h.log.Error(err)

		res.Message = "id tidak valid"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	payload := CustomerUpdateModel{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		h.log.Error(err)

		res.Message = http.StatusText(http.StatusBadRequest)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	err = h.repo.UpdateSingleById(r.Context(), id, payload)
	if err != nil {
		h.log.Error(err)

		mysql_error, ok := err.(*mysql.MySQLError)
		if ok {
			if mysql_error.Number == 1292 {
				res.Message = http.StatusText(http.StatusBadRequest)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(res.ToJson())
			}

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))

		return
	}

	res.Message = "berhasil memperbarui data customer"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res.ToJson())
}

func (h *handler) DeleteSingleById(w http.ResponseWriter, r *http.Request) {
	res := response.New()

	id, err := parseIdStringToUint(r.PathValue("id"))
	if err != nil {
		res.Message = "id tidak valid"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	err = h.repo.DeleteSingleById(r.Context(), id)
	if err != nil {
		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) GetSingleAndUpdateAddressById(w http.ResponseWriter, r *http.Request) {
	res := response.New()
	var err error

	id, err := parseIdStringToUint(r.PathValue("id"))
	if err != nil {
		h.log.Error(err)

		res.Message = "id tidak valid"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	customer, err := h.repo.SelectSingleById(r.Context(), id)
	if err != nil {
		h.log.Error(err)

		res.Message = http.StatusText(http.StatusNotFound)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(res.ToJson())

		return
	}

	payload := address.AddressUpdateModel{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		h.log.Error(err)

		res.Message = http.StatusText(http.StatusBadRequest)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	err = h.repo.UpdateSingleAddressByCustomerId(r.Context(), customer.Address.Id, payload)
	if err != nil {
		h.log.Error(err)

		mysql_error, ok := err.(*mysql.MySQLError)
		if ok {
			if mysql_error.Number == 1292 {
				res.Message = http.StatusText(http.StatusBadRequest)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(res.ToJson())
			}

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))

		return
	}

	res.Message = "berhasil memperbarui data customer"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res.ToJson())
}

func parseIdStringToUint(id_str string) (uint, error) {
	id_uint64, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id_uint64), nil
}
