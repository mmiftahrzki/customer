package app

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/mmiftahrzki/customer/auth"
	"github.com/mmiftahrzki/customer/customer"
	"github.com/mmiftahrzki/customer/docs"
	"github.com/mmiftahrzki/customer/middleware"
)

func newMux(db *sql.DB) *http.ServeMux {
	jwtSigningKey := make([]byte, 256)
	rand.Read(jwtSigningKey[:])
	base64Encoded := base64.StdEncoding.EncodeToString(jwtSigningKey)
	fmt.Println(base64Encoded)

	appHandler := handler{}
	mux := http.NewServeMux()
	auth := auth.New(jwtSigningKey)
	customer := customer.New(db)
	doc := docs.New()

	deleteSingleById := middleware.ChainMiddleware(customer.Handler.DeleteSingleById, auth.Middleware.VerifyJWT)
	postSingle := middleware.ChainMiddleware(customer.Handler.PostSingle, auth.Middleware.VerifyJWT)
	putSingleById := middleware.ChainMiddleware(customer.Handler.PutSingleById, auth.Middleware.VerifyJWT)
	getSingleAndUpdateAddressById := middleware.ChainMiddleware(customer.Handler.GetSingleAndUpdateAddressById, auth.Middleware.VerifyJWT)

	mux.Handle("GET /{$}", appHandler)

	mux.HandleFunc("GET /swagger-css", doc.Handler.SwaggerCSS)
	mux.HandleFunc("GET /swagger-js", doc.Handler.SwaggerJS)
	mux.HandleFunc("GET /swagger", doc.Handler.SwaggerJson)
	mux.HandleFunc("GET /restful-api", doc.Handler.Swagger)

	mux.HandleFunc("POST /api/auth", auth.Handler.CreateAuthToken)
	mux.HandleFunc("POST /api/auth/", auth.Handler.CreateAuthToken)

	mux.HandleFunc("GET /api/customer", customer.Handler.GetMultiple)
	mux.HandleFunc("GET /api/customer/", customer.Handler.GetMultiple)
	mux.HandleFunc("GET /api/customer/{id}", customer.Handler.GetSingleById)
	mux.HandleFunc("GET /api/customer/{id}/prev", customer.Handler.GetMultiplePrev)
	mux.HandleFunc("GET /api/customer/{id}/prev/", customer.Handler.GetMultiplePrev)
	mux.HandleFunc("GET /api/customer/{id}/next", customer.Handler.GetMultipleNext)
	mux.HandleFunc("GET /api/customer/{id}/next/", customer.Handler.GetMultipleNext)
	mux.HandleFunc("POST /api/customer", postSingle)
	mux.HandleFunc("POST /api/customer/", postSingle)
	mux.HandleFunc("PUT /api/customer/{id}", putSingleById)
	mux.HandleFunc("PATCH /api/customer/{customer_id}/address/{address_id}", getSingleAndUpdateAddressById)
	mux.HandleFunc("DELETE /api/customer/{id}", deleteSingleById)

	return mux
}
