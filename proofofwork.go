package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var maxNonece = math.MaxInt64

const targetBits = 24

type ProofOfWork struct {
	block  *Block
	target *big.Int // requirement of the output hash
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	// shift left by (256 - targetBits)
	// => 0000010000000000000000000000000000000000000000000000000000000000
	return &ProofOfWork{block: b, target: target}
}

func (p *ProofOfWork) concatData(nonce int) []byte {
	data := bytes.Join([][]byte{
		p.block.PrevBlockHash,
		p.block.Data,
		[]byte(fmt.Sprintf("%x", p.block.Timestamp)),
		[]byte(fmt.Sprintf("%x", targetBits)),
		[]byte(fmt.Sprintf("%x", nonce)),
	}, []byte{})
	return data
}

// HashCash algorithm
func (p *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", p.block.Data)
	for nonce < maxNonece {
		data := p.concatData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(p.target) == -1 { // hashInt < target
			break
		}
		nonce++
	}
	fmt.Println()
	return nonce, hash[:]
}

func (p *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := p.concatData(p.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(p.target) == -1 // hash is valid if hashInt < target
}
