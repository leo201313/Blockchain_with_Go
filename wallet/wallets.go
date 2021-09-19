package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "./tmp/wallets.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

func (ws *Wallets) SaveFile() {
	var content bytes.Buffer
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

func (ws *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		return err
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
	err = decoder.Decode(&wallets)

	if err != nil {
		return err
	}

	ws.Wallets = wallets.Wallets

	return nil
}

func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFile()

	return &wallets, err
}

func (ws *Wallets) AddWallet(name string) string {
	wallet := MakeWallet(name)
	address := fmt.Sprintf("%s", wallet.Address())

	ws.Wallets[address] = wallet

	return address
}

func (ws Wallets) GetWalletByAddress(address string) Wallet {
	aimwallet, ok := ws.Wallets[address]
	if ok != true {
		log.Panic("Error: No Wallet with such Address!")

	}

	return *aimwallet
}

func (ws *Wallets) GetAllAddresses() ([]string, []string) {
	var addresses []string
	var nickname []string

	for address, wallet := range ws.Wallets {
		addresses = append(addresses, address)
		nickname = append(nickname, wallet.RefName)
	}
	return addresses, nickname

}

func (ws *Wallets) GetWalletByName(name string) Wallet {
	var aimwallet Wallet

	empty := true
	for _, wallet := range ws.Wallets {
		if wallet.RefName == name {
			aimwallet = *wallet
			empty = false
			break
		}
	}
	if empty {
		log.Panic("Error: No Wallet with such Name!")
	}
	return aimwallet
}
