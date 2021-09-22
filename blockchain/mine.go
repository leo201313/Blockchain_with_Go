package blockchain

import "fmt"

func (chain *BlockChain) RunMine() {
	candidateBlock, err := CreateCandidateBlock()
	Handle(err)
	currentHeight := chain.GetCurrentBlock().Height + 1
	chain.AddBlock(candidateBlock.PubTx, currentHeight)
	err = RemoveCandidateBlockFile()
	Handle(err)

	currentBlock := chain.GetCurrentBlock()                                                                            //Test SPV
	route, hashroute := currentBlock.MTree.BackValidationRoute(candidateBlock.PubTx[0].ID)                             //Test SPV
	SPVwork := SimplePaymentValidation(candidateBlock.PubTx[0].ID, currentBlock.MTree.RootNode.Data, route, hashroute) //Test SPV
	fmt.Println("Whether SPV works: ", SPVwork)                                                                        //Test SPV

}
