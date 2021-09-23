package blockchain

import (
	"bytes"
	"encoding/gob"

	"github.com/leo201313/Blockchain_with_Go/constcoe"
	"github.com/leo201313/Blockchain_with_Go/utils"
)

type TxOutput struct {
	Value      int
	HashPubKey []byte
}

type TxInput struct {
	TxID   []byte
	Out    int
	Sig    []byte
	PubKey []byte
}

type UTXO struct {
	TxID   []byte
	OutIdx int
	Value  int
}

func (u *UTXO) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(u)
	utils.Handle(err)

	return res.Bytes()
}

func DeserializeUTXO(data []byte) *UTXO {
	var utxo UTXO
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&utxo)
	utils.Handle(err)
	return &utxo
}

func (in *TxInput) CanUnlock(data string) bool {
	return bytes.Compare(in.Sig, []byte(data)) == 0
}

func (out *TxOutput) Lock(address []byte) {
	out.HashPubKey = Address2PubHash(address)
}

func Address2PubHash(address []byte) []byte {
	pubKeyHash := utils.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-constcoe.ChecksumLength]
	return pubKeyHash
}

func (out *TxOutput) CanBeUnlocked(pubkeyhash []byte) bool {
	return bytes.Compare(out.HashPubKey, pubkeyhash) == 0
}
