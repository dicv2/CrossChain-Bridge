package tokens

import (
	"fmt"
	"math/big"
)

// SwapType type
type SwapType uint32

// SwapType constants
const (
	NoSwapType SwapType = iota
	SwapinType
	SwapoutType
)

func (s SwapType) String() string {
	switch s {
	case NoSwapType:
		return "noswap"
	case SwapinType:
		return "swapin"
	case SwapoutType:
		return "swapout"
	default:
		return fmt.Sprintf("unknown swap type %d", s)
	}
}

// SwapTxType type
type SwapTxType uint32

// SwapTxType constants
const (
	SwapinTx     SwapTxType = iota // 0
	SwapoutTx                      // 1
	P2shSwapinTx                   // 2
)

func (s SwapTxType) String() string {
	switch s {
	case SwapinTx:
		return "swapintx"
	case SwapoutTx:
		return "swapouttx"
	case P2shSwapinTx:
		return "p2shswapintx"
	default:
		return fmt.Sprintf("unknown swaptx type %d", s)
	}
}

// TxSwapInfo struct
type TxSwapInfo struct {
	PairID    string   `json:"pairid"`
	Hash      string   `json:"hash"`
	Height    uint64   `json:"height"`
	Timestamp uint64   `json:"timestamp"`
	From      string   `json:"from"`
	TxTo      string   `json:"txto"`
	To        string   `json:"to"`
	Bind      string   `json:"bind"`
	Value     *big.Int `json:"value"`
}

// TxStatus struct
type TxStatus struct {
	Receipt       interface{} `json:"receipt,omitempty"`
	Confirmations uint64      `json:"confirmations"`
	BlockHeight   uint64      `json:"block_height"`
	BlockHash     string      `json:"block_hash"`
	BlockTime     uint64      `json:"block_time"`
}

// SwapInfo struct
type SwapInfo struct {
	PairID     string     `json:"pairid,omitempty"`
	SwapID     string     `json:"swapid,omitempty"`
	SwapType   SwapType   `json:"swaptype,omitempty"`
	TxType     SwapTxType `json:"txtype,omitempty"`
	Bind       string     `json:"bind,omitempty"`
	Identifier string     `json:"identifier,omitempty"`
	Reswapping bool       `json:"reswapping,omitempty"`
}

// IsSwapin is swapin type
func (s *SwapInfo) IsSwapin() bool {
	return s.SwapType == SwapinType
}

// BuildTxArgs struct
type BuildTxArgs struct {
	SwapInfo    `json:"swapInfo,omitempty"`
	From        string     `json:"from,omitempty"`
	To          string     `json:"to,omitempty"`
	OriginFrom  string     `json:"originFrom,omitempty"`
	OriginTxTo  string     `json:"originTxTo,omitempty"`
	Value       *big.Int   `json:"value,omitempty"`
	OriginValue *big.Int   `json:"originValue,omitempty"`
	SwapValue   *big.Int   `json:"swapvalue,omitempty"`
	Memo        string     `json:"memo,omitempty"`
	Input       *[]byte    `json:"input,omitempty"`
	Extra       *AllExtras `json:"extra,omitempty"`
}

// GetReplaceNum get rplace swap count
func (args *BuildTxArgs) GetReplaceNum() uint64 {
	if args.Extra != nil {
		return args.Extra.ReplaceNum
	}
	return 0
}

// GetExtraArgs get extra args
func (args *BuildTxArgs) GetExtraArgs() *BuildTxArgs {
	return &BuildTxArgs{
		SwapInfo: args.SwapInfo,
		Extra:    args.Extra,
	}
}

// GetTxGasPrice get tx gas price
func (args *BuildTxArgs) GetTxGasPrice() *big.Int {
	if args.Extra != nil && args.Extra.EthExtra != nil {
		return args.Extra.EthExtra.GasPrice
	}
	return nil
}

// SetTxGasPrice set tx gas price
func (args *BuildTxArgs) SetTxGasPrice(gasPrice *big.Int) {
	if args != nil && args.Extra != nil && args.Extra.EthExtra != nil {
		args.Extra.EthExtra.GasPrice = gasPrice
	}
}

// GetTxNonce get tx nonce
func (args *BuildTxArgs) GetTxNonce() uint64 {
	if args.Extra == nil {
		return 0
	}
	if args.Extra.EthExtra != nil && args.Extra.EthExtra.Nonce != nil {
		return *args.Extra.EthExtra.Nonce
	}
	if args.Extra.RippleExtra != nil && args.Extra.RippleExtra.Sequence != nil {
		return uint64(*args.Extra.RippleExtra.Sequence)
	}
	if args.Extra.TerraExtra != nil && args.Extra.TerraExtra.Sequence != nil {
		return uint64(*args.Extra.TerraExtra.Sequence)
	}
	return 0
}

// SetTxNonce set tx nonce
func (args *BuildTxArgs) SetTxNonce(nonce uint64) {
	switch {
	case args == nil || args.Extra == nil:
		return
	case args.Extra.EthExtra != nil:
		args.Extra.EthExtra.Nonce = &nonce
	case args.Extra.RippleExtra != nil:
		seq := uint32(nonce)
		args.Extra.RippleExtra.Sequence = &seq
	case args.Extra.TerraExtra != nil:
		args.Extra.TerraExtra.Sequence = &nonce
	}
}

func (args *BuildTxArgs) SetReplaceNum(replaceNum uint64) {
	if args != nil && args.Extra != nil {
		args.Extra.ReplaceNum = replaceNum
	}
}

// AllExtras struct
type AllExtras struct {
	ReplaceNum  uint64        `json:"replaceNum,omitempty"`
	BtcExtra    *BtcExtraArgs `json:"btcExtra,omitempty"`
	EthExtra    *EthExtraArgs `json:"ethExtra,omitempty"`
	RippleExtra *RippleExtra  `json:"rippleExtra,omitempty"`
	TerraExtra  *TerraExtra   `json:"terraExtra,omitempty"`
}

// EthExtraArgs struct
type EthExtraArgs struct {
	Gas       *uint64  `json:"gas,omitempty"`
	GasPrice  *big.Int `json:"gasPrice,omitempty"`
	GasTipCap *big.Int `json:"gasTipCap,omitempty"`
	GasFeeCap *big.Int `json:"gasFeeCap,omitempty"`
	Nonce     *uint64  `json:"nonce,omitempty"`
}

// RippleExtra struct
type RippleExtra struct {
	Sequence *uint32 `json:"sequence,omitempty"`
	Fee      *int64  `json:"fee,omitempty"`
}

// TerraExtra struct
type TerraExtra struct {
	Sequence *uint64 `json:"sequence,omitempty"`
	Fees     *string `json:"fees,omitempty"`
	Gas      *uint64 `json:"gas,omitempty"`
}

// BtcOutPoint struct
type BtcOutPoint struct {
	Hash  string `json:"hash"`
	Index uint32 `json:"index"`
}

// BtcExtraArgs struct
type BtcExtraArgs struct {
	RelayFeePerKb     *int64         `json:"relayFeePerKb,omitempty"`
	ChangeAddress     *string        `json:"-"`
	PreviousOutPoints []*BtcOutPoint `json:"previousOutPoints,omitempty"`
}

// P2shAddressInfo struct
type P2shAddressInfo struct {
	BindAddress        string
	P2shAddress        string
	RedeemScript       string
	RedeemScriptDisasm string
}
