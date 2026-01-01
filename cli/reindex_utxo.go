package cli

import (
	"fmt"
	"simplified-blockchain/core"
)

func (cli *CLI) reindexUTXO() {
	bc := core.GetBlockchain()
	UTXOSet := core.UTXOSet{bc}
	UTXOSet.Reindex()

	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}