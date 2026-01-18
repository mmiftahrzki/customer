package customer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
		log:     logger.GetLogger().WithField("component", "customerHandler"),
	}

	return handler
}

func timeoutMiddleware(w http.ResponseWriter, r *http.Request) {
	var queryTimeout string
	var timeoutValue int
	var timeoutDuration time.Duration

	queryTimeout = r.URL.Query().Get("timeout")

	if queryTimeout != "" {
		var strConvErr error

		timeoutValue, strConvErr = strconv.Atoi(queryTimeout)
		if strConvErr != nil {
			responses.Error(w, http.StatusUnprocessableEntity, "invalid timeout duration value")

			return
		}

		timeoutDuration = time.Duration(timeoutValue) * time.Millisecond

		if timeoutDuration > time.Duration(0) {
			if timeoutDuration > 30000*time.Millisecond {
				timeoutDuration = 30000 * time.Millisecond
			}

			ctxWithTimeout, cancel := context.WithTimeout(r.Context(), timeoutDuration)
			defer cancel()

			r = r.WithContext(ctxWithTimeout)
		}
	}
}

func (h *handler) PostSingle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	contentLenStr := r.Header.Get("Content-Length")
	contentLen, err := strconv.Atoi(contentLenStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if contentLen == 0 {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if contentLen > 2048 {
		responses.Error(w, http.StatusRequestEntityTooLarge, "content length cannot be more than 2048")

		return
	}

	payload := modelCreate{}
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
	var res responses.GetMultipleResponse[modelRead]

	customers, svcErr := h.service.GetMultiple(r.Context())
	if svcErr != nil {
		if errors.Is(svcErr, context.DeadlineExceeded) {
			responses.Error(w, http.StatusServiceUnavailable, "server took too long to respond")

			return
		}

		h.log.Error(svcErr)

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
	var res responses.GetSingleResponse[modelRead]

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
	var res responses.GetMultipleResponse[modelRead]

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
	var res responses.GetMultipleResponse[modelRead]

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
	defer r.Body.Close()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.log.Error(err)

		responses.Error(w, http.StatusBadRequest, "invalid id")

		return
	}

	payload := modelUpdate{}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.log.Error(err)

		w.WriteHeader(http.StatusBadRequest)

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
	defer r.Body.Close()

	var err error

	customerId, err := strconv.Atoi(r.PathValue("customer_id"))
	if err != nil {
		h.log.Error(err)

		responses.Error(w, http.StatusBadRequest, "invalid id")

		return
	}

	addressId, err := strconv.Atoi(r.PathValue("address_id"))
	if err != nil {
		h.log.Error(err)

		responses.Error(w, http.StatusBadRequest, "invalid id")

		return
	}

	payload := address.ModelUpdate{}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.log.Error(err)

		responses.Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))

		return
	}

	err = payload.Validate()
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())

		return
	}

	err = h.service.ModifySingleAddressById(r.Context(), customerId, uint16(addressId), payload)
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
