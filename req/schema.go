package req

type IDReq[T uint64 | string] struct {
	ID T `json:"id" yaml:"id" form:"id" binding:"required"`
}

type NameReq struct {
	Name string `form:"name" json:"name" binding:"required"`
}

type PaginationReq struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type LikeReq struct {
	Keyword string `json:"keyword" form:"keyword"`
}

type Captcha struct {
	UUID string `json:"uuid" form:"uuid" binding:"required"`
	Code string `json:"code" form:"code" binding:"required"`
}

func (p *PaginationReq) FindPage(sizeLimit ...int) (offset int, limit int) {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Size == 0 {
		p.Size = 10
	}
	if len(sizeLimit) > 0 && p.Size > sizeLimit[0] {
		p.Size = sizeLimit[0]
	}
	return (p.Page - 1) * p.Size, p.Size
}

func (p *PaginationReq) StartEnd() (start, end int) {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Size == 0 {
		p.Size = 10
	}
	start = (p.Page - 1) * p.Size
	end = p.Size + start
	return
}
