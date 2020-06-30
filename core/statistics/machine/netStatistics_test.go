package machine

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAcquireNetStatistics_RelativelyToPreviousValues(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	var mutex sync.Mutex
	netStats := NetStatistics{}
	netStats.BpsSentPeak = 41
	netStats.BpsReceivedPeak = 41

	// Many routines in launched in parallel
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			mutex.Lock()
			if i == 3 {
				netStats.BpsSentPeak = 42
				netStats.BpsReceivedPeak = 42
			}

			netStats = AcquireNetStatistics(netStats)
			mutex.Unlock()

			time.Sleep(10 * time.Millisecond)
			wg.Done()
		}(i)
	}

	wg.Wait()

	require.Equal(t, uint64(42), netStats.BpsReceivedPeak)
	require.Equal(t, uint64(42), netStats.BpsSentPeak)
}
