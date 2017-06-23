package api

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/zonesan/clog"
	kapi "k8s.io/kubernetes/pkg/api/v1"
)

type Container struct {
	BaseURL string
	Params  interface{}
}

const BROKER_CONTAINER_NS = "service-brokers"

func (c *Container) UsageAmount(svc string, bsi *BackingServiceInstance) (*svcAmountList, error) {

	k, v := c.findPodLabel(bsi.Spec.Creds)
	if len(k) == 0 || len(v) == 0 {
		return nil, fmt.Errorf("can't locate pod due to an empty label.")
	}

	clog.Debug("label:", k, v)

	pods, err := c.findPodsByLabelSelector(k, v)
	if err != nil {
		clog.Error(err)
		return nil, err
	}

	clog.Debug("pods:", pods)

	amounts := &svcAmountList{Items: []svcAmount{
		{Name: "volume", Used: "30", Size: "50"},
		{Name: svc, Used: bsi.Spec.BackingServiceName, Desc: "faked response from container."}}}
	return amounts, nil
}

func (c *Container) findPodLabel(creds map[string]string) (k, v string) {
	svcname := ""

	if host, ok := creds["host"]; !ok {
		clog.Error("can't find 'host' in credentials.")
		return
	} else {
		clog.Debug("vhost:", host)
		s := strings.Split(host, ".")
		if len(s) > 0 && len(s[0]) > 0 {
			svcname = s[0]
		} else {
			return
		}
	}

	svc, err := c.getService(svcname)
	if err != nil {
		clog.Error(err)
		return
	}

	for k, v = range svc.Labels {

	}
	if len(k) == 0 || len(v) == 0 {
		clog.Error("can't find label.")
	}
	return k, v
}

func (c *Container) getService(name string) (*kapi.Service, error) {
	oc := DFClient()
	return oc.GetService(BROKER_CONTAINER_NS, name)
}

func (c *Container) findPodsByLabelSelector(k, v string) (*kapi.PodList, error) {
	labelSelector := k + "=" + v

	values := url.Values{}
	values.Set("labelSelector", labelSelector)
	encodedLabel := values.Encode()
	clog.Debug(encodedLabel)

	oc := DFClient()

	return oc.ListPods(BROKER_CONTAINER_NS, encodedLabel)
}

func init() {
	services := []string{"neo4j", "rabbitmq"}
	container := &Container{}
	register("container", services, container)
}
