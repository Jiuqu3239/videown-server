package paging

type Paging interface {
	PageNo() uint
	PageSize() uint
}

type PagingResulter interface {
	TotalPages() uint
	TotalRecords() uint
	DataList() any
}

type PageRequest struct {
	Pn   uint `json:"pn" form:"pn"`
	Size uint `json:"size" form:"size"`
}

func (self PageRequest) PageNo() uint {
	if self.Pn == 0 {
		return 1
	}
	return self.Pn
}

func (self PageRequest) PageSize() uint {
	if self.Size == 0 {
		return 10
	} else {
		return self.Size
	}
}

type PageResult struct {
	Pn      uint `json:"pn"`
	Size    uint `json:"size"`
	Pages   uint `json:"tps"`
	Records uint `json:"trs"`
	List    any  `json:"list"`
}

func (self PageResult) TotalPages() uint {
	return self.Pages
}

func (self PageResult) TotalRecords() uint {
	return self.Records
}

func (self PageResult) DataList() any {
	return self.List
}

func NewPagingResult(pr Paging, totalRecords uint, totalPages uint, list any) PagingResulter {
	var r PageResult
	r.Pn = pr.PageNo()
	r.Size = pr.PageSize()
	r.Pages = totalPages
	r.Records = totalRecords
	r.List = list
	return r
}
