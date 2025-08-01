package models

// Package представляет пакет Go с его интерфейсами и типами
type Package struct {
	ModuleName  string              `json:"moduleName"`
	PackagePath string              `json:"packagePath"`
	Annotations Annotations         `json:"annotations"`
	Interfaces  []Interface         `json:"interfaces"`
	Types       map[string]TypeInfo `json:"types"`
}
