package cli

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/leo201313/Blockchain_with_Go/blockchain"
	"github.com/leo201313/Blockchain_with_Go/wallet"
)

type CommandLine struct{}

func (cli *CommandLine) printUsage() {
	fmt.Println("Welcome to Leo Cao's tiny blockchain system, usage is as follows:")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("All you need is to first create a wallet, and use the wallet address to init a blockchain.")
	fmt.Println("Then, you can add more wallets and try to make some transactions use 'send'.")
	fmt.Println("At this version, one block only contains one transaction. Also the mining is not considered yet.")
	fmt.Println("The nickname of wallet is just used to refer the wallet instead of typing the full address.")
	fmt.Println("Nothing left to say, wish you good luck.")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("createwallet -nickname NICKNAME               ----> Creates a new wallet with the nick name")
	fmt.Println("listwallets                                   ----> Lists all the wallets in the wallet file")
	fmt.Println("walletinfo -nickname NICKNAME                 ----> Back all the information of a wallet")
	fmt.Println("createblockchain -nickname NICKNAME           ----> Creates a blockchain using the name of the user's wallet")
	fmt.Println("blockchaininfo                                ----> Prints the blocks in the chain")
	fmt.Println("send -from FROMNAME -to TONAME -amount AMOUNT ----> Send amount of coins from one wallet address to another")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
}

func (cli *CommandLine) createWallet(nickname string) {
	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet(nickname)
	wallets.SaveFile()
	fmt.Printf("Owner:%s, New wallet address is:%s\n", nickname, address)
}

func (cli *CommandLine) listWallets() {
	wallets, _ := wallet.CreateWallets()
	address, names := wallets.GetAllAddresses()
	for i := 0; i < len(address); i++ {
		fmt.Printf("Wallet Address: %s , Reffered Name:%s\n", address[i], names[i])
	}
}

func (cli *CommandLine) getWalletInfo(nickname string) {
	wallets, _ := wallet.CreateWallets()
	aimwallet := wallets.GetWalletByName(nickname)
	address := aimwallet.Address()

	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	balance := 0
	UTXOs := chain.FindUTXO(address)
	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("Owner:%s, Address:%s, Balance:%d \n", nickname, string(address), balance)
}

func (cli *CommandLine) createBlockChain(nickname string) {
	wallets, _ := wallet.CreateWallets()
	aimwallet := wallets.GetWalletByName(nickname)
	address := aimwallet.Address()
	newChain := blockchain.InitBlockChain(address)
	newChain.Database.Close()
	fmt.Println("Finished creating chain")
}

func (cli *CommandLine) getBlockChainInfo() {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	iterator := chain.Iterator()
	ogprevhash := chain.BackOgPrevHash()
	for {
		block := iterator.Next()
		fmt.Printf("Previous hash:%x\n", block.PrevHash)
		fmt.Printf("Transactions:%v\n", block.Transactions)
		fmt.Printf("hash:%x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if bytes.Equal(block.PrevHash, ogprevhash) {
			break
		}
	}
}

func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	wallets, _ := wallet.CreateWallets()
	fromwallet := wallets.GetWalletByName(from)
	towallet := wallets.GetWalletByName(to)
	tx := blockchain.NewTransaction(fromwallet.Address(), towallet.Address(), amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listWalletsCmd := flag.NewFlagSet("listwallets", flag.ExitOnError)
	walletInfoCmd := flag.NewFlagSet("walletinfo", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	blockChainInfoCmd := flag.NewFlagSet("blockchaininfo", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	createWalletName := createWalletCmd.String("nickname", "", "The name to refer the wallet")
	walletInfoName := walletInfoCmd.String("nickname", "", "The name to refer the wallet")
	createBlockChainName := createBlockChainCmd.String("nickname", "", "The name of wallet to find address to init a blockchain")
	sendFrom := sendCmd.String("from", "", "Source wallet nickname")
	sendTo := sendCmd.String("to", "", "Destination wallet nickname")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listwallets":
		err := listWalletsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "walletinfo":
		err := walletInfoCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "blockchaininfo":
		err := blockChainInfoCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if createWalletCmd.Parsed() {
		if *createWalletName == "" {
			createWalletCmd.Usage()
			runtime.Goexit()
		}
		cli.createWallet(*createWalletName)
	}

	if walletInfoCmd.Parsed() {
		if *walletInfoName == "" {
			walletInfoCmd.Usage()
			runtime.Goexit()
		}
		cli.getWalletInfo(*walletInfoName)
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainName == "" {
			createBlockChainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockChainName)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if listWalletsCmd.Parsed() {
		cli.listWallets()
	}

	if blockChainInfoCmd.Parsed() {
		cli.getBlockChainInfo()
	}
}
