package customer

import (
	"net/http"
)

func newMux(h *handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/customer/", h.GetMultiple)
	mux.HandleFunc("POST /api/customer/", h.PostSingle)
	mux.HandleFunc("GET /api/customer/{id}", h.GetSingleById)
	mux.HandleFunc("GET /api/customer/{id}/prev", h.GetMultiplePrev)
	mux.HandleFunc("GET /api/customer/{id}/next", h.GetMultipleNext)
	mux.HandleFunc("PUT /api/customer/{id}", h.PutSingleById)
	mux.HandleFunc("DELETE /api/customer/{id}", h.DeleteSingleById)
	mux.HandleFunc("PATCH /api/customer/{id}/address", h.GetSingleAndUpdateAddressById)

	return mux
}
