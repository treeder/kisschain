package main

import "fmt"

func main() {
	blockchain := []*Block{CreateGenesisBlock()}
	prevBlock := blockchain[0]

	// How many blocks should we add to the chain
	// after the genesis block
	numBlocksToAdd := 20

	// Add blocks to the chain
	for i := 0; i < numBlocksToAdd; i++ {
		b := NextBlock(prevBlock)
		blockchain = append(blockchain, b)
		prevBlock = b
		// Tell everyone about it!
		fmt.Printf("Block %v has been added to the blockchain!\n", b.Index)
		fmt.Printf("Hash: %v\n", b.HashHex())
	}
}
