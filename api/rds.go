package api

type Rds struct {
	BaseURL string
	Params  interface{}
}

func (r *Rds) UsageAmount(params interface{}) *svcAmountList {
	return nil
}
