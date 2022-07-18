package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"strings"
)

func get_rds_list(sess *session.Session) []*rds.DBInstance {
	svc := rds.New(sess)
	params := &rds.DescribeDBInstancesInput{}
	resp, err := svc.DescribeDBInstances(params)

	if err != nil {
		fmt.Printf("[ERROR] Fail DescribeDBInstances API call: %s \n", err.Error())
		return nil
	}
	return resp.DBInstances
}

func (z *Zaws) ShowRdsList() {
	list := make([]Data, 0)
	rdsInstances := get_rds_list(z.AwsSession)
	for _, rdsInstance := range rdsInstances {
		data := Data{
			InstanceName: *rdsInstance.DBInstanceIdentifier,
			RdsName:      *rdsInstance.DBInstanceIdentifier,
			RdsDnsName:   *rdsInstance.Endpoint.Address,
			InstanceType: *rdsInstance.DBInstanceClass}
		list = append(list, data)
	}
	fmt.Printf(convert_to_lldjson_string(list))
}

func (z *Zaws) ShowRdsCloudwatchMetricsList() {
	list := make([]Data, 0)
	metrics := get_metric_list(z.AwsSession, "DBInstanceIdentifier", z.TargetId)
	for _, metric := range metrics {
		datapoints := get_metric_stats(z.AwsSession, "DBInstanceIdentifier", z.TargetId, *metric.MetricName, *metric.Namespace)
		metricName := strings.Replace(*metric.MetricName, "%", "Perc", -1)

		data := Data{MetricName: metricName, MetricNamespace: *metric.Namespace}
		if len(datapoints) > 0 {
			data.MetricUnit = *datapoints[0].Unit
		}
		list = append(list, data)
	}

	fmt.Printf(convert_to_lldjson_string(list))
}

func (z *Zaws) SendRdsMetricStats() {
	z.SendMetricStats("DBInstanceIdentifier")
}
