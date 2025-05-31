package docs

import (
	"net/http"
)

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /swagger-css", swaggerCSS)
	mux.HandleFunc("GET /swagger-js", swaggerJS)
	mux.HandleFunc("GET /swagger", swaggerJson)
	mux.HandleFunc("GET /restful-api", swagger)

	return mux
}
