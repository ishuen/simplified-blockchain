package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CLI struct{}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  printchain - print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
}

func (cli *CLI) printChain() {
	bc := NewBlockchain("")
	defer bc.db.Close()
	bci := bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		// break if it is the genesis block
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) getBalance(address string) {
	bc := NewBlockchain(address)
	defer bc.db.Close()
	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) createBlockchain(address string) {
	bc := CreateBlockchain(address)
	defer bc.db.Close()
	fmt.Println("Done!")
}

func (cli *CLI) send(from, to string, amount int) {
	bc := NewBlockchain(from)
	defer bc.db.Close()
	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.AddBlock([]*Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceCommand := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCommand := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCommand := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCommand := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCommand.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCommand.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCommand.String("from", "", "Source wallet address")
	sendTo := sendCommand.String("to", "", "Destination wallet address")
	sendAmount := sendCommand.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCommand.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
	case "createblockchain":
		err := createBlockchainCommand.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
	case "send":
		err := sendCommand.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
	case "printchain":
		err := printChainCommand.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCommand.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCommand.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCommand.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCommand.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if sendCommand.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCommand.Usage()
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if printChainCommand.Parsed() {
		cli.printChain()
	}
}
