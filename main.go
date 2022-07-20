package main

import (
	"fmt"

	"github.com/manishmeganathan/essensio/core"
)

func main() {
	chain := core.NewBlockChain("Genesis Information")

	chain.AddBlock("Data 1")
	chain.AddBlock("Data 2")
	chain.AddBlock("Data 3")

	fmt.Println("====")
	for _, block := range chain.Blocks {
		fmt.Println(block)
		fmt.Println("====")
	}
}
