package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

func get_elb_list(sess *session.Session) []*elb.LoadBalancerDescription {
	svc := elb.New(sess)
	params := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{},
	}
	resp, err := svc.DescribeLoadBalancers(params)

	if err != nil {
		fmt.Printf("[ERROR] Fail DescribeLoadBalancers API call: %s \n", err.Error())
		return nil
	}
	return resp.LoadBalancerDescriptions
}

func (z *Zaws) ShowElbList() {
	list := make([]Data, 0)
	elbs := get_elb_list(z.AwsSession)
	for _, elb := range elbs {
		data := Data{ElbName: *elb.LoadBalancerName, ElbDnsName: *elb.DNSName}
		list = append(list, data)
	}
	fmt.Printf(convert_to_lldjson_string(list))
}

func (z *Zaws) ShowELBCloudwatchMetricsList() {
	list := make([]Data, 0)
	metrics := get_metric_list(z.AwsSession, "LoadBalancerName", z.TargetId)
	for _, metric := range metrics {
		datapoints := get_metric_stats(z.AwsSession, "LoadBalancerName", z.TargetId, *metric.MetricName, *metric.Namespace)
		metric_name := *metric.MetricName
		for _, dimension := range metric.Dimensions {
			if *dimension.Name == "AvailabilityZone" {
				metric_name = *metric.MetricName + "." + *dimension.Value
				break
			}
		}
		data := Data{MetricName: metric_name, MetricNamespace: *metric.Namespace}
		if len(datapoints) > 0 {
			data.MetricUnit = *datapoints[0].Unit
		}
		list = append(list, data)
	}

	fmt.Printf(convert_to_lldjson_string(list))
}

func (z *Zaws) SendElbMetricStats() {
	z.SendMetricStats("LoadBalancerName")
}
