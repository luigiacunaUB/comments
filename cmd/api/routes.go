package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) routes() http.Handler {
	//setup a new router
	router := httprouter.New()
	//handle 404
	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	//handle 405s
	router.NotFound = http.HandlerFunc(a.methodNotAllowedResponse)
	//setup routes
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthcheckHandler)
	//return router
	return a.recoverPanic(router)
}
