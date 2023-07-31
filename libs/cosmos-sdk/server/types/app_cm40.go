package app

import (
	cliContext "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/client/context"
	gogogrpc "github.com/gogo/protobuf/grpc"
)

type ApplicationAdapter interface {
	RegisterGRPCServer(gogogrpc.Server)
	RegisterTxService(clientCtx cliContext.CLIContext)
}
