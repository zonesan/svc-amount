package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/zonesan/clog"
)

type Mux struct{}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)
	//RespError(w, ErrorNew(ErrCodeNotFound), http.StatusNotFound)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))
}

func AmountInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	ns := ps.ByName("name")
	instance := ps.ByName("instance_name")

	if oClient == nil {
		oClient = NewDataFoundryTokenClient()
	}
	bsi, err := oClient.GetServiceInstance(ns, instance)
	if err != nil {
		clog.Error(err)
		RespError(w, err)
	} else {
		amounts, err := DoomServiceInstance(bsi)
		if err != nil {
			clog.Error(err)
			RespError(w, err)
		} else {
			// amounts = &svcAmountList{Items: []svcAmount{
			// 	{Name: "RegionsQuota", Used: "300", Size: "500"},
			// 	{Name: "TablesQuotaa", Used: "20", Size: "100", Desc: "HBase命名空间的表数目"},
			// 	{Name: ns, Used: instance, Size: r.URL.RequestURI(), Desc: "faked response."}}}
			// fmt.Fprint(w, "Welcome!\n")
			RespOK(w, amounts)
		}
	}
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)
	for driver, services := range drivers {
		clog.Infof("driver[%v]: %v", driver, services)
	}
	fmt.Fprint(w, "Welcome!\n")
}
