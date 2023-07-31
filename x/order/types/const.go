package types

import (
	sdk "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/types"
	"github.com/fibonacci-chain/fbc/x/common"
)

// nolint
const (
	FeeTypeOrderNew     = "new"
	FeeTypeOrderCancel  = "cancel"
	FeeTypeOrderExpire  = "expire"
	FeeTypeOrderDeal    = "deal"
	FeeTypeOrderReceive = "receive"
	TestTokenPair       = common.TestToken + "_" + sdk.DefaultBondDenom
	BuyOrder            = "BUY"
	SellOrder           = "SELL"
)
