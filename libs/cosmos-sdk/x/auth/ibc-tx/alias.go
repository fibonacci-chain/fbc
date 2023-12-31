package ibc_tx

import (
	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/x/auth/ibc-tx/internal/adapter"
	ibccodec "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/x/auth/ibc-tx/internal/pb-codec"
)

var (
	PubKeyRegisterInterfaces = ibccodec.RegisterInterfaces
	LagacyKey2PbKey          = adapter.LagacyPubkey2ProtoBuffPubkey
)
