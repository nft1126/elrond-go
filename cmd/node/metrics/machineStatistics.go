package metrics

import (
	"errors"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/appStatusPolling"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/core/statistics/machine"
)

// StartMachineStatisticsPolling will start read information about current  running machini
func StartMachineStatisticsPolling(ash core.AppStatusHandler, pollingInterval time.Duration) error {
	if check.IfNil(ash) {
		return errors.New("nil AppStatusHandler")
	}

	appStatusPollingHandler, err := appStatusPolling.NewAppStatusPolling(ash, pollingInterval)
	if err != nil {
		return errors.New("cannot init AppStatusPolling")
	}

	err = registerCpuStatistics(appStatusPollingHandler)
	if err != nil {
		return err
	}

	err = registerMemStatistics(appStatusPollingHandler)
	if err != nil {
		return err
	}

	err = registerNetStatistics(appStatusPollingHandler)
	if err != nil {
		return err
	}

	appStatusPollingHandler.Poll()

	return nil
}

func registerMemStatistics(appStatusPollingHandler *appStatusPolling.AppStatusPolling) error {
	return appStatusPollingHandler.RegisterPollingFunc(func(appStatusHandler core.AppStatusHandler) {
		mem := machine.AcquireMemStatistics()

		appStatusHandler.SetUInt64Value(core.MetricMemLoadPercent, mem.PercentUsed)
		appStatusHandler.SetUInt64Value(core.MetricMemTotal, mem.Total)
		appStatusHandler.SetUInt64Value(core.MetricMemUsedGolang, mem.UsedByGolang)
		appStatusHandler.SetUInt64Value(core.MetricMemUsedSystem, mem.UsedBySystem)
		appStatusHandler.SetUInt64Value(core.MetricMemHeapInUse, mem.HeapInUse)
		appStatusHandler.SetUInt64Value(core.MetricMemStackInUse, mem.StackInUse)
	})
}

func registerNetStatistics(appStatusPollingHandler *appStatusPolling.AppStatusPolling) error {
	var mutex sync.Mutex
	var netStats machine.NetStatistics

	return appStatusPollingHandler.RegisterPollingFunc(func(appStatusHandler core.AppStatusHandler) {
		mutex.Lock()
		defer mutex.Unlock()
		netStats = machine.AcquireNetStatistics(netStats)

		appStatusHandler.SetUInt64Value(core.MetricNetworkRecvBps, netStats.BpsReceived)
		appStatusHandler.SetUInt64Value(core.MetricNetworkRecvBpsPeak, netStats.BpsReceivedPeak)
		appStatusHandler.SetUInt64Value(core.MetricNetworkRecvPercent, netStats.ReceivedPercent)

		appStatusHandler.SetUInt64Value(core.MetricNetworkSentBps, netStats.BpsSent)
		appStatusHandler.SetUInt64Value(core.MetricNetworkSentBpsPeak, netStats.BpsSentPeak)
		appStatusHandler.SetUInt64Value(core.MetricNetworkSentPercent, netStats.SentPercent)
	})
}

func registerCpuStatistics(appStatusPollingHandler *appStatusPolling.AppStatusPolling) error {
	provider, err := machine.NewCpuStatisticsProvider()
	if err != nil {
		return err
	}

	return appStatusPollingHandler.RegisterPollingFunc(func(appStatusHandler core.AppStatusHandler) {
		cpuStats := provider.AcquireStatistics()
		appStatusHandler.SetUInt64Value(core.MetricCpuLoadPercent, cpuStats.CpuPercentUsage)
	})
}
