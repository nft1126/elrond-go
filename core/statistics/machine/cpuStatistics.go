package machine

import (
	"errors"
	"fmt"

	"github.com/shirou/gopsutil/cpu"
)

var errCpuCount = errors.New("cpu count is zero")

// CpuStatistics holds CPU statistics
type CpuStatistics struct {
	CpuUsagePercent uint64
}

func (stats *CpuStatistics) String() string {
	return fmt.Sprintf("usage:%d%%", stats.CpuUsagePercent)
}

// CpuStatisticsProvider provides CPU statistics
type CpuStatisticsProvider struct {
	numCpu int
}

// NewCpuStatisticsProvider creates a CpuStatisticsProvider
func NewCpuStatisticsProvider() (*CpuStatisticsProvider, error) {
	numCpu, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}
	if numCpu == 0 {
		return nil, errCpuCount
	}

	return &CpuStatisticsProvider{
		numCpu: numCpu,
	}, nil
}

// AcquireStatistics acquires CPU statistics
func (provider *CpuStatisticsProvider) AcquireStatistics() CpuStatistics {
	currentProcess, err := GetCurrentProcess()
	if err != nil {
		return CpuStatistics{}
	}

	percent, err := currentProcess.Percent(0)
	if err != nil {
		return CpuStatistics{}
	}

	result := CpuStatistics{
		CpuUsagePercent: uint64(percent / float64(provider.numCpu)),
	}

	log.Trace("CpuStatisticsProvider.AcquireStatistics", "stats", result.String())
	return result
}
