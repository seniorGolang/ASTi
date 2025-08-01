package models

type Position struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}
