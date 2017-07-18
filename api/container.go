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

const (
	BROKER_CONTAINER_NS = "service-brokers"
	DISKFREE_CMD        = "df"
)

func (c *Container) UsageAmount(svc string, bsi *BackingServiceInstance) (*svcAmountList, error) {

	k, v := c.findPodLabel(bsi.Spec.Creds)
	if len(k) == 0 || len(v) == 0 {
		return nil, fmt.Errorf("can't locate pod due to an empty label.")
	}

	pods, err := c.findPodsByLabelSelector(k, v)
	if err != nil {
		clog.Error(err)
		return nil, err
	}

	var podsname []string
	for _, v := range pods.Items {
		podsname = append(podsname, v.Name)
	}
	clog.Debug("pods list:", podsname)

	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("can't list pods by specified label.")
	}

	pod := &pods.Items[0]
	mountPath, err := c.findVolumeMountPath(pod)
	if err != nil {
		clog.Error(err)
		return nil, err
	}

	amount, err := c.getVolumeAmount(pod.Name, mountPath)
	if err != nil {
		clog.Error(err)
		return nil, err
	}

	amounts := &svcAmountList{Items: []svcAmount{*amount}}

	// if amount != nil {
	// 	amounts.Items = append(amounts.Items, *amount)
	// }
	return amounts, nil
}

func (c *Container) getVolumeAmount(podName, mountPath string) (*svcAmount, error) {
	oc := DFClient()
	res, err := oc.ExecCommand(BROKER_CONTAINER_NS, podName, DISKFREE_CMD, mountPath)
	if err != nil {
		clog.Error(err)
		return nil, err
	}
	amount, ok := res.(*svcAmount)
	if !ok {
		return nil, fmt.Errorf("unknown error..")
	}
	return amount, nil
}

func (c *Container) findVolumeMountPath(pod *kapi.Pod) (string, error) {
	volumes := pod.Spec.Volumes
	var volumeName string
	var mountPath string
	for _, volume := range volumes {
		if volume.PersistentVolumeClaim != nil {
			volumeName = volume.Name
			break
		}
	}
	if len(volumeName) == 0 {
		return "", fmt.Errorf("can't locate pvc in pod %v.", pod.Name)
	}

	container := pod.Spec.Containers[0]
	for _, mounts := range container.VolumeMounts {
		if mounts.Name == volumeName {
			mountPath = mounts.MountPath
			break
		}
	}
	if len(mountPath) == 0 {
		return "", fmt.Errorf("can't find mount point of volume '%v' in pod '%v'", volumeName, pod.Name)
	}

	clog.Debugf("volume '%v' mounted to '%v'", volumeName, mountPath)
	return mountPath, nil
}

func (c *Container) findPodLabel(creds map[string]string) (k, v string) {
	svcname := ""

	if host, ok := creds["vhost"]; !ok {
		clog.Error("can't find 'vhost' in credentials.")
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
	clog.Debugf("label: %v=%v", k, v)
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

func (c *Container) execCommand(cmd string, args ...string) (interface{}, error) {
	return nil, nil
}

func init() {
	services := []string{"neo4j", "rabbitmq"}
	container := &Container{}
	register("container", services, container)
}
