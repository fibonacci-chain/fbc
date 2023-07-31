package client

import (
	"github.com/fibonacci-chain/fbc/x/feesplit/client/cli"
	"github.com/fibonacci-chain/fbc/x/feesplit/client/rest"
	govcli "github.com/fibonacci-chain/fbc/x/gov/client"
)

var (
	// FeeSplitSharesProposalHandler alias gov NewProposalHandler
	FeeSplitSharesProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdFeeSplitSharesProposal,
		rest.FeeSplitSharesProposalRESTHandler,
	)
)
