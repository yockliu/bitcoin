package main

import (
	"fmt"

	"github.com/yockliu/bitcoin/lib"
)

func main() {
	bitcoin := lib.InitBitcoin()

	for true {
		fmt.Println("--------------------------------")
		fmt.Println("Please Input Command:")
		args, _ := interfaceScanln(10)
		bitcoin.HandleCommand(args)
		fmt.Println("")
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
