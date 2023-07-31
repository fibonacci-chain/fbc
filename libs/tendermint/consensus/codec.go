package consensus

import (
	amino "github.com/tendermint/go-amino"

	"github.com/fibonacci-chain/fbc/libs/tendermint/types"
)

var cdc = amino.NewCodec()

func init() {
	RegisterMessages(cdc)
	RegisterWALMessages(cdc)
	types.RegisterBlockAmino(cdc)
}
