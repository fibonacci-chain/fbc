package vmbridge

import (
	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/codec"
	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/types/module"
	"github.com/fibonacci-chain/fbc/x/vmbridge/keeper"
	"github.com/fibonacci-chain/fbc/x/wasm"
)

func RegisterServices(cfg module.Configurator, keeper keeper.Keeper) {
	RegisterMsgServer(cfg.MsgServer(), NewMsgServerImpl(keeper))
}

func GetWasmOpts(cdc *codec.ProtoCodec) wasm.Option {
	return wasm.WithMessageEncoders(RegisterSendToEvmEncoder(cdc))
}
