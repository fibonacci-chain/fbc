package continuousauction

import (
	sdk "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/types"

	"github.com/fibonacci-chain/fbc/x/order/keeper"
)

// nolint
type CaEngine struct {
}

// nolint
func (e *CaEngine) Run(ctx sdk.Context, keeper keeper.Keeper) {
	// TODO
}
