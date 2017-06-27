package api

import (
	"net/url"
	"os"
	"sync/atomic"

	"fmt"

	"github.com/zonesan/clog"
	"k8s.io/kubernetes/pkg/api/unversioned"
	kapi "k8s.io/kubernetes/pkg/api/v1"
)

var (
	dataFoundryHostAddr string
	dataFoundryToken    string
	oClient             *DataFoundryClient
)

type DataFoundryClient struct {
	host        string
	oapiURL     string
	kapiURL     string
	bearerToken atomic.Value
	token       atomic.Value
}

// BackingServiceInstance describe a BackingServiceInstance
type BackingServiceInstance struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	kapi.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the behavior of the Namespace.
	Spec BackingServiceInstanceSpec `json:"spec,omitempty" description:"spec defines the behavior of the BackingServiceInstance"`

	// Status describes the current status of a Namespace
}

// BackingServiceInstanceSpec describes the attributes on a BackingServiceInstance
type BackingServiceInstanceSpec struct {
	// description of an instance.
	InstanceProvisioning `json:"provisioning, omitempty"`
	// id of an instance
	InstanceID string `json:"instance_id, omitempty"`
	// tags of an instance
}

// InstanceProvisioning describe an InstanceProvisioning detail
type InstanceProvisioning struct {
	// dashboard url of an instance
	DashboardUrl string `json:"dashboard_url, omitempty"`
	// bs name of an instance
	BackingServiceName string `json:"backingservice_name, omitempty"`
	// bs id of an instance
	BackingServiceSpecID string `json:"backingservice_spec_id, omitempty"`
	// bs plan id of an instance
	BackingServicePlanGuid string `json:"backingservice_plan_guid, omitempty"`
	// bs plan name of an instance
	BackingServicePlanName string `json:"backingservice_plan_name, omitempty"`
	// parameters of an instance
	Parameters map[string]string `json:"parameters, omitempty"`
	// credentials of an instance
	Creds map[string]string `json:"credentials, omitempty"`
}

func DFClient() *DataFoundryClient {
	return NewDataFoundryTokenClient(dataFoundryToken)
}

func NewDataFoundryTokenClient(token string) *DataFoundryClient {

	if oClient != nil {
		return oClient
	}
	// host = setBaseUrl(host)
	oClient = &DataFoundryClient{
		host:    dataFoundryHostAddr,
		oapiURL: dataFoundryHostAddr + "/oapi/v1",
		kapiURL: dataFoundryHostAddr + "/api/v1",
	}

	oClient.setBearerToken("Bearer " + token)
	oClient.setToken(token)

	return oClient
}

func (c *DataFoundryClient) setBearerToken(token string) {
	c.bearerToken.Store(token)
}

func (c *DataFoundryClient) setToken(token string) {
	c.token.Store(token)
}

func (c *DataFoundryClient) BearerToken() string {
	//return oc.bearerToken
	return c.bearerToken.Load().(string)
}

func (c *DataFoundryClient) Token() string {
	return c.token.Load().(string)
}

func (c *DataFoundryClient) GetServiceInstance(ns, name string) (*BackingServiceInstance, error) {
	uri := "/namespaces/" + ns + "/backingserviceinstances/" + name
	bsi := new(BackingServiceInstance)
	err := c.OGet(uri, bsi)
	clog.Trace(bsi)
	return bsi, err
}

func (c *DataFoundryClient) GetService(ns, name string) (*kapi.Service, error) {
	uri := "/namespaces/" + ns + "/services/" + name
	svc := new(kapi.Service)
	err := c.KGet(uri, svc)
	clog.Trace(svc)
	return svc, err
}

func (c *DataFoundryClient) ListPods(ns, queryParam string) (*kapi.PodList, error) {
	uri := fmt.Sprintf("/namespaces/%s/pods?%s", ns, queryParam)
	pods := &kapi.PodList{}
	err := c.KGet(uri, pods)
	if err != nil {
		clog.Error(err)
		return nil, err
	}

	return pods, err

}

func (c *DataFoundryClient) ExecCommand(ns, pod, cmd string, args ...string) (interface{}, error) {
	// wsd -insecureSkipVerify \
	// -url \
	// 'wss://10.247.32.17/api/v1/namespaces/service-brokers/pods/sb-ceeajasbecimq-rbbtmq-a4zgh/exec?command=df&command=%2fvar%2flib%2frabbitmq&acss_token=22736IIO7vk_lD1_Bq_rktRQzP7JZzDgjk66-4DjLHk' \
	// -origin https://10.247.32.17
	values := url.Values{}
	values.Set("command", cmd)
	command := values.Encode()

	for _, arg := range args {
		values.Set("command", arg)
		cmdArg := values.Encode()
		command = command + "&" + cmdArg
	}
	uri := fmt.Sprintf("%s/namespaces/%s/pods/%s/exec?%s&access_token=%s",
		c.kapiURL, ns, pod, command, c.Token())
	clog.Debug(uri)

	u, err := url.Parse(uri)
	if err != nil {
		clog.Error(err)
		return nil, err
	}
	u.Scheme = "wss"
	url := u.String()
	origin := c.host
	clog.Debugf("url: %s, origin: %s", url, origin)

	ws(url, origin)

	return nil, nil
}

func (c *DataFoundryClient) OGet(uri string, into interface{}) error {
	return doRequest("GET", c.oapiURL+uri, nil, into, c.BearerToken())
}

func (c *DataFoundryClient) OPost(uri string, body, into interface{}) error {
	return doRequest("POST", c.oapiURL+uri, body, into, c.BearerToken())
}

func (c *DataFoundryClient) KGet(uri string, into interface{}) error {
	return doRequest("GET", c.kapiURL+uri, nil, into, c.BearerToken())
}

func (c *DataFoundryClient) KPost(uri string, body, into interface{}) error {
	return doRequest("POST", c.kapiURL+uri, body, into, c.BearerToken())
}

func init() {
	dataFoundryHostAddr = os.Getenv("DATAFOUNDRY_API_SERVER")
	if len(dataFoundryHostAddr) == 0 {
		clog.Fatal("DATAFOUNDRY_API_SERVER must be specified.")
	}
	dataFoundryHostAddr = httpsAddr(dataFoundryHostAddr)
	clog.Debug("datafoundry api server:", dataFoundryHostAddr)

	dataFoundryToken = os.Getenv("DATAFOUNDRY_API_TOKEN")
	if len(dataFoundryToken) == 0 {
		clog.Fatal("DATAFOUNDRY_API_TOKEN must be specified.")
	}
	clog.Debug("datafoundry api token:", "*HIDDEN*") // dataFoundryToken)

	oClient = NewDataFoundryTokenClient(dataFoundryToken)
}
