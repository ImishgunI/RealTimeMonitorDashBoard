package metrics

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type CPUMetrics struct {
	Name      string
	Cores     int8
	Threads   int8
	Frequency float32
	Temreture float32
	Workload  float32
}

func New() *CPUMetrics {
	return &CPUMetrics{
		Name:      "",
		Cores:     0,
		Threads:   0,
		Frequency: 0.0,
		Temreture: 0.0,
		Workload:  0.0,
	}
}

func getDataFromCpuInfo() (string, error) {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func parseData(data string) [][]string {
	parsed := [][]string{}
	scanner := bufio.NewScanner(strings.NewReader(data))
	for range 13 {
		scanner.Scan()
		parsed = append(parsed, strings.Fields(scanner.Text()))
	}
	return parsed
}

func GetValueForCPUMetrics() (Name string, Cores int8, Threads int8, Frequency float32) {
	data, err := getDataFromCpuInfo()
	if err != nil {
		return "Unknown processor", 0, 0, 0.0
	}
	parsed := parseData(data)
	Name = getName(parsed)
	Cores = getCores(parsed)
	Threads = getThreads(parsed)
	Frequency = getFrequency(parsed)
	return Name, Cores, Threads, Frequency
}

func getName(parsed [][]string) string {
	builder := strings.Builder{}
	for j := 3; j < len(parsed[4]); j++ {
		builder.WriteString(parsed[4][j])
		if j+1 < len(parsed[4]) {
			builder.WriteString(" ")
		}
	}
	return builder.String()
}

func getCores(parsed [][]string) int8 {
	core, err := strconv.ParseInt(parsed[12][3], 10, 8)
	if err != nil {
		fmt.Printf("Failed to parse int for cores: %v", err)
		return 0
	}
	return int8(core)
}

func getThreads(parsed [][]string) int8 {
	threads, err := strconv.ParseInt(parsed[10][2], 10, 8)
	if err != nil {
		fmt.Printf("Failed to parse int for threads: %v", err)
		return 0 
	}
	return int8(threads)
}

func getFrequency(parsed [][]string) float32 {
	freq, err := strconv.ParseFloat(parsed[7][3], 10)
	if err != nil {
		fmt.Printf("Failed to parse float for frequency: %v", err)
		return 0.0
	}
	return float32(freq)
}

func getdataFromThermalZone() string {
	data, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		fmt.Printf("%v\n", err)
		return ""
	}
	builder := strings.Builder{}
	for _, c := range data {
		if c != '\n' {
			builder.WriteByte(c)
		}
	}
	return builder.String()
}

func GetTemretureForCPU() float32 {
	temp :=  getdataFromThermalZone()
	temreture, err := strconv.ParseFloat(temp, 32)
	if err != nil {
		fmt.Printf("%v", err)
		return 0
	}
	return float32(temreture / 1000)
}

func getDataForWorkload() string {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return ""
	}
	return string(data)
}

func getLine(data string) []string {
	scanner := bufio.NewScanner(strings.NewReader(data))
	scanner.Scan()
	fields := strings.Fields(scanner.Text())
	return fields
}

func getCPUMeasurement(line []string) (int64, int64) {

	var values []int64
	for _, f := range line[1:] {
		v, _ := strconv.ParseInt(f, 10, 64)
		values = append(values, v)
	}

	user := values[0]
	nice := values[1]
	system := values[2]
	idle := values[3]
	iowait := values[4]
	irq := values[5]
	softirq := values[6]
	steal := values[7]

	total := user + nice + system + idle + iowait + irq + softirq + steal

	idle_total := idle + iowait

	return total, idle_total
}

func GetWorkload() float32 {
	total1, idle_total1 := getCPUMeasurement(getLine(getDataForWorkload()))
	time.Sleep(1 * time.Second)
	total2, idle_total2 := getCPUMeasurement(getLine(getDataForWorkload()))

	totalDiff := total2 - total1
	idleDiff := idle_total2 - idle_total1

	if totalDiff == 0 {
		return 0.0
	}
	load := 100.0 * float32(totalDiff-idleDiff) / float32(totalDiff)
	return load
}
