package core

import (
	"encoding/json"
	"os"
)

const walletFile = "wallet.dat"

type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFromFile()
	return &wallets, err
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := string(wallet.GetAddress())
	ws.Wallets[address] = wallet
	return address
}

func (ws *Wallets) GetAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}
	fileContent, err := os.ReadFile(walletFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileContent, ws)
	if err != nil {
		panic(err)
	}
	return nil
}

func (ws Wallets) SaveToFile() error {
	jsonData, err := json.Marshal(ws)
	if err != nil {
		return err
	}

	// file mode 0644: owner can read/write, group and others can read
	err = os.WriteFile(walletFile, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}
