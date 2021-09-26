package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"runtime"

	"github.com/dgraph-io/badger"
	"github.com/leo201313/Blockchain_with_Go/constcoe"
	"github.com/leo201313/Blockchain_with_Go/utils"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (chain *BlockChain) GetCurrentBlock() *Block {
	var block *Block
	err := chain.Database.View(func(txn *badger.Txn) error {

		item, err := txn.Get(chain.LastHash)
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})
		utils.Handle(err)
		return err
	})

	utils.Handle(err)
	return block
}

func (chain *BlockChain) AddBlock(transactions []*Transaction, height uint32) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		utils.Handle(err)

		return err
	})

	utils.Handle(err)

	newBlock := CreateBlock(transactions, lastHash, height) // doing PoW

	err = chain.Database.Update(func(transaction *badger.Txn) error {
		err := transaction.Set(newBlock.Hash, newBlock.Serialize())
		utils.Handle(err)
		err = transaction.Set([]byte("lh"), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
	utils.Handle(err)
}

func InitBlockChain(address []byte) *BlockChain {
	var lastHash []byte

	if utils.DBexists(constcoe.DbFile) {
		fmt.Println("blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(constcoe.DbPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	utils.Handle(err)

	err = db.Update(func(txn *badger.Txn) error {

		cbtx := CoinbaseTx(address, []byte(constcoe.GenesisData), []byte{})
		genesis := Genesis(cbtx)
		fmt.Println("Genesis Created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		utils.Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)
		utils.Handle(err)
		err = txn.Set([]byte("ogprevhash"), genesis.PrevHash)
		utils.Handle(err)
		lastHash = genesis.Hash

		return err

	})
	utils.Handle(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func ContinueBlockChain() *BlockChain {
	if utils.DBexists(constcoe.DbFile) == false {
		fmt.Println("No blockchain found, please create one first")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(constcoe.DbPath)
	opts.Logger = nil
	db, err := badger.Open(opts)
	utils.Handle(err)

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	chain := BlockChain{lastHash, db}
	return &chain

}

func (chain *BlockChain) BackOgPrevHash() []byte {
	var ogprevhash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("ogprevhash"))
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			ogprevhash = val
			return nil
		})

		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	return ogprevhash
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{chain.LastHash, chain.Database}
	return &iterator
}

func (iterator *BlockChainIterator) Next() *Block {
	var block *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	iterator.CurrentHash = block.PrevHash

	return block
}

func (chain *BlockChain) FindUnspentTransactions(address []byte) []Transaction {
	var unspentTxs []Transaction

	spentTXNs := make(map[string][]int) // can't use type []byte as key value
	iter := chain.Iterator()

all:
	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXNs[txID] != nil {
					for _, spentOut := range spentTXNs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(Address2PubHash(address)) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					if in.CanUnlock(string(address)) {
						inTxID := hex.EncodeToString(in.TxID)
						spentTXNs[inTxID] = append(spentTXNs[inTxID], in.Out)
					}
				}
			}
			if bytes.Equal(block.PrevHash, chain.BackOgPrevHash()) {
				break all
			}
		}

	}
	return unspentTxs

}

func (chain *BlockChain) FindUTXO(address []byte) []UTXO {
	var UTXOs []UTXO
	unspentTransactions := chain.FindUnspentTransactions(address)
	for _, tx := range unspentTransactions {
		for outidx, out := range tx.Outputs {
			if out.CanBeUnlocked(Address2PubHash(address)) {
				tempUTXO := UTXO{tx.ID, outidx, out.Value}
				UTXOs = append(UTXOs, tempUTXO)
			}
		}
	}
	return UTXOs
}

func (chain *BlockChain) FindSpendableOutputs(address []byte, amount int) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := chain.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(Address2PubHash(address)) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = outIdx
				if accumulated >= amount {
					break Work
				}
				continue Work //one transaction can only have one output fit.
			}
		}
	}

	return accumulated, unspentOuts
}

func (bc *BlockChain) FindTransaction(txID []byte) (Transaction, error) {
	bIter := bc.Iterator()
	for {
		block := bIter.Next()

		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, txID) {
				return *tx, nil
			}
		}

		if bytes.Equal(block.PrevHash, bc.BackOgPrevHash()) {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)
	for _, input := range tx.Inputs {
		prevTX, err := bc.FindTransaction(input.TxID)
		utils.Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	tx.Sign(privKey, prevTXs)
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	prevTXs := make(map[string]Transaction)
	for _, input := range tx.Inputs {
		prevTX, err := bc.FindTransaction(input.TxID)
		utils.Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}
