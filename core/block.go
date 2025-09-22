package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

/*
For simplicity, block headers (previous block hash, timestamp, hash)
and data are all stored in the Block struct.
*/
type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Hash          []byte
	Transactions  []*Transaction
	Nonce         int
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result) // buffer pointer implements io.Writer
	err := encoder.Encode(b)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		panic(err)
	}
	return &block
}

func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Transactions:  transactions,
		Nonce:         0,
		Hash:          []byte{},
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}
