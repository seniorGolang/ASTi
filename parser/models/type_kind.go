package models

// TypeKind представляет тип Go типа
type TypeKind int

const (
	TypeStruct TypeKind = iota
	TypeInterface
	TypeEnum
	TypeAlias
	TypeGeneric
	TypeFunction
	TypeChannel
	TypeMap
	TypeSlice
	TypeArray
	TypePointer
	TypeBasic
)

// String возвращает строковое представление TypeKind
func (k TypeKind) String() (str string) {
	switch k {
	case TypeStruct:
		str = "struct"
	case TypeInterface:
		str = "interface"
	case TypeEnum:
		str = "enum"
	case TypeAlias:
		str = "alias"
	case TypeGeneric:
		str = "generic"
	case TypeFunction:
		str = "function"
	case TypeChannel:
		str = "channel"
	case TypeMap:
		str = "map"
	case TypeSlice:
		str = "slice"
	case TypeArray:
		str = "array"
	case TypePointer:
		str = "pointer"
	case TypeBasic:
		str = "basic"
	default:
		str = "unknown"
	}
	return
} 