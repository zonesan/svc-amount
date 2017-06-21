package api

import (
	"net/url"

	"github.com/zonesan/clog"
)

type Container struct {
	BaseURL string
	Params  interface{}
}

func (c *Container) UsageAmount(svc string, bsi *BackingServiceInstance) *svcAmountList {

	c.findLabel(bsi.Spec.Creds)

	values := url.Values{}
	values.Set("labelSelector", "servicebroker=sb-3rxdog4nyzr3c-neo4j")

	amounts := &svcAmountList{Items: []svcAmount{
		{Name: "volume", Used: "30", Size: "50"},
		{Name: svc, Used: bsi.Spec.BackingServiceName, Desc: "faked response from container."}}}
	return amounts
}

func (c *Container) findLabel(creds map[string]string) (label string) {
	if v, ok := creds["host"]; !ok {
		clog.Error("can't find 'host' in credentials.")
		return ""
	} else {
		clog.Debug(v)
		label = v
	}
	return label
}

func init() {
	services := []string{"neo4j", "rabbitmq", "etcd"}
	container := &Container{}
	register("container", services, container)
}
