package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
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

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

type TXInput struct {
	Txid      []byte // ID of the transaction containing the output we want to reference
	Vout      int    // index of the output in the previous transaction
	Signature string // in a real blockchain, this would be a digital signature
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.Signature == unlockingData
}

type TXOutput struct {
	Value  int    // amount of coins
	PubKey string // in a real blockchain, this would be a hash of the recipient's public key
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.PubKey == unlockingData
}

/*
Coinbase transaction is a special transaction that has no inputs.
It is used to reward miners for mining a new block.
*/
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, data} // coinbase has no previous transaction
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()
	return &tx
}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			panic(err)
		}
		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from}) // a change output
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}
