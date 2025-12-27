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
	_ = metrics.ParseData(data)
}

func TestGetDataFromStringSlice(t *testing.T) {
	data := []string{"processor", "0"}
	name, value := metrics.GetDataFromStringSlice(data)
	if name == "" {
		t.Errorf("%v", fmt.Errorf("Name is empty"))
	}
	fmt.Printf("%s %s", name, value)
}

func TestGetDataForCPUMetrics(t *testing.T) {
	name, cores, threads, freq := metrics.GetValueForCPUMetrics()
	fmt.Printf("%s %d %d %.3f", name, cores, threads, freq/1000)
}

func TestGetTempForCPU(t *testing.T) {
	temp := metrics.GetTemretureForCPU()
	fmt.Println(temp)
}
