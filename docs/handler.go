package docs

import "net/http"

type handler struct {
}

func (h handler) SwaggerYAML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/yaml")
	w.WriteHeader(http.StatusOK)
	w.Write(swaggerYAML)
}

func (h handler) SwaggerJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/yaml")
	w.WriteHeader(http.StatusOK)
	w.Write(swaggerJSON)
}

func (h handler) SwaggerCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	w.WriteHeader(http.StatusOK)
	w.Write(swaggerUICSS)
}

func (h handler) SwaggerJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")
	w.WriteHeader(http.StatusOK)
	w.Write(swaggerUIJS)
}

func (h handler) Swagger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(swaggerUIHTML)
}
