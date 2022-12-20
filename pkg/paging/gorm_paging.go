package paging

import (
	"math"

	"gorm.io/gorm"
)

func QueryByPaging(pr Paging, listResult any, listQ *gorm.DB, cntQ *gorm.DB) (PagingResulter, error) {
	var totalRecords int64 = -1
	if cntQ != nil {
		if err := cntQ.Count(&totalRecords).Error; err != nil {
			return nil, err
		}
	}
	size := pr.PageSize()
	if size > 1000 {
		size = 1000
	}
	totalPages := uint(math.Ceil(float64(totalRecords) / float64(size)))
	var err error
	if totalRecords != -1 && pr.PageNo() <= totalPages {
		offset := (pr.PageNo() - 1) * size
		err = listQ.Offset(int(offset)).Limit(int(size)).Scan(listResult).Error
	}
	return NewPagingResult(pr, uint(totalRecords), totalPages, listResult), err
}
