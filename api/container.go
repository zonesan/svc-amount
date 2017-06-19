package api

type Container struct {
	BaseURL string
	Params  interface{}
}

func (c *Container) UsageAmount(params interface{}) *svcAmountList {
	return nil
}
