package config

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/server"
	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/store/iavl"
	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/store/types"
	tmiavl "github.com/fibonacci-chain/fbc/libs/iavl"
	iavlconfig "github.com/fibonacci-chain/fbc/libs/iavl/config"
	"github.com/fibonacci-chain/fbc/libs/system"
	"github.com/fibonacci-chain/fbc/libs/system/trace"
	tmconfig "github.com/fibonacci-chain/fbc/libs/tendermint/config"
	"github.com/fibonacci-chain/fbc/libs/tendermint/consensus"
	"github.com/fibonacci-chain/fbc/libs/tendermint/libs/log"
	"github.com/fibonacci-chain/fbc/libs/tendermint/state"
	tmtypes "github.com/fibonacci-chain/fbc/libs/tendermint/types"

	"github.com/spf13/viper"
)

var _ tmconfig.IDynamicConfig = &FecConfig{}
var _ iavlconfig.IDynamicConfig = &FecConfig{}

type FecConfig struct {
	// mempool.recheck
	mempoolRecheck bool
	// mempool.force_recheck_gap
	mempoolForceRecheckGap int64
	// mempool.size
	mempoolSize int
	// mempool.cache_size
	mempoolCacheSize int
	// mempool.flush
	mempoolFlush bool
	// mempool.max_tx_num_per_block
	maxTxNumPerBlock int64
	// mempool.enable_delete_min_gp_tx
	enableDeleteMinGPTx bool
	// mempool.max_gas_used_per_block
	maxGasUsedPerBlock int64
	// mempool.enable-pgu
	enablePGU bool
	// mempool.pgu-adjustment
	pguAdjustment float64
	//mempool.pgu-concurrency
	pguConcurrency int
	// mempool.node_key_whitelist
	nodeKeyWhitelist []string
	//mempool.check_tx_cost
	mempoolCheckTxCost bool
	// p2p.sentry_addrs
	sentryAddrs []string

	// gas-limit-buffer
	gasLimitBuffer uint64

	// enable-dynamic-gp
	enableDynamicGp bool
	// dynamic-gp-weight
	dynamicGpWeight int
	// dynamic-gp-check-blocks
	dynamicGpCheckBlocks int
	// dynamic-gp-coefficient
	dynamicGpCoefficient int
	// dynamic-gp-max-gas-used
	dynamicGpMaxGasUsed int64
	// dynamic-gp-max-tx-num
	dynamicGpMaxTxNum int64
	// dynamic-gp-mode
	dynamicGpMode int

	// consensus.timeout_propose
	csTimeoutPropose time.Duration
	// consensus.timeout_propose_delta
	csTimeoutProposeDelta time.Duration
	// consensus.timeout_prevote
	csTimeoutPrevote time.Duration
	// consensus.timeout_prevote_delta
	csTimeoutPrevoteDelta time.Duration
	// consensus.timeout_precommit
	csTimeoutPrecommit time.Duration
	// consensus.timeout_precommit_delta
	csTimeoutPrecommitDelta time.Duration
	// consensus.timeout_commit
	csTimeoutCommit time.Duration

	// iavl-cache-size
	iavlCacheSize int
	// commit-gap-height
	commitGapHeight int64

	// iavl-fast-storage-cache-size
	iavlFSCacheSize int64

	// enable-wtx
	enableWtx bool

	// enable-analyzer
	enableAnalyzer bool

	deliverTxsMode int

	// active view change
	activeVC bool

	blockPartSizeBytes int
	blockCompressType  int
	blockCompressFlag  int

	// enable broadcast hasBlockPartMsg
	enableHasBlockPartMsg bool
	gcInterval            int

	iavlAcNoBatch bool

	//
	commitGapOffset int64
}

const (
	FlagEnableDynamic = "config.enable-dynamic"

	FlagMempoolRecheck         = "mempool.recheck"
	FlagMempoolForceRecheckGap = "mempool.force_recheck_gap"
	FlagMempoolSize            = "mempool.size"
	FlagMempoolCacheSize       = "mempool.cache_size"
	FlagMempoolFlush           = "mempool.flush"
	FlagMaxTxNumPerBlock       = "mempool.max_tx_num_per_block"
	FlagMaxGasUsedPerBlock     = "mempool.max_gas_used_per_block"
	FlagEnablePGU              = "mempool.enable-pgu"
	FlagPGUAdjustment          = "mempool.pgu-adjustment"
	FlagPGUConcurrency         = "mempool.pgu-concurrency"

	FlagNodeKeyWhitelist           = "mempool.node_key_whitelist"
	FlagMempoolCheckTxCost         = "mempool.check_tx_cost"
	FlagMempoolEnableDeleteMinGPTx = "mempool.enable_delete_min_gp_tx"
	FlagGasLimitBuffer             = "gas-limit-buffer"
	FlagEnableDynamicGp            = "enable-dynamic-gp"
	FlagDynamicGpMode              = "dynamic-gp-mode"
	FlagDynamicGpWeight            = "dynamic-gp-weight"
	FlagDynamicGpCheckBlocks       = "dynamic-gp-check-blocks"
	FlagDynamicGpCoefficient       = "dynamic-gp-coefficient"
	FlagDynamicGpMaxGasUsed        = "dynamic-gp-max-gas-used"
	FlagDynamicGpMaxTxNum          = "dynamic-gp-max-tx-num"
	FlagEnableWrappedTx            = "enable-wtx"
	FlagSentryAddrs                = "p2p.sentry_addrs"
	FlagCsTimeoutPropose           = "consensus.timeout_propose"
	FlagCsTimeoutProposeDelta      = "consensus.timeout_propose_delta"
	FlagCsTimeoutPrevote           = "consensus.timeout_prevote"
	FlagCsTimeoutPrevoteDelta      = "consensus.timeout_prevote_delta"
	FlagCsTimeoutPrecommit         = "consensus.timeout_precommit"
	FlagCsTimeoutPrecommitDelta    = "consensus.timeout_precommit_delta"
	FlagCsTimeoutCommit            = "consensus.timeout_commit"
	FlagEnableHasBlockPartMsg      = "enable-blockpart-ack"
	FlagDebugGcInterval            = "debug.gc-interval"
	FlagCommitGapOffset            = "commit-gap-offset"
)

var (
	testnetNodeIdWhitelist = []string{
		// RPC nodes for users
		"3a339568305c5aff58a1f134437b608490e2ec6d",
	}

	mainnetNodeIdWhitelist = []string{}

	fecConfig  *FecConfig
	once       sync.Once
	confLogger log.Logger
)

func GetFecConfig() *FecConfig {
	once.Do(func() {
		fecConfig = NewFecConfig()
	})
	return fecConfig
}

func NewFecConfig() *FecConfig {
	c := defaultFecConfig()
	c.loadFromConfig()

	if viper.GetBool(FlagEnableDynamic) {
		if viper.IsSet(FlagApollo) {
			loaded := c.loadFromApollo()
			if !loaded {
				panic("failed to connect apollo or no config items in apollo")
			}
		} else {
			ok, err := c.loadFromLocal()
			if err != nil {
				confLogger.Error("failed to load config from local", "err", err)
			}
			if !ok {
				confLogger.Error("failed to load config from local")
			} else {
				confLogger.Info("load config from local success")
			}
		}
	}

	return c
}

func defaultFecConfig() *FecConfig {
	return &FecConfig{
		mempoolRecheck:         false,
		mempoolForceRecheckGap: 2000,
		commitGapHeight:        iavlconfig.DefaultCommitGapHeight,
		iavlFSCacheSize:        tmiavl.DefaultIavlFastStorageCacheSize,
	}
}

func RegisterDynamicConfig(logger log.Logger) {
	confLogger = logger
	// set the dynamic config
	fecConfig := GetFecConfig()
	tmconfig.SetDynamicConfig(fecConfig)
	iavlconfig.SetDynamicConfig(fecConfig)
	trace.SetDynamicConfig(fecConfig)
}

func (c *FecConfig) loadFromConfig() {
	c.SetMempoolRecheck(viper.GetBool(FlagMempoolRecheck))
	c.SetMempoolForceRecheckGap(viper.GetInt64(FlagMempoolForceRecheckGap))
	c.SetMempoolSize(viper.GetInt(FlagMempoolSize))
	c.SetMempoolCacheSize(viper.GetInt(FlagMempoolCacheSize))
	c.SetMempoolFlush(viper.GetBool(FlagMempoolFlush))
	c.SetMempoolCheckTxCost(viper.GetBool(FlagMempoolCheckTxCost))
	c.SetMaxTxNumPerBlock(viper.GetInt64(FlagMaxTxNumPerBlock))
	c.SetEnableDeleteMinGPTx(viper.GetBool(FlagMempoolEnableDeleteMinGPTx))
	c.SetMaxGasUsedPerBlock(viper.GetInt64(FlagMaxGasUsedPerBlock))
	c.SetEnablePGU(viper.GetBool(FlagEnablePGU))
	c.SetPGUAdjustment(viper.GetFloat64(FlagPGUAdjustment))
	c.SetPGUConcurrency(viper.GetInt(FlagPGUConcurrency))
	c.SetGasLimitBuffer(viper.GetUint64(FlagGasLimitBuffer))

	c.SetEnableDynamicGp(viper.GetBool(FlagEnableDynamicGp))
	c.SetDynamicGpWeight(viper.GetInt(FlagDynamicGpWeight))
	c.SetDynamicGpCheckBlocks(viper.GetInt(FlagDynamicGpCheckBlocks))
	c.SetDynamicGpCoefficient(viper.GetInt(FlagDynamicGpCoefficient))
	c.SetDynamicGpMaxGasUsed(viper.GetInt64(FlagDynamicGpMaxGasUsed))
	c.SetDynamicGpMaxTxNum(viper.GetInt64(FlagDynamicGpMaxTxNum))

	c.SetDynamicGpMode(viper.GetInt(FlagDynamicGpMode))
	c.SetCsTimeoutPropose(viper.GetDuration(FlagCsTimeoutPropose))
	c.SetCsTimeoutProposeDelta(viper.GetDuration(FlagCsTimeoutProposeDelta))
	c.SetCsTimeoutPrevote(viper.GetDuration(FlagCsTimeoutPrevote))
	c.SetCsTimeoutPrevoteDelta(viper.GetDuration(FlagCsTimeoutPrevoteDelta))
	c.SetCsTimeoutPrecommit(viper.GetDuration(FlagCsTimeoutPrecommit))
	c.SetCsTimeoutPrecommitDelta(viper.GetDuration(FlagCsTimeoutPrecommitDelta))
	c.SetCsTimeoutCommit(viper.GetDuration(FlagCsTimeoutCommit))
	c.SetIavlCacheSize(viper.GetInt(iavl.FlagIavlCacheSize))
	c.SetIavlFSCacheSize(viper.GetInt64(tmiavl.FlagIavlFastStorageCacheSize))
	c.SetCommitGapHeight(viper.GetInt64(server.FlagCommitGapHeight))
	c.SetSentryAddrs(viper.GetString(FlagSentryAddrs))
	c.SetNodeKeyWhitelist(viper.GetString(FlagNodeKeyWhitelist))
	c.SetEnableWtx(viper.GetBool(FlagEnableWrappedTx))
	c.SetEnableAnalyzer(viper.GetBool(trace.FlagEnableAnalyzer))
	c.SetDeliverTxsExecuteMode(viper.GetInt(state.FlagDeliverTxsExecMode))
	c.SetCommitGapOffset(viper.GetInt64(FlagCommitGapOffset))
	c.SetBlockPartSize(viper.GetInt(server.FlagBlockPartSizeBytes))
	c.SetEnableHasBlockPartMsg(viper.GetBool(FlagEnableHasBlockPartMsg))
	c.SetGcInterval(viper.GetInt(FlagDebugGcInterval))
	c.SetIavlAcNoBatch(viper.GetBool(tmiavl.FlagIavlCommitAsyncNoBatch))
}

func resolveNodeKeyWhitelist(plain string) []string {
	if len(plain) == 0 {
		return []string{}
	}
	return strings.Split(plain, ",")
}

func resolveSentryAddrs(plain string) []string {
	if len(plain) == 0 {
		return []string{}
	}
	return strings.Split(plain, ";")
}

func (c *FecConfig) loadFromApollo() bool {
	client := NewApolloClient(c)
	return client.LoadConfig()
}

func (c *FecConfig) loadFromLocal() (bool, error) {
	var err error
	rootDir := viper.GetString("home")
	configPath := path.Join(rootDir, "config", LocalDynamicConfigPath)
	configPath, err = filepath.Abs(configPath)
	if err != nil {
		return false, err
	}
	client, err := NewLocalClient(configPath, c, confLogger)
	if err != nil {
		return false, err
	}
	ok := client.LoadConfig()
	err = client.Enable()
	return ok, err
}

func (c *FecConfig) format() string {
	return fmt.Sprintf(`%s config:
	mempool.recheck: %v
	mempool.force_recheck_gap: %d
	mempool.size: %d
	mempool.cache_size: %d

	mempool.flush: %v
	mempool.max_tx_num_per_block: %d
	mempool.enable_delete_min_gp_tx: %v
	mempool.max_gas_used_per_block: %d
	mempool.check_tx_cost: %v

	gas-limit-buffer: %d
	dynamic-gp-weight: %d
	dynamic-gp-check-blocks: %d
	dynamic-gp-coefficient: %d
	dynamic-gp-max-gas-used: %d
	dynamic-gp-max-tx-num: %d
	dynamic-gp-mode: %d

	consensus.timeout_propose: %s
	consensus.timeout_propose_delta: %s
	consensus.timeout_prevote: %s
	consensus.timeout_prevote_delta: %s
	consensus.timeout_precommit: %s
	consensus.timeout_precommit_delta: %s
	consensus.timeout_commit: %s
	
	iavl-cache-size: %d
    iavl-fast-storage-cache-size: %d
    commit-gap-height: %d
	enable-analyzer: %v
    iavl-commit-async-no-batch: %v
	active-view-change: %v`, system.ChainName,
		c.GetMempoolRecheck(),
		c.GetMempoolForceRecheckGap(),
		c.GetMempoolSize(),
		c.GetMempoolCacheSize(),
		c.GetMempoolFlush(),
		c.GetMaxTxNumPerBlock(),
		c.GetEnableDeleteMinGPTx(),
		c.GetMaxGasUsedPerBlock(),
		c.GetMempoolCheckTxCost(),
		c.GetGasLimitBuffer(),
		c.GetDynamicGpWeight(),
		c.GetDynamicGpCheckBlocks(),
		c.GetDynamicGpCoefficient(),
		c.GetDynamicGpMaxGasUsed(),
		c.GetDynamicGpMaxTxNum(),
		c.GetDynamicGpMode(),
		c.GetCsTimeoutPropose(),
		c.GetCsTimeoutProposeDelta(),
		c.GetCsTimeoutPrevote(),
		c.GetCsTimeoutPrevoteDelta(),
		c.GetCsTimeoutPrecommit(),
		c.GetCsTimeoutPrecommitDelta(),
		c.GetCsTimeoutCommit(),
		c.GetIavlCacheSize(),
		c.GetIavlFSCacheSize(),
		c.GetCommitGapHeight(),
		c.GetEnableAnalyzer(),
		c.GetIavlAcNoBatch(),
		c.GetActiveVC(),
	)
}

func (c *FecConfig) update(key, value interface{}) {
	k, v := key.(string), value.(string)
	c.updateFromKVStr(k, v)
}

func (c *FecConfig) updateFromKVStr(k, v string) {
	switch k {
	case FlagMempoolRecheck:
		r, err := strconv.ParseBool(v)
		if err != nil {
			return
		}
		c.SetMempoolRecheck(r)
	case FlagMempoolForceRecheckGap:
		r, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return
		}
		c.SetMempoolForceRecheckGap(r)
	case FlagMempoolSize:
		r, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		c.SetMempoolSize(r)
	case FlagMempoolCacheSize:
		r, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		c.SetMempoolCacheSize(r)
	case FlagMempoolFlush:
		r, err := strconv.ParseBool(v)
		if err != nil {
			return
		}
		c.SetMempoolFlush(r)
	case FlagMaxTxNumPerBlock:
		r, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return
		}
		c.SetMaxTxNumPerBlock(r)
	case FlagMempoolEnableDeleteMinGPTx:
		r, err := strconv.ParseBool(v)
		if err != nil {
			return
		}
		c.SetEnableDeleteMinGPTx(r)
	case FlagNodeKeyWhitelist:
		c.SetNodeKeyWhitelist(v)
	case FlagMempoolCheckTxCost:
		r, err := strconv.ParseBool(v)
		if err != nil {
			return
		}
		c.SetMempoolCheckTxCost(r)
	case FlagSentryAddrs:
		c.SetSentryAddrs(v)
	case FlagMaxGasUsedPerBlock:
		r, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return
		}
		c.SetMaxGasUsedPerBlock(r)
	case FlagEnablePGU:
		r, err := strconv.ParseBool(v)
		if err != nil {
			return
		}
		c.SetEnablePGU(r)
	case FlagPGUAdjustment:
		r, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return
		}
		c.SetPGUAdjustment(r)
	case FlagPGUConcurrency:
		r, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		c.SetPGUConcurrency(r)
	case FlagGasLimitBuffer:
		r, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return
		}
		c.SetGasLimitBuffer(r)
	case FlagEnableDynamicGp:
		r, err := strconv.ParseBool(v)
		if err != nil {
			return
		}
		c.SetEnableDynamicGp(r)
	case FlagDynamicGpWeight:
		r, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		c.SetDynamicGpWeight(r)
	case FlagDynamicGpCheckBlocks:
		r, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		c.SetDynamicGpCheckBlocks(r)
	case FlagDynamicGpCoefficient:
		r, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		c.SetDynamicGpCoefficient(r)
	case FlagDynamicGpMaxGasUsed:
		r, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return
		}
		c.SetDynamicGpMaxGasUsed(r)
	case FlagDynamicGpMaxTxNum:
		r, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return
		}
		c.SetDynamicGpMaxTxNum(r)
	case FlagDynamicGpMode:
		r, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		c.SetDynamicGpMode(r)
	case FlagCsTimeoutPropose:
		r, err := time.ParseDuration(v)
		if err != nil {
			return
		}
		c.SetCsTimeoutPropose(r)
	case FlagCsTimeoutProposeDelta:
		r, err := time.ParseDuration(v)
		if err != nil {
			return
		}
		c.SetCsTimeoutProposeDelta(r)
	case FlagCsTimeoutPrevote:
		r, err := time.ParseDuration(v)
		if err != nil {
			return
		}
		c.SetCsTimeoutPrevote(r)
	case FlagCsTimeoutPrevoteDelta:
		r, err := time.ParseDuration(v)
		if err != nil {
			return
		}
		c.SetCsTimeoutPrevoteDelta(r)
	case FlagCsTimeoutPrecommit:
		r, err := time.ParseDuration(v)
		if err != nil {
			return
		}
		c.SetCsTimeoutPrecommit(r)
	case FlagCsTimeoutPrecommitDelta:
		r, err := time.ParseDuration(v)
		if err != nil {
			return
		}
		c.SetCsTimeoutPrecommitDelta(r)
	case FlagCsTimeoutCommit:
		r, err := time.ParseDuration(v)
		if err != nil {
			return
		}
		c.SetCsTimeoutCommit(r)
	case iavl.FlagIavlCacheSize:
		r, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		c.SetIavlCacheSize(r)
	case tmiavl.FlagIavlFastStorageCacheSize:
		r, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		c.SetIavlFSCacheSize(int64(r))
	case server.FlagCommitGapHeight:
		r, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return
		}
		c.SetCommitGapHeight(r)
	case trace.FlagEnableAnalyzer:
		r, err := strconv.ParseBool(v)
		if err != nil {
			return
		}
		c.SetEnableAnalyzer(r)
	case state.FlagDeliverTxsExecMode:
		r, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		c.SetDeliverTxsExecuteMode(r)
	case server.FlagActiveViewChange:
		r, err := strconv.ParseBool(v)
		if err != nil {
			return
		}
		c.SetActiveVC(r)
	case server.FlagBlockPartSizeBytes:
		r, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		c.SetBlockPartSize(r)
	case tmtypes.FlagBlockCompressType:
		r, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		c.SetBlockCompressType(r)
	case tmtypes.FlagBlockCompressFlag:
		r, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		c.SetBlockCompressFlag(r)
	case FlagEnableHasBlockPartMsg:
		r, err := strconv.ParseBool(v)
		if err != nil {
			return
		}
		c.SetEnableHasBlockPartMsg(r)
	case FlagDebugGcInterval:
		r, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		c.SetGcInterval(r)
	case tmiavl.FlagIavlCommitAsyncNoBatch:
		r, err := strconv.ParseBool(v)
		if err != nil {
			return
		}
		c.SetIavlAcNoBatch(r)
	case FlagCommitGapOffset:
		r, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return
		}
		c.SetCommitGapOffset(r)
	}

}

func (c *FecConfig) GetEnableAnalyzer() bool {
	return c.enableAnalyzer
}
func (c *FecConfig) SetEnableAnalyzer(value bool) {
	c.enableAnalyzer = value
}

func (c *FecConfig) GetMempoolRecheck() bool {
	return c.mempoolRecheck
}
func (c *FecConfig) SetMempoolRecheck(value bool) {
	c.mempoolRecheck = value
}

func (c *FecConfig) GetMempoolForceRecheckGap() int64 {
	return c.mempoolForceRecheckGap
}
func (c *FecConfig) SetMempoolForceRecheckGap(value int64) {
	if value <= 0 {
		return
	}
	c.mempoolForceRecheckGap = value
}

func (c *FecConfig) GetMempoolSize() int {
	return c.mempoolSize
}
func (c *FecConfig) SetMempoolSize(value int) {
	if value < 0 {
		return
	}
	c.mempoolSize = value
}

func (c *FecConfig) GetMempoolCacheSize() int {
	return c.mempoolCacheSize
}
func (c *FecConfig) SetMempoolCacheSize(value int) {
	if value < 0 {
		return
	}
	c.mempoolCacheSize = value
}

func (c *FecConfig) GetMempoolFlush() bool {
	return c.mempoolFlush
}
func (c *FecConfig) SetMempoolFlush(value bool) {
	c.mempoolFlush = value
}

func (c *FecConfig) GetEnableWtx() bool {
	return c.enableWtx
}

func (c *FecConfig) SetDeliverTxsExecuteMode(mode int) {
	c.deliverTxsMode = mode
}

func (c *FecConfig) GetDeliverTxsExecuteMode() int {
	return c.deliverTxsMode
}

func (c *FecConfig) SetEnableWtx(value bool) {
	c.enableWtx = value
}

func (c *FecConfig) GetNodeKeyWhitelist() []string {
	return c.nodeKeyWhitelist
}

func (c *FecConfig) GetMempoolCheckTxCost() bool {
	return c.mempoolCheckTxCost
}
func (c *FecConfig) SetMempoolCheckTxCost(value bool) {
	c.mempoolCheckTxCost = value
}

func (c *FecConfig) SetNodeKeyWhitelist(value string) {
	idList := resolveNodeKeyWhitelist(value)

	for _, id := range idList {
		if id == "testnet-node-ids" {
			c.nodeKeyWhitelist = append(c.nodeKeyWhitelist, testnetNodeIdWhitelist...)
		} else if id == "mainnet-node-ids" {
			c.nodeKeyWhitelist = append(c.nodeKeyWhitelist, mainnetNodeIdWhitelist...)
		} else {
			c.nodeKeyWhitelist = append(c.nodeKeyWhitelist, id)
		}
	}
}

func (c *FecConfig) GetSentryAddrs() []string {
	return c.sentryAddrs
}

func (c *FecConfig) SetSentryAddrs(value string) {
	addrs := resolveSentryAddrs(value)
	for _, addr := range addrs {
		c.sentryAddrs = append(c.sentryAddrs, strings.TrimSpace(addr))
	}
}

func (c *FecConfig) GetMaxTxNumPerBlock() int64 {
	return c.maxTxNumPerBlock
}
func (c *FecConfig) SetMaxTxNumPerBlock(value int64) {
	if value < 0 {
		return
	}
	c.maxTxNumPerBlock = value
}

func (c *FecConfig) GetEnableDeleteMinGPTx() bool {
	return c.enableDeleteMinGPTx
}

func (c *FecConfig) SetEnableDeleteMinGPTx(enable bool) {
	c.enableDeleteMinGPTx = enable
}

func (c *FecConfig) GetMaxGasUsedPerBlock() int64 {
	return c.maxGasUsedPerBlock
}

func (c *FecConfig) SetMaxGasUsedPerBlock(value int64) {
	if value < -1 {
		return
	}
	c.maxGasUsedPerBlock = value
}

func (c *FecConfig) GetEnablePGU() bool {
	return c.enablePGU
}

func (c *FecConfig) SetEnablePGU(value bool) {
	c.enablePGU = value
}

func (c *FecConfig) GetPGUAdjustment() float64 {
	return c.pguAdjustment
}

func (c *FecConfig) SetPGUAdjustment(value float64) {
	c.pguAdjustment = value
}

func (c *FecConfig) GetGasLimitBuffer() uint64 {
	return c.gasLimitBuffer
}
func (c *FecConfig) SetGasLimitBuffer(value uint64) {
	c.gasLimitBuffer = value
}

func (c *FecConfig) GetEnableDynamicGp() bool {
	return c.enableDynamicGp
}

func (c *FecConfig) SetEnableDynamicGp(value bool) {
	c.enableDynamicGp = value
}

func (c *FecConfig) GetDynamicGpWeight() int {
	return c.dynamicGpWeight
}

func (c *FecConfig) SetDynamicGpWeight(value int) {
	if value <= 0 {
		value = 1
	} else if value > 100 {
		value = 100
	}
	c.dynamicGpWeight = value
}

func (c *FecConfig) GetDynamicGpCoefficient() int {
	return c.dynamicGpCoefficient
}
func (c *FecConfig) SetDynamicGpCoefficient(value int) {
	if value <= 0 {
		value = 1
	} else if value > 100 {
		value = 100
	}
	c.dynamicGpCoefficient = value
}

func (c *FecConfig) GetDynamicGpMaxGasUsed() int64 {
	return c.dynamicGpMaxGasUsed
}

func (c *FecConfig) SetDynamicGpMaxGasUsed(value int64) {
	if value < -1 {
		return
	}
	c.dynamicGpMaxGasUsed = value
}

func (c *FecConfig) GetDynamicGpMaxTxNum() int64 {
	return c.dynamicGpMaxTxNum
}

func (c *FecConfig) SetDynamicGpMaxTxNum(value int64) {
	if value < 0 {
		return
	}
	c.dynamicGpMaxTxNum = value
}

func (c *FecConfig) GetDynamicGpMode() int {
	return c.dynamicGpMode
}

func (c *FecConfig) SetDynamicGpMode(value int) {
	if value < 0 || value > 2 {
		return
	}
	c.dynamicGpMode = value
}

func (c *FecConfig) GetDynamicGpCheckBlocks() int {
	return c.dynamicGpCheckBlocks
}

func (c *FecConfig) SetDynamicGpCheckBlocks(value int) {
	if value <= 0 {
		value = 1
	} else if value > 100 {
		value = 100
	}
	c.dynamicGpCheckBlocks = value
}

func (c *FecConfig) GetCsTimeoutPropose() time.Duration {
	return c.csTimeoutPropose
}
func (c *FecConfig) SetCsTimeoutPropose(value time.Duration) {
	if value < 0 {
		return
	}
	c.csTimeoutPropose = value
}

func (c *FecConfig) GetCsTimeoutProposeDelta() time.Duration {
	return c.csTimeoutProposeDelta
}
func (c *FecConfig) SetCsTimeoutProposeDelta(value time.Duration) {
	if value < 0 {
		return
	}
	c.csTimeoutProposeDelta = value
}

func (c *FecConfig) GetCsTimeoutPrevote() time.Duration {
	return c.csTimeoutPrevote
}
func (c *FecConfig) SetCsTimeoutPrevote(value time.Duration) {
	if value < 0 {
		return
	}
	c.csTimeoutPrevote = value
}

func (c *FecConfig) GetCsTimeoutPrevoteDelta() time.Duration {
	return c.csTimeoutPrevoteDelta
}
func (c *FecConfig) SetCsTimeoutPrevoteDelta(value time.Duration) {
	if value < 0 {
		return
	}
	c.csTimeoutPrevoteDelta = value
}

func (c *FecConfig) GetCsTimeoutPrecommit() time.Duration {
	return c.csTimeoutPrecommit
}
func (c *FecConfig) SetCsTimeoutPrecommit(value time.Duration) {
	if value < 0 {
		return
	}
	c.csTimeoutPrecommit = value
}

func (c *FecConfig) GetCsTimeoutPrecommitDelta() time.Duration {
	return c.csTimeoutPrecommitDelta
}
func (c *FecConfig) SetCsTimeoutPrecommitDelta(value time.Duration) {
	if value < 0 {
		return
	}
	c.csTimeoutPrecommitDelta = value
}

func (c *FecConfig) GetCsTimeoutCommit() time.Duration {
	return c.csTimeoutCommit
}
func (c *FecConfig) SetCsTimeoutCommit(value time.Duration) {
	if value < 0 {
		return
	}
	c.csTimeoutCommit = value
}

func (c *FecConfig) GetIavlCacheSize() int {
	return c.iavlCacheSize
}
func (c *FecConfig) SetIavlCacheSize(value int) {
	c.iavlCacheSize = value
}

func (c *FecConfig) GetIavlFSCacheSize() int64 {
	return c.iavlFSCacheSize
}

func (c *FecConfig) SetIavlFSCacheSize(value int64) {
	c.iavlFSCacheSize = value
}

func (c *FecConfig) GetCommitGapHeight() int64 {
	return atomic.LoadInt64(&c.commitGapHeight)
}
func (c *FecConfig) SetCommitGapHeight(value int64) {
	if IsPruningOptionNothing() { // pruning nothing the gap should 1
		value = 1
	}
	if value <= 0 {
		return
	}
	atomic.StoreInt64(&c.commitGapHeight, value)
}

func IsPruningOptionNothing() bool {
	strategy := strings.ToLower(viper.GetString(server.FlagPruning))
	if strategy == types.PruningOptionNothing {
		return true
	}
	return false
}

func (c *FecConfig) GetActiveVC() bool {
	return c.activeVC
}
func (c *FecConfig) SetActiveVC(value bool) {
	c.activeVC = value
	consensus.SetActiveVC(value)
}

func (c *FecConfig) GetBlockPartSize() int {
	return c.blockPartSizeBytes
}
func (c *FecConfig) SetBlockPartSize(value int) {
	c.blockPartSizeBytes = value
	tmtypes.UpdateBlockPartSizeBytes(value)
}

func (c *FecConfig) GetBlockCompressType() int {
	return c.blockCompressType
}
func (c *FecConfig) SetBlockCompressType(value int) {
	c.blockCompressType = value
	tmtypes.BlockCompressType = value
}

func (c *FecConfig) GetBlockCompressFlag() int {
	return c.blockCompressFlag
}
func (c *FecConfig) SetBlockCompressFlag(value int) {
	c.blockCompressFlag = value
	tmtypes.BlockCompressFlag = value
}

func (c *FecConfig) GetGcInterval() int {
	return c.gcInterval
}

func (c *FecConfig) SetGcInterval(value int) {
	// close gc for debug
	if value > 0 {
		debug.SetGCPercent(-1)
	} else {
		debug.SetGCPercent(100)
	}
	c.gcInterval = value

}

func (c *FecConfig) GetCommitGapOffset() int64 {
	return c.commitGapOffset
}

func (c *FecConfig) SetCommitGapOffset(value int64) {
	if value < 0 {
		value = 0
	}
	c.commitGapOffset = value
}

func (c *FecConfig) GetEnableHasBlockPartMsg() bool {
	return c.enableHasBlockPartMsg
}

func (c *FecConfig) SetEnableHasBlockPartMsg(value bool) {
	c.enableHasBlockPartMsg = value
}

func (c *FecConfig) GetIavlAcNoBatch() bool {
	return c.iavlAcNoBatch
}

func (c *FecConfig) SetIavlAcNoBatch(value bool) {
	c.iavlAcNoBatch = value
}

func (c *FecConfig) GetPGUConcurrency() int {
	return c.pguConcurrency
}
func (c *FecConfig) SetPGUConcurrency(value int) {
	c.pguConcurrency = value
}
