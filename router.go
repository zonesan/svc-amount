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
	// router.GET(API_PREFIX+"/namespaces/:name/serviceinstances/:instance_name", openshift.ListMembers)

	router.NotFound = &api.Mux{}

	return router
}
