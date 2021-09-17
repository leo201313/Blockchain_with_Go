package wallet

import (
	"log"

	"github.com/mr-tron/base58"
)

func base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

func base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	if err != nil {
		log.Panic(err)
	}
	return decode
}
