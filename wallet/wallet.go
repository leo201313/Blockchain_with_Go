package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const (
	ChecksumLength = 4
	Version        = byte(0x00)
)

type Wallet struct {
	RefName    string // The nick name of your wallet, just for local test.
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pub

}

func PublicKeyHash(publicKey []byte) []byte {
	hashedPublicKey := sha256.Sum256(publicKey)
	hasher := ripemd160.New()
	_, err := hasher.Write(hashedPublicKey[:])
	if err != nil {
		log.Panic(err)
	}
	publicRipeMd := hasher.Sum(nil)
	return publicRipeMd
}

func Checksum(ripeMdHash []byte) []byte {
	firstHash := sha256.Sum256(ripeMdHash)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:ChecksumLength]
}

func (w *Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)
	versionedHash := append([]byte{Version}, pubHash...)
	checksum := Checksum(versionedHash)
	finalHash := append(versionedHash, checksum...)
	address := base58Encode(finalHash)
	return address
}

func MakeWallet(name string) *Wallet {
	privateKey, publicKey := NewKeyPair()
	wallet := Wallet{name, privateKey, publicKey}
	return &wallet
}
