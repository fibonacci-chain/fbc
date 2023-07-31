package client

import (
	"github.com/fibonacci-chain/fbc/x/dex/client/cli"
	"github.com/fibonacci-chain/fbc/x/dex/client/rest"
	govclient "github.com/fibonacci-chain/fbc/x/gov/client"
)

// param change proposal handler
var (
	// DelistProposalHandler alias gov NewProposalHandler
	DelistProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitDelistProposal, rest.DelistProposalRESTHandler)
)
