package params

import (
	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/codec"
	"github.com/fibonacci-chain/fbc/x/params/types"
)

// ModuleCdc is the codec of module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers all necessary param module types with a given codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(types.ParameterChangeProposal{}, "fbexchain/params/ParameterChangeProposal", nil)
	cdc.RegisterConcrete(types.UpgradeProposal{}, "fbexchain/params/UpgradeProposal", nil)
	cdc.RegisterConcrete(types.UpgradeInfo{}, "fbexchain/params/UpgradeInfo", nil)
}
