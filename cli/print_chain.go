package cli

import (
	"fmt"
	"simplified-blockchain/core"
	"strconv"
)

func (cli *CLI) printChain() {
	bc := core.GetBlockchain("")
	defer bc.Db.Close()
	bci := bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := core.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		// break if it is the genesis block
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
