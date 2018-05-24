package httphandler

import (
	"net/http"
	"os"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
	"gitlab.brickchain.com/libs/go-httphandler.v2/middleware"
)

var allowedHeaders = []string{"Accept", "Accept-Language", "Content-Language", "Content-Type", "Origin", "Authorization", "X-Auth-Token"}
var allowedMethods = []string{"GET", "POST", "PUT", "DELETE"}

// SetAllowedHeaders sets the headers we allow in CORS
func SetAllowedHeaders(h []string) {
	allowedHeaders = h
}

// SetAllowedMethods sets the methods we allow in CORS
func SetAllowedMethods(m []string) {
	allowedMethods = m
}

// NewRouter returns a new httprouter.Router object
func NewRouter() *httprouter.Router {
	return httprouter.New()
}

// LoadMiddlewares adds the middlewares for CORS and Server and proxy headers
func LoadMiddlewares(router *httprouter.Router, version string) http.Handler {
	headersOk := gorillaHandlers.AllowedHeaders(allowedHeaders)
	methodsOk := gorillaHandlers.AllowedMethods(allowedMethods)
	handler := gorillaHandlers.CORS(headersOk, methodsOk)(router)
	handler = middleware.ResponseHeader(handler, "Server", os.Args[0]+"/"+version)
	handler = gorillaHandlers.ProxyHeaders(handler)

	return handler
}
