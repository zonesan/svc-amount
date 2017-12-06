package api

type Rds struct {
	ServiceDefault
	BaseURL string
	Params  interface{}
}

// func (r *Rds) UsageAmount(svc string, bsi *BackingServiceInstance, req *http.Request) (*svcAmountList, error) {
// 	amounts := &svcAmountList{Items: []svcAmount{
// 		{Name: "dbsize", Used: "30", Size: "50"},
// 		{Name: svc, Used: bsi.Spec.BackingServiceName, Desc: "faked response from rds."}}}
// 	return amounts, nil
// }

func init() {
	// services := []string{"mongodb", "greenplum"}
	// rds := &Rds{}
	// register("rds", services, rds)
}
