package transfer

import (
	"github.com/fibonacci-chain/fbc/libs/ibc-go/modules/apps/transfer/keeper"
	"github.com/fibonacci-chain/fbc/libs/ibc-go/modules/apps/transfer/types"
)

var (
	NewKeeper  = keeper.NewKeeper
	ModuleCdc  = types.ModuleCdc
	SetMarshal = types.SetMarshal
)
