package ethsecp256k1

import (
	cryptoamino "github.com/fibonacci-chain/fbc/libs/tendermint/crypto/encoding/amino"

	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/codec"
	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/crypto/keys"
)

// CryptoCodec is the default amino codec used by ethermint
var CryptoCodec = codec.New()

func init() {
	// replace the keyring codec with the ethermint crypto codec to prevent
	// amino panics because of unregistered Priv/PubKey
	keys.CryptoCdc = CryptoCodec
	keys.RegisterCodec(CryptoCodec)
	cryptoamino.RegisterAmino(CryptoCodec)
	RegisterCodec(CryptoCodec)
}

// RegisterCodec registers all the necessary types with amino for the given
// codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(PubKey{}, PubKeyName, nil)
	cdc.RegisterConcrete(PrivKey{}, PrivKeyName, nil)
}
