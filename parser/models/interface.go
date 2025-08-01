package models

// Interface представляет интерфейс Go
type Interface struct {
	Name        string      `json:"name"`
	ID          string      `json:"id"`
	Package     string      `json:"package"`
	Import      string      `json:"import,omitempty"`
	Methods     []Method    `json:"methods"`
	Annotations Annotations `json:"annotations,omitempty"`
	Position    Position    `json:"position"`
}
