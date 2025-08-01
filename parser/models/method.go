package models

// Method представляет метод интерфейса
type Method struct {
	MethodInfo
	ID           string `json:"id"`
	Serializable bool   `json:"serializable"`
} 