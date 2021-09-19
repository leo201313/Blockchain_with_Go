package blockchain

import (
	"bytes"
	"log"

	"github.com/leo201313/Blockchain_with_Go/wallet"
	"github.com/mr-tron/base58/base58"
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

func (in *TxInput) CanUnlock(data string) bool {
	return bytes.Compare(in.Sig, []byte(data)) == 0
}

func (in *TxInput) RightKey(pubkeyhash []byte) bool {
	hashpubkey := wallet.PublicKeyHash([]byte(in.PubKey))
	return bytes.Compare(hashpubkey, pubkeyhash) == 0

}

func (out *TxOutput) Lock(address []byte) {
	out.HashPubKey = Address2PubHash(address)
}

func Address2PubHash(address []byte) []byte {
	pubKeyHash, err := base58.Decode(string(address))
	if err != nil {
		log.Panic(err)
	}
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-wallet.ChecksumLength]
	return pubKeyHash
}

func (out *TxOutput) CanBeUnlocked(pubkeyhash []byte) bool {
	return bytes.Compare(out.HashPubKey, pubkeyhash) == 0
}
