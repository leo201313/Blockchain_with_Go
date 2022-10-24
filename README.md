# Blockchain_with_Go
**Use golang to reproduce the basic blockchain**

Here is a tutorial to realize a simple basic [blockchain using golang](https://www.krad.top/goblockchain01/). Any you'll find much more systematical and clear codes at [here](https://github.com/leo201313/goblockchain_tutorial).

### Update Panel
* V0.1 No transactions but noly blocks are allowed.
* V0.5 Transactions are now supported. One block can have multiple transactions.
* V1.0 The UTXO module has been just supported. Users can now publish transactions to reallocate the coins.
* V1.1 Wallet module has been added, but it is not fully supported by the blockchain yet.
* V1.2 Now you can use the wallet address to make TxOutputs. Also you can use wallet address to refer the transactions.
* V1.5 There is a big jump in this version. Wallet module can be fully supported (Signature and Validation has been done), and even an API for the future mining functionality has been created.
* V1.6 Now supports the Merkle Tree and SPV.
* V1.7 Add UTXO sets to accelerate finding the spendable outputs of a wallet instead of go through the whole blockchain.


### Insight Future
<font color='red'> **Time is now to go straight forward to construct the P2P network step by step!**</font>
* Write server and client programs. It should follow the P2P protocal and seperate the nodes into full nodes and others.
* Activate the mining mechanism of the network.
* Realize the self-adaption of difficulty (using RPC).

### How to test and use
Recently I have made a .bat to test my program. If you want to know how to learn from this half-way program, just run the test.bat and see what I have done at this stage.

### Requirements
module github.com/leo201313/Blockchain_with_Go

go 1.17

require github.com/dgraph-io/badger v1.6.2

    require (
        github.com/AndreasBriese/bbloom v0.0.0-20190825152654-46b345b51c96 // indirect
        github.com/cespare/xxhash v1.1.0 // indirect
        github.com/dgraph-io/ristretto v0.0.2 // indirect
        github.com/dustin/go-humanize v1.0.0 // indirect
        github.com/golang/protobuf v1.3.1 // indirect
        github.com/mr-tron/base58 v1.2.0
        github.com/pkg/errors v0.8.1 // indirect
        golang.org/x/crypto v0.0.0-20210915214749-c084706c2272
        golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
        golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
    )

