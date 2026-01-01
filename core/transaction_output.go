package core

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TXOutput struct {
	Value      int // amount of coins
	PubKeyHash []byte
}

type TXOutputs struct {
	Outputs []TXOutput
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

func (outs TXOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DeserializeOutputs(data []byte) TXOutputs {
	var outputs TXOutputs
	
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)	
	if err != nil {
		log.Panic(err)
	}

	return outputs
}