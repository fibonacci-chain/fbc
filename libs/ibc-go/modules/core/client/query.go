package client

import (
	"fmt"

	clictx "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/client/context"
	clienttypes "github.com/fibonacci-chain/fbc/libs/ibc-go/modules/core/02-client/types"
	commitmenttypes "github.com/fibonacci-chain/fbc/libs/ibc-go/modules/core/23-commitment/types"
	host "github.com/fibonacci-chain/fbc/libs/ibc-go/modules/core/24-host"
	abci "github.com/fibonacci-chain/fbc/libs/tendermint/abci/types"
)

// QueryTendermintProof performs an ABCI query with the given key and returns
// the value of the query, the proto encoded merkle proof, and the height of
// the Tendermint block containing the state root. The desired tendermint height
// to perform the query should be set in the client context. The query will be
// performed at one below this height (at the IAVL version) in order to obtain
// the correct merkle proof. Proof queries at height less than or equal to 2 are
// not supported. Queries with a client context height of 0 will perform a query
// at the lastest state available.
// Issue: https://github.com/cosmos/cosmos-sdk/issues/6567
func QueryTendermintProof(clientCtx clictx.CLIContext, key []byte) ([]byte, []byte, clienttypes.Height, error) {
	height := clientCtx.Height

	// ABCI queries at heights 1, 2 or less than or equal to 0 are not supported.
	// Base app does not support queries for height less than or equal to 1.
	// Therefore, a query at height 2 would be equivalent to a query at height 3.
	// A height of 0 will query with the lastest state.
	if height != 0 && height <= 2 {
		return nil, nil, clienttypes.Height{}, fmt.Errorf("proof queries at height <= 2 are not supported")
	}

	// Use the IAVL height if a valid tendermint height is passed in.
	// A height of 0 will query with the latest state.
	if height != 0 {
		height--
	}

	req := abci.RequestQuery{
		Path:   fmt.Sprintf("store/%s/key", host.StoreKey),
		Height: height,
		Data:   key,
		Prove:  true,
	}

	res, err := clientCtx.QueryABCI(req)
	if err != nil {
		return nil, nil, clienttypes.Height{}, err
	}

	merkleProof, err := commitmenttypes.ConvertProofs(res.Proof)
	if err != nil {
		return nil, nil, clienttypes.Height{}, err
	}

	proofBz, err := clientCtx.CodecProy.GetProtocMarshal().MarshalBinaryBare(&merkleProof)
	if err != nil {
		return nil, nil, clienttypes.Height{}, err
	}

	revision := clienttypes.ParseChainID(clientCtx.ChainID)
	return res.Value, proofBz, clienttypes.NewHeight(revision, uint64(res.Height)+1), nil
}
