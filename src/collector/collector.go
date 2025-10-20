package collector

import "time"

type MetricCollector struct {
	CPU       float64   `json:"cpu_usage"`
	MEM       float64   `json:"memory_usage"`
	MEMMB     uint64    `json:"memory_mb"`
	Timestamp time.Time `json:"timestamp"`
}

func NewMetricCollector() *MetricCollector {
	return &MetricCollector{
		CPU:       0.0,
		MEM:       0.0,
		MEMMB:     0,
		Timestamp: time.Time{},
	}
}
