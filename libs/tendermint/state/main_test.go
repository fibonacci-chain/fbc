package state_test

import (
	"os"
	"testing"

	"github.com/fibonacci-chain/fbc/libs/tendermint/types"
)

func TestMain(m *testing.M) {
	types.RegisterMockEvidencesGlobal()
	os.Exit(m.Run())
}
