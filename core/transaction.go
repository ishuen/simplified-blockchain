package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
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

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}
	for _, in := range tx.Vin {
		if prevTXs[hex.EncodeToString(in.Txid)].ID == nil {
			panic("ERROR: Previous transaction is not correct")
		}
	}
	txCopy := tx.TrimmedCopy()

	for inIdx, in := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(in.Txid)]
		txCopy.Vin[inIdx].Signature = nil
		txCopy.Vin[inIdx].PubKey = prevTx.Vout[in.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inIdx].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.ID)
		if err != nil {
			panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vin[inIdx].Signature = signature
	}
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, in := range tx.Vin {
		inputs = append(inputs, TXInput{in.Txid, in.Vout, nil, nil})
	}
	for _, out := range tx.Vout {
		outputs = append(outputs, TXOutput{out.Value, out.PubKeyHash})
	}
	txCopy := Transaction{tx.ID, inputs, outputs}
	return txCopy
}

// String returns a human-readable representation of a transaction
func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))
	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("       Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Vout))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("       Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       PubKeyHash: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	for _, in := range tx.Vin {
		if prevTXs[hex.EncodeToString(in.Txid)].ID == nil {
			panic("ERROR: Previous transaction is not correct")
		}
	}
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inIdx, in := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(in.Txid)]
		txCopy.Vin[inIdx].Signature = nil
		txCopy.Vin[inIdx].PubKey = prevTx.Vout[in.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inIdx].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen / 2)])
		s.SetBytes(in.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.PubKey)
		x.SetBytes(in.PubKey[:(keyLen / 2)])
		y.SetBytes(in.PubKey[(keyLen / 2):])

		publicKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		log.Println("Public Key", publicKey)
		if !ecdsa.Verify(&publicKey, txCopy.ID, &r, &s) {
			log.Println("WARNING: Signature verification failed")
			return false
		}
	}
	return true
}

/*
Coinbase transaction is a special transaction that has no inputs.
It is used to reward miners for mining a new block.
*/
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		log.Println("Creating coinbase transaction...")
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}
	txin := TXInput{[]byte{}, -1, nil, []byte(data)} // coinbase has no previous transaction
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	tx.ID = tx.Hash()
	return &tx
}

// UTXO = Unspent Transaction Output
func NewUTXOTransaction(from, to string, amount int, UTXOSet *UTXOSet) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	wallets, err := NewWallets()
	if err != nil {
		panic(err)
	}

	wallet := wallets.GetWallet(from)
	publicKeyHash := HashPubKey(wallet.PublicKey)
	acc, validOutputs := UTXOSet.FindSpendableOutputs(publicKeyHash, amount)
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
	UTXOSet.Blockchain.SignTransaction(&tx, wallet.PrivateKey)
	return &tx
}
