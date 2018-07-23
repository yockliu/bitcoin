package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/yockliu/bitcoin/lib"
)

var interval = 10

func main() {
	bitcoin := lib.InitBitcoin()

	bitcoin.HandleCommand([]string{"key", "new"})
	bitcoin.HandleCommand([]string{"key", "new"})
	fmt.Println("")

	go loopBitcoin(bitcoin)

	for true {
		fmt.Printf(":")
		args, _ := interfaceScanln(10)

		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "interval":
			if len(args) < 2 {
				continue
			}
			i, err := strconv.Atoi(args[1])
			if err != nil {
				continue
			}
			interval = i
		default:
			bitcoin.HandleCommand(args)
		}

		fmt.Printf("")
	}
}

func interfaceScanln(n int) ([]string, error) {
	x := make([]string, n)
	y := make([]interface{}, len(x))
	for i := range x {
		y[i] = &x[i]
	}
	n, err := fmt.Scanln(y...)
	x = x[:n]
	return x, err
}

func loopBitcoin(bitcoin *lib.Bitcoin) {
	for true {
		bitcoin.CaculateNewBlock()

		time.Sleep(time.Duration(interval) * time.Second)
	}
}
