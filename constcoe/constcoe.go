package constcoe

const (
	ChecksumLength     = 4
	Version            = byte(0x00)
	Reward             = 100
	WalletFile         = "./tmp/wallets/wallets.data"
	WalletPath         = "./tmp/wallets"
	CandidateBlockFile = "./tmp/transactions.data"
	DbPath             = "./tmp/blocks"
	DbFile             = "./tmp/blocks/MANIFEST"
	GenesisData        = "First Transaction in Genesis"
)
