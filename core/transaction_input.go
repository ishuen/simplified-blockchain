package core

import "bytes"

type TXInput struct {
	Txid      []byte // ID of the transaction containing the output we want to reference
	Vout      int    // index of the output in the previous transaction
	Signature []byte
	PublicKey []byte
}

func (in *TXInput) UsesKey(publicKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PublicKey)
	return bytes.Equal(lockingHash, publicKeyHash)
}
