package simulation

import (
	"math/rand"

	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/x/simulation"
	"github.com/fibonacci-chain/fbc/libs/ibc-go/modules/core/04-channel/types"
)

// GenChannelGenesis returns the default channel genesis state.
func GenChannelGenesis(_ *rand.Rand, _ []simulation.Account) types.GenesisState {
	return types.DefaultGenesisState()
}
