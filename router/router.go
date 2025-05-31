package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mmiftahrzki/customer/middleware"
	"github.com/mmiftahrzki/customer/response"
)

type Endpoint struct {
	Middlewares []middleware.Middleware
	Method      string
	Path        string
}

type Router struct {
	endpoints  map[string]Endpoint
	httprouter *httprouter.Router
}

func New() *Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := response.New()
		response.Message = "Resources not found!"

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write(response.ToJson())
	})
	router.MethodNotAllowed = router.NotFound

	return &Router{
		httprouter: router,
		endpoints:  map[string]Endpoint{},
	}
}

func (router *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	router.httprouter.ServeHTTP(writer, request)
}

func (router *Router) Handle(endpoint Endpoint, handle httprouter.Handle) {
	var handlers httprouter.Handle = handle

	for i := len(endpoint.Middlewares) - 1; i >= 0; i-- {
		handlers = endpoint.Middlewares[i](handlers)
	}

	router.httprouter.Handle(endpoint.Method, endpoint.Path, handlers)
}
