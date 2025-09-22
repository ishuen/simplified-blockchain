package cli

import (
	"fmt"
	"simplified-blockchain/core"
)

func (cli *CLI) createBlockchain(address string) {
	bc := core.CreateBlockchain(address)
	defer bc.Db.Close()
	fmt.Println("Done!")
}
