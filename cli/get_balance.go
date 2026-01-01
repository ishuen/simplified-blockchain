package cli

import (
	"fmt"
	"simplified-blockchain/core"
)

func (cli *CLI) getBalance(address string) {
	if !core.ValidateAddress(address) {
		panic("ERROR: Address is not valid")
	}
	bc := core.GetBlockchain()
	UTXOSet := core.UTXOSet{bc}
	defer bc.Db.Close()
	balance := 0
	publicKeyHash := core.Base58Decode([]byte(address))
	publicKeyHash = publicKeyHash[1 : len(publicKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(publicKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
