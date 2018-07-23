package lib

import (
	"fmt"
	"strconv"
)

// HandleCommand handle commaned
func (bitcoin *Bitcoin) HandleCommand(args []string) {
	argsLen := len(args)
	if argsLen == 0 {
		fmt.Println("command is empty!\nThe Commands: \n key -- new key, address\n tx -- new transaction")
		return
	}

	cmd := args[0]

	switch cmd {
	case "key":
		bitcoin.handleKeyCommand(args[1:])
	case "tx":
		bitcoin.handleTxCommand(args[1:])
	case "block":
		bitcoin.handleBlockCommand(args[1:])
	}
}

func (bitcoin *Bitcoin) handleKeyCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("key's operations:\n new \n list")
	}

	op := args[0]

	switch op {
	case "new":
		address := bitcoin.newKey()
		fmt.Printf("New Address is %s\n", address)
	case "list":
		list := bitcoin.listAdress()
		for index, address := range list {
			fmt.Printf("%d --- %s\n", index, address)
		}
	case "default":
		if len(args) != 2 {
			fmt.Println("set default address arg err: address")
		}
		address := args[1]
		err := bitcoin.wallet.setDefaultAddress(address)
		if err != nil {
			fmt.Println("set default address error: ", err)
		} else {
			fmt.Println("set default address ok: ", address)
		}
	}
}

func (bitcoin *Bitcoin) handleTxCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("tx's operations:\n new, utxo")
	}

	op := args[0]

	switch op {
	case "new":
		if len(args) != 4 {
			fmt.Println("arguments should be: value fromAddress toAddress")
			return
		}

		value, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("value error: %s\n", err)
			return
		}

		fromAddress := args[2]
		toAddress := args[3]
		err = bitcoin.addTransaction(value, fromAddress, toAddress)

		if err != nil {
			fmt.Printf("addTransaction error: %s\n", err)
		}
	case "utxo":
		if len(args) != 2 {
			fmt.Println("arguments should be: address")
			return
		}

		fmt.Println("utxo storage: ", bitcoin.storage.utxoIndex)

		if utxos, ok := bitcoin.storage.utxoIndex[args[1]]; ok {
			total := uint64(0)
			for _, utxo := range utxos {
				total += utxo.Utxo.Value
			}
			fmt.Println("utxo size = : ", len(utxos))
			fmt.Println("utxo total value = : ", total)
		}
	}
}

func (bitcoin *Bitcoin) handleBlockCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("block's operations:\n new")
	}

	op := args[0]
	blockChain := bitcoin.blockChain

	switch op {
	case "info":
		fmt.Println("Height: ", blockChain.Height())
		fmt.Printf("Current: %x\n", blockChain.Current().Hash())
	case "index":
		if len(args) != 2 {
			fmt.Println("arguments should be: index")
			return
		}
		index, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("argument error: ", err)
			return
		}
		block := blockChain.BlockOfHeight(index)
		fmt.Printf("Block %x: %x", index, block.Hash())
	}
}
