package types

import (
	"fmt"

	amino "github.com/tendermint/go-amino"

	abci "github.com/fibonacci-chain/fbc/libs/tendermint/abci/types"
	tmpubsub "github.com/fibonacci-chain/fbc/libs/tendermint/libs/pubsub"
	tmquery "github.com/fibonacci-chain/fbc/libs/tendermint/libs/pubsub/query"
)

// Reserved event types (alphabetically sorted).
const (
	// Block level events for mass consumption by users.
	// These events are triggered from the state package,
	// after a block has been committed.
	// These are also used by the tx indexer for async indexing.
	// All of this data can be fetched through the rpc.
	EventNewBlock            = "NewBlock"
	EventNewBlockHeader      = "NewBlockHeader"
	EventTx                  = "Tx"
	EventPendingTx           = "PendingTx"
	EventRmPendingTx         = "RmPendingTx"
	EventValidatorSetUpdates = "ValidatorSetUpdates"
	EventBlockTime           = "BlockTime"
	EventTxs                 = "Txs"

	// Internal consensus events.
	// These are used for testing the consensus state machine.
	// They can also be used to build real-time consensus visualizers.
	EventCompleteProposal = "CompleteProposal"
	EventLock             = "Lock"
	EventNewRound         = "NewRound"
	EventNewRoundStep     = "NewRoundStep"
	EventPolka            = "Polka"
	EventRelock           = "Relock"
	EventTimeoutPropose   = "TimeoutPropose"
	EventTimeoutWait      = "TimeoutWait"
	EventUnlock           = "Unlock"
	EventValidBlock       = "ValidBlock"
	EventVote             = "Vote"
	EventSignVote         = "SignVote"
	EventBlockPart        = "BlockPart"
	EventProposeRequest   = "ProposeRequest"
)

type RmPendingTxReason int

const (
	Recheck RmPendingTxReason = iota
	MinGasPrice
	Confirmed
)

var EnableEventBlockTime = false

///////////////////////////////////////////////////////////////////////////////
// ENCODING / DECODING
///////////////////////////////////////////////////////////////////////////////

// TMEventData implements events.EventData.
type TMEventData interface {
	// empty interface
}

func RegisterEventDatas(cdc *amino.Codec) {
	cdc.RegisterInterface((*TMEventData)(nil), nil)
	cdc.RegisterConcrete(CM40EventDataNewBlock{}, "tendermint/event/NewBlock", nil)
	cdc.RegisterConcrete(EventDataNewBlockHeader{}, "tendermint/event/NewBlockHeader", nil)
	cdc.RegisterConcrete(EventDataTx{}, "tendermint/event/Tx", nil)
	cdc.RegisterConcrete(EventDataRoundState{}, "tendermint/event/RoundState", nil)
	cdc.RegisterConcrete(EventDataNewRound{}, "tendermint/event/NewRound", nil)
	cdc.RegisterConcrete(EventDataCompleteProposal{}, "tendermint/event/CompleteProposal", nil)
	cdc.RegisterConcrete(EventDataVote{}, "tendermint/event/Vote", nil)
	cdc.RegisterConcrete(EventDataValidatorSetUpdates{}, "tendermint/event/ValidatorSetUpdates", nil)
	cdc.RegisterConcrete(EventDataString(""), "tendermint/event/ProposalString", nil)
}

// Most event messages are basic types (a block, a transaction)
// but some (an input to a call tx or a receive) are more exotic

type EventDataNewBlock struct {
	Block *Block `json:"block"`

	ResultBeginBlock abci.ResponseBeginBlock `json:"result_begin_block"`
	ResultEndBlock   abci.ResponseEndBlock   `json:"result_end_block"`
}

func (e EventDataNewBlock) Upgrade() interface{} {
	ret := CM40EventDataNewBlock{}
	return ret.From(e)
}

type EventDataNewBlockHeader struct {
	Header Header `json:"header"`

	NumTxs           int64                   `json:"num_txs"` // Number of txs in a block
	ResultBeginBlock abci.ResponseBeginBlock `json:"result_begin_block"`
	ResultEndBlock   abci.ResponseEndBlock   `json:"result_end_block"`
}

// All txs fire EventDataTx
type EventDataTx struct {
	TxResult
	Nonce uint64
}

type EventDataTxs struct {
	Height int64
	//Txs     Txs
	Results []*abci.ResponseDeliverTx
}

type EventDataRmPendingTx struct {
	Hash   []byte
	From   string
	Nonce  uint64
	Reason RmPendingTxReason
}

// latest blockTime
type EventDataBlockTime struct {
	Height    int64
	TimeNow   int64
	TxNum     int
	Available bool
}

// NOTE: This goes into the replay WAL
type EventDataRoundState struct {
	Height int64  `json:"height"`
	Round  int    `json:"round"`
	Step   string `json:"step"`
}

type ValidatorInfo struct {
	Address Address `json:"address"`
	Index   int     `json:"index"`
}

type EventDataNewRound struct {
	Height int64  `json:"height"`
	Round  int    `json:"round"`
	Step   string `json:"step"`

	Proposer ValidatorInfo `json:"proposer"`
}

type EventDataCompleteProposal struct {
	Height int64  `json:"height"`
	Round  int    `json:"round"`
	Step   string `json:"step"`

	BlockID BlockID `json:"block_id"`
}

type EventDataVote struct {
	Vote *Vote
}

type EventDataString string

type EventDataValidatorSetUpdates struct {
	ValidatorUpdates []*Validator `json:"validator_updates"`
}

///////////////////////////////////////////////////////////////////////////////
// PUBSUB
///////////////////////////////////////////////////////////////////////////////

const (
	// EventTypeKey is a reserved composite key for event name.
	EventTypeKey = "tm.event"
	// TxHashKey is a reserved key, used to specify transaction's hash.
	// see EventBus#PublishEventTx
	TxHashKey = "tx.hash"
	// TxHeightKey is a reserved key, used to specify transaction block's height.
	// see EventBus#PublishEventTx
	TxHeightKey = "tx.height"

	// BlockHeightKey is a reserved key used for indexing BeginBlock and Endblock
	// events.
	BlockHeightKey = "block.height"
)

var (
	EventQueryCompleteProposal    = QueryForEvent(EventCompleteProposal)
	EventQueryLock                = QueryForEvent(EventLock)
	EventQueryNewBlock            = QueryForEvent(EventNewBlock)
	EventQueryNewBlockHeader      = QueryForEvent(EventNewBlockHeader)
	EventQueryNewRound            = QueryForEvent(EventNewRound)
	EventQueryNewRoundStep        = QueryForEvent(EventNewRoundStep)
	EventQueryPolka               = QueryForEvent(EventPolka)
	EventQueryRelock              = QueryForEvent(EventRelock)
	EventQueryTimeoutPropose      = QueryForEvent(EventTimeoutPropose)
	EventQueryTimeoutWait         = QueryForEvent(EventTimeoutWait)
	EventQueryTx                  = QueryForEvent(EventTx)
	EventQueryUnlock              = QueryForEvent(EventUnlock)
	EventQueryValidatorSetUpdates = QueryForEvent(EventValidatorSetUpdates)
	EventQueryValidBlock          = QueryForEvent(EventValidBlock)
	EventQueryVote                = QueryForEvent(EventVote)
)

func EventQueryTxFor(tx Tx, height int64) tmpubsub.Query {
	return tmquery.MustParse(fmt.Sprintf("%s='%s' AND %s='%X'", EventTypeKey, EventTx, TxHashKey, tx.Hash(height)))
}

func QueryForEvent(eventType string) tmpubsub.Query {
	return tmquery.MustParse(fmt.Sprintf("%s='%s'", EventTypeKey, eventType))
}

// BlockEventPublisher publishes all block related events
type BlockEventPublisher interface {
	PublishEventNewBlock(block EventDataNewBlock) error
	PublishEventNewBlockHeader(header EventDataNewBlockHeader) error
	PublishEventTx(EventDataTx) error
	PublishEventTxs(EventDataTxs) error
	PublishEventPendingTx(EventDataTx) error
	PublishEventValidatorSetUpdates(EventDataValidatorSetUpdates) error
	PublishEventLatestBlockTime(time EventDataBlockTime) error
	PublishEventRmPendingTx(EventDataRmPendingTx) error
}

type TxEventPublisher interface {
	PublishEventTx(EventDataTx) error
	PublishEventPendingTx(EventDataTx) error
	PublishEventRmPendingTx(EventDataRmPendingTx) error
}
