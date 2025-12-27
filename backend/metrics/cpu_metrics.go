package metrics

import (
	"fmt"
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
	name, cores, threads, freq := GetValueForCPUMetrics()
	return &CPUMetrics{
		Name:      name,
		Cores:     int8(cores),
		Threads:   int8(threads),
		Frequency: freq / 1000,
		Temreture: GetTemretureForCPU(),
		Workload:  0,
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

func GetValueForCPUMetrics() (Name string, Cores int, Threads int, Frequency float32) {
	data, err := GetDataFromCpuInfo()
	if err != nil {
		return "Uknown Processor", 0, 0, 0
	}
	mp := ParseData(data)
	Cores, err = strconv.Atoi(mp["cpu cores"])
	if err != nil {
		return "Uknown Processor", 0, 0, 0
	}
	Name = mp["model name"]
	Threads, err = strconv.Atoi(mp["siblings"])
	if err != nil {
		return "Uknown Processor", 0, 0, 0
	}
	freq, err := strconv.ParseFloat(mp["cpu MHz"], 32)
	if err != nil {
		return "Uknown Processor", 0, 0, 0
	}
	Frequency = float32(freq)
	return Name, Cores, Threads, Frequency
}

func GetTemretureForCPU() float32 {
	data, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		fmt.Printf("%v\n", err)
		return 0
	}
	builder := strings.Builder{}
	for _, c := range data {
		if c != '\n' {
			builder.WriteByte(c)
		}
	}
	temreture, err := strconv.ParseFloat(builder.String(), 32)
	if err != nil {
		fmt.Printf("%v", err)
		return 0
	}
	return float32(temreture / 1000)
}
