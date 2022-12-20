package model

import (
	"gorm.io/gorm"
)

type ActivityState int32

const (
	SUCCESS ActivityState = iota
	FAILED
	LISTENING
	WITHDRAW
)

type EventType int32

const (
	ACT_CREATE EventType = iota //nft metadata create
	ACT_MINT                    //nft mint
	ACT_TS                      //nft transfer
	ACT_TX                      //nft transaction
	ACT_MELT                    //nft melt
	ACT_ALT                     //status alter(price change or list and     )
	ACT_FPG                     //file storage progress
)

type Activity struct {
	Id        int64  `gorm:"primary_key;auto_increment" json:"-"`
	EventType string `json:"eventType"`
	Creator   string `json:"creator"`
	Source    string `json:"from"`
	Target    string `json:"to"`
	FileHash  string `gorm:"not null" json:"fileHash"`
	NftToken  string `json:"nftToken,omitempty"`
	Price     string `json:"price,omitempty"`
	State     string `json:"state"`
	TxHash    string `json:"txhash,omitempty"`
	Gas       string `json:"gas,omitempty"`
	StartDate string `gorm:"not null" json:"-"`
	EndDate   string `gorm:"not null" json:"date"`
}

func (t ActivityState) String() string {
	switch t {
	case SUCCESS:
		return "success"
	case FAILED:
		return "failed"
	case LISTENING:
		return "listening"
	case WITHDRAW:
		return "withdraw"
	}
	return "unknow"
}

func (t EventType) String() string {
	switch t {
	case ACT_CREATE:
		return "create"
	case ACT_MINT:
		return "mint"
	case ACT_TS:
		return "ts"
	case ACT_TX:
		return "tx"
	case ACT_MELT:
		return "melt"
	case ACT_ALT:
		return "alt"
	case ACT_FPG:
		return "fpg"
	}
	return "unknow"
}

func (t *Activity) Create(db *gorm.DB) error {
	return db.Create(t).Error
}

func (t *Activity) Get(db *gorm.DB) (res []Activity, err error) {
	tx := db.Where(t).Find(&res)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return
}

func (t *Activity) GetLatest(db *gorm.DB) (res Activity, err error) {
	tx := db.Where(t).Last(&res)
	return res, tx.Error
}

func (t *Activity) IsExist(db *gorm.DB) (bool, error) {
	var count int64
	tx := db.Model(t).Where(t).Count(&count)
	return count >= 1, tx.Error
}

func (t *Activity) Update(db *gorm.DB) error {
	tx := db.Updates(t)
	return tx.Error
}

func QueryNftEvents(db *gorm.DB, query string, args ...any) (res []Activity, err error) {
	tx := db.Where(query, args...).Find(&res)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return
}
