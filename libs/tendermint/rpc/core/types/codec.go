package coretypes

import (
	amino "github.com/tendermint/go-amino"

	"github.com/fibonacci-chain/fbc/libs/tendermint/types"
)

func RegisterAmino(cdc *amino.Codec) {
	types.RegisterEventDatas(cdc)
	types.RegisterBlockAmino(cdc)
}
