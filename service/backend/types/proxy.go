package types

type ResizeTermRequest struct {
	Columns int `form:"cols" json:"cols" xml:"cols" binding:"required"`
	Rows    int `form:"rows" json:"rows" xml:"rows" binding:"required"`
}
