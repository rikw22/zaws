package main

import (
	"os"
	"testing"
)

func TestZaws_ShowEc2List(t *testing.T) {
	zaws := NewZaws()
	zaws.ShowEc2List()
}

func TestZaws_ShowEC2CloudwatchMetricsList(t *testing.T) {
	zaws := NewZaws()
	zaws.TargetId = os.Getenv("AWS_EC2_TEST_INSTANCE")
	zaws.ShowEC2CloudwatchMetricsList()
}

func TestZaws_ShowRdsList(t *testing.T) {
	zaws := NewZaws()
	zaws.ShowRdsList()
}

func TestZaws_ShowRdsCloudwatchMetricsList(t *testing.T) {
	zaws := NewZaws()
	zaws.TargetId = os.Getenv("AWS_RDS_TEST_INSTANCE")
	zaws.ShowRdsCloudwatchMetricsList()
}
