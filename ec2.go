package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
)

func get_ec2_list(sess *session.Session) []*ec2.Instance {
	var instances []*ec2.Instance
	svc := ec2.New(sess)
	resp, err := svc.DescribeInstances(nil)

	if err != nil {
		fmt.Printf("[ERROR] Fail DescribeInstances API call: %s \n", err.Error())
		os.Exit(1)
	}
	for _, reservation := range resp.Reservations {
		instances = append(instances, reservation.Instances...)
	}
	return instances
}

func (z *Zaws) ShowEc2List() {
	list := make([]Data, 0)
	instances := get_ec2_list(z.AwsSession)
	for _, instance := range instances {
		data := Data{InstanceType: *instance.InstanceType, InstanceId: *instance.InstanceId}
		if instance.PrivateIpAddress != nil {
			data.InstancePrivateAddr = *instance.PrivateIpAddress
		}
		for _, tag := range instance.Tags {
			if *tag.Key == "Name" {
				data.InstanceName = *tag.Value
			}
		}
		if data.InstanceName == "" {
			data.InstanceName = *instance.InstanceId
		}
		list = append(list, data)
	}
	fmt.Printf(convert_to_lldjson_string(list))
}

func (z *Zaws) ShowEC2CloudwatchMetricsList() {
	list := make([]Data, 0)
	metrics := get_metric_list(z.AwsSession, "InstanceId", z.TargetId)
	for _, metric := range metrics {
		datapoints := get_metric_stats(z.AwsSession, "InstanceId", z.TargetId, *metric.MetricName, *metric.Namespace)
		data := Data{MetricName: *metric.MetricName, MetricNamespace: *metric.Namespace}
		if len(datapoints) > 0 {
			data.MetricUnit = *datapoints[0].Unit
		}
		list = append(list, data)
	}

	fmt.Printf(convert_to_lldjson_string(list))
}

func (z *Zaws) SendEc2MetricStats() {
	z.SendMetricStats("InstanceId")
}
