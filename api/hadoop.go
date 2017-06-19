package api

type Hadoop struct {
	BaseURL string
	Params  interface{}
}

func (h *Hadoop) UsageAmount(svc, name string, params interface{}) *svcAmountList {
	amounts := &svcAmountList{Items: []svcAmount{
		{Name: "RegionsQuota", Used: "300", Size: "500"},
		{Name: "TablesQuotaa", Used: "20", Size: "100", Desc: "HBase命名空间的表数目"},
		{Name: svc, Used: name, Desc: "faked response from hadoop."}}}
	return amounts
}

func init() {
	services := []string{"hbase", "hive", "hdfs", "kafka", "spark", "mapreduce"}
	hadoop := &Hadoop{}
	register("hadoop", services, hadoop)
}
