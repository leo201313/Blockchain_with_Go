package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"math/big"

	"github.com/leo201313/Blockchain_with_Go/constcoe"
)

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

	txOut := TxOutput{constcoe.Reward, Address2PubHash(toAddress)}

	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}

	return &tx
}

func (tx *Transaction) TxHash() []byte {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	return hash[:]
}

func (tx *Transaction) SetID() {
	tx.ID = tx.TxHash()
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].TxID) == 0 && tx.Inputs[0].Out == -1
}

func NewTransaction(from, to []byte, amount int, chain *BlockChain, privkey ecdsa.PrivateKey) *Transaction {
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
	chain.SignTransaction(&tx, privkey)
	return &tx
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}
	txCopy := tx.PlainCopy()

	for idx, input := range txCopy.Inputs {
		plainhash := txCopy.PlainHash(idx, prevTXs[hex.EncodeToString(input.TxID)].Outputs[input.Out].HashPubKey) // This is because we want to sign the inputs seperately!
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, plainhash)
		Handle(err)
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Inputs[idx].Sig = signature
	}

}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	txCopy := tx.PlainCopy()
	curve := elliptic.P256()

	for idx, input := range tx.Inputs {
		plainhash := txCopy.PlainHash(idx, prevTXs[hex.EncodeToString(input.TxID)].Outputs[input.Out].HashPubKey)

		r := big.Int{}
		s := big.Int{}
		sigLen := len(input.Sig)
		r.SetBytes(input.Sig[:(sigLen / 2)])
		s.SetBytes(input.Sig[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(input.PubKey)
		x.SetBytes(input.PubKey[:(keyLen / 2)])
		y.SetBytes(input.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, plainhash, &r, &s) == false {
			return false
		}

	}
	return true
}

func (tx *Transaction) PlainCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, txin := range tx.Inputs {
		inputs = append(inputs, TxInput{txin.TxID, txin.Out, nil, txin.PubKey})
	}

	for _, txout := range tx.Outputs {
		outputs = append(outputs, TxOutput{txout.Value, txout.HashPubKey})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

func (tx *Transaction) PlainHash(inidx int, prevPubKeyHash []byte) []byte {
	txCopy := tx.PlainCopy()
	txCopy.Inputs[inidx].PubKey = prevPubKeyHash
	return txCopy.TxHash()
}
