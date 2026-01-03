package core

import "bytes"

type TXInput struct {
	Txid      []byte // ID of the transaction containing the output we want to reference
	Vout      int    // index of the output in the previous transaction
	Signature []byte
	PubKey []byte
}

func (in *TXInput) UsesKey(publicKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	return bytes.Equal(lockingHash, publicKeyHash)
}
