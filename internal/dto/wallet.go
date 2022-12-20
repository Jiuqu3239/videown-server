package dto

import (
	"math/big"
	"time"
	"videown-server/pkg/paging"
)

type WalletState uint8

const (
	WalletState_Ok = WalletState(1) + iota
	WalletState_Frozen
)

type AppealAuditResult uint8

const (
	AAR_Pending = AppealAuditResult(0) + iota
	AAR_Passed
	AAR_Refuse
)

type (
	WalletAttachReq struct {
		WalletSignWrapper
	}

	WalletsQueryReq struct {
		paging.PageRequest
		Address  string `json:"address" form:"address"`
		OnlyZero int32  `json:"onlyZero" form:"onlyZero"`
	}

	AppealReq struct {
		WalletSignWrapper
		Cause string `binding:"required,min=1,max=120" json:"cause" form:"cause"`
	}

	AppealsQueryReq struct {
		paging.PageRequest
		Address string `json:"address" form:"address"`
	}

	WalletFreezeReq struct {
		OnlyWalletAddress
		Reasons string `binding:"required,min=1,max=120" json:"reasons" form:"reasons"`
	}

	WalletFrozenInfo struct {
		FrozenAt      time.Time  `json:"frozenAt"`
		FrozenReasons string     `json:"frozenReasons"`
		AppealAt      *time.Time `json:"appealAt,omitempty"`
		AppealCause   *string    `json:"appealCause,omitempty"`
		AuditResult   *uint8     `json:"auditResult,omitempty"`
	}

	WalletStateResp struct {
		WalletState WalletState       `json:"walletState"`
		FrozenInfo  *WalletFrozenInfo `json:"frozenInfo,omitempty"`
	}

	AppealAuditReq struct {
		AppealId uint64 `binding:"required" json:"appealId"`
		Passed   *bool  `binding:"required" json:"passed"`
	}

	TransferIntoReq struct {
		OnlyWalletAddress
		Amount          *big.Int `binding:"required" json:"amount"`
		ConfirmPassword string   `binding:"required" json:"confirmPassword"`
	}

	MakeEthSignReq struct {
		PrivateKey string `binding:"required" json:"prk" form:"prk"`
		Message    string `binding:"required" json:"msg" form:"msg"`
	}
)

type (
	WalletAppealWithFrozenView struct {
		Id            uint64     `json:"id"`
		FrozenId      uint64     `json:"frozenId"`
		WalletAddress string     `json:"walletAddress"`
		FrozenAt      time.Time  `json:"frozenAt"`
		FrozenReasons string     `json:"frozenReasons"`
		AppealAt      time.Time  `json:"appealAt"`
		AppealCause   string     `json:"appealCause"`
		AuditResult   uint8      `json:"auditResult"`
		AuditAt       *time.Time `json:"auditAt"`
	}
)
