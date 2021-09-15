package main

import (
	"fmt"
	"strconv"

	"github.com/leo201313/Blockchain_with_Go/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()

	chain.AddBlock("first block after genesis")
	chain.AddBlock("second block after genesis")
	chain.AddBlock("third block after genesis")

	for _, block := range chain.Blocks {
		fmt.Printf("Previous hash:%x\n", block.PrevHash)
		fmt.Printf("data:%s\n", block.Data)
		fmt.Printf("hash:%x\n", block.Hash)

		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
