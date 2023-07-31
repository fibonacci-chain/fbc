package bank

import (
	"strings"
	"testing"

	abci "github.com/fibonacci-chain/fbc/libs/tendermint/abci/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/types"
	sdkerrors "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/types/errors"
)

func TestInvalidMsg(t *testing.T) {
	h := NewHandler(nil)

	res, err := h(sdk.NewContext(nil, abci.Header{}, false, nil), sdk.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)

	_, _, log := sdkerrors.ABCIInfo(err, false)
	require.True(t, strings.Contains(log, "unrecognized bank message type"))
}
