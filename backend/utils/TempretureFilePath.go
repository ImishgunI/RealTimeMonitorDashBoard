package utils

import (
	"os"
	"strings"
)

func FindFilePathForCPUTemreture() string {
	if isExist("/sys/class/thermal/thermal_zone0/temp") {
		return "/sys/class/thermal/thermal_zone0/temp"
	}
	name := "/sys/class/hwmon/"
	dirs, _ := os.ReadDir(name)
	c := 0
	names := []string{}
	for _, entry := range dirs {
		names = append(names, name+entry.Name()+"/")
		for c < len(names) {
			if isExist(names[c] + "name") {
				dev, err := readFile(names[c] + "name")
				if err != nil {
					break
				}
				if dev == "k10temp" || dev == "coretemp" {
					return names[c] + "temp1_input"
				} else {
					break
				}
			}
		}
		c++
	}
	return ""
}

func isExist(filepath string) bool {
	_, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return true
}

func readFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	str := strings.Trim(string(data), "\n")
	return str, nil
}
