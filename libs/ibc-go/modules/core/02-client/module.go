package client

import (
	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/codec"
	interfacetypes "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/codec/types"
	"github.com/fibonacci-chain/fbc/libs/ibc-go/modules/core/02-client/client/cli"
	"github.com/fibonacci-chain/fbc/libs/ibc-go/modules/core/02-client/types"
	"github.com/gogo/protobuf/grpc"
	"github.com/spf13/cobra"
)

// Name returns the IBC client name
func Name() string {
	return types.SubModuleName
}

// GetQueryCmd returns no root query command for the IBC client
func GetQueryCmd(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	return cli.GetQueryCmd(cdc, reg)
}

// GetTxCmd returns the root tx command for 02-client.
func GetTxCmd(cdc *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	return cli.NewTxCmd(cdc, reg)
}

// RegisterQueryService registers the gRPC query service for IBC client.
func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
	types.RegisterQueryServer(server, queryServer)
}
