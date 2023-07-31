package keeper_test

import (
	"testing"

	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/store/mpt"
	"github.com/fibonacci-chain/fbc/libs/tendermint/types"
	"github.com/stretchr/testify/suite"
)

type KeeperMptTestSuite struct {
	KeeperTestSuite
}

func (suite *KeeperMptTestSuite) SetupTest() {
	mpt.TrieWriteAhead = true
	types.UnittestOnlySetMilestoneMarsHeight(1)

	suite.KeeperTestSuite.SetupTest()
}

func TestKeeperMptTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperMptTestSuite))
}
