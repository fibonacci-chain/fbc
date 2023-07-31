package simulation

import (
	"fmt"

	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/codec"
	clientsim "github.com/fibonacci-chain/fbc/libs/ibc-go/modules/core/02-client/simulation"
	connectionsim "github.com/fibonacci-chain/fbc/libs/ibc-go/modules/core/03-connection/simulation"
	channelsim "github.com/fibonacci-chain/fbc/libs/ibc-go/modules/core/04-channel/simulation"
	host "github.com/fibonacci-chain/fbc/libs/ibc-go/modules/core/24-host"
	"github.com/fibonacci-chain/fbc/libs/ibc-go/modules/core/keeper"
	"github.com/fibonacci-chain/fbc/libs/tendermint/libs/kv"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding ibc type.
func NewDecodeStore(k keeper.Keeper) func(cdc *codec.Codec, kvA, kvB kv.Pair) string {
	return func(cdc *codec.Codec, kvA, kvB kv.Pair) string {
		if res, found := clientsim.NewDecodeStore(k.ClientKeeper, kvA, kvB); found {
			return res
		}

		if res, found := connectionsim.NewDecodeStore(k.Codec(), kvA, kvB); found {
			return res
		}

		if res, found := channelsim.NewDecodeStore(k.Codec(), kvA, kvB); found {
			return res
		}

		panic(fmt.Sprintf("invalid %s key prefix: %s", host.ModuleName, string(kvA.Key)))
	}
}
