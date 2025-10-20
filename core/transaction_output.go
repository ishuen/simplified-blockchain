package core

import "bytes"

type TXOutput struct {
	Value      int // amount of coins
	PubKeyHash []byte
}

func (out *TXOutput) Lock(address []byte) {
	publicKeyHash := Base58Decode(address)
	// remove version and checksum
	publicKeyHash = publicKeyHash[1 : len(publicKeyHash)-addressChecksumLen]
	out.PubKeyHash = publicKeyHash
}

// checks if the output can be used by the owner of the publicKeyHash
func (out *TXOutput) isLockedWithKey(publicKeyHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, publicKeyHash)
}

func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}
