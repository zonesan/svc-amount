package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ocmanager/svc-amount/api"
)

const (
	API_PREFIX = "/sapi/v1"
)

func createRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", api.Index)

	router.GET(API_PREFIX+"/namespaces/:name/instances/:instance_name", api.AmountInfo)

	router.PUT(API_PREFIX+"/namespaces/:name/instances/:instance_name", api.RestartInstance)

	router.PUT(API_PREFIX+"/debug", api.EnableDebug)
	router.DELETE(API_PREFIX+"/debug", api.DisableDebug)

	router.GET(API_PREFIX+"/ns/:ns/cmd/:pod", api.Command)

	debug := true
	if debug {
		router.GET("/debug/pprof/", api.DebugIndex)
		router.GET("/debug/pprof/:name", api.DebugIndex)
	}

	router.NotFound = &api.Mux{}

	return router
}

func init() {
	api.ListDrivers()
}
