package customer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mmiftahrzki/customer/app/model"
	"github.com/mmiftahrzki/customer/customer/address"
	"github.com/mmiftahrzki/customer/logger"
	"github.com/sirupsen/logrus"
)

type handler struct {
	service *service
	log     *logrus.Entry
}

func NewHandler(svc *service) *handler {
	return &handler{
		service: svc,
		log:     logger.GetLogger().WithField("component", "customer_handler"),
	}
}

func (h *handler) PostSingle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := model.NewReadModel()
	payload := CustomerCreateModel{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		h.log.Error(err)

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err = h.service.CreateNewSingle(r.Context(), payload); err != nil {
		switch err {
		case ErrCustomerAlreadyExists:
			response.Message = fmt.Sprintf("customer with email: %s already exists", payload.Email)

			w.WriteHeader(http.StatusConflict)
			w.Write(response.ToJson())

			return
		default:
			h.log.Error(err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))

			return
		}
	}

	response.Message = "success create new customer"

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response.ToJson()))
}

func (h *handler) GetMultiple(w http.ResponseWriter, r *http.Request) {
	if fmt.Sprintf("%s %s", r.Method, r.RequestURI) != r.Pattern {
		http.NotFound(w, r)

		return
	}

	response := model.NewReadModel()

	customers, err := h.service.GetMultiple(r.Context())
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
	res := model.NewReadModel()

	id, err := parseIdStringToInt(r.PathValue("id"))
	if err != nil {
		res.Message = "invalid id"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	// customer, err := h.repo.SelectSingleById(r.Context(), id)
	customer, err := h.service.GetSingleById(r.Context(), id)
	if err != nil {
		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))

		return
	}

	empty_customer := customerReadModel{}
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
	res := model.NewReadModel()

	id, err := parseIdStringToInt(r.PathValue("id"))
	if err != nil {
		h.log.Error(err)

		res.Message = "invalid id"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	customers, err := h.service.GetMultipleNext(r.Context(), id)
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

// func (h *handler) GetMultiplePrev(w http.ResponseWriter, r *http.Request) {
// 	res := model.NewReadModel()

// 	id, err := parseIdStringToInt(r.PathValue("id"))
// 	if err != nil {
// 		h.log.Error(err)

// 		res.Message = "invalid id"

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write(res.ToJson())

// 		return
// 	}

// 	customers, err := h.service.GetPrev(r.Context(), id)
// 	if err != nil {
// 		h.log.Error(err)

// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))

// 		return
// 	}

// 	if len(customers) == LIMIT+1 {
// 		res.Data["__prev"] = fmt.Sprintf("/api/customer/%d/prev", customers[1].Id)

// 		customers = customers[1 : LIMIT+1]
// 	}

// 	if len(customers) > 0 {
// 		res.Data["__next"] = fmt.Sprintf("/api/customer/%d/next", customers[len(customers)-1].Id)
// 	}

// 	res.Data["customers"] = customers
// 	res.Message = "berhasil mendapatkan data customer"

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(res.ToJson())
// }

func (h *handler) PutSingleById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var err error = nil
	res := model.NewReadModel()

	id, err := parseIdStringToInt(r.PathValue("id"))
	if err != nil {
		h.log.Error(err)

		res.Message = "invalid id"

		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	payload := CustomerUpdateModel{}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.log.Error(err)

		res.Message = http.StatusText(http.StatusBadRequest)

		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	err = ValidateCustomerUpdateModel(payload)
	if err != nil {
		res.Message = err.Error()

		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	err = h.service.ModifySingleById(r.Context(), id, payload)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			res.Message = http.StatusText(http.StatusUnprocessableEntity)

			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write(res.ToJson())

			return
		} else {
			h.log.Error(err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))

			return
		}
	}

	res.Message = "success modify customer"
	w.WriteHeader(http.StatusOK)
	w.Write(res.ToJson())
}

func (h *handler) DeleteSingleById(w http.ResponseWriter, r *http.Request) {
	res := model.NewReadModel()

	id, err := parseIdStringToInt(r.PathValue("id"))
	if err != nil {
		res.Message = "invalid id"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	err = h.service.DeleteSingleById(r.Context(), id)
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
	w.Header().Set("Content-Type", "application/json")

	res := model.NewReadModel()
	var err error

	customer_id, err := parseIdStringToInt(r.PathValue("customer_id"))
	if err != nil {
		h.log.Error(err)

		res.Message = "invalid id"

		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	address_id, err := parseIdStringToUint16(r.PathValue("address_id"))
	if err != nil {
		h.log.Error(err)

		res.Message = "invalid id"

		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	payload := address.AddressUpdateModel{}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.log.Error(err)

		res.Message = http.StatusText(http.StatusBadRequest)

		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	err = address.ValidateAddressUpdateModel(payload)
	if err != nil {
		res.Message = err.Error()

		w.WriteHeader(http.StatusBadRequest)
		w.Write(res.ToJson())

		return
	}

	err = h.service.ModifySingleAddressById(r.Context(), customer_id, address_id, payload)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) || errors.Is(err, ErrInvalidCustomerAddressMismatch) {
			res.Message = http.StatusText(http.StatusUnprocessableEntity)

			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write(res.ToJson())

			return
		}

		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))

		return
	}

	res.Message = "success modify customer's address"

	w.WriteHeader(http.StatusOK)
	w.Write(res.ToJson())
}

func parseIdStringToInt(id_str string) (int, error) {
	id_int, err := strconv.ParseInt(id_str, 10, 0)
	if err != nil {
		return 0, err
	}

	return int(id_int), nil
}

func parseIdStringToUint16(id_str string) (uint16, error) {
	id_int, err := strconv.ParseInt(id_str, 10, 0)
	if err != nil {
		return 0, err
	}

	return uint16(id_int), nil
}
