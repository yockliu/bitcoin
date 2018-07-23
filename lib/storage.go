package lib

import (
	"fmt"

	. "github.com/yockliu/bitcoinlib"
	. "github.com/yockliu/bitcointx"
	. "github.com/yockliu/blockchain"
)

// Storage something that store
type Storage struct {
	blockIndex map[HashCode]*Block
	utxoIndex  map[string][]*Outpoint
}

// NewStorage create new Storage instance pointer
func NewStorage() *Storage {
	storage := Storage{}

	storage.blockIndex = map[HashCode]*Block{}
	storage.utxoIndex = map[string][]*Outpoint{}

	return &storage
}

func (storage *Storage) findBlockByHash(hash HashCode) *Block {
	return storage.blockIndex[hash]
}

func (storage *Storage) saveBlock(block *Block) {
	storage.blockIndex[block.Hash()] = block
}

func (storage *Storage) findUTXO(address string) []*Outpoint {
	if outpoints, ok := storage.utxoIndex[address]; ok {
		return outpoints
	}
	return []*Outpoint{}
}

func (storage *Storage) saveUTXO(address string, outpoint *Outpoint) {
	if outpoints, ok := storage.utxoIndex[address]; ok {
		outpoints = append(outpoints, outpoint)
		storage.utxoIndex[address] = outpoints
	} else {
		outpoints = []*Outpoint{outpoint}
		storage.utxoIndex[address] = outpoints
	}
}

func (storage *Storage) removeUTXO(address string, outpoint *Outpoint) {
	if outpoints, ok := storage.utxoIndex[address]; ok {
		removedIndex := -1

		for i, element := range outpoints {
			if outpoint == element || (outpoint.TxHash == element.TxHash && outpoint.N == element.N) {
				removedIndex = i
				break
			}
		}

		var newOutpoints []*Outpoint

		if removedIndex >= 0 && removedIndex < len(outpoints) {
			head := outpoints[:removedIndex]
			tail := outpoints[removedIndex+1 : len(outpoints)]
			newOutpoints = append(head, tail...)
		}

		storage.utxoIndex[address] = newOutpoints
	} else {
		fmt.Println("removeUTXO find outpoint not ok")
	}
}
