package customer

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/mmiftahrzki/customer/logger"
	"github.com/mmiftahrzki/customer/responses"
	"github.com/sirupsen/logrus"
)

type handler struct {
	service Service
	log     *logrus.Entry
}

func newHandler(svc Service) handler {
	handler := handler{
		service: svc,
		log:     logger.GetLogger().WithField("component", "customerHandler"),
	}

	return handler
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

	err = h.service.RegisterCustomer(r.Context(), payload)
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
	var res responses.GetMultipleResponse[ModelRead]
	timeout := 15 * time.Second

	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()

	queryParams := r.URL.Query()

	sortBy, err := ParseOrderColumn(queryParams.Get("sortBy"))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())

		return
	}

	sortDirection, err := ParseOrderDirection(queryParams.Get("order"))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())

		return
	}

	take, err := ParseTake(queryParams.Get("take"))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())

		return
	}

	page, err := ParsePage(queryParams.Get("page"))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())

		return
	}

	orderBy := OrderBy{
		Column:    sortBy,
		Direction: sortDirection,
	}

	customers, err := h.service.CustomerList(ctx, take, page, orderBy)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			responses.Error(w, http.StatusServiceUnavailable, "server took too long to respond")

			return
		}

		h.log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	res.Data = customers

	responses.WithJson(w, http.StatusOK, res)

	h.log.Info("customers data retrieved successfully")
}

func (h *handler) GetSingleById(w http.ResponseWriter, r *http.Request) {
	var res responses.GetSingleResponse[ModelRead]

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.log.Error(err)

		responses.Error(w, http.StatusBadRequest, "invalid id")

		return
	}

	customer, err := h.service.CustomerDetails(r.Context(), id)
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

func (h *handler) PutSingleById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.log.Error(err)

		responses.Error(w, http.StatusBadRequest, "invalid id")

		return
	}

	var payload modelUpdate
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.log.Error(err)

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = payload.Validate()
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
