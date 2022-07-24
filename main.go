package main

import (
	"fmt"
	"log"
	"os"

	"github.com/manishmeganathan/essensio/core"
)

func main() {
	// Check if a command has been entered
	if len(os.Args) < 2 {
		log.Fatalln("Command Not Found")
	}

	// Read the command and check its value
	switch cmd := os.Args[1]; cmd {
	// AddBlock command
	case "addblock":
		// Check if an input has been provided for the block data
		if len(os.Args) < 3 {
			log.Fatalln("Missing Input for 'AddBlock' command")
		}

		// Load up the BlockChain
		chain, err := core.NewBlockChain()
		if err != nil {
			log.Fatalln("Failed to Start Blockchain:", err)
		}

		// Add the Block with the given data to the chain
		if err := chain.AddBlock(os.Args[2]); err != nil {
			log.Fatalln("Failed to Add Block to Chain:", err)
		}

	// ShowChain command
	case "showchain":
		// Load up the BlockChain
		chain, err := core.NewBlockChain()
		if err != nil {
			log.Fatalln("Failed to Start Blockchain:", err)
		}

		// Create a new iterator for the chain
		iterator := chain.NewIterator()
		fmt.Print("\nFull Chain Information\n\n")

		// Iterate until the chain is done
		for !iterator.Done() {
			// Get the next block
			block, err := iterator.Next()
			if err != nil {
				log.Fatalln("Iterator Error:", err)
			}

			// Print the Block
			fmt.Println(block)
		}

	default:
		log.Fatalf("Unsupported Command: %v\n", cmd)
	}
}
