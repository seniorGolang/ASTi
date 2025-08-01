package models

// TypeInfo представляет информацию о типе
type TypeInfo struct {
	Name        string         `json:"name"`
	Package     string         `json:"package"`
	Import      string         `json:"import,omitempty"`
	Kind        TypeKind       `json:"kind"`
	Fields      []FieldInfo    `json:"fields,omitempty"`
	Methods     []MethodInfo   `json:"methods,omitempty"`
	Annotations Annotations    `json:"annotations,omitempty"`
	Position    Position       `json:"position"`
	Generic     *GenericInfo   `json:"generic,omitempty"`
	Underlying  string         `json:"underlying,omitempty"`
	Constants   []ConstantInfo `json:"constants,omitempty"`

	Pointer     bool `json:"pointer,omitempty"`
	Slice       bool `json:"slice,omitempty"`
	Map         bool `json:"map,omitempty"`
	Channel     bool `json:"channel,omitempty"`
	GenericType bool `json:"genericType,omitempty"`
	Array       bool `json:"array,omitempty"`
	ArrayLen    int  `json:"arrayLen,omitempty"`
	Interface   bool `json:"interface,omitempty"`
	Function    bool `json:"function,omitempty"`
} 