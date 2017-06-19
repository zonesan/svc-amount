package api

import (
	"fmt"
	"strings"
)

type AmountDriver interface {
	UsageAmount(svc, name string, params interface{}) *svcAmountList
}

type Agent struct {
	driver  AmountDriver
	service string
}

type AmountAgent struct {
	driver   AmountDriver
	services []string
}

var drivers = make(map[string]*AmountAgent)

func register(name string, services []string, driver AmountDriver) {
	if driver == nil {
		panic("driver: Register driver is nil!")
	}
	if _, dup := drivers[name]; dup {
		panic("driver: Register called twice for driver" + name)
	}
	drivers[name] = &AmountAgent{driver: driver, services: services}
}

func findDriver(service string) (*Agent, error) {
	for driver, agent := range drivers {
		for _, svc := range agent.services {
			if strings.ToLower(svc) == strings.ToLower(service) {
				return NewAgent(svc, driver)
			}
		}
	}
	return nil, fmt.Errorf("unsupported service '%s'", service)
}

func NewAgent(svc, name string) (*Agent, error) {
	driver, ok := drivers[name]
	if !ok {
		return nil, fmt.Errorf("Can't find agent %s", name)
	}
	return &Agent{driver: driver.driver, service: svc}, nil
}

func (agent *Agent) GetAmount(name string, params interface{}) *svcAmountList {
	return agent.driver.UsageAmount(agent.service, name, params)
}
