package customer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mmiftahrzki/customer/customer/address"
	"github.com/mmiftahrzki/customer/logger"
	"github.com/mmiftahrzki/customer/responses"
	"github.com/sirupsen/logrus"
)

type handler struct {
	service service
	log     *logrus.Entry
}

func newHandler(svc service) handler {
	handler := handler{
		service: svc,
		log:     logger.GetLogger().WithField("component", "customer_handler"),
	}

	return handler
}

func (h *handler) PostSingle(w http.ResponseWriter, r *http.Request) {
	content_length_str := r.Header.Get("Content-Length")
	content_length, err := strconv.Atoi(content_length_str)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if content_length == 0 {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if content_length > 2048 {
		responses.Error(w, http.StatusRequestEntityTooLarge, "content length cannot be more than 2048")

		return
	}

	payload := customerCreateModel{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		h.log.Error(err)

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = h.service.CreateNewSingle(r.Context(), payload)
	if err != nil {
		if errors.Is(err, errCustomerAlreadyExists) {
			responses.WithJson(w, http.StatusConflict, err.Error())

			return
		}

		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handler) GetMultiple(w http.ResponseWriter, r *http.Request) {
	var res responses.GetMultipleResponse[customerReadModel]

	if fmt.Sprintf("%s %s", r.Method, r.RequestURI) != r.Pattern {
		http.NotFound(w, r)

		return
	}

	customers, err := h.service.GetMultiple(r.Context())
	if err != nil {
		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if len(customers) == limit+1 {
		res.Next = fmt.Sprintf("/api/customer/%d/next", customers[limit-1].Id)

		customers = customers[:limit]
	}
	res.Data = customers

	responses.WithJson(w, http.StatusOK, res)

	h.log.Info("customers data retrieved successfully")
}

func (h *handler) GetSingleById(w http.ResponseWriter, r *http.Request) {
	var res responses.GetSingleResponse[customerReadModel]

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.log.Error(err)

		responses.Error(w, http.StatusBadRequest, "invalid id")

		return
	}

	customer, err := h.service.GetSingleById(r.Context(), id)
	if err != nil {
		if errors.Is(err, errCustomerNotFound) {
			responses.Error(w, http.StatusNotFound, err.Error())

			return
		}

		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	res.Data = customer

	responses.WithJson(w, http.StatusOK, res)
}

func (h *handler) GetMultipleNext(w http.ResponseWriter, r *http.Request) {
	var res responses.GetMultipleResponse[customerReadModel]

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.log.Error(err)

		responses.Error(w, http.StatusBadRequest, "invalid id")

		return
	}

	customers, err := h.service.GetMultipleNext(r.Context(), id)
	if err != nil {
		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if len(customers) == limit+1 {
		res.Prev = fmt.Sprintf("/api/customer/%d/prev", customers[0].Id)
		res.Next = fmt.Sprintf("/api/customer/%d/next", customers[limit-1].Id)

		customers = customers[:limit]
	}
	res.Data = customers

	responses.WithJson(w, http.StatusOK, res)
}

func (h *handler) GetMultiplePrev(w http.ResponseWriter, r *http.Request) {
	var res responses.GetMultipleResponse[customerReadModel]

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.log.Error(err)

		responses.Error(w, http.StatusBadRequest, "invalid id")

		return
	}

	customers, err := h.service.GetMultiplePrev(r.Context(), id)
	if err != nil {
		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if len(customers) == limit {
		res.Prev = fmt.Sprintf("/api/customer/%d/prev", customers[0].Id)
	}

	if len(customers) > 0 {
		res.Next = fmt.Sprintf("/api/customer/%d/next", customers[len(customers)-1].Id)
	}

	if len(customers) < limit {
		http.Redirect(w, r, "/api/customer", http.StatusSeeOther)

		return
	}

	res.Data = customers

	responses.WithJson(w, http.StatusOK, res)
}

func (h *handler) PutSingleById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.log.Error(err)

		responses.Error(w, http.StatusBadRequest, "invalid id")

		return
	}

	payload := customerUpdateModel{}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.log.Error(err)

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = validatecustomerUpdateModel(payload)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())

		return
	}

	err = h.service.ModifySingleById(r.Context(), id, payload)
	if err != nil {
		if errors.Is(err, errCustomerNotFound) {
			responses.Error(w, http.StatusUnprocessableEntity, err.Error())

			return
		}

		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) DeleteSingleById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid id")

		return
	}

	err = h.service.DeleteSingleById(r.Context(), id)
	if err != nil {
		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) GetSingleAndUpdateAddressById(w http.ResponseWriter, r *http.Request) {
	var err error

	customer_id, err := strconv.Atoi(r.PathValue("customer_id"))
	if err != nil {
		h.log.Error(err)

		responses.Error(w, http.StatusBadRequest, "invalid id")

		return
	}

	address_id, err := strconv.Atoi(r.PathValue("address_id"))
	if err != nil {
		h.log.Error(err)

		responses.Error(w, http.StatusBadRequest, "invalid id")

		return
	}

	payload := address.AddressUpdateModel{}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.log.Error(err)

		responses.Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))

		return
	}

	err = address.ValidateAddressUpdateModel(payload)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())

		return
	}

	err = h.service.ModifySingleAddressById(r.Context(), customer_id, uint16(address_id), payload)
	if err != nil {
		if errors.Is(err, errCustomerNotFound) || errors.Is(err, errInvalidCustomerAddressMismatch) {
			responses.WithJson(w, http.StatusUnprocessableEntity, err.Error())

			return
		}

		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
