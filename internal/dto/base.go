package dto

type OnlyWalletAddress struct {
	Address string `binding:"required,min=40,max=42" json:"address" form:"address"`
}

type WalletSign struct {
	OnlyWalletAddress
	Nonce int64  `form:"nonce" binding:"required" json:"nonce"`
	Sign  string `form:"sign" binding:"required" json:"sign"`
}

type WalletSignWrapper struct {
	WalletSign
}

func (t WalletSignWrapper) Address() string {
	return t.WalletSign.Address
}

func (t WalletSignWrapper) Nonce() int64 {
	return t.WalletSign.Nonce
}

func (t WalletSignWrapper) Sign() string {
	return t.WalletSign.Sign
}
