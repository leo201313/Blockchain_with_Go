package constcoe

type configuration struct {
	Version  byte   // the version of protocal
	Mode     string // Full node or SPV
	Address  string // the Ip address of the node, e.g. 127.0.0.1:9600
	MaxWaitInterval int    // the Max interval for waiting
	PoolSize int    //how many transactions can be hold

}
