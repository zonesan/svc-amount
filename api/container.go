package api

type Container struct {
	BaseURL string
	Params  interface{}
}

func (c *Container) UsageAmount(svc, name string, params interface{}) *svcAmountList {
	amounts := &svcAmountList{Items: []svcAmount{
		{Name: "volume", Used: "30", Size: "50"},
		{Name: svc, Used: name, Desc: "faked response from container."}}}
	return amounts
}

func init() {
	services := []string{"neo4j", "rabbitmq"}
	container := &Container{}
	register("container", services, container)
}
