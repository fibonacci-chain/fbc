package token

import (
	sdk "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/types"
	"github.com/fibonacci-chain/fbc/x/common/perf"
	"github.com/fibonacci-chain/fbc/x/token/types"
)

// BeginBlocker is called when dapp handles with abci::BeginBlock
func beginBlocker(ctx sdk.Context, keeper Keeper) {
	seq := perf.GetPerf().OnBeginBlockEnter(ctx, types.ModuleName)
	defer perf.GetPerf().OnBeginBlockExit(ctx, types.ModuleName, seq)

	keeper.ResetCache(ctx)
}
