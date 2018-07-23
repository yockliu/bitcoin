package lib

import (
	"fmt"

	. "github.com/yockliu/bitcoinlib"
)

// Wallet some place that hold the keys
type Wallet struct {
	defaultAddress string
	keysMap        map[string]*KeyPair
}

// NewWallet create Wallet instance pointer
func NewWallet() *Wallet {
	wallet := Wallet{}
	wallet.defaultAddress = ""
	wallet.keysMap = make(map[string]*KeyPair)
	return &wallet
}

func (wallet *Wallet) keyCount() int {
	return len(wallet.keysMap)
}

func (wallet *Wallet) newAddress() string {
	keypair := GenKeyPair()
	address := fmt.Sprintf("%s", keypair.Address)
	wallet.keysMap[address] = keypair
	if len(wallet.keysMap) == 1 {
		wallet.defaultAddress = address
	}
	return address
}

func (wallet *Wallet) addressList() []string {
	addressList := []string{}

	for k := range wallet.keysMap {
		addressList = append(addressList, k)
	}

	return addressList
}

func (wallet *Wallet) setDefaultAddress(address string) error {
	if _, ok := wallet.keysMap[address]; ok {
		wallet.defaultAddress = address
		return nil
	}
	return fmt.Errorf("no such address")
}

func (wallet *Wallet) getDefaultAddress() string {
	return wallet.defaultAddress
}
