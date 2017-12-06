package api

import (
	"fmt"
	"net/http"

	"net/http/pprof"

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
	clog.Debug("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	ns := ps.ByName("name")
	instance := ps.ByName("instance_name")

	if oClient == nil {
		oClient = DFClient()
	}
	bsi, err := oClient.GetServiceInstance(ns, instance)
	if err != nil {
		clog.Error(err)
		RespError(w, err)
	} else {
		amounts, err := DoomServiceInstance(r, bsi)
		if err != nil {
			clog.Error(err)
			RespError(w, err)
		} else {
			clog.Tracef("%#v", amounts)
			RespOK(w, amounts)
		}
	}
}

func RestartInstance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Debug("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	ns := ps.ByName("name")
	instance := ps.ByName("instance_name")

	if oClient == nil {
		oClient = DFClient()
	}
	bsi, err := oClient.GetServiceInstance(ns, instance)
	if err != nil {
		clog.Error(err)
		RespError(w, err)
	} else {
		agent, err := findDriver(bsi.Spec.BackingServiceName)
		if err != nil {
			clog.Error(err)
			RespError(w, err)
			return
		}
		agent.req = r
		err = agent.RestartInstance(bsi)
		if err != nil {
			RespError(w, err)
			return
		}
		RespOK(w, nil)
	}
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)
	for driver, services := range drivers {
		clog.Infof("driver[%v]: %v", driver, services)
	}
	fmt.Fprint(w, "Welcome!\n")
}

func EnableDebug(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)
	clog.SetLogLevel(clog.LOG_LEVEL_TRACE)
	clog.Debug("DEBUG MODE ENABLED")
	fmt.Fprintf(w, "DEBUG MODE ENABLED")
}

func DisableDebug(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	clog.Debug("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)
	clog.SetLogLevel(clog.LOG_LEVEL_INFO)
	clog.Info("DEBUG MODE DISABLED")
	fmt.Fprintf(w, "DEBUG MODE DISABLED")
}

func DebugIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)
	pprof.Index(w, r)
}

func Command(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	pod := ps.ByName("pod")
	ns := ps.ByName("ns")
	if ns == "" {
		ns = "datafoundry"
	}

	oc := DFClient()
	oc.ExecCommand(ns, pod, "df", "/run/secrets")
}
