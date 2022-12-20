package dto

type CreateReq struct {
	Creator     string `json:"creator"`
	FileHash    string `json:"filehash"`
	FileName    string `json:"filename"`
	FileSize    int64  `json:"filesize"`
	Description string `json:"description"`
	CoverImage  string `json:"image"`
	Length      string `json:"length"`
	Label       string `json:"label"`
}

type NftReq struct {
	FileHash string `json:"filehash"`
	TxHash   string `json:"txhash,omitempty"`
	Token    string `json:"token,omitempty"`
	Signtx   string `json:"signtx,omitempty"`
	To       string `json:"to,omitempty"`
	Price    string `json:"price,omitempty"`
	Status   string `json:"status,omitempty"`
}
