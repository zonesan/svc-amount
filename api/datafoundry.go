package api

import (
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

	return oClient
}

func (c *DataFoundryClient) setBearerToken(token string) {
	c.bearerToken.Store(token)
}

func (c *DataFoundryClient) BearerToken() string {
	//return oc.bearerToken
	return c.bearerToken.Load().(string)
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
