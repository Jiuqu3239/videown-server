package nft

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
	"videown-server/global"
	"videown-server/internal/dto"
	"videown-server/internal/model"
	"videown-server/pkg/utils"

	"github.com/pkg/errors"
)

var threshold = 40

const (
	FILE_STATUS_PENDING = "pending"
	FILE_STATUS_ACTIVE  = "active"
	FILE_STATUS_CANCEL  = "cancel"
)

func QueryVideosCount() (int64, error) {
	var query model.VideoMetadata
	return query.GetCount(global.GormDb)
}

func QueryVideosByQuerier(querier dto.Querier) (dto.Responder, error) {
	var (
		responder dto.Responder
		list      []model.VideoMetadata
		count     int64
	)
	err := dto.UniversalQuery(global.GormDb, querier, &list, &count)
	if err != nil {
		return responder, errors.Wrap(err, "query videos error")
	}
	if count != int64(len(list)) {
		count = int64(len(list))
	}
	responder, err = dto.UniversalLoader(querier.Field, list, count)
	if err != nil {
		return responder, errors.Wrap(err, "query videos error")
	}
	return responder, nil
}

func QueryRelatedVideos(searcher dto.Searcher) (dto.Responder, error) {
	var (
		v     model.VideoMetadata
		res   []model.VideoMetadata
		hits  []int
		text  string
		count int
	)
	if searcher.PageSize <= 0 {
		searcher.PageSize = len(res)
	}
	if searcher.PageIndex < 1 {
		searcher.PageIndex = 1
	}
	videos, err := v.Get(global.GormDb)
	if err != nil {
		return dto.Responder{}, err
	}
	if len(videos) <= 0 {
		return dto.Responder{}, errors.New("empty table")
	}
	matcher := utils.NewStringMatcher(strings.Split(searcher.Key, " "))
	for _, v := range videos {
		switch searcher.Type {
		case "name":
			text = v.FileName
		case "description":
			text = v.Description
		default:
			text = v.FileName
		}
		hits = matcher.Match([]byte(text))
		if len(hits) >= len(strings.Split(text, " "))*threshold/100 {
			count++
			if (searcher.PageIndex-1)*searcher.PageSize < count && count < searcher.PageIndex*searcher.PageSize {
				res = append(res, v)
			}
		}
	}
	return dto.UniversalLoader(searcher.Field, res, int64(len(res)))
}

func UpdateVideoStatus(filehash, status string) {
	//check metadata is exsited
	video := &model.VideoMetadata{FileHash: filehash}
	res, err := video.Get(global.GormDb)
	if err != nil {
		return
	}
	video = &res[0]
	switch status {
	case FILE_STATUS_PENDING:
		status = model.SCHEDULE.String()
	case FILE_STATUS_ACTIVE:
		status = model.STORAGE.String()
	default:
		status = model.SCHEDULE.String()
	}
	if status == video.FileStatus {
		return
	}
	//create video file event
	videoEvent := model.Activity{
		EventType: model.ACT_FPG.String(),
		Creator:   video.Creator,
		FileHash:  filehash,
		Source:    video.FileStatus,
		Target:    status,
		State:     model.SUCCESS.String(),
		StartDate: time.Now().Local().Format(global.Time_FMT),
		EndDate:   time.Now().Local().Format(global.Time_FMT),
	}
	err = videoEvent.Create(global.GormDb)
	if err != nil {
		return
	}
	//update video metadata
	video.FileStatus = status
	video.Update(global.GormDb)
}

func SaveCoverImg(file *multipart.FileHeader, filename string) error {
	fpath := path.Join(
		global.COVER_IMAGE_PATH,
		filename[:4],
		filename[4:8],
		filename[8:12],
	)
	os.MkdirAll(fpath, os.ModeDir)
	fpath = path.Join(fpath, filename)
	if _, err := os.Stat(fpath); err == nil {
		return nil
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func GetCoverImagePath(filename string) string {
	return path.Join(
		global.COVER_IMAGE_PATH,
		filename[:4],
		filename[4:8],
		filename[8:12],
		filename,
	)
}
