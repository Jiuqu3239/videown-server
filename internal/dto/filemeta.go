package dto

type UserBrief struct {
	User       string `json:"user,omitempty"`
	FileName   string `json:"fileName,omitempty"`
	BucketName string `json:"bucketName,omitempty"`
}

type FileBlock struct {
	MinerId   int    `json:"minerId,omitempty"`
	BlockSize int    `json:"blockSize,omitempty"`
	BlockNum  int    `json:"blockNum,omitempty"`
	BlockId   string `json:"blockId,omitempty"`
	MinerIp   string `json:"minerIp,omitempty"`
	MinerAcc  string `json:"minerAcc,omitempty"`
}

type FileMeta struct {
	State      string      `json:"state,omitempty"`
	Size       int         `json:"size,omitempty"`
	Index      int         `json:"index,omitempty"`
	Blocks     []FileBlock `json:"blocks,omitempty"`
	UserBriefs []UserBrief `json:"userBriefs,omitempty"`
}

type QueryResponse struct {
	Ok FileMeta `json:"ok,omitempty"`
}
