package cli

import (
	"fmt"
	"simplified-blockchain/core"
)

func (cli *CLI) createBlockchain(address string) {
	if !core.ValidateAddress(address) {
		fmt.Println("ERROR: Address is not valid")
		return
	}
	bc := core.CreateBlockchain(address)
	defer bc.Db.Close()

	UTXOSet := core.UTXOSet{bc}
	UTXOSet.Reindex()
	fmt.Println("Done!")
}
