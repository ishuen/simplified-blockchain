package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
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
	Data          []byte
	Nonce         int
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Data:          []byte(data),
		Nonce:         0,
		Hash:          []byte{},
	}
	block.SetHash()
	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
