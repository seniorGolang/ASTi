// @asti name="NestedTypes" version=1.8
package nestedtypes

import (
	"context"
	"time"
)

// @asti type=DeepNestedStruct validation=strict
type DeepNestedStruct struct {
	// @asti field=Level1 required=true
	Level1 struct {
		// @asti field=Level2 required=true
		Level2 struct {
			// @asti field=Level3 required=true
			Level3 struct {
				// @asti field=Level4 required=true
				Level4 struct {
					// @asti field=Level5 required=true
					Level5 struct {
						// @asti field=Data required=true
						Data map[string]interface{} `json:"data"`

						// @asti field=Timestamp validation=timestamp
						Timestamp time.Time `json:"timestamp"`
					} `json:"level5"`
				} `json:"level4"`
			} `json:"level3"`
		} `json:"level2"`
	} `json:"level1"`
}

// @asti type=ComplexSliceStruct validation=strict
type ComplexSliceStruct struct {
	// @asti field=Matrix required=true
	Matrix [][][][]int `json:"matrix"`

	// @asti field=SliceMap required=true
	SliceMap []map[string][]interface{} `json:"slice_map"`

	// @asti field=MapSlice required=true
	MapSlice map[string][][]string `json:"map_slice"`

	// @asti field=DeepSlice required=true
	DeepSlice [][][][][]*DeepNestedStruct `json:"deep_slice"`
}

// @asti type=ComplexMapStruct validation=strict
type ComplexMapStruct struct {
	// @asti field=DeepMap required=true
	DeepMap map[string]map[string]map[string]interface{} `json:"deep_map"`

	// @asti field=MapSliceMap required=true
	MapSliceMap map[string][]map[string]interface{} `json:"map_slice_map"`

	// @asti field=SliceMapSlice required=true
	SliceMapSlice []map[string][]map[string]string `json:"slice_map_slice"`
}

// @asti type=InterfaceEmbedding validation=strict
type InterfaceEmbedding struct {
	// @asti field=Base required=true
	Base interface {
		// @asti method=GetID
		GetID() string

		// @asti method=GetName
		GetName() string
	} `json:"base"`

	// @asti field=Extended required=true
	Extended interface {
		// @asti method=GetID
		GetID() string

		// @asti method=GetName
		GetName() string

		// @asti method=GetMetadata
		GetMetadata() map[string]interface{}
	} `json:"extended"`
}

// @asti type=FunctionType validation=strict
type FunctionType struct {
	// @asti field=Handler required=true
	Handler func(context.Context, string) (result interface{}, err error) `json:"handler"`

	// @asti field=Processor required=true
	Processor func([]interface{}) func(map[string]interface{}) error `json:"processor"`

	// @asti field=Validator required=true
	Validator func(interface{}) (bool, error) `json:"validator"`
}

// @asti type=ChannelStruct validation=strict
type ChannelStruct struct {
	// @asti field=DataChan required=true
	DataChan chan interface{} `json:"data_chan"`

	// @asti field=ResultChan required=true
	ResultChan chan<- string `json:"result_chan"`

	// @asti field=ErrorChan required=true
	ErrorChan <-chan error `json:"error_chan"`

	// @asti field=BufferedChan required=true
	BufferedChan chan []interface{} `json:"buffered_chan"`
}

// @asti name="NestedTypeService" timeout=150
type NestedTypeService interface {
	// @asti method=ProcessDeepNested retry=3
	ProcessDeepNested(ctx context.Context, data DeepNestedStruct) (result DeepNestedStruct, err error)

	// @asti method=ProcessComplexSlice timeout=60
	ProcessComplexSlice(ctx context.Context, data ComplexSliceStruct) (result ComplexSliceStruct, err error)

	// @asti method=ProcessComplexMap retry=5
	ProcessComplexMap(ctx context.Context, data ComplexMapStruct) (result ComplexMapStruct, err error)

	// @asti method=ProcessInterfaceEmbedding timeout=45
	ProcessInterfaceEmbedding(ctx context.Context, data InterfaceEmbedding) (result InterfaceEmbedding, err error)

	// @asti method=ProcessFunctionType retry=2
	ProcessFunctionType(ctx context.Context, data FunctionType) (result FunctionType, err error)

	// @asti method=ProcessChannelStruct timeout=30
	ProcessChannelStruct(ctx context.Context, data ChannelStruct) (result ChannelStruct, err error)
}
