//go:build darwin

package main

import (
	"os/exec"
	"strconv"
	"strings"
)

// Thermal state thresholds based on kernel_task CPU usage
const (
	ThermalThresholdIdle      = 5.0
	ThermalThresholdLight     = 20.0
	ThermalThresholdHeavy     = 50.0
	ThermalThresholdThrottle  = 100.0
)

// GetKernelTaskCPU returns kernel_task's CPU usage percentage.
// It uses the 'top' command with two samples to get current CPU usage
// rather than cumulative usage since boot. This takes approximately
// 1 second to execute due to the sampling interval.
func GetKernelTaskCPU() (float64, error) {
	// Use top -l 2 -s 1 to get current CPU usage (not cumulative)
	// The first sample shows cumulative CPU since boot, 
	// the second sample shows current CPU usage
	cmd := exec.Command("top", "-l", "2", "-s", "1", "-pid", "0")
	output, err := cmd.Output()
	if err != nil {
		return 0.0, err
	}
	
	// Extract CPU% from the second sample (current usage)
	lines := strings.Split(string(output), "\n")
	var lastKernelLine string
	for _, line := range lines {
		if strings.Contains(line, "kernel_task") {
			lastKernelLine = line
		}
	}
	
	if lastKernelLine != "" {
		fields := strings.Fields(lastKernelLine)
		if len(fields) >= 3 {
			cpuStr := strings.TrimSuffix(fields[2], "%")
			cpu, err := strconv.ParseFloat(cpuStr, 64)
			if err != nil {
				return 0.0, err
			}
			return cpu, nil
		}
	}
	
	return 0.0, nil
}

// GetThermalState returns a human-readable thermal state based on kernel_task CPU usage.
// The states are:
//   - "Idle": CPU <= 5% (system is cool)
//   - "Light Load": 5% < CPU <= 20% (normal operation)
//   - "Heavy Load": 20% < CPU <= 50% (system under load)
//   - "Throttling": 50% < CPU <= 100% (thermal throttling active)
//   - "Heavy Throttling": CPU > 100% (severe thermal throttling)
func GetThermalState(cpu float64) string {
	switch {
	case cpu <= ThermalThresholdIdle:
		return "Idle"
	case cpu <= ThermalThresholdLight:
		return "Light Load"
	case cpu <= ThermalThresholdHeavy:
		return "Heavy Load"
	case cpu <= ThermalThresholdThrottle:
		return "Throttling"
	default:
		// 100%+ indicates serious thermal throttling
		return "Heavy Throttling"
	}
}