package types

type CreateTermProcessRequest struct {
	Columns int `form:"cols" json:"cols" xml:"cols"`
	Rows    int `form:"rows" json:"rows" xml:"rows"`
}

type UpdateTermSizeRequest struct {
	Columns int `form:"cols" json:"cols" xml:"cols" binding:"required"`
	Rows    int `form:"rows" json:"rows" xml:"rows" binding:"required"`
}

type ResizeTermRequest struct {
	Columns int `form:"cols" json:"cols" xml:"cols" binding:"required"`
	Rows    int `form:"rows" json:"rows" xml:"rows" binding:"required"`
}
