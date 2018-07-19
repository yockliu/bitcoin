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
	txIndex    map[HashCode]*Transaction
	utxoIndex  map[string][]*Outpoint
}

// NewStorage create new Storage instance pointer
func NewStorage() *Storage {
	storage := Storage{}

	return &storage
}

func (storage *Storage) findBlockByHash(hash HashCode) *Block {
	return storage.blockIndex[hash]
}

func (storage *Storage) saveBlock(block *Block) error {
	if existBlock, ok := storage.blockIndex[block.Hash()]; ok {
		if existBlock.MerkleRoot.Compare(&block.MerkleRoot) != 0 {
			return fmt.Errorf("the hash has another block that merkleroot conflict")
		}
	}
	storage.blockIndex[block.Hash()] = block
	return nil
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
			if outpoint == element {
				removedIndex = i
				break
			}
		}

		var newOutpoints []*Outpoint

		// if removedIndex == 0 {
		// 	newOutpoints = outpoints[1:]
		// } else if removedIndex == len(outpoints)-1 {
		// 	newOutpoints = outpoints[:len(outpoints)-1]
		// } else if removedIndex > 0 && removedIndex < len(outpoints)-1 {
		// 	head := outpoints[:removedIndex]
		// 	tail := outpoints[removedIndex+1 : len(outpoints)]
		// 	newOutpoints = append(head, tail...)
		// }

		if removedIndex >= 0 && removedIndex < len(outpoints) {
			head := outpoints[:removedIndex]
			tail := outpoints[removedIndex+1 : len(outpoints)]
			newOutpoints = append(head, tail...)
		}

		storage.utxoIndex[address] = newOutpoints
	}
}
