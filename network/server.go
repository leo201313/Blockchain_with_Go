package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"sync"

	"github.com/leo201313/Blockchain_with_Go/blockchain"
	"github.com/leo201313/Blockchain_with_Go/utils"
)

var localAddress string
var knownNodes = make([]string, 0)
var bcReadMutex sync.RWMutex

func addrIsKnown(addr string) bool {
	for _, tempaddr := range knownNodes {
		if tempaddr == addr {
			return true
		}
	}
	return false
}

func StartServer(localaddr string, knownodes []string) {
	localAddress = localaddr
	knownNodes = append(knownNodes, knownodes...)

	listen, err := net.Listen("tcp", localAddress)
	utils.Handle(err)
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		utils.Handle(err)
		go handleConnection(conn)
	}

	// bc := blockchain.ContinueBlockChain()

}

func handleConnection(conn net.Conn) {
	request, err := ioutil.ReadAll(conn)
	utils.Handle(err)
	command := Bytes2String(request[:commandLength])

	fmt.Printf("Received %s command\n", command)

	switch command {
	case versionCommand:
		handleVersion(conn, request)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

func handleVersion(conn net.Conn, request []byte) {
	var buff bytes.Buffer
	var msg MSGversion

	buff.Write(request[commandLength:])
	decoder := gob.NewDecoder(&buff)
	err := decoder.Decode(&msg)
	utils.Handle(err)

	bc := blockchain.ContinueBlockChain()
	myNowHeight := bc.GetCurrentBlock().Height
	foreignHeight := msg.NowHeight

	if myNowHeight < foreignHeight {
		sendGetBlocks(conn, localAddress)
	} else if myNowHeight > foreignHeight {
		sendVersion(conn, localAddress)
	}

	if !addrIsKnown(msg.FromAddress) {
		knownNodes = append(knownNodes, msg.FromAddress)
	}

}

func handleGetBlocks(conn net.Conn, request []byte) {

}
