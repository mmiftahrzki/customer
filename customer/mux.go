package customer

import (
	"net/http"

	"github.com/mmiftahrzki/customer/auth"
)

func newMux(h *handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("DELETE /api/customer/{id}", auth.Verify(http.HandlerFunc(h.DeleteSingleById)))
	mux.HandleFunc("POST /api/customer/", h.PostSingle)
	mux.HandleFunc("GET /api/customer/", h.GetMultiple)
	mux.HandleFunc("GET /api/customer/{id}", h.GetSingleById)
	// mux.HandleFunc("GET /api/customer/{id}/prev", h.GetMultiplePrev)
	mux.HandleFunc("GET /api/customer/{id}/next", h.GetMultipleNext)
	mux.HandleFunc("PUT /api/customer/{id}", h.PutSingleById)
	mux.HandleFunc("PATCH /api/customer/{customer_id}/address/{address_id}", h.GetSingleAndUpdateAddressById)

	return mux
}
