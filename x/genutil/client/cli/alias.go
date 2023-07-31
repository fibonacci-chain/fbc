package cli

import (
	genutilcli "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/x/genutil/client/cli"
)

type (
	stakingMsgBuildingHelpers = genutilcli.StakingMsgBuildingHelpers
)

var (
	// nolint
	ValidateGenesisCmd = genutilcli.ValidateGenesisCmd
)
