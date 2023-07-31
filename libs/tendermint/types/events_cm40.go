package types

import (
	abci "github.com/fibonacci-chain/fbc/libs/tendermint/abci/types"
)

type CM40EventDataNewBlock struct {
	Block *CM40Block `json:"block"`

	ResultBeginBlock abci.ResponseBeginBlock `json:"result_begin_block"`
	ResultEndBlock   abci.ResponseEndBlock   `json:"result_end_block"`
}
