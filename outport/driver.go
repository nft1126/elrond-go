package outport

import (
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/data"
)

var log = logger.GetOrCreate("outport")

var _ Driver = (*OutportDriver)(nil)

type OutportDriver struct {
	config        config.OutportConfig
	txCoordinator TransactionCoordinator
	logsProcessor TransactionLogProcessor
}

func NewOutportDriver(config config.OutportConfig, txCoordinator TransactionCoordinator, logsProcessor TransactionLogProcessor) (*OutportDriver, error) {
	if check.IfNil(txCoordinator) {
		return nil, ErrNilTxCoordinator
	}
	if check.IfNil(logsProcessor) {
		return nil, ErrNilLogsProcessor
	}

	return &OutportDriver{
		config:        config,
		txCoordinator: txCoordinator,
		logsProcessor: logsProcessor,
	}, nil
}

// DigestBlock digests a block
func (driver *OutportDriver) DigestCommittedBlock(header data.HeaderHandler, body data.BodyHandler) {
	if check.IfNil(header) {
		return
	}
	if check.IfNil(body) {
		return
	}

	// txPool := txCoordinator.GetAllCurrentUsedTxs(block.TxBlock)
	// scPool := txCoordinator.GetAllCurrentUsedTxs(block.SmartContractResultBlock)
	// rewardPool := txCoordinator.GetAllCurrentUsedTxs(block.RewardsBlock)
	// invalidPool := txCoordinator.GetAllCurrentUsedTxs(block.InvalidBlock)
	// receiptPool := txCoordinator.GetAllCurrentUsedTxs(block.ReceiptBlock)

	// fmt.Println("txPool", txPool)
	// fmt.Println("scPool", scPool)
	// fmt.Println("rewardPool", rewardPool)
	// fmt.Println("invalidPool", invalidPool)
	// fmt.Println("receiptPool", receiptPool)

	// Write to files (streams)
	// Example: https://github.com/ElrondNetwork/arwen-wasm-vm/pull/78/commits/3e23f1c44625363816cd3584fc64f01345be94b2
}

// IsInterfaceNil returns true if there is no value under the interface
func (driver *OutportDriver) IsInterfaceNil() bool {
	return driver == nil
}

// // TransactionsToDigest holds current transactions to digest
// type TransactionsToDigest struct {
// 	RegularTxs  map[string]data.TransactionHandler
// 	RewardTxs   map[string]data.TransactionHandler
// 	ScResults   map[string]data.TransactionHandler
// 	InvalidTxs  map[string]data.TransactionHandler
// 	ReceiptsTxs map[string]data.TransactionHandler
// }