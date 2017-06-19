package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/zonesan/clog"
)

func ParseRequestBody(r *http.Request, v interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	clog.Debug("Request Body:", string(b))
	if err := json.Unmarshal(b, v); err != nil {
		return err
	}

	return nil
}

func RespOK(w http.ResponseWriter, data interface{}) {
	if data == nil {
		data = genRespJSON(nil)
	}

	if body, err := json.MarshalIndent(data, "", "  "); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func RespError(w http.ResponseWriter, err error) {
	resp := genRespJSON(err)

	if body, err := json.MarshalIndent(resp, "", "  "); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.status)
		w.Write(body)
	}

}

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Reason  string `json:"reason,omitempty"`
	status  int    `json:"status,omitempty"`
	//Data    interface{} `json:"data,omitempty"`
}

func genRespJSON(err error) *APIResponse {
	resp := new(APIResponse)

	if err == nil {
		resp.Code = http.StatusOK
		resp.status = http.StatusOK
	} else {
		if e, ok := err.(*StatusError); ok {
			resp.Code = int(e.ErrStatus.Code)

			// frontend can't handle 403/401, he will panic...
			{
				if resp.Code == http.StatusForbidden || resp.Code == http.StatusUnauthorized {
					resp.Code = http.StatusBadRequest
				}
			}
			resp.status = resp.Code
			resp.Message = e.ErrStatus.Message

		} else {
			resp.Code = http.StatusBadRequest
			resp.Message = err.Error()
			resp.status = http.StatusBadRequest //http.StatusBadRequest
		}
	}

	resp.Reason = http.StatusText(resp.status)

	return resp
}
