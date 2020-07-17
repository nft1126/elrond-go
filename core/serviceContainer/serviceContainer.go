package serviceContainer

import (
	"github.com/ElrondNetwork/elrond-go/core/indexer"
	"github.com/ElrondNetwork/elrond-go/core/statistics"
	"github.com/ElrondNetwork/elrond-go/scwatcher"
)

var _ Core = (*serviceContainer)(nil)

type serviceContainer struct {
	indexer         indexer.Indexer
	tpsBenchmark    statistics.TPSBenchmark
	scWatcherDriver scwatcher.Driver
}

// Option represents a functional configuration parameter that
//  can operate over the serviceContainer struct
type Option func(container *serviceContainer) error

// NewServiceContainer creates a new serviceContainer responsible in
//  providing access to all injected core features
func NewServiceContainer(opts ...Option) (Core, error) {
	sc := &serviceContainer{
		scWatcherDriver: &scwatcher.DisabledScWatcherDriver{},
	}
	for _, opt := range opts {
		err := opt(sc)
		if err != nil {
			return nil, err
		}
	}
	return sc, nil
}

// Indexer returns the core package's indexer
func (sc *serviceContainer) Indexer() indexer.Indexer {
	return sc.indexer
}

// TPSBenchmark returns the core package's tpsBenchmark
func (sc *serviceContainer) TPSBenchmark() statistics.TPSBenchmark {
	return sc.tpsBenchmark
}

// ScWatcherDriver returns the ScWatcher driver
func (sc *serviceContainer) ScWatcherDriver() scwatcher.Driver {
	return sc.scWatcherDriver
}

// WithIndexer sets up the database indexer
func WithIndexer(indexer indexer.Indexer) Option {
	return func(sc *serviceContainer) error {
		sc.indexer = indexer
		return nil
	}
}

// WithTPSBenchmark sets up the tpsBenchmark object
func WithTPSBenchmark(tpsBenchmark statistics.TPSBenchmark) Option {
	return func(sc *serviceContainer) error {
		sc.tpsBenchmark = tpsBenchmark
		return nil
	}
}

// WithScWatcherDriver sets up the Smart Contracts Driver
func WithScWatcherDriver(driver scwatcher.Driver) Option {
	return func(sc *serviceContainer) error {
		sc.scWatcherDriver = driver
		return nil
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (sc *serviceContainer) IsInterfaceNil() bool {
	return sc == nil
}
