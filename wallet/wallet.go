package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
	"github.com/leo201313/Blockchain_with_Go/blockchain"
	"github.com/leo201313/Blockchain_with_Go/constcoe"
	"github.com/leo201313/Blockchain_with_Go/utils"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	RefName    string // The nick name of your wallet, just for local test.
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	utils.Handle(err)
	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pub

}

func PublicKeyHash(publicKey []byte) []byte {
	hashedPublicKey := sha256.Sum256(publicKey)
	hasher := ripemd160.New()
	_, err := hasher.Write(hashedPublicKey[:])
	utils.Handle(err)
	publicRipeMd := hasher.Sum(nil)
	return publicRipeMd
}

func Checksum(ripeMdHash []byte) []byte {
	firstHash := sha256.Sum256(ripeMdHash)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:constcoe.ChecksumLength]
}

func (w *Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)
	versionedHash := append([]byte{constcoe.Version}, pubHash...)
	checksum := Checksum(versionedHash)
	finalHash := append(versionedHash, checksum...)
	address := utils.Base58Encode(finalHash)
	return address
}

func MakeWallet(name string) *Wallet {
	privateKey, publicKey := NewKeyPair()
	wallet := Wallet{name, privateKey, publicKey}
	return &wallet
}

func (wt *Wallet) MakeTransaction(toaddress []byte, amount int, chain *blockchain.BlockChain) {
	candidateblock, err := blockchain.CreateCandidateBlock()
	utils.Handle(err)
	candidateblock.AddTransaction(blockchain.NewTransaction(wt.Address(), toaddress, amount, chain, wt.PrivateKey))
	candidateblock.SaveFile()
}

func (wt *Wallet) GetFileAddress() string {
	strAddress := string(wt.Address())
	fileAddress := constcoe.WalletPath + "/" + strAddress + "/" + "MANIFEST"
	return fileAddress
}

func (wt *Wallet) GetFilePath() string {
	strPath := string(wt.Address())
	filePath := constcoe.WalletPath + "/" + strPath
	return filePath
}

func (wt *Wallet) LoadUTXOSet() *UTXOSet {
	if utils.DBexists(wt.GetFileAddress()) == false {
		fmt.Println("No UTXOSet found for this wallet, please create one first")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(wt.GetFilePath())
	opts.Logger = nil
	db, err := badger.Open(opts)
	utils.Handle(err)
	utxoSet := UTXOSet{wt.Address(), db}

	return &utxoSet
}

func (wt *Wallet) CreateUTXOSet(chain *blockchain.BlockChain) *UTXOSet {
	if utils.DBexists(wt.GetFileAddress()) {
		fmt.Println("UTXOSet has already existed, now rebuild it.")
		err := os.RemoveAll(wt.GetFilePath())
		utils.Handle(err)
	}

	opts := badger.DefaultOptions(wt.GetFilePath())
	opts.Logger = nil
	db, err := badger.Open(opts)
	utils.Handle(err)

	utxoSet := UTXOSet{wt.Address(), db}

	UTXOs := chain.FindUTXO(wt.Address())

	err = db.Update(func(txn *badger.Txn) error {
		for _, utxo := range UTXOs {
			err = txn.Set(utxo.TxID, utxo.Serialize())
			utils.Handle(err)
		}
		return err
	})
	return &utxoSet
}

func (us *UTXOSet) GetBalance() int {
	amount := 0
	err := us.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				tmpUTXO := blockchain.DeserializeUTXO(v)
				amount += tmpUTXO.Value
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	utils.Handle(err)
	return amount
}

func (us *UTXOSet) AddUTXO(utxo *blockchain.UTXO) {
	err := us.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(utxo.TxID, utxo.Serialize())
		utils.Handle(err)
		return err
	})
	utils.Handle(err)
}

func (us *UTXOSet) DelUTXO(utxo *blockchain.UTXO) {
	err := us.DB.Update(func(txn *badger.Txn) error {
		err := txn.Delete(utxo.TxID)
		utils.Handle(err)
		return err
	})
	utils.Handle(err)

}
