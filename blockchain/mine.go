package blockchain

import (
	"fmt"

	"github.com/leo201313/Blockchain_with_Go/utils"
)

func (chain *BlockChain) RunMine() {
	candidateBlock, err := CreateCandidateBlock()
	utils.Handle(err)
	currentHeight := chain.GetCurrentBlock().Height + 1
	chain.AddBlock(candidateBlock.PubTx, currentHeight)
	err = RemoveCandidateBlockFile()
	utils.Handle(err)

	currentBlock := chain.GetCurrentBlock()                                                                            //Test SPV
	route, hashroute := currentBlock.MTree.BackValidationRoute(candidateBlock.PubTx[0].ID)                             //Test SPV
	SPVwork := SimplePaymentValidation(candidateBlock.PubTx[0].ID, currentBlock.MTree.RootNode.Data, route, hashroute) //Test SPV
	fmt.Println("Whether SPV works: ", SPVwork)                                                                        //Test SPV

}
