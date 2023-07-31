package client

import (
	"github.com/fibonacci-chain/fbc/x/farm/client/cli"
	"github.com/fibonacci-chain/fbc/x/farm/client/rest"
	govcli "github.com/fibonacci-chain/fbc/x/gov/client"
)

var (
	// ManageWhiteListProposalHandler alias gov NewProposalHandler
	ManageWhiteListProposalHandler = govcli.NewProposalHandler(cli.GetCmdManageWhiteListProposal, rest.ManageWhiteListProposalRESTHandler)
)
