package req

type IDReq struct {
	ID uint64 `json:"id" yaml:"id" form:"id" binding:"required"`
}

type UUIDReq struct {
	ID string `json:"id" yaml:"id" form:"id" binding:"required"`
}

type NameReq struct {
	Name string `form:"name" uri:"name" json:"name" binding:"required"`
}

type PaginationReq struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
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

type LikeReq struct {
	Keyword string `json:"keyword" form:"keyword"`
}

type ConfirmReq struct {
	Confirm bool `json:"confirm" form:"confirm" binding:"required"`
}

type Captcha struct {
	UUID string `json:"uuid" form:"uuid" binding:"required"`
	Code string `json:"code" form:"code" binding:"required"`
}

type AuthLoginReq struct {
	Captcha
	Account  string `json:"account" form:"account" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type AuthRegisterReq struct {
	Captcha
	Account  string `json:"account" form:"account" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Confirm  string `json:"confirm" form:"confirm" binding:"required"`
}

type AuthPasswordReq struct {
	OldPassword string `json:"old_password" form:"old_password" binding:"required"`
	NewPassword string `json:"new_password" form:"new_password" binding:"required"`
	Confirm     string `json:"confirm" form:"confirm" binding:"required"`
}
