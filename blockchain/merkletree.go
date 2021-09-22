package blockchain

import (
	"bytes"
	"crypto/sha256"
	"log"
)

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	LeftNode  *MerkleNode
	RightNode *MerkleNode
	Data      []byte // Hash
}

func CreateMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	tempNode := MerkleNode{}

	if left == nil && right == nil { //The leaf

		tempNode.Data = data
	} else {
		catenateHash := append(left.Data, right.Data...)
		hash := sha256.Sum256(catenateHash)
		tempNode.Data = hash[:]
	}

	tempNode.LeftNode = left
	tempNode.RightNode = right

	return &tempNode

}

func CrateMerkleTree(txs []*Transaction) *MerkleTree {
	txslen := len(txs)
	if txslen%2 != 0 {
		txs = append(txs, txs[txslen-1])
	}

	var nodePool []*MerkleNode

	for _, tx := range txs {
		nodePool = append(nodePool, CreateMerkleNode(nil, nil, tx.ID))
	}

	for len(nodePool) > 1 {
		var tempNodePool []*MerkleNode
		poolLen := len(nodePool)
		if poolLen%2 != 0 {
			tempNodePool = append(tempNodePool, nodePool[poolLen-1])
		}
		for i := 0; i < poolLen/2; i++ {
			tempNodePool = append(tempNodePool, CreateMerkleNode(nodePool[2*i], nodePool[2*i+1], nil))
		}
		nodePool = tempNodePool
	}

	merkleTree := MerkleTree{nodePool[0]}

	return &merkleTree
}

func (mn *MerkleNode) Find(data []byte, route []int, hashroute [][]byte) (bool, []int, [][]byte) {
	findFlag := false

	if bytes.Equal(mn.Data, data) {
		findFlag = true
		return findFlag, route, hashroute

	} else {
		if mn.LeftNode != nil {
			route_t := append(route, 0)
			hashroute_t := append(hashroute, mn.RightNode.Data)
			findFlag, route_t, hashroute_t = mn.LeftNode.Find(data, route_t, hashroute_t)
			if findFlag {
				return findFlag, route_t, hashroute_t

			} else {
				if mn.RightNode != nil {
					route_t = append(route, 1)
					hashroute_t = append(hashroute, mn.LeftNode.Data)
					findFlag, route_t, hashroute_t = mn.RightNode.Find(data, route_t, hashroute_t)
					if findFlag {
						return findFlag, route_t, hashroute_t
					} else {
						return findFlag, route, hashroute
					}

				}
			}
		} else {
			return findFlag, route, hashroute
		}
	}

	return findFlag, route, hashroute
}

func (mt *MerkleTree) BackValidationRoute(txid []byte) ([]int, [][]byte) {

	ok, route, hashroute := mt.RootNode.Find(txid, []int{}, [][]byte{})
	if !ok {
		log.Panic("Error in BackValidationRoute: No such transaction.")

	}
	return route, hashroute
}

func SimplePaymentValidation(txid, mtroothash []byte, route []int, hashroute [][]byte) bool {
	routeLen := len(route)
	var tempHash []byte
	tempHash = txid

	for i := routeLen - 1; i >= 0; i-- {
		if route[i] == 0 {
			catenateHash := append(tempHash, hashroute[i]...)
			hash := sha256.Sum256(catenateHash)
			tempHash = hash[:]
		} else if route[i] == 1 {
			catenateHash := append(hashroute[i], tempHash...)
			hash := sha256.Sum256(catenateHash)
			tempHash = hash[:]
		} else {
			log.Panic("Error in SimplePaymentValidation.")
		}
	}

	return bytes.Equal(tempHash, mtroothash)

}
