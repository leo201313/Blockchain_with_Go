package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math"
	"math/big"

	"github.com/leo201313/Blockchain_with_Go/constcoe"
)

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-constcoe.Difficulty))
	pow := &ProofOfWork{b, target}
	return pow
}

func ToHexInt(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func (pow *ProofOfWork) InitNonce(nonce int64) []byte {
	data := bytes.Join([][]byte{
		pow.Block.PrevHash,
		pow.Block.HashTransactions(),
		pow.Block.MTree.RootNode.Data,
		ToHexInt(int64(nonce)),
		ToHexInt(int64(constcoe.Difficulty)),
	},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Run() (int64, []byte) {
	var intHash big.Int
	var hash [32]byte
	var nonce int64
	nonce = 0

	for nonce < math.MaxInt64 {
		data := pow.InitNonce(nonce)
		hash = sha256.Sum256(data)

		// fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	// fmt.Println()
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := pow.InitNonce(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}
