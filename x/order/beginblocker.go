package order

import (
	sdk "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/types"

	"github.com/fibonacci-chain/fbc/x/common/perf"
	"github.com/fibonacci-chain/fbc/x/order/keeper"
	"github.com/fibonacci-chain/fbc/x/order/types"
	//"github.com/fibonacci-chain/fbc/x/common/version"
)

// BeginBlocker runs the logic of BeginBlocker with version 0.
// BeginBlocker resets keeper cache.
func BeginBlocker(ctx sdk.Context, keeper keeper.Keeper) {
	seq := perf.GetPerf().OnBeginBlockEnter(ctx, types.ModuleName)
	defer perf.GetPerf().OnBeginBlockExit(ctx, types.ModuleName, seq)

	keeper.ResetCache(ctx)
}
