package tests

import (
	"fmt"
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
	mp := metrics.ParseData(data)
	fmt.Print(mp)
}

func TestGetDataFromStringSlice(t *testing.T) {
	data := []string{"processor", "0"}
	name, value := metrics.GetDataFromStringSlice(data)
	if name == "" {
		t.Errorf("%v", fmt.Errorf("Name is empty"))
	}
	fmt.Println(name, value)
}
