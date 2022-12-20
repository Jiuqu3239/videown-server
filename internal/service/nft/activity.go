package nft

import (
	"strings"
	"videown-server/global"
	"videown-server/internal/dto"
	"videown-server/internal/model"

	"github.com/pkg/errors"
)

func QueryActivitiesByQuerier(querier dto.Querier) (dto.Responder, error) {
	var (
		responder dto.Responder
		list      []model.Activity
		res       []dto.EventResp
		count     int64
	)
	err := dto.EventQuery(global.GormDb, querier, &list, &count)
	if err != nil {
		return responder, errors.Wrap(err, "query videos error")
	}
	if count != int64(len(list)) {
		count = int64(len(list))
	}
	video := &model.VideoMetadata{}
	for _, v := range list {
		if v.EventType == model.ACT_FPG.String() {
			continue
		}
		r := dto.EventResp{
			EventType: v.EventType,
			FileHash:  v.FileHash,
			Price:     v.Price,
			State:     v.State,
			Date:      v.EndDate,
			From:      v.Source,
			To:        v.Target,
		}
		if v.EventType != model.ACT_TS.String() &&
			v.EventType != model.ACT_TX.String() {
			r.From = model.NULL
			r.To = v.Creator
		}
		r.EventType = getEventType(v)
		video.FileHash = v.FileHash
		videos, err := video.Get(global.GormDb)
		if err != nil || len(videos) <= 0 {
			continue
		}
		r.FileName = videos[0].FileName
		r.CoverImg = videos[0].CoverImg
		res = append(res, r)
	}
	//response
	responder, err = dto.UniversalLoader(querier.Field, res, count)
	if err != nil {
		return responder, errors.Wrap(err, "query videos error")
	}
	return responder, nil
}

func getEventType(act model.Activity) string {
	eventType := ""
	switch act.EventType {
	case model.ACT_ALT.String():
		if strings.ToLower(act.Target) == "list" {
			eventType = "list"
		} else if strings.ToLower(act.Source) == "list" {
			eventType = "unlist"
		}
	case model.ACT_TS.String():
		eventType = "transfer"
	case model.ACT_TX.String():
		eventType = "transaction"
	default:
		eventType = act.EventType
	}
	return eventType
}
