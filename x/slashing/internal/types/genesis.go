package types

import (
	"fmt"
	"time"

	sdk "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/types"
)

// GenesisState - all slashing state that must be provided at genesis
type GenesisState struct {
	Params       Params                          `json:"params" yaml:"params"`
	SigningInfos map[string]ValidatorSigningInfo `json:"signing_infos" yaml:"signing_infos"`
	MissedBlocks map[string][]MissedBlock        `json:"missed_blocks" yaml:"missed_blocks"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params, signingInfos map[string]ValidatorSigningInfo, missedBlocks map[string][]MissedBlock,
) GenesisState {

	return GenesisState{
		Params:       params,
		SigningInfos: signingInfos,
		MissedBlocks: missedBlocks,
	}
}

// MissedBlock
type MissedBlock struct {
	Index  int64 `json:"index" yaml:"index"`
	Missed bool  `json:"missed" yaml:"missed"`
}

// NewMissedBlock creates a new MissedBlock instance
func NewMissedBlock(index int64, missed bool) MissedBlock {
	return MissedBlock{
		Index:  index,
		Missed: missed,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:       DefaultParams(),
		SigningInfos: make(map[string]ValidatorSigningInfo),
		MissedBlocks: make(map[string][]MissedBlock),
	}
}

// ValidateGenesis validates the slashing genesis parameters
func ValidateGenesis(data GenesisState) error {
	downtime := data.Params.SlashFractionDowntime
	if downtime.IsNegative() || downtime.GT(sdk.OneDec()) {
		return fmt.Errorf("slashing fraction downtime should be less than or equal to one and greater than zero, is %s", downtime.String())
	}

	dblSign := data.Params.SlashFractionDoubleSign
	if dblSign.IsNegative() || dblSign.GT(sdk.OneDec()) {
		return fmt.Errorf("slashing fraction double sign should be less than or equal to one and greater than zero, is %s", dblSign.String())
	}

	minSign := data.Params.MinSignedPerWindow
	if minSign.IsNegative() || minSign.GT(sdk.OneDec()) {
		return fmt.Errorf("min signed per window should be less than or equal to one and greater than zero, is %s", minSign.String())
	}

	downtimeJail := data.Params.DowntimeJailDuration
	if downtimeJail < 1*time.Minute {
		return fmt.Errorf("downtime unblond duration must be at least 1 minute, is %s", downtimeJail.String())
	}

	signedWindow := data.Params.SignedBlocksWindow
	if signedWindow < 10 {
		return fmt.Errorf("signed blocks window must be at least 10, is %d", signedWindow)
	}

	return nil
}
