package blockchain

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"

	"github.com/leo201313/Blockchain_with_Go/constcoe"
)

type CandidateBlock struct {
	PubTx []*Transaction
}

func (cb *CandidateBlock) SaveFile() {
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(cb)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(constcoe.CandidateBlockFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

func (cb *CandidateBlock) LoadFile() error {
	if _, err := os.Stat(constcoe.CandidateBlockFile); os.IsNotExist(err) {
		return err
	}

	var candidateBlock CandidateBlock

	fileContent, err := ioutil.ReadFile(constcoe.CandidateBlockFile)
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
	err = decoder.Decode(&candidateBlock)

	if err != nil {
		return err
	}

	cb.PubTx = candidateBlock.PubTx
	return nil
}

func CreateCandidateBlock() (*CandidateBlock, error) {
	candidateblock := CandidateBlock{}

	err := candidateblock.LoadFile()
	if os.IsNotExist(err) {
		return &candidateblock, nil
	}
	return &candidateblock, err
}

func (cb *CandidateBlock) AddTransaction(transaction *Transaction) {
	cb.PubTx = append(cb.PubTx, transaction)
}

func RemoveCandidateBlockFile() error {
	err := os.Remove(constcoe.CandidateBlockFile)
	return err
}
