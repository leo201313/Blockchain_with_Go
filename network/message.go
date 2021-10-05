package network

type MSGversion struct {
	Version     byte
	NowHeight   uint32
	FromAddress string
}

type MSGgetblocks struct {
	FromAddress string
}

const (
	commandLength    = 3
	versionCommand   = "000"
	getblocksCommand = "001"
)

func String2Bytes(str string) []byte {
	resBytes := []byte(str)
	return resBytes
}

func Bytes2String(byt []byte) string {
	return string(byt)
}
