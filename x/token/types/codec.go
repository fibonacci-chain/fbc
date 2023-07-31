package types

import (
	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgTokenIssue{}, "fbexchain/token/MsgIssue", nil)
	cdc.RegisterConcrete(MsgTokenBurn{}, "fbexchain/token/MsgBurn", nil)
	cdc.RegisterConcrete(MsgTokenMint{}, "fbexchain/token/MsgMint", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "fbexchain/token/MsgMultiTransfer", nil)
	cdc.RegisterConcrete(MsgSend{}, "fbexchain/token/MsgTransfer", nil)
	cdc.RegisterConcrete(MsgTransferOwnership{}, "fbexchain/token/MsgTransferOwnership", nil)
	cdc.RegisterConcrete(MsgConfirmOwnership{}, "fbexchain/token/MsgConfirmOwnership", nil)
	cdc.RegisterConcrete(MsgTokenModify{}, "fbexchain/token/MsgModify", nil)

	// for test
	//cdc.RegisterConcrete(MsgTokenDestroy{}, "fbexchain/token/MsgDestroy", nil)
}

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
