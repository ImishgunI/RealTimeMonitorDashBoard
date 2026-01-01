package metrics

import (
	"testing"
	"github.com/go-playground/assert/v2"
)

func TestGetDataFromCpuInfo(t *testing.T) {
	data, err := getDataFromCpuInfo()
	if err != nil {
		t.Fatalf("Error from getDataFromCpuInfo() function: %v", err)
	}
	assert.NotEqual(t, data, "")
}
