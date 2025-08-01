package models

// GenericInfo представляет информацию о generic типах
type GenericInfo struct {
	TypeParams  []string          `json:"typeParams"`
	Constraints []string          `json:"constraints,omitempty"`
	Bounds      map[string]string `json:"bounds,omitempty"`
} 