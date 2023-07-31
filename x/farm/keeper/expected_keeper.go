package keeper

import (
	govtypes "github.com/fibonacci-chain/fbc/x/gov/types"

	sdk "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/types"
	"time"
)

// GovKeeper defines the expected gov Keeper
type GovKeeper interface {
	RemoveFromActiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time)
	GetDepositParams(ctx sdk.Context) govtypes.DepositParams
	GetVotingParams(ctx sdk.Context) govtypes.VotingParams
}
