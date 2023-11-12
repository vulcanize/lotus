package ethtypes

import (
	"encoding/json"
	"math/big"

	logger "github.com/ipfs/go-log/v2"
)

// TraceConfig holds extra parameters to trace functions.
type TraceConfig struct {
	*logger.Config
	Tracer  *string
	Timeout *string
	Reexec  *uint64
	// Config specific to given tracer. Note struct logger
	// config are historically embedded in main object.
	TracerConfig json.RawMessage
}

// TraceCallConfig is the config for traceCall API. It holds one more
// field to override the state for tracing.
type TraceCallConfig struct {
	TraceConfig
	StateOverrides *StateOverride
	BlockOverrides *BlockOverrides
}

// OverrideAccount indicates the overriding fields of account during the execution
// of a message call.
// Note, state and stateDiff can't be specified at the same time. If state is
// set, message execution will only use the data in the given state. Otherwise
// if statDiff is set, all diff will be applied first and then execute the call
// message.
type OverrideAccount struct {
	Nonce     *EthUint64           `json:"nonce"`
	Code      *EthBytes            `json:"code"`
	Balance   **EthBigInt          `json:"balance"`
	State     *map[EthHash]EthHash `json:"state"`
	StateDiff *map[EthHash]EthHash `json:"stateDiff"`
}

// StateOverride is the collection of overridden accounts.
type StateOverride map[EthHash]OverrideAccount

// BlockOverrides is a set of header fields to override.
type BlockOverrides struct {
	Number     *EthBigInt
	Difficulty *EthBigInt
	Time       *EthUint64
	GasLimit   *EthUint64
	Coinbase   *EthAddress
	Random     *EthHash
	BaseFee    *EthUint64
}

// Apply overrides the given header fields into the given block context.
func (diff *BlockOverrides) Apply(blockCtx *BlockContext) {
	if diff == nil {
		return
	}
	if diff.Number != nil {
		blockCtx.BlockNumber = diff.Number.Int
	}
	if diff.Difficulty != nil {
		blockCtx.Difficulty = diff.Difficulty.Int
	}
	if diff.Time != nil {
		blockCtx.Time = uint64(*diff.Time)
	}
	if diff.GasLimit != nil {
		blockCtx.GasLimit = uint64(*diff.GasLimit)
	}
	if diff.Coinbase != nil {
		blockCtx.Coinbase = *diff.Coinbase
	}
	if diff.Random != nil {
		blockCtx.Random = diff.Random
	}
	if diff.BaseFee != nil {
		blockCtx.BaseFee = new(big.Int).SetUint64((uint64)(*diff.BaseFee))
	}
}

type (
	// CanTransferFunc is the signature of a transfer guard function
	CanTransferFunc func(StateDB, EthAddress, *big.Int) bool
	// TransferFunc is the signature of a transfer function
	TransferFunc func(StateDB, EthAddress, EthAddress, *big.Int)
	// GetHashFunc returns the n'th block hash in the blockchain
	// and is used by the BLOCKHASH EVM op code.
	GetHashFunc func(uint64) EthAddress
)

// BlockContext provides the EVM with auxiliary information. Once provided
// it shouldn't be modified.
type BlockContext struct {
	// CanTransfer returns whether the account contains
	// sufficient ether to transfer the value
	CanTransfer CanTransferFunc
	// Transfer transfers ether from one account to the other
	Transfer TransferFunc
	// GetHash returns the hash corresponding to n
	GetHash GetHashFunc

	// Block information
	Coinbase    EthAddress // Provides information for COINBASE
	GasLimit    uint64     // Provides information for GASLIMIT
	BlockNumber *big.Int   // Provides information for NUMBER
	Time        uint64     // Provides information for TIME
	Difficulty  *big.Int   // Provides information for DIFFICULTY
	BaseFee     *big.Int   // Provides information for BASEFEE
	Random      *EthHash   // Provides information for PREVRANDAO
}

// TxTraceResult is the result of a single transaction trace.
type TxTraceResult struct {
	Result interface{} `json:"result,omitempty"` // Trace results produced by the tracer
	Error  string      `json:"error,omitempty"`  // Trace failure produced by the tracer
}
