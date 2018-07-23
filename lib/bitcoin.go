package lib

import (
	"fmt"

	. "github.com/yockliu/bitcoinlib"
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

	// fmt.Println(bitcoin.txPool)

	return nil
}

// CaculateNewBlock caculate new block
func (bitcoin *Bitcoin) CaculateNewBlock() {
	defaultAddress := bitcoin.wallet.defaultAddress
	if len(defaultAddress) == 0 {
		return
	}

	height := bitcoin.blockChain.Height()

	coinbase := NewCoinbase(defaultAddress, uint32(height+1), 0x12A05F200)

	txes := bitcoin.popTx()

	contents := []Cell{}
	contents = append(contents, coinbase)
	for _, cell := range txes {
		contents = append(contents, cell)
	}

	block := bitcoin.blockChain.GenerateBlock(contents, 20)

	// fmt.Printf("\nbitcoin caculateNewBlock: %x, \n:", block.Hash())

	bitcoin.restorage(block)
}

func (bitcoin *Bitcoin) popTx() []*Transaction {
	tx := bitcoin.txPool
	bitcoin.txPool = []*Transaction{}
	return tx
}

func (bitcoin *Bitcoin) restorage(block *Block) {
	storage := bitcoin.storage

	// storage block
	storage.saveBlock(block)

	var cell interface{}
	for _, cell = range block.Contents {

		switch v := cell.(type) {
		case *Coinbase:
			// only add a new UTX
			coinbase, _ := cell.(*Coinbase)
			outpoint := Outpoint{}
			outpoint.TxHash = *coinbase.Hash()
			outpoint.N = 0
			outpoint.Utxo = coinbase.Out
			outpoint.Lock = false
			storage.saveUTXO(coinbase.Out.Address, &outpoint)

		case *Transaction:
			// remove input UTXO & add output UTXO
			tx, _ := cell.(*Transaction)
			ins := tx.GetIns()
			for _, in := range ins {
				outpoint := in.GetOutpoint()
				outpoint.Lock = false
				storage.removeUTXO(outpoint.Utxo.Address, outpoint)
			}
			outs := tx.GetOuts()
			for index, out := range outs {
				outpoint := Outpoint{}
				outpoint.TxHash = *tx.Hash()
				outpoint.N = uint32(index)
				outpoint.Utxo = out
				outpoint.Lock = false
				storage.saveUTXO(out.Address, &outpoint)
			}

		default:
			fmt.Println("default cell.(type): ", v)
		}
	}

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
