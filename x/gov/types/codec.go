package types

import (
	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/codec"
)

// module codec
var ModuleCdc = codec.New()

// RegisterCodec registers all the necessary types and interfaces for
// governance.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Content)(nil), nil)

	cdc.RegisterConcrete(MsgSubmitProposal{}, "fbexchain/gov/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "fbexchain/gov/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "fbexchain/gov/MsgVote", nil)

	cdc.RegisterConcrete(TextProposal{}, "fbexchain/gov/TextProposal", nil)
	cdc.RegisterConcrete(SoftwareUpgradeProposal{}, "fbexchain/gov/SoftwareUpgradeProposal", nil)
}

// RegisterProposalTypeCodec registers an external proposal content type defined
// in another module for the internal ModuleCdc. This allows the MsgSubmitProposal
// to be correctly Amino encoded and decoded.
func RegisterProposalTypeCodec(o interface{}, name string) {
	ModuleCdc.RegisterConcrete(o, name, nil)
}

// TODO determine a good place to seal this codec
func init() {
	RegisterCodec(ModuleCdc)
}
