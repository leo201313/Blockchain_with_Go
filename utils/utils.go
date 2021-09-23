package utils

import (
	"log"
	"os"

	"github.com/mr-tron/base58"
)

func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)
	return []byte(encode)
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	Handle(err)
	return decode
}

func DBexists(db string) bool {
	if _, err := os.Stat(db); os.IsNotExist(err) {
		return false
	}
	return true
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
