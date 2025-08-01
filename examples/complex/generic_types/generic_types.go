// @asti name="GenericTypes" version=2.1
package generictypes

import (
	"context"
	"time"
)

// @asti type=GenericContainer validation=strict
type GenericContainer[T any] struct {
	// @asti field=Data required=true
	Data T `json:"data"`

	// @asti field=CreatedAt validation=timestamp
	CreatedAt time.Time `json:"created_at"`

	// @asti field=Metadata validation=optional
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// @asti type=ComplexGeneric validation=strict
type ComplexGeneric[K comparable, V any, T ~[]V] struct {
	// @asti field=Key required=true
	Key K `json:"key"`

	// @asti field=Value required=true
	Value V `json:"value"`

	// @asti field=Collection validation=required
	Collection T `json:"collection"`

	// @asti field=Map validation=optional
	Map map[K]V `json:"map,omitempty"`

	// @asti field=Pointer validation=optional
	Pointer *V `json:"pointer,omitempty"`
}

// @asti type=GenericNode validation=strict
type GenericNode[T any] struct {
	// @asti field=ID required=true
	ID string `json:"id"`

	// @asti field=Data required=true
	Data T `json:"data"`

	// Дженерик слайс
	Children []*GenericNode[T] `json:"children,omitempty"`

	// Дженерик карта
	Attributes map[string]T `json:"attributes,omitempty"`

	// Дженерик указатель
	Parent *GenericNode[T] `json:"parent,omitempty"`
}

// @asti type=MultiGeneric validation=strict
type MultiGeneric[K comparable, V any, T ~[]K, U ~map[K]V] struct {
	// @asti field=Keys required=true
	Keys T `json:"keys"`

	// @asti field=Values required=true
	Values U `json:"values"`

	// @asti field=Combined validation=optional
	Combined map[K][]V `json:"combined,omitempty"`
}

// @asti name="GenericTypeService" timeout=90
type GenericTypeService interface {
	// @asti method=CreateContainer retry=3
	CreateContainer(ctx context.Context, data interface{}) (container *GenericContainer[interface{}], err error)

	// @asti method=ProcessComplexGeneric timeout=45
	ProcessComplexGeneric(ctx context.Context, data ComplexGeneric[string, interface{}, []interface{}]) (result ComplexGeneric[string, interface{}, []interface{}], err error)

	// @asti method=BuildGenericTree retry=5
	BuildGenericTree(ctx context.Context, nodes []*GenericNode[interface{}]) (root *GenericNode[interface{}], err error)

	// @asti method=ProcessMultiGeneric timeout=60
	ProcessMultiGeneric(ctx context.Context, data MultiGeneric[string, interface{}, []string, map[string]interface{}]) (result MultiGeneric[string, interface{}, []string, map[string]interface{}], err error)

	// @asti method=GenericMapOperation retry=2
	GenericMapOperation(ctx context.Context, input map[string]interface{}) (result map[string][]interface{}, err error)

	// @asti method=GenericSliceOperation timeout=30
	GenericSliceOperation(ctx context.Context, input []interface{}) (result []*interface{}, err error)
}
