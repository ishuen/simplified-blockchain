package cli

import (
	"fmt"
	"simplified-blockchain/core"
)

func (cli *CLI) send(from, to string, amount int) {
	bc := core.NewBlockchain(from)
	defer bc.Db.Close()
	tx := core.NewUTXOTransaction(from, to, amount, bc)
	bc.AddBlock([]*core.Transaction{tx})
	fmt.Println("Success!")
}
