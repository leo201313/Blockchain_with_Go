package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp    int64
	Height       uint32
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int64
}

func CreateBlock(txs []*Transaction, prevHash []byte, height uint32) *Block {
	block := &Block{time.Now().Unix(), height, []byte{}, txs, prevHash, 0}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func Genesis(coinbase *Transaction) *Block {
	originalByte := []byte("Leo Cao is awesome!")
	return CreateBlock([]*Transaction{coinbase}, originalByte, 0)
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	Handle(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	Handle(err)
	return &block
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

// func BytesToUint32(inbytes []byte) uint32 {
// 	bytesBuffer := bytes.NewBuffer(inbytes)
// 	var outint uint32
// 	err := binary.Read(bytesBuffer, binary.BigEndian, &outint)
// 	Handle(err)
// 	return outint
// }

// func Uint32ToBytes(inint uint32) []byte {
// 	bytesBuffer := bytes.NewBuffer([]byte{})
// 	err := binary.Write(bytesBuffer, binary.BigEndian, &inint)
// 	Handle(err)
// 	outbytes := bytesBuffer.Bytes()
// 	return outbytes
// }
