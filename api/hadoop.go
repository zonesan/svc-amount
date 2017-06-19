package api

type Hadoop struct {
	BaseURL string
	Params  interface{}
}

func (h *Hadoop) UsageAmount(params interface{}) *svcAmountList {
	return nil
}
