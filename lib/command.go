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
	}

	cmd := args[0]

	switch cmd {
	case "key":
		bitcoin.handleKeyCommand(args[1:])
	case "tx":
		bitcoin.handleTxCommand(args[1:])
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
	}
}

func (bitcoin *Bitcoin) handleTxCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("tx's operations:\n new")
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
	}
}
