package api

import (
	"videown-server/global"
	"videown-server/internal/dto"
	"videown-server/internal/ginlet/resp"
	"videown-server/internal/service/auth"

	"github.com/gin-gonic/gin"
)

type SysAPI struct{}

func NewSysAPI() SysAPI {
	return SysAPI{}
}

type LoginReq struct {
	Username string `form:"username" binding:"required" json:"username"`
	Password string `form:"password" binding:"required" json:"password"`
}

type ChangePwdReq struct {
	NewPwd string `form:"newPwd" binding:"required,min=8" json:"newPwd"`
	OldPwd string `form:"oldPwd" binding:"required" json:"oldPwd"`
}

func (t SysAPI) UserLogin(c *gin.Context) {
	var f LoginReq
	if err := c.ShouldBind(&f); err != nil {
		resp.Error(c, err)
		return
	}
	ar, err := auth.RequestAuth(f.Username, f.Password, auth.UserVerifyProvider{})
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.Ok(c, ar)
}

func (t SysAPI) WalletLogin(c *gin.Context) {
	var w dto.WalletAttachReq
	if err := c.ShouldBind(&w); err != nil {
		resp.Error(c, err)
		return
	}
	ar, err := auth.RequestWalletAuth(w)
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.Ok(c, ar)
}

func (t SysAPI) AdminLogin(c *gin.Context) {
	var f LoginReq
	if err := c.ShouldBind(&f); err != nil {
		resp.Error(c, err)
		return
	}
	ar, err := auth.RequestAuth(f.Username, f.Password, auth.GetAdminVerifier())
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.Ok(c, ar)
}

func (t SysAPI) ChangePassword(c *gin.Context) {
	var f ChangePwdReq
	if err := c.ShouldBind(&f); err != nil {
		resp.Error(c, err)
		return
	}
	if err := global.Settings.ChangePassword(f.NewPwd, f.OldPwd); err != nil {
		resp.Error(c, err)
		return
	}
	resp.Ok(c, nil)
}
