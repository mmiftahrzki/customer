package customer

import (
	"database/sql"
	"net/http"

	"github.com/mmiftahrzki/customer/middleware"
)

type customer struct {
	http.Handler
}

func New(db *sql.DB, middlewares ...middleware.Middleware) customer {
	handler := newHandler(newService(newRepo(db)))
	mux := http.NewServeMux()

	deleteSingleById := middleware.ChainMiddleware(handler.DeleteSingleById, middlewares...)
	postSingle := middleware.ChainMiddleware(handler.PostSingle, middlewares...)
	putSingleById := middleware.ChainMiddleware(handler.PutSingleById, middlewares...)
	getSingleAndUpdateAddressById := middleware.ChainMiddleware(handler.GetSingleAndUpdateAddressById, middlewares...)

	mux.HandleFunc("GET /api/customer", handler.GetMultiple)
	mux.HandleFunc("GET /api/customer/{id}", handler.GetSingleById)
	mux.HandleFunc("GET /api/customer/{id}/prev", handler.GetMultiplePrev)
	mux.HandleFunc("GET /api/customer/{id}/next", handler.GetMultipleNext)
	mux.HandleFunc("POST /api/customer", postSingle)
	mux.HandleFunc("PUT /api/customer/{id}", putSingleById)
	mux.HandleFunc("PATCH /api/customer/{customer_id}/address/{address_id}", getSingleAndUpdateAddressById)
	mux.HandleFunc("DELETE /api/customer/{id}", deleteSingleById)

	return customer{Handler: mux}
}
