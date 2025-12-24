package metrics

import (
	"os"
	"strings"
)

type CPUMetrics struct {
	Name      string
	Cores     int8
	Threads   int8
	Frequency float32
	Temreture float32
	Workload  int8
}

func New() *CPUMetrics {
	return &CPUMetrics{}
}

func GetDataFromCpuInfo() (string, error) {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ParseData(data string) {
	for line := range strings.Lines(data) {
		_ = strings.Split(line, ":")
	}
}
