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
	name, cores, threads, freq := GetValueForCPUMetrics()
	return &CPUMetrics{
		Name:      name,
		Cores:     int8(cores),
		Threads:   int8(threads),
		Frequency: freq / 1000,
		Temreture: GetTemretureForCPU(),
		Workload:  GetWorkload(),
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

func GetDataForWorkload() string {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return ""
	}
	return string(data)
}

func GetLine(data string) []string {
	scanner := bufio.NewScanner(strings.NewReader(data))
	scanner.Scan()
	fields := strings.Fields(scanner.Text())
	return fields
}

func GetCPUMeasurement(line []string) (int64, int64) {
	
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
	total1, idle_total1 := GetCPUMeasurement(GetLine(GetDataForWorkload()))
	time.Sleep(1 * time.Second)
	total2, idle_total2 := GetCPUMeasurement(GetLine(GetDataForWorkload()))

	totalDiff := total2 - total1
	idleDiff := idle_total2 - idle_total1

	if totalDiff == 0 {
		return 0.0
	}
	load := 100.0 * float32(totalDiff-idleDiff) / float32(totalDiff)
	return load
}
