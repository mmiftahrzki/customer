package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mmiftahrzki/customer/logger"
	"github.com/mmiftahrzki/customer/responses"
	"github.com/sirupsen/logrus"
)

type handler struct {
	service service
	log     *logrus.Entry
}

func newHandler(svc service) handler {
	return handler{
		service: svc,
		log:     logger.GetLogger().WithField("component", "auth/handler"),
	}
}

func (h *handler) CreateAuthToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	payload := AuthCreateModel{}

	json_decoder := json.NewDecoder(r.Body)
	err := json_decoder.Decode(&payload)
	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	token, err := h.service.generateJWT(payload)
	if err != nil {
		log.Println(err)

		responses.Error(w, http.StatusInternalServerError, "errors occured when generating JWT")

		return
	}

	res := AuthReadModel{Token: token}

	responses.WithJson(w, http.StatusOK, res)
}
