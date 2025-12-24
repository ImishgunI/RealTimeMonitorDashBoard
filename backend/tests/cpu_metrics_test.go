package tests

import (
	"real_time_monitor_dashboard/backend/metrics"
	"testing"
)

func TestCpuGetDataFromCpuInfo(t *testing.T) {
	data, err := metrics.GetDataFromCpuInfo()
	if err != nil {
		t.Errorf("%v", err)
	}
	if data == "" {
		t.Errorf("data is empty")
	}
}

func TestParseData(t *testing.T) {
	data, err := metrics.GetDataFromCpuInfo()
	if err != nil {
		t.Errorf("%v", err)
	}
	metrics.ParseData(data)
}
