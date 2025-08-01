package models

// MethodInfo представляет метод типа
type MethodInfo struct {
	Name        string      `json:"name"`
	Parameters  []Variable  `json:"parameters"`
	Results     []Variable  `json:"results"`
	Annotations Annotations `json:"annotations,omitempty"`
	Position    Position    `json:"position"`
} 