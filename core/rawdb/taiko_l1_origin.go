package rawdb

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	// Database key prefix for L2 block's L1Origin.
	l1OriginPrefix  = []byte("TKO:L1O")
	headL1OriginKey = []byte("TKO:LastL1O")
)

// l1OriginKey calculates the L1Origin key.
// l1OriginPrefix + l2HeaderHash -> l1OriginKey
func l1OriginKey(blockID *big.Int) []byte {
	data, _ := (*math.HexOrDecimal256)(blockID).MarshalText()
	return append(l1OriginPrefix, data...)
}

//go:generate go run github.com/fjl/gencodec -type L1Origin -field-override l1OriginMarshaling -out gen_taiko_l1_origin.go

// L1Origin represents a L1Origin of a L2 block.
type L1Origin struct {
	BlockID       *big.Int    `json:"blockID" gencodec:"required"`
	L2BlockHash   common.Hash `json:"l2BlockHash"`
	L1BlockHeight *big.Int    `json:"l1BlockHeight" gencodec:"required"`
	L1BlockHash   common.Hash `json:"l1BlockHash" gencodec:"required"`
}

type l1OriginMarshaling struct {
	BlockID       *math.HexOrDecimal256
	L1BlockHeight *math.HexOrDecimal256
}

// WriteL1Origin stores a L1Origin into the database.
func WriteL1Origin(db ethdb.KeyValueWriter, blockID *big.Int, l1Origin *L1Origin) {
	data, err := rlp.EncodeToBytes(l1Origin)
	if err != nil {
		log.Crit("Failed to encode L1Origin", "err", err)
	}

	if err := db.Put(l1OriginKey(blockID), data); err != nil {
		log.Crit("Failed to store L1Origin", "err", err)
	}
}

// ReadL1Origin retrieves the given L2 block's L1Origin from database.
func ReadL1Origin(db ethdb.KeyValueReader, blockID *big.Int) (*L1Origin, error) {
	data, _ := db.Get(l1OriginKey(blockID))
	if len(data) == 0 {
		return nil, nil
	}

	l1Origin := new(L1Origin)
	if err := rlp.Decode(bytes.NewReader(data), l1Origin); err != nil {
		return nil, fmt.Errorf("invalid L1Origin RLP bytes: %w", err)
	}

	return l1Origin, nil
}

// WriteHeadL1Origin stores the given L1Origin as the last L1Origin.
func WriteHeadL1Origin(db ethdb.KeyValueWriter, blockID *big.Int) {
	data, _ := (*math.HexOrDecimal256)(blockID).MarshalText()
	if err := db.Put(headL1OriginKey, data); err != nil {
		log.Crit("Failed to store head L1Origin", "error", err)
	}
}

// ReadHeadL1Origin retrieves the last L1Origin from database.
func ReadHeadL1Origin(db ethdb.KeyValueReader) (*big.Int, error) {
	data, _ := db.Get(headL1OriginKey)
	if len(data) == 0 {
		return nil, nil
	}

	blockID := new(math.HexOrDecimal256)
	if err := blockID.UnmarshalText(data); err != nil {
		log.Error("Unmarshal L1Origin unmarshal error", "error", err)
		return nil, fmt.Errorf("invalid L1Origin unmarshal: %w", err)
	}

	return (*big.Int)(blockID), nil
}

// DeleteL1OriginByBlockNumber deletes L1Origin for the given block number.
// This is used during blockchain reorg/setHead operations.
func DeleteL1OriginByBlockNumber(db ethdb.KeyValueWriter, blockNumber uint64) {
	// Convert block number to BlockID (assuming BlockID == block number for simplicity)
	blockID := new(big.Int).SetUint64(blockNumber)

	// Delete the L1Origin entry
	if err := db.Delete(l1OriginKey(blockID)); err != nil {
		log.Warn("Failed to delete L1Origin", "blockNumber", blockNumber, "err", err)
	}
}

// UpdateHeadL1OriginAfterSetHead updates the head L1Origin pointer after a setHead operation.
// It finds the latest valid L1Origin at or before the given block number.
func UpdateHeadL1OriginAfterSetHead(db ethdb.Database, newHeadBlockNumber uint64) {
	// Search backwards from the new head to find the latest valid L1Origin
	for blockNum := newHeadBlockNumber; blockNum > 0; blockNum-- {
		blockID := new(big.Int).SetUint64(blockNum)
		if l1Origin, err := ReadL1Origin(db, blockID); err == nil && l1Origin != nil {
			// Found a valid L1Origin, update the head pointer
			WriteHeadL1Origin(db, blockID)
			log.Info("Updated head L1Origin after setHead", "newHead", newHeadBlockNumber, "l1OriginBlock", blockNum)
			return
		}
	}

	// If no L1Origin found, clear the head pointer
	if err := db.Delete(headL1OriginKey); err != nil {
		log.Warn("Failed to clear head L1Origin", "err", err)
	} else {
		log.Info("Cleared head L1Origin after setHead (no valid L1Origin found)", "newHead", newHeadBlockNumber)
	}
}
