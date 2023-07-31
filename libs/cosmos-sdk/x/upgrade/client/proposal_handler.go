package client

import (
	govclient "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/x/gov/client"
	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/x/upgrade/client/cli"
	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/x/upgrade/client/rest"
)

var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitUpgradeProposal, rest.ProposalRESTHandler)
