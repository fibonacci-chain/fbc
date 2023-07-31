package core

import (
	"errors"
	"fmt"
	"sort"

	sdkerrors "github.com/fibonacci-chain/fbc/libs/cosmos-sdk/types/errors"

	nullblockindexer "github.com/fibonacci-chain/fbc/libs/tendermint/state/indexer/block/null"

	tmquery "github.com/fibonacci-chain/fbc/libs/tendermint/libs/pubsub/query"

	tmmath "github.com/fibonacci-chain/fbc/libs/tendermint/libs/math"
	ctypes "github.com/fibonacci-chain/fbc/libs/tendermint/rpc/core/types"
	rpctypes "github.com/fibonacci-chain/fbc/libs/tendermint/rpc/jsonrpc/types"
	sm "github.com/fibonacci-chain/fbc/libs/tendermint/state"
	"github.com/fibonacci-chain/fbc/libs/tendermint/types"
)

// BlockchainInfo gets block headers for minHeight <= height <= maxHeight.
// Block headers are returned in descending order (highest first).
// More: https://docs.tendermint.com/master/rpc/#/Info/blockchain
func BlockchainInfo(ctx *rpctypes.Context, minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
	// maximum 20 block metas
	const limit int64 = 20
	var err error
	minHeight, maxHeight, err = filterMinMax(
		env.BlockStore.Base(),
		env.BlockStore.Height(),
		minHeight,
		maxHeight,
		limit)
	if err != nil {
		return nil, err
	}
	env.Logger.Debug("BlockchainInfoHandler", "maxHeight", maxHeight, "minHeight", minHeight)

	blockMetas := make([]*types.BlockMeta, 0, maxHeight-minHeight+1)
	for height := maxHeight; height >= minHeight; height-- {
		blockMeta := env.BlockStore.LoadBlockMeta(height)
		blockMetas = append(blockMetas, blockMeta)
	}

	return &ctypes.ResultBlockchainInfo{
		LastHeight: env.BlockStore.Height(),
		BlockMetas: blockMetas}, nil
}

func LatestBlockNumber() (int64, error) {
	return env.BlockStore.Height(), nil
}

// error if either min or max are negative or min > max
// if 0, use blockstore base for min, latest block height for max
// enforce limit.
func filterMinMax(base, height, min, max, limit int64) (int64, int64, error) {
	// filter negatives
	if min < 0 || max < 0 {
		return min, max, fmt.Errorf("heights must be non-negative")
	}

	// adjust for default values
	if min == 0 {
		min = 1
	}
	if max == 0 {
		max = height
	}

	// limit max to the height
	max = tmmath.MinInt64(height, max)

	// limit min to the base
	min = tmmath.MaxInt64(base, min)

	// limit min to within `limit` of max
	// so the total number of blocks returned will be `limit`
	min = tmmath.MaxInt64(min, max-limit+1)

	if min > max {
		return min, max, fmt.Errorf("min height %d can't be greater than max height %d", min, max)
	}
	return min, max, nil
}

// Block gets block at a given height.
// If no height is provided, it will fetch the latest block.
// More: https://docs.tendermint.com/master/rpc/#/Info/block
func Block(ctx *rpctypes.Context, heightPtr *int64) (*ctypes.ResultBlock, error) {
	height, err := getHeight(env.BlockStore.Height(), heightPtr)
	if err != nil {
		return nil, err
	}
	block := env.BlockStore.LoadBlock(height)
	blockMeta := env.BlockStore.LoadBlockMeta(height)
	if blockMeta == nil {
		return &ctypes.ResultBlock{BlockID: types.BlockID{}, Block: block}, nil
	}
	return &ctypes.ResultBlock{BlockID: blockMeta.BlockID, Block: block}, nil
}

// BlockByHash gets block by hash.
// More: https://docs.tendermint.com/master/rpc/#/Info/block_by_hash
func BlockByHash(ctx *rpctypes.Context, hash []byte) (*ctypes.ResultBlock, error) {
	block := env.BlockStore.LoadBlockByHash(hash)
	if block == nil {
		return &ctypes.ResultBlock{BlockID: types.BlockID{}, Block: nil}, nil
	}
	// If block is not nil, then blockMeta can't be nil.
	blockMeta := env.BlockStore.LoadBlockMeta(block.Height)
	return &ctypes.ResultBlock{BlockID: blockMeta.BlockID, Block: block}, nil
}

// Commit gets block commit at a given height.
// If no height is provided, it will fetch the commit for the latest block.
// More: https://docs.tendermint.com/master/rpc/#/Info/commit
func Commit(ctx *rpctypes.Context, heightPtr *int64) (*ctypes.ResultCommit, error) {
	height, err := getHeight(env.BlockStore.Height(), heightPtr)
	if err != nil {
		return nil, err
	}

	blockMeta := env.BlockStore.LoadBlockMeta(height)
	if blockMeta == nil {
		return nil, nil
	}
	header := blockMeta.Header

	// If the next block has not been committed yet,
	// use a non-canonical commit
	if height == env.BlockStore.Height() {
		commit := env.BlockStore.LoadSeenCommit(height)
		return ctypes.NewResultCommit(&header, commit, false), nil
	}

	// Return the canonical commit (comes from the block at height+1)
	commit := env.BlockStore.LoadBlockCommit(height)
	return ctypes.NewResultCommit(&header, commit, true), nil
}

// BlockResults gets ABCIResults at a given height.
// If no height is provided, it will fetch results for the latest block.
//
// Results are for the height of the block containing the txs.
// Thus response.results.deliver_tx[5] is the results of executing
// getBlock(h).Txs[5]
// More: https://docs.tendermint.com/master/rpc/#/Info/block_results
func BlockResults(ctx *rpctypes.Context, heightPtr *int64) (*ctypes.ResultBlockResults, error) {
	height, err := getHeight(env.BlockStore.Height(), heightPtr)
	if err != nil {
		return nil, err
	}

	results, err := sm.LoadABCIResponses(env.StateDB, height)
	if err != nil {
		return nil, err
	}

	return &ctypes.ResultBlockResults{
		Height:                height,
		TxsResults:            results.DeliverTxs,
		BeginBlockEvents:      results.BeginBlock.Events,
		EndBlockEvents:        results.EndBlock.Events,
		ValidatorUpdates:      results.EndBlock.ValidatorUpdates,
		ConsensusParamUpdates: results.EndBlock.ConsensusParamUpdates,
	}, nil
}

// Header gets block header at a given height.
// If no height is provided, it will fetch the latest header.
// More: https://docs.tendermint.com/master/rpc/#/Info/header
func BlockInfo(ctx *rpctypes.Context, heightPtr *int64) (*types.BlockMeta, error) {
	height, err := getHeight(env.BlockStore.Height(), heightPtr)
	if err != nil {
		return nil, err
	}

	return env.BlockStore.LoadBlockMeta(height), nil
}

func BlockSearch(
	ctx *rpctypes.Context,
	query string,
	pagePtr, perPagePtr *int,
	orderBy string,
) (*ctypes.ResultBlockSearch, error) {

	// if index is disabled, return error
	if _, ok := env.BlockIndexer.(*nullblockindexer.BlockerIndexer); ok {
		return nil, errors.New("indexing is disabled")
	}

	q, err := tmquery.New(query)
	if err != nil {
		return nil, err
	}

	results, err := env.BlockIndexer.Search(ctx.Context(), q)
	if err != nil {
		return nil, err
	}

	// sort results (must be done before pagination)
	switch orderBy {
	case "desc", "":
		sort.Slice(results, func(i, j int) bool { return results[i] > results[j] })

	case "asc":
		sort.Slice(results, func(i, j int) bool { return results[i] < results[j] })

	default:
		return nil, fmt.Errorf("expected order_by to be either `asc` or `desc` or empty: %w", sdkerrors.ErrInvalidRequest)
	}

	// paginate results
	totalCount := len(results)
	perPage := validatePerPage(*perPagePtr)

	page, err := validatePage(*pagePtr, perPage, totalCount)
	if err != nil {
		return nil, err
	}

	skipCount := validateSkipCount(page, perPage)
	pageSize := tmmath.MinInt(perPage, totalCount-skipCount)

	apiResults := make([]*ctypes.ResultBlock, 0, pageSize)
	for i := skipCount; i < skipCount+pageSize; i++ {
		block := env.BlockStore.LoadBlock(results[i])
		if block != nil {
			blockMeta := env.BlockStore.LoadBlockMeta(block.Height)
			if blockMeta != nil {
				apiResults = append(apiResults, &ctypes.ResultBlock{
					Block:   block,
					BlockID: blockMeta.BlockID,
				})
			}
		}
	}

	return &ctypes.ResultBlockSearch{Blocks: apiResults, TotalCount: totalCount}, nil
}
