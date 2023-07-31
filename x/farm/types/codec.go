package types

import (
	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreatePool{}, "fbexchain/farm/MsgCreatePool", nil)
	cdc.RegisterConcrete(MsgDestroyPool{}, "fbexchain/farm/MsgDestroyPool", nil)
	cdc.RegisterConcrete(MsgLock{}, "fbexchain/farm/MsgLock", nil)
	cdc.RegisterConcrete(MsgUnlock{}, "fbexchain/farm/MsgUnlock", nil)
	cdc.RegisterConcrete(MsgClaim{}, "fbexchain/farm/MsgClaim", nil)
	cdc.RegisterConcrete(MsgProvide{}, "fbexchain/farm/MsgProvide", nil)
	cdc.RegisterConcrete(ManageWhiteListProposal{}, "fbexchain/farm/ManageWhiteListProposal", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
