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
)

func newMux(db *sql.DB) *http.ServeMux {
	jwtSigningKey := make([]byte, 256)
	rand.Read(jwtSigningKey[:])
	base64Encoded := base64.StdEncoding.EncodeToString(jwtSigningKey)
	fmt.Println(base64Encoded)

	appHandler := handler{}

	auth := auth.New(jwtSigningKey)
	customer := customer.New(db)
	doc := docs.New()

	customerMux := http.NewServeMux()
	mux := http.NewServeMux()

	testThenVerifyAuth := pipe(Test, auth.Middleware.VerifyJWT)
	deleteSingleById := add(auth.Middleware.VerifyJWT, customer.Handler.DeleteSingleById)
	postSingle := add(auth.Middleware.VerifyJWT, customer.Handler.PostSingle)
	putSingleById := add(auth.Middleware.VerifyJWT, customer.Handler.PutSingleById)
	getSingleAndUpdateAddressById := add(auth.Middleware.VerifyJWT, customer.Handler.GetSingleAndUpdateAddressById)

	customerMux.HandleFunc("GET /api/customer/{$}", customer.Handler.GetMultiple)
	customerMux.HandleFunc("GET /api/customer/{id}", customer.Handler.GetSingleById)
	customerMux.HandleFunc("GET /api/customer/{id}/prev/{$}", customer.Handler.GetMultiplePrev)
	customerMux.HandleFunc("GET /api/customer/{id}/next/{$}", customer.Handler.GetMultipleNext)
	customerMux.HandleFunc("POST /api/customer/{$}", postSingle)
	customerMux.HandleFunc("PUT /api/customer/{id}", putSingleById)
	customerMux.HandleFunc("PATCH /api/customer/{customer_id}/address/{address_id}", getSingleAndUpdateAddressById)
	customerMux.HandleFunc("DELETE /api/customer/{id}", deleteSingleById)

	mux.Handle("GET /{$}", appHandler)
	mux.HandleFunc("GET /swagger-css", doc.Handler.SwaggerCSS)
	mux.HandleFunc("GET /swagger-js", doc.Handler.SwaggerJS)
	mux.HandleFunc("GET /swagger", doc.Handler.SwaggerJson)
	mux.HandleFunc("GET /restful-api", doc.Handler.Swagger)

	mux.HandleFunc("POST /api/test", testThenVerifyAuth(customer.Handler.GetMultiple))

	mux.HandleFunc("POST /api/auth/{$}", auth.Handler.CreateAuthToken)

	mux.HandleFunc("/api/customer/", customerMux.ServeHTTP)

	return mux
}
