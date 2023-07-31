package keeper

import (
	sdk "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/types"
	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/x/mint/internal/types"
	"github.com/pkg/errors"
	"time"
)

func (k Keeper) AddYieldFarming(ctx sdk.Context, yieldAmt sdk.Coins) error {
	// todo: verify farmModuleName
	if len(k.farmModuleName) == 0 {
		return nil
	}
	return k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.farmModuleName, yieldAmt)
}

// get the minter custom
func (k Keeper) GetMinterCustom(ctx sdk.Context) (minter types.MinterCustom) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.MinterKey)
	if b != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &minter)
	}
	return
}

// set the minter custom
func (k Keeper) SetMinterCustom(ctx sdk.Context, minter types.MinterCustom) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(minter)
	store.Set(types.MinterKey, b)
}

// UpdateMinterCustomToRightJupiter occur 2023year
func (k Keeper) UpdateMinterCustomToRightJupiter(ctx sdk.Context, minter *types.MinterCustom, params types.Params) {
	//update param to right
	params.DeflationEpoch = types.FiboChainDeflationEpoch
	params.DeflationRate = types.FiboChainDeflationRate

	firstYearPerBlock := sdk.MustNewDecFromStr("1.2506231842725684") //first year mint per block
	genesisTimeStr := "2022-02-01 00:58:12"
	rightNextBlockToUpdateTimeStr := "2024-02-01 00:58:12"

	genesisTime, err := time.Parse("2006-04-02 15:04:05", genesisTimeStr)
	if err != nil {
		panic("parse genesisTime err: " + err.Error())
	}

	rightNextBlockToUpdateTime, err := time.Parse("2006-04-02 15:04:05", rightNextBlockToUpdateTimeStr)
	if err != nil {
		panic("parse rightNextBlockToUpdateTime err: " + err.Error())
	}

	shouldNextBlockToUpdateJupiterHeight := uint64(rightNextBlockToUpdateTime.Sub(genesisTime).Seconds() / 3) //assume 3s per block
	minter.MintedPerBlock = sdk.NewDecCoinsFromDec(params.MintDenom, firstYearPerBlock)

	var provisionAmtPerBlock sdk.Dec
	if ctx.BlockHeight() == 0 || minter.NextBlockToUpdate == 0 {
		provisionAmtPerBlock = k.GetOriginalMintedPerBlock()
	} else {
		provisionAmtPerBlock = minter.MintedPerBlock.AmountOf(params.MintDenom).Mul(params.DeflationRate)
	}

	// update new MinterCustom
	minter.MintedPerBlock = sdk.NewDecCoinsFromDec(params.MintDenom, provisionAmtPerBlock)
	minter.NextBlockToUpdate = shouldNextBlockToUpdateJupiterHeight

	k.SetParams(ctx, params)
	k.SetMinterCustom(ctx, *minter)
}

// UpdateMinterCustom every year deflation rate 20%, assume year blockNumber 10519200
// total mint:  7777_7777 - 1200_0000 = 65777777
// 2022 1year: 65777777 * 0.2 = 13155555.4 / 10519200 = 1.2506231842725684
// 2023 2year: (65777777 - 13155555.4) * 0.2 = 10524444.32 / 10519200 = 1.0004985474180546
// 2024 3year: (65777777 - 13155555.4-  10524444.32) * 0.2 = 8419555.456 / 10519200 = 0.8003988379344437
// 2025 4year (65777777 - 13155555.4-  10524444.32 - 8419555.456) * 0.2 = 6735644.364800001 / 10519200 = 0.640319070347555
// 2026 4year (65777777 - 13155555.4-  10524444.32 - 8419555.456 - 6735644.364800001) * 0.2 = 5388515.491840001 / 10519200 = 0.512255256278044
// ...
func (k Keeper) UpdateMinterCustom(ctx sdk.Context, minter *types.MinterCustom, params types.Params) {
	maxMintYear := types.FiboChainMaxMintYear

	//stop minting
	maxMintBlock := int64(maxMintYear * params.BlocksPerYear) // 77 * 10519200 = 809978400

	if ctx.BlockHeight() > maxMintBlock {
		minter.MintedPerBlock = sdk.NewDecCoinsFromDec(params.MintDenom, sdk.NewDec(0))
		minter.NextBlockToUpdate += params.DeflationEpoch * params.BlocksPerYear
		k.SetMinterCustom(ctx, *minter)
		return
	}

	var provisionAmtPerBlock sdk.Dec
	if ctx.BlockHeight() == 0 || minter.NextBlockToUpdate == 0 {
		provisionAmtPerBlock = k.GetOriginalMintedPerBlock()
	} else {
		provisionAmtPerBlock = minter.MintedPerBlock.AmountOf(params.MintDenom).Mul(params.DeflationRate)
	}

	// update new MinterCustom
	minter.MintedPerBlock = sdk.NewDecCoinsFromDec(params.MintDenom, provisionAmtPerBlock)
	minter.NextBlockToUpdate += params.DeflationEpoch * params.BlocksPerYear

	//if tmtypes.HigherThanVenus5(ctx.BlockHeight()) {
	//	minter.NextBlockToUpdate += params.DeflationEpoch * params.BlocksPerYear / 12
	//} else {
	//	minter.NextBlockToUpdate += params.DeflationEpoch * params.BlocksPerYear
	//}

	k.SetMinterCustom(ctx, *minter)
}

//______________________________________________________________________

// GetOriginalMintedPerBlock returns the init tokens per block.
func (k Keeper) GetOriginalMintedPerBlock() sdk.Dec {
	return k.originalMintedPerBlock
}

// SetOriginalMintedPerBlock sets the init tokens per block.
func (k Keeper) SetOriginalMintedPerBlock(originalMintedPerBlock sdk.Dec) {
	k.originalMintedPerBlock = originalMintedPerBlock
}

// ValidateMinterCustom validate minter
func ValidateOriginalMintedPerBlock(originalMintedPerBlock sdk.Dec) error {
	if originalMintedPerBlock.IsNegative() {
		panic("init tokens per block must be non-negative")
	}

	return nil
}

// SetTreasures set the treasures to db
func (k Keeper) SetTreasures(ctx sdk.Context, treasures []types.Treasure) {
	store := ctx.KVStore(k.storeKey)
	types.SortTreasures(treasures)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(treasures)
	store.Set(types.TreasuresKey, b)
}

// GetTreasures get the treasures from db
func (k Keeper) GetTreasures(ctx sdk.Context) (treasures []types.Treasure) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.TreasuresKey)
	if b != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &treasures)
	}
	return
}

// AllocateTokenToTreasure allocate token to treasure and return remain
func (k Keeper) AllocateTokenToTreasure(ctx sdk.Context, fees sdk.Coins) (remain sdk.Coins, err error) {
	treasures := k.GetTreasures(ctx)
	remain = sdk.NewCoins()
	remain = remain.Add(fees...)
	for i, _ := range treasures {
		allocated := fees.MulDecTruncate(treasures[i].Proportion)
		if err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, treasures[i].Address, allocated); err != nil {
			return
		}
		remain = remain.Sub(allocated)
		if remain.IsAnyNegative() {
			return remain, errors.New("allocate coin is more than mint coin")
		}
		k.Logger(ctx).Debug("allocate treasure", "addr", treasures[i].Address, "proportion", treasures[i].Proportion, "sum coins", fees.String(), "allocated", allocated.String(), "remain", remain.String())
	}
	return
}

func (k Keeper) UpdateTreasures(ctx sdk.Context, treasures []types.Treasure) error {
	src := k.GetTreasures(ctx)
	result := types.InsertAndUpdateTreasures(src, treasures)
	if err := types.ValidateTreasures(result); err != nil {
		return err
	}
	k.SetTreasures(ctx, result)
	return nil
}

func (k Keeper) DeleteTreasures(ctx sdk.Context, treasures []types.Treasure) error {
	src := k.GetTreasures(ctx)
	result, err := types.DeleteTreasures(src, treasures)
	if err != nil {
		return err
	}
	if err := types.ValidateTreasures(result); err != nil {
		return err
	}
	k.SetTreasures(ctx, result)
	return nil
}
