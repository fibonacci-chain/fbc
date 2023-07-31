package types

import "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/codec"

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
