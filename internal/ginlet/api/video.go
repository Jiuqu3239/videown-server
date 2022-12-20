package api

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"videown-server/global"
	"videown-server/internal/dto"
	"videown-server/internal/ginlet/resp"
	"videown-server/internal/model"
	"videown-server/internal/service/nft"
	"videown-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type VideoAPI struct{}

func NewVideoAPI() VideoAPI {
	return VideoAPI{}
}

func (v VideoAPI) QueryVideos(c *gin.Context) {
	var querier dto.Querier
	err := c.BindJSON(&querier)
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 400, "Invalid form data"))
		return
	}
	videos, err := nft.QueryVideosByQuerier(querier)
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 500, "Query videos service error"))
		return
	}
	resp.Ok(c, videos)
}

func (v VideoAPI) SearchVideos(c *gin.Context) {
	var searcher dto.Searcher
	err := c.BindJSON(&searcher)
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 500, "bind json data error"))
		return
	}
	if searcher.Type == "" || (searcher.Type != "name" && searcher.Type != "description") {
		searcher.Type = "name"
	}
	res, err := nft.QueryRelatedVideos(searcher)
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 500, "search video service error"))
		return
	}
	resp.Ok(c, res)
}

func (v VideoAPI) UploadVideoCoverImg(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil || !utils.IsImage(path.Ext(fileHeader.Filename)) {
		resp.Error(c, resp.NewErrorWraper(err, 400, "invalid image file"))
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 500, "open image file error"))
		return
	}
	h := sha256.New()
	_, err = io.Copy(h, file)
	file.Close()
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 500, "create image file hash error"))
		return
	}
	hash := h.Sum(nil)
	fileName := hex.EncodeToString(hash) + path.Ext(fileHeader.Filename)
	err = nft.SaveCoverImg(fileHeader, fileName)
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 500, "save image file error"))
		return
	}
	resp.Ok(c, fileName)
}

func (v VideoAPI) DownloadVideoCoverImg(c *gin.Context) {
	filename, ok := c.GetQuery("filename")
	if !ok || len(filename) < 12 {
		resp.ErrorWithHttpStatus(c, errors.New("invalid query params"), 400)
		return
	}
	imgType := path.Ext(filename)
	if imgType != "" {
		imgType = "image/" + strings.ToLower(imgType[1:])
	}
	c.Writer.Header().Add("Content-Type", imgType)
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename)) //attachment
	if _, err := os.Stat(nft.GetCoverImagePath(filename)); err != nil {
		c.File(global.DEFAULT_COVER_IMAGE)
		return
	}
	c.File(nft.GetCoverImagePath(filename))
}

func (v VideoAPI) AddVideoViews(c *gin.Context) {
	filehash, ok := c.GetQuery("filehash")
	if !ok {
		resp.ErrorWithHttpStatus(c, errors.New("invalid query params"), 400)
		return
	}
	video := model.VideoMetadata{FileHash: filehash}
	res, err := video.Get(global.GormDb)
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 500, "query video info error"))
		return
	}
	if len(res) <= 0 {
		resp.ErrorWithHttpStatus(c, errors.New("video metadata not found"), 404)
		return
	}
	video = res[0]
	video.Views++
	err = video.Update(global.GormDb)
	if err != nil {
		resp.Error(c, resp.NewErrorWraper(err, 500, "update views error"))
		return
	}
	resp.Ok(c, video.Views)
}
