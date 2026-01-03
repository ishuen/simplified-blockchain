package cli

import (
	"fmt"
	"simplified-blockchain/core"
)

func (cli *CLI) send(from, to string, amount int) {
	if !core.ValidateAddress(from) {
		fmt.Println("ERROR: Sender address is not valid")
		return
	}
	if !core.ValidateAddress(to) {
		fmt.Println("ERROR: Recipient address is not valid")
		return
	}
	bc := core.GetBlockchain()
	UTXOSet := core.UTXOSet{bc}
	defer bc.Db.Close()
	tx := core.NewUTXOTransaction(from, to, amount, &UTXOSet)
	cbTx := core.NewCoinbaseTX(from, "")
	txs := []*core.Transaction{cbTx, tx}

	fmt.Println("transactions: ", txs)

	newBlock := bc.MineBlock(txs)
	UTXOSet.Update(newBlock)
	fmt.Println("Success!")
}
