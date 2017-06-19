package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"k8s.io/kubernetes/pkg/api/unversioned"

	"github.com/zonesan/clog"
)

func doRequest(method, url string, bodyParams, v interface{}, token string) (err error) {
	clog.Debug(method, url)
	var reqbody []byte
	if bodyParams != nil {
		reqbody, err = json.Marshal(bodyParams)
		if err != nil {
			return err
		}
	}

	resp, err := request(method, url, reqbody, token)

	if err != nil {
		return err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	err = checkRespStatus(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
			//clog.Tracef("%#v", v)
		}
	}
	return err

}

func request(method string, url string, body []byte, token string) (*http.Response, error) {
	// token := c.BearerToken()
	// if token == "" {
	// 	return nil, errors.New("token is blank")
	// }

	clog.Trace("request url:", url)
	var req *http.Request
	var err error
	if len(body) == 0 {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, bytes.NewReader(body))
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// token is used to request df api.
	if len(token) > 0 {
		req.Header.Set("Authorization", token)
	}

	transCfg := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transCfg,
		// Timeout:   timeout,
	}
	return client.Do(req)
}

// StatusError is an error intended for consumption by a REST API server; it can also be
// reconstructed by clients from a REST response. Public to allow easy type switches.
type StatusError struct {
	ErrStatus unversioned.Status
}

// Error implements the Error interface.
func (e *StatusError) Error() string {
	return e.ErrStatus.Message
}
func checkRespStatus(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	// openshift returns 401 with a plain text but not ErrStatus json, so we hacked this response text.
	if r.StatusCode == http.StatusUnauthorized {
		errorResponse := &StatusError{}
		errorResponse.ErrStatus.Code = http.StatusUnauthorized
		errorResponse.ErrStatus.Message = http.StatusText(http.StatusUnauthorized)
		return errorResponse
	}

	errorResponse := &StatusError{}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		clog.Errorf("%v,%s", r.StatusCode, data)
		json.Unmarshal(data, &errorResponse.ErrStatus)
	}

	return errorResponse
}

func setBaseURL(urlStr string) string {
	// Make sure the given URL end with a slash
	if strings.HasSuffix(urlStr, "/") {
		return setBaseURL(strings.TrimSuffix(urlStr, "/"))
	}
	return urlStr
}

func httpsAddr(addr string) string {

	if !strings.HasPrefix(strings.ToLower(addr), "http://") &&
		!strings.HasPrefix(strings.ToLower(addr), "https://") {
		return fmt.Sprintf("https://%s", addr)
	}

	return setBaseURL(addr)
}
