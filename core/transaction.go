package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
)

const subsidy = 10 // amound of reward for mining a block

type Transaction struct {
	ID []byte
	/*
	   Inputs of a new transaction reference outputs of a previous transaction.
	   In one transaction, inputs can reference outputs from multiple transactions.
	*/
	Vin  []TXInput
	Vout []TXOutput
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		panic(err)
	}
	return encoded.Bytes()
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	(*tx).ID = []byte{}
	hash = sha256.Sum256(tx.Serialize())
	return hash[:]
}

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

func (out *TXOutput) isLockedWithKey(publicKeyHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, publicKeyHash)
}

func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}

/*
Coinbase transaction is a special transaction that has no inputs.
It is used to reward miners for mining a new block.
*/
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, nil, []byte(data)} // coinbase has no previous transaction
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	tx.ID = tx.Hash()
	return &tx
}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	wallets, err := NewWallets()
	if err != nil {
		panic(err)
	}

	wallet := wallets.GetWallet(from)
	publicKeyHash := HashPubKey(wallet.PublicKey)
	acc, validOutputs := bc.FindSpendableOutputs(publicKeyHash, amount)
	if acc < amount {
		panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			panic(err)
		}
		for _, out := range outs {
			input := TXInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from)) // a change output
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	return &tx
}
