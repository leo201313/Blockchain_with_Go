package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"

	"github.com/leo201313/Blockchain_with_Go/blockchain"
	"github.com/leo201313/Blockchain_with_Go/constcoe"
	"github.com/leo201313/Blockchain_with_Go/utils"
)

func StartClient() {
	exitFlag := false
	if !exitFlag {
	}
}

func sendVersion(conn net.Conn, addr string) {
	bc := blockchain.ContinueBlockChain()
	nowHeight := bc.GetCurrentBlock().Height
	var encoded bytes.Buffer

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(MSGversion{constcoe.Version, nowHeight, addr})
	utils.Handle(err)

	msg := append(String2Bytes(versionCommand), encoded.Bytes()...)

	err = SendData(conn, msg)
	if err != nil {
		fmt.Printf("send failed,err:%v\n", err)
		log.Panic(err)
	}
}

func sendGetBlocks(conn net.Conn, addr string) {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(MSGgetblocks{addr})
	utils.Handle(err)

	msg := append(String2Bytes(getblocksCommand), encoded.Bytes()...)

	err = SendData(conn, msg)
	if err != nil {
		fmt.Printf("send failed,err:%v\n", err)
		log.Panic(err)
	}

}

func SendData(conn net.Conn, msg []byte) error {
	_, err := conn.Write(msg)
	return err
}
