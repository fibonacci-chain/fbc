package keeper_test

import (
	"testing"
	"time"

	"github.com/fibonacci-chain/fbc/app"
	"github.com/fibonacci-chain/fbc/app/crypto/ethsecp256k1"
	sdk "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/types"
	abci "github.com/fibonacci-chain/fbc/libs/tendermint/abci/types"
	"github.com/fibonacci-chain/fbc/x/feesplit/keeper"
	"github.com/fibonacci-chain/fbc/x/feesplit/types"
	"github.com/stretchr/testify/suite"
)

var (
	contract = ethsecp256k1.GenerateAddress()
	deployer = sdk.AccAddress(ethsecp256k1.GenerateAddress().Bytes())
	withdraw = sdk.AccAddress(ethsecp256k1.GenerateAddress().Bytes())
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.FBChainApp

	querier sdk.Querier
}

func (suite *KeeperTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(checkTx)
	suite.ctx = suite.app.NewContext(checkTx, abci.Header{
		Height:  1,
		ChainID: "ethermint-3",
		Time:    time.Now().UTC(),
	})
	suite.querier = keeper.NewQuerier(suite.app.FeeSplitKeeper)
	suite.app.FeeSplitKeeper.SetParams(suite.ctx, types.DefaultParams())
}
