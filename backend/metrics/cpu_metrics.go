package metrics

import (
	"os"
	"strconv"
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
	name, cores, threads := GetValueForCPUMetrics()
	return &CPUMetrics{
		Name:    name,
		Cores:   int8(cores),
		Threads: int8(threads),
	}
}

func GetDataFromCpuInfo() (string, error) {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ParseData(data string) map[string]string {
	mp := make(map[string]string)
	for line := range strings.Lines(data) {
		slice := strings.Split(line, ":")
		name, value := GetDataFromStringSlice(slice)
		if _, ok := mp[name]; !ok {
			mp[name] = value
		}
	}
	return mp
}

func GetDataFromStringSlice(data []string) (string, string) {
	var (
		name  string
		value string
	)
	if len(data) > 1 {
		name = strings.Trim(data[0], "\t\n ")
		value = strings.Trim(data[1], " \n\t")
	}
	return name, value
}

func GetValueForCPUMetrics() (Name string, Cores int, Threads int) {
	data, err := GetDataFromCpuInfo()
	if err != nil {
		return "Uknown Processor", 0, 0
	}
	mp := ParseData(data)
	Cores, err = strconv.Atoi(mp["cpu cores"])
	if err != nil {
		return "Uknown Processor", 0, 0
	}
	Name = mp["model name"]
	Threads, err = strconv.Atoi(mp["siblings"])
	if err != nil {
		return "Uknown Processor", 0, 0
	}
	return Name, Cores, Threads
}
