package machine

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCpuStatisticsUsagePercent(t *testing.T) {
	t.Parallel()

	provider, err := NewCpuStatisticsProvider()
	require.Nil(t, err)

	stats := provider.AcquireStatistics()
	require.True(t, stats.CpuUsagePercent <= 100)
}
