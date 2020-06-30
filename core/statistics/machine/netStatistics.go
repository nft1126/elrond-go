package machine

import (
	"fmt"
	"time"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/shirou/gopsutil/net"
)

// NetStatistics holds network statistics
type NetStatistics struct {
	timestamp     time.Time
	totalReceived uint64
	totalSent     uint64

	BpsReceived     uint64
	BpsReceivedPeak uint64
	ReceivedPercent uint64
	BpsSent         uint64
	BpsSentPeak     uint64
	SentPercent     uint64
}

func (stats *NetStatistics) String() string {
	return fmt.Sprintf(
		"received: %s/s (peak=%s/s, usage=%d%%), sent: %s/s (peak=%s/s, usage=%d%%)",
		core.ConvertBytes(stats.BpsReceived),
		core.ConvertBytes(stats.BpsReceivedPeak),
		stats.ReceivedPercent,
		core.ConvertBytes(stats.BpsSent),
		core.ConvertBytes(stats.BpsSentPeak),
		stats.SentPercent,
	)
}

// AcquireNetStatistics acquires the current network statistics (usage), relatively to previously acquired statistics
func AcquireNetStatistics(previous NetStatistics) NetStatistics {
	timestamp := time.Now()
	timeDelta := timestamp.Sub(previous.timestamp).Seconds()

	netCounters, err := net.IOCounters(false)
	if err != nil {
		return NetStatistics{}
	}
	if len(netCounters) == 0 {
		return NetStatistics{}
	}

	totalReceived := netCounters[0].BytesRecv
	totalSent := netCounters[0].BytesSent

	isLessReceived := totalReceived < previous.totalReceived
	isLessSent := totalSent < previous.totalSent
	if isLessReceived || isLessSent {
		return NetStatistics{}
	}

	bpsReceived := uint64(float64(totalReceived-previous.totalReceived) / timeDelta)
	bpsSent := uint64(float64(totalSent-previous.totalSent) / timeDelta)

	bpsReceivedPeak := previous.BpsReceivedPeak
	if bpsReceived > bpsReceivedPeak {
		bpsReceivedPeak = bpsReceived
	}

	bpsSentPeak := previous.BpsSentPeak
	if bpsSent > bpsSentPeak {
		bpsSentPeak = bpsSent
	}

	receivedPercent := uint64(0)
	if bpsReceivedPeak != 0 {
		receivedPercent = bpsReceived * 100 / bpsReceivedPeak
	}

	sentPercent := uint64(0)
	if bpsSentPeak != 0 {
		sentPercent = bpsSent * 100 / bpsSentPeak
	}

	result := NetStatistics{
		timestamp:     timestamp,
		totalReceived: totalReceived,
		totalSent:     totalSent,

		BpsReceived:     bpsReceived,
		BpsReceivedPeak: bpsReceivedPeak,
		ReceivedPercent: receivedPercent,
		BpsSent:         bpsSent,
		BpsSentPeak:     bpsSentPeak,
		SentPercent:     sentPercent,
	}

	log.Trace("AcquireNetStatistics", "stats", result.String())
	return result
}
