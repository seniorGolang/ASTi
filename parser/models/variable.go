package models

// Variable представляет переменную в сигнатуре метода или функции
type Variable struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Variadic bool   `json:"variadic,omitempty"`
	Pointer  bool   `json:"pointer,omitempty"`
	Slice    bool   `json:"slice,omitempty"`
	Map      bool   `json:"map,omitempty"`
	Channel  bool   `json:"channel,omitempty"`
	Generic  bool   `json:"generic,omitempty"`
	Array    bool   `json:"array,omitempty"`
	ArrayLen int    `json:"arrayLen,omitempty"`
} 