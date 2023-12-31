package types

import (
	"encoding/json"

	abci "github.com/fibonacci-chain/fbc/libs/tendermint/abci/types"

	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/codec"
	sdk "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/types"
	authexported "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/x/auth/exported"
)

// StakingKeeper defines the expected staking keeper (noalias)
type StakingKeeper interface {
	ApplyAndReturnValidatorSetUpdates(sdk.Context) (updates []abci.ValidatorUpdate)
}

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	NewAccount(sdk.Context, authexported.Account) authexported.Account
	SetAccount(sdk.Context, authexported.Account)
	IterateAccounts(ctx sdk.Context, process func(authexported.Account) (stop bool))
}

// GenesisAccountsIterator defines the expected iterating genesis accounts object (noalias)
type GenesisAccountsIterator interface {
	IterateGenesisAccounts(
		cdc *codec.Codec,
		appGenesis map[string]json.RawMessage,
		iterateFn func(authexported.Account) (stop bool),
	)
}
