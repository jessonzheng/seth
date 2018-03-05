package types

import (
	"math/big"
	"seth/common"
	"seth/crypto"
	"sync/atomic"
	"time"
)

// BlockNonce Block nonce
type BlockNonce [8]byte

// Header Block Header
type Header struct {
	ParentHash common.Hash    `json:"parentHash"       gencodec:"required"`
	Coinbase   common.Address `json:"miner"            gencodec:"required"`
	Root       common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash     common.Hash    `json:"transactionsRoot" gencodec:"required"`
	Difficulty *big.Int       `json:"difficulty"       gencodec:"required"`
	Number     *big.Int       `json:"number"           gencodec:"required"`
	Time       *big.Int       `json:"timestamp"        gencodec:"required"`
	Extra      []byte         `json:"extraData"        gencodec:"required"`
	MixDigest  common.Hash    `json:"mixHash"          gencodec:"required"`
	Nonce      BlockNonce     `json:"nonce"            gencodec:"required"`
}

// Hash return the block hash of the heade
func (h *Header) Hash() common.Hash {
	return crypto.RlpHash(h)
}

// Body block body struct
type Body struct {
	Transactions []*Transaction
}

// Block block struct define
type Block struct {
	header       *Header
	transactions Transactions

	// caches
	hash atomic.Value
	size atomic.Value

	// Td is used by package core to store the total difficulty
	// of the chain up to and including the block.
	td *big.Int

	// These fields are used by package eth to track
	// inter-peer block relay.
	ReceivedAt   time.Time
	ReceivedFrom interface{}
}

// Hash returns the keccak256 hash of block's header.
func (b *Block) Hash() common.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := b.header.Hash()
	b.hash.Store(v)
	return v
}

// FindTransaction find transaction with hash value
func (b *Block) FindTransaction(hash common.Hash) *Transaction {
	for _, transaction := range b.transactions {
		if transaction.Hash() == hash {
			return transaction
		}
	}
	return nil
}