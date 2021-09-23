package wallet

import "github.com/dgraph-io/badger"

type UTXOSet struct {
	WalletAddress []byte
	DB            *badger.DB
}
