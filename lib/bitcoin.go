package lib

import (
	"fmt"

	. "github.com/yockliu/bitcointx"
	. "github.com/yockliu/blockchain"
)

// Bitcoin The bitcoin client core class
type Bitcoin struct {
	wallet     *Wallet
	blockChain *BlockChain
	storage    *Storage
	txPool     []*Transaction
}

// InitBitcoin init bitcoin instance
func InitBitcoin() *Bitcoin {
	wallet := NewWallet()
	blockChain := NewBlockChain()
	storage := NewStorage()

	txPool := []*Transaction{}

	bitcoin := Bitcoin{}
	bitcoin.wallet = wallet
	bitcoin.blockChain = blockChain
	bitcoin.storage = storage
	bitcoin.txPool = txPool

	return &bitcoin
}

func (bitcoin *Bitcoin) newKey() string {
	return bitcoin.wallet.newAddress()
}

func (bitcoin *Bitcoin) listAdress() []string {
	return bitcoin.wallet.addressList()
}

func (bitcoin *Bitcoin) addTransaction(value int, fromAddress string, toAddress string) error {
	outpoints, totalInput := bitcoin.findUTXO(value, fromAddress)

	if len(outpoints) == 0 {
		return fmt.Errorf("find not UTXO")
	}

	keypair := bitcoin.wallet.keysMap[fromAddress]

	inputs := []*TXIn{}
	outputs := []*TXOut{}

	for _, outpoint := range outpoints {
		outpoint.Lock = true
		in := NewTXIn(outpoint, keypair)
		inputs = append(inputs, in)
	}

	toOutput := NewTXOut(uint64(value), toAddress)
	outputs = append(outputs, toOutput)
	if totalInput > value {
		changeOutput := NewTXOut(uint64(totalInput-value), fromAddress)
		outputs = append(outputs, changeOutput)
	}

	tx := NewTransaction(inputs, outputs)
	bitcoin.txPool = append(bitcoin.txPool, tx)

	fmt.Println(bitcoin.txPool)

	return nil
}

func (bitcoin *Bitcoin) findUTXO(value int, address string) ([]*Outpoint, int) {
	if value <= 0 {
		return []*Outpoint{}, 0
	}

	outpoints := bitcoin.storage.findUTXO(address)

	total := 0
	outs := []*Outpoint{}
	for _, outpoint := range outpoints {
		if outpoint.Lock {
			continue
		}

		total = total + int(outpoint.Utxo.Value)
		outs = append(outs, outpoint)

		if total >= value {
			return outs, total
		}
	}

	return []*Outpoint{}, 0
}
