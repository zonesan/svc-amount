package api

import (
	"os"

	"fmt"

	"github.com/zonesan/clog"
)

type Hadoop struct {
	BaseURL string
	Params  interface{}
}

var hadoopBaseURL string

func (h *Hadoop) UsageAmount(svc string, bsi *BackingServiceInstance) (*svcAmountList, error) {
	uri := fmt.Sprintf("%s/%s/%s", h.BaseURL, svc, bsi.Spec.InstanceID)

	amounts, err := h.getAmountFromRemote(uri)
	if err != nil {
		clog.Error(err)
	}
	return amounts, err

	// amounts := &svcAmountList{Items: []svcAmount{
	// 	{Name: "RegionsQuota", Used: "300", Size: "500"},
	// 	{Name: "TablesQuotaa", Used: "20", Size: "100", Desc: "HBase命名空间的表数目"},
	// 	{Name: svc, Used: bsi.Spec.BackingServiceName, Desc: "faked response from hadoop."}}}

	// return amounts
}

func (h *Hadoop) getAmountFromRemote(uri string) (*svcAmountList, error) {
	result := new(svcAmountList)
	err := doRequest("GET", uri, nil, result, "")
	return result, err
}

func init() {

	hadoopBaseURL = os.Getenv("HADOOP_AMOUNT_BASEURL")
	if len(hadoopBaseURL) == 0 {
		clog.Fatal("HADOOP_AMOUNT_BASEURL must be specified.")
	}
	hadoopBaseURL = httpsAddr(hadoopBaseURL)
	clog.Debug("hadoop amount base url:", hadoopBaseURL)

	services := []string{"hbase", "hive", "hdfs", "kafka", "spark", "mapreduce"}
	hadoop := &Hadoop{BaseURL: hadoopBaseURL}
	register("hadoop", services, hadoop)

	// since hadoop and rds is the same api.
	hdpservices := []string{"mongodb", "greenplum", "mysql"}
	register("rds", hdpservices, hadoop)
}
