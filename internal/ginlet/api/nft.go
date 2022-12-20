package api

import (
	"errors"
	"videown-server/internal/dto"
	"videown-server/internal/ginlet/resp"
	"videown-server/internal/service/nft"

	"github.com/gin-gonic/gin"
)

type NftAPI struct{}

func NewNftAPI() NftAPI {
	return NftAPI{}
}

func (v NftAPI) CreateVideoMetadata(c *gin.Context) {
	var req dto.CreateReq
	err := c.BindJSON(&req)
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 500, "bind json data error"))
		return
	}
	result, err := nft.CreateVideoMetadata(req)
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 500, "create video metadata service error"))
		return
	}
	resp.Ok(c, result)
}

func (n NftAPI) QueryActivities(c *gin.Context) {
	var querier dto.Querier
	err := c.BindJSON(&querier)
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 500, "bind json data error"))
		return
	}
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 400, "invalid form data"))
		return
	}
	activities, err := nft.QueryActivitiesByQuerier(querier)
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 500, "query activities service error"))
		return
	}
	resp.Ok(c, activities)
}

func (n NftAPI) MintNFT(c *gin.Context) {
	act, ok := c.Params.Get("act")
	if !ok {
		resp.Error(c, resp.NewErrorWraper(errors.New("empty act"), 400, ""))
	}
	var req dto.NftReq
	err := c.BindJSON(&req)
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 500, "bind json data error"))
		return
	}
	if act == "update" {
		data := nft.TxData{IsSent: true, Data: req.TxHash}
		res, err := nft.UpdateForMint(req.FileHash, data)
		if err != nil {
			resp.Error(c, resp.NewErrorWraper(err, 500, "service error"))
			return
		}
		resp.Ok(c, res)
		return
	}
	if act == "send" {
		data := nft.TxData{IsSent: false, Data: req.Signtx}
		res, err := nft.UpdateForMint(req.FileHash, data)
		if err != nil {
			resp.Error(c, resp.NewErrorWraper(err, 500, "service error"))
			return
		}
		resp.Ok(c, res)
		return
	}
}

func (n NftAPI) BuyNFT(c *gin.Context) {
	act, ok := c.Params.Get("act")
	if !ok {
		resp.Error(c, resp.NewErrorWraper(errors.New("empty act"), 400, ""))
	}
	var req dto.NftReq
	err := c.BindJSON(&req)
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 500, "bind json data error"))
		return
	}
	// if err := utils.VerityAddress(req.To, utils.CessPrefix); err != nil {
	// 	resp.Error(c, resp.NewErrorWraper(err, 400, "invalid address"))
	// 	return
	// }
	if act == "update" {
		data := nft.TxData{IsSent: true, Data: req.TxHash}
		res, err := nft.UpdateForPurchase(req.FileHash, req.To, data)
		if err != nil {
			resp.Error(c, resp.NewErrorWraper(err, 400, "service error"))
			return
		}
		resp.Ok(c, res)
		return
	}
	if act == "send" {
		data := nft.TxData{IsSent: false, Data: req.Signtx}
		res, err := nft.UpdateForPurchase(req.FileHash, req.To, data)
		if err != nil {
			resp.Error(c, resp.NewErrorWraper(err, 400, "service error"))
			return
		}
		resp.Ok(c, res)
		return
	}
}

func (n NftAPI) TransferNFT(c *gin.Context) {
	act, ok := c.Params.Get("act")
	if !ok {
		resp.Error(c, resp.NewErrorWraper(errors.New("empty act"), 400, ""))
	}
	var req dto.NftReq
	err := c.BindJSON(&req)
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 500, "bind json data error"))
		return
	}
	// if err := utils.VerityAddress(req.To, utils.CessPrefix); err != nil {
	// 	resp.Error(c, resp.NewErrorWraper(err, 400, "invalid address"))
	// 	return
	// }
	if act == "update" {
		txhash := c.PostForm("txhash")
		data := nft.TxData{IsSent: true, Data: txhash}
		res, err := nft.UpdateForTransfer(req.FileHash, req.To, data)
		if err != nil {
			resp.Error(c, resp.NewErrorWraper(err, 400, "service error"))
			return
		}
		resp.Ok(c, res)
		return
	}
	if act == "send" {
		data := nft.TxData{IsSent: false, Data: req.Signtx}
		res, err := nft.UpdateForTransfer(req.FileHash, req.To, data)
		if err != nil {
			resp.Error(c, resp.NewErrorWraper(err, 400, "service error"))
			return
		}
		resp.Ok(c, res)
		return
	}
}

func (n NftAPI) ChangeSellingStatus(c *gin.Context) {
	act, ok := c.Params.Get("act")
	if !ok {
		resp.Error(c, resp.NewErrorWraper(errors.New("empty act"), 400, ""))
	}
	var req dto.NftReq
	err := c.BindJSON(&req)
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 500, "bind json data error"))
		return
	}
	if act == "update" {
		data := nft.TxData{IsSent: true, Data: req.TxHash}
		res, err := nft.ChangeStatus(req.FileHash, req.Status, req.Price, data)
		if err != nil {
			resp.Error(c, resp.NewErrorWraper(err, 400, "service error"))
			return
		}
		resp.Ok(c, res)
		return
	}
	if act == "send" {
		data := nft.TxData{IsSent: false, Data: req.Signtx}
		res, err := nft.ChangeStatus(req.FileHash, req.Status, req.Price, data)
		if err != nil {
			resp.Error(c, resp.NewErrorWraper(err, 400, "service error"))
			return
		}
		resp.Ok(c, res)
		return
	}
}

func (n NftAPI) ChangeSellingPrice(c *gin.Context) {
	act, ok := c.Params.Get("act")
	if !ok {
		resp.Error(c, resp.NewErrorWraper(errors.New("empty act"), 400, ""))
	}
	var req dto.NftReq
	err := c.BindJSON(&req)
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 500, "bind json data error"))
		return
	}
	if act == "update" {
		data := nft.TxData{IsSent: true, Data: req.TxHash}
		res, err := nft.ChangePrice(req.FileHash, req.Price, data)
		if err != nil {
			resp.Error(c, resp.NewErrorWraper(err, 400, "service error"))
			return
		}
		resp.Ok(c, res)
		return
	}
	if act == "send" {
		data := nft.TxData{IsSent: false, Data: req.Signtx}
		res, err := nft.ChangePrice(req.FileHash, req.Price, data)
		if err != nil {
			resp.Error(c, resp.NewErrorWraper(err, 400, "service error"))
			return
		}
		resp.Ok(c, res)
		return
	}
}
