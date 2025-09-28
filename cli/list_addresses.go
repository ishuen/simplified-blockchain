package cli

import (
	"fmt"
	"simplified-blockchain/core"
)

func (cli *CLI) listAddresses() {
	wallets, err := core.NewWallets()
	if err != nil {
		panic(err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}
