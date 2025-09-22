package cli

import (
	"fmt"
	"simplified-blockchain/core"
)

func (cli *CLI) getBalance(address string) {
	bc := core.GetBlockchain(address)
	defer bc.Db.Close()
	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
