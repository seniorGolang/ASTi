package models

// FieldInfo представляет поле структуры
type FieldInfo struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Tags        map[string]string `json:"tags,omitempty"`
	Annotations Annotations       `json:"annotations,omitempty"`
	Position    Position          `json:"position"`
	Embedded    bool              `json:"embedded,omitempty"`
	Pointer     bool              `json:"pointer,omitempty"`
	Slice       bool              `json:"slice,omitempty"`
	Map         bool              `json:"map,omitempty"`
	Channel     bool              `json:"channel,omitempty"`
	Generic     bool              `json:"generic,omitempty"`
	Array       bool              `json:"array,omitempty"`
	ArrayLen    int               `json:"arrayLen,omitempty"`
} 