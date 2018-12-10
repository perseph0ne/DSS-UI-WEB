package model

type Document struct {
	ID string `json:"id" form:"id"`
	Name string `json: "name" form:"name"`
	Size int64 `json:"size" form:"size"`
	File [] byte `json:"file" form:"file"`
}
