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
	bc := core.GetBlockchain(from)
	defer bc.Db.Close()
	tx := core.NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*core.Transaction{tx})
	fmt.Println("Success!")
}
