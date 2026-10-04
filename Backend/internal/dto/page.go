package dto

// PageQuery: ?page=&limit= untuk endpoint daftar.
type PageQuery struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

func (q *PageQuery) Normalize() {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 20
	}
}

func (q PageQuery) Offset() int {
	return (q.Page - 1) * q.Limit
}
