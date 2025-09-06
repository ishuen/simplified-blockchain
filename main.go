package main

import "fmt"

func main() {
	chain := NewBlockchain()

	chain.AddBlock(("Send 1 BTC to Ivan"))
	chain.AddBlock(("Send 1 more BTC to Ivan"))

	for _, block := range chain.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash) // hexadecimal format
		fmt.Println()
	}
}
