package types

import (
	sdk "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/types"
	sdkerrors "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/types/errors"
)

// for denom convert wei to fibo and reject fibo direct
func (m *MsgTransfer) RulesFilter() (sdk.Msg, error) {
	if m.Token.Denom == sdk.DefaultBondDenom {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "ibc MsgTransfer not support fibo denom")
	}
	return m, nil
}
