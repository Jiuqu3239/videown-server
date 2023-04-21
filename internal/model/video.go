package model

import (
	"errors"

	"gorm.io/gorm"
)

type FileStatus int32

const (
	UPLOAD FileStatus = iota
	SCHEDULE
	STORAGE
	DELETE
)

// NFT Status
type NftStatus int32

const (
	CREATE NftStatus = iota
	MINT
	LIST
	MELT
)
const (
	NULL          = "--"
	DEFAULT_CHAIN = "CESS"
)

type VideoMetadata struct {
	Id           int64  `gorm:"primary_key;auto_increment" json:"-"`
	FileName     string `gorm:"not null;comment:文件名" json:"fileName"`
	FileHash     string `gorm:"not null;unique;" json:"fileHash"`
	Description  string `gorm:"type:text;not null;" json:"description"`
	CoverImg     string `json:"coverImg"`
	Length       string `json:"length"`
	Views        int64  `json:"views"`
	Label        string `json:"label"`
	Size         int64  `gorm:"not null" json:"size"`
	FileStatus   string `json:"fileStatus"`
	Creator      string `gorm:"not null" json:"creator"`
	Owner        string `gorm:"not null" json:"owner"`
	NftToken     string `gorm:"unique;default:NULL" json:"nftToken"`
	Price        string `gorm:"default:NULL" json:"price"`
	NftStatus    string `json:"nftStatus"`
	Chain        string `json:"chain"`
	ContractAddr string `gorm:"default:NULL" json:"contractAddr"`
}

func (t FileStatus) String() string {
	switch t {
	case UPLOAD:
		return "Upload"
	case SCHEDULE:
		return "Schedule"
	case STORAGE:
		return "Storage"
	case DELETE:
		return "Delete"
	}
	return "unknow"
}

func (t NftStatus) String() string {
	switch t {
	case CREATE:
		return "Create"
	case MINT:
		return "Mint"
	case LIST:
		return "List"
	case MELT:
		return "Melt"
	}
	return "Unknow"
}

func (t *VideoMetadata) Create(db *gorm.DB) error {
	return db.Create(t).Error
}

func (t *VideoMetadata) Get(db *gorm.DB) (res []VideoMetadata, err error) {
	tx := db.Where(t).Find(&res)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return
}

func (t *VideoMetadata) IsExist(db *gorm.DB) (bool, error) {
	var count int64
	tx := db.Model(t).Where(t).Count(&count)
	return count >= 1, tx.Error
}

func (t *VideoMetadata) Update(db *gorm.DB) error {
	tx := db.Updates(t)
	return tx.Error
}

func (t *VideoMetadata) Delete(db *gorm.DB) error {
	if t.FileHash == "" {
		return errors.New("empty key")
	}
	tx := db.Delete(t)
	return tx.Error
}

func (t *VideoMetadata) GetCount(db *gorm.DB) (int64, error) {
	var count int64
	tx := db.Model(t).Count(&count)
	return count, tx.Error
}
