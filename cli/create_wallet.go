package cli

import (
	"fmt"
	"simplified-blockchain/core"
)

func (cli *CLI) createWallet() {
	wallets, _ := core.NewWallets()
	address := wallets.CreateWallet()
	err := wallets.SaveToFile()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Your new address: %s\n", address)
}
