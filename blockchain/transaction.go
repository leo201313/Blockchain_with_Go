package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
)

const reward = 100

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func ToHexString(str string) []byte {
	return []byte(str)
}

func CoinbaseTx(toAddress, signature, publickey []byte) *Transaction {
	if bytes.Compare(signature, nil) == 0 {
		signature = bytes.Join([][]byte{
			ToHexString("Coins to "),
			toAddress,
		},
			[]byte{},
		)
	}

	txIn := TxInput{[]byte{}, -1, signature, publickey}

	txOut := TxOutput{reward, Address2PubHash(toAddress)}

	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}

	return &tx
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].TxID) == 0 && tx.Inputs[0].Out == -1
}

func NewTransaction(from, to []byte, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := chain.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("Error: Not enough funds!")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)

		for _, out := range outs {
			input := TxInput{txID, out, []byte{}, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, Address2PubHash(to)})

	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, Address2PubHash(from)})
	}

	tx := Transaction{nil, inputs, outputs}

	tx.SetID()

	return &tx
}
