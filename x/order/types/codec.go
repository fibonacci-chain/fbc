package types

import (
	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgNewOrders{}, "fbexchain/order/MsgNew", nil)
	cdc.RegisterConcrete(MsgCancelOrders{}, "fbexchain/order/MsgCancel", nil)
}

// ModuleCdc generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
