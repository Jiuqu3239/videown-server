package dto

import (
	"encoding/json"
	"fmt"
	"strings"
	"videown-server/internal/model"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Searcher struct {
	Key       string   `json:"key"`
	Type      string   `json:"type"`
	PageIndex int      `json:"pageindex"`
	PageSize  int      `json:"pagesize"`
	Field     []string `json:"field"`
	OnlyTotal string   `json:"onlytotal"`
}

type FilterItem struct {
	Column string `json:"column"`
	Sign   string `json:"sign"`
	Values []any  `json:"values"`
}

type SorterItem struct {
	Column string `json:"column"`
	Type   string `json:"type"`
}

type Querier struct {
	PageIndex  int          `json:"pageindex"`
	PageSize   int          `json:"pagesize"`
	Field      []string     `json:"field"`
	FilterType string       `json:"filterType"`
	Filter     []FilterItem `json:"filter"`
	Sorter     []SorterItem `json:"sorter"`
	Groupby    string       `json:"groupby"`
	OnlyTotal  bool         `json:"onlytotal"`
}

type Responder struct {
	Header []string `json:"header"`
	List   any      `json:"list"`
	Total  int64    `json:"total"`
}

func (t *Querier) GetPointer(key string) any {
	switch key {
	case "pageindex":
		return &t.PageIndex
	case "pagesize":
		return &t.PageSize
	case "field":
		return &t.Field
	case "filterType":
		return &t.FilterType
	case "filter":
		return &t.Filter
	case "sorter":
		return &t.Sorter
	case "groupby":
		return &t.Groupby
	case "onlytotal":
		return &t.OnlyTotal
	}
	return nil
}

func UniversalLoader(header []string, data any, total int64) (Responder, error) {
	var responder Responder
	if len(header) == 0 {
		responder.List = data
		responder.Total = total
		return responder, nil
	}
	var mapDatas []map[string]any
	bytes, err := json.Marshal(data)
	if err != nil {
		return responder, errors.Wrap(err, "load data to responder error")
	}
	err = json.Unmarshal(bytes, &mapDatas)
	if err != nil {
		return responder, errors.Wrap(err, "load data to responder error")
	}
	responder.Total = total
	responder.Header = header
	mapHeader := make(map[string]struct{}, len(header))
	for _, h := range header {
		mapHeader[h] = struct{}{}
	}
	for _, item := range mapDatas {
		for key := range item {
			if _, ok := mapHeader[key]; !ok {
				delete(item, key)
			}
		}
	}
	responder.List = mapDatas
	return responder, nil
}

func UniversalQuery(db *gorm.DB, querier Querier, dest any, count *int64) error {
	if db == nil {
		return errors.New("empty DB pointer")
	}
	for _, item := range querier.Filter {
		cond := getCondStr(item)
		if strings.ToLower(querier.FilterType) == "and" {
			db = db.Where(cond, item.Values...)
		} else {
			db = db.Or(cond, item.Values...)
		}
		if db.Error != nil {
			return db.Error
		}
	}
	for _, item := range querier.Sorter {
		if item.Column == "" {
			continue
		}
		if item.Type == "" {
			item.Type = "asc"
		}
		if strings.ToLower(item.Type) != "asc" && strings.ToLower(item.Type) != "desc" {
			continue
		}
		db = db.Order(item.Column + " " + item.Type)
	}
	if querier.Groupby != "" {
		db = db.Group(querier.Groupby)
	}
	if querier.OnlyTotal {
		return db.Count(count).Error
	}
	if querier.PageIndex <= 0 {
		querier.PageIndex = 1
	}
	if querier.PageSize <= 0 {
		querier.PageSize = 20
	}
	db = db.Limit(querier.PageSize).Offset((querier.PageIndex - 1) * querier.PageSize)
	return db.Find(dest).Count(count).Error
}

func getCondStr(filter FilterItem) string {
	cond := ""
	switch filter.Sign {
	case "=", "!=", "<", ">", "<=", ">=", "like", "not like":
		cond = fmt.Sprintf(" %s %s ? ", filter.Column, filter.Sign)
	case "between":
		cond = fmt.Sprintf(" %s %s ? and ? ", filter.Column, filter.Sign)
	case "in", "not in":
		cond = fmt.Sprintf(" %s %s (", filter.Column, filter.Sign)
		for i := 0; i < len(filter.Values)-1; i++ {
			cond += " ? ,"
		}
		cond += " ? )"
	}
	return cond
}

func EventQuery(db *gorm.DB, querier Querier, dest any, count *int64) error {
	var accounts []string
	for i, v := range querier.Filter {
		if strings.ToLower(v.Column) == "creator" {
			acc, ok := v.Values[0].(string)
			if !ok {
				return errors.New("value type error")
			}
			accounts = append(accounts, acc)
			querier.Filter = append(querier.Filter[:i], querier.Filter[i+1:]...)
		}
	}
	var sdb = db.Model(&model.Activity{})
	for _, v := range accounts {
		sdb = sdb.Or("creator = ?", v).Or("source = ?", v).Or("target = ?", v)
	}
	return UniversalQuery(db.Table("(?) as act", sdb), querier, dest, count)
}
