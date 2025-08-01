// @asti name="CornerCases" version=1.0
// @asti category=testing model=corner
package cornercases

import (
	"context"
	"time"
)

// @asti type=EmptyStruct validation=strict
type EmptyStruct struct{}

// @asti type=SingleFieldStruct validation=strict
type SingleFieldStruct struct {
	// @asti field=Value required=true
	Value string `json:"value"`
}

// @asti type=PointerStruct validation=strict
type PointerStruct struct {
	// @asti field=StringPtr validation=optional
	StringPtr *string `json:"string_ptr,omitempty"`

	// @asti field=IntPtr validation=optional
	IntPtr *int `json:"int_ptr,omitempty"`

	// @asti field=BoolPtr validation=optional
	BoolPtr *bool `json:"bool_ptr,omitempty"`

	// @asti field=TimePtr validation=optional
	TimePtr *time.Time `json:"time_ptr,omitempty"`

	// @asti field=StructPtr validation=optional
	StructPtr *SingleFieldStruct `json:"struct_ptr,omitempty"`
}

// @asti type=ArrayStruct validation=strict
type ArrayStruct struct {
	// @asti field=FixedArray required=true
	FixedArray [5]int `json:"fixed_array"`

	// @asti field=StringArray required=true
	StringArray [10]string `json:"string_array"`

	// @asti field=StructArray required=true
	StructArray [3]SingleFieldStruct `json:"struct_array"`

	// @asti field=PointerArray required=true
	PointerArray [2]*SingleFieldStruct `json:"pointer_array"`
}

// @asti type=SliceStruct validation=strict
type SliceStruct struct {
	// @asti field=IntSlice required=true
	IntSlice []int `json:"int_slice"`

	// @asti field=StringSlice required=true
	StringSlice []string `json:"string_slice"`

	// @asti field=StructSlice required=true
	StructSlice []SingleFieldStruct `json:"struct_slice"`

	// @asti field=PointerSlice required=true
	PointerSlice []*SingleFieldStruct `json:"pointer_slice"`

	// @asti field=EmptySlice required=true
	EmptySlice []interface{} `json:"empty_slice"`
}

// @asti type=MapStruct validation=strict
type MapStruct struct {
	// @asti field=StringMap required=true
	StringMap map[string]string `json:"string_map"`

	// @asti field=IntMap required=true
	IntMap map[string]int `json:"int_map"`

	// @asti field=StructMap required=true
	StructMap map[string]SingleFieldStruct `json:"struct_map"`

	// @asti field=PointerMap required=true
	PointerMap map[string]*SingleFieldStruct `json:"pointer_map"`

	// @asti field=EmptyMap required=true
	EmptyMap map[string]interface{} `json:"empty_map"`

	// @asti field=NestedMap required=true
	NestedMap map[string]map[string]interface{} `json:"nested_map"`
}

// @asti type=ChannelStruct validation=strict
type ChannelStruct struct {
	// @asti field=StringChan required=true
	StringChan chan string `json:"string_chan"`

	// @asti field=IntChan required=true
	IntChan chan int `json:"int_chan"`

	// @asti field=StructChan required=true
	StructChan chan SingleFieldStruct `json:"struct_chan"`

	// @asti field=PointerChan required=true
	PointerChan chan *SingleFieldStruct `json:"pointer_chan"`

	// @asti field=SendOnlyChan required=true
	SendOnlyChan chan<- string `json:"send_only_chan"`

	// @asti field=ReceiveOnlyChan required=true
	ReceiveOnlyChan <-chan string `json:"receive_only_chan"`

	// @asti field=BufferedChan required=true
	BufferedChan chan interface{} `json:"buffered_chan"`
}

// @asti type=InterfaceStruct validation=strict
type InterfaceStruct struct {
	// @asti field=EmptyInterface required=true
	EmptyInterface interface{} `json:"empty_interface"`

	// @asti field=StringInterface required=true
	StringInterface interface {
		String() string
	} `json:"string_interface"`

	// @asti field=ComplexInterface required=true
	ComplexInterface interface {
		// @asti method=GetID
		GetID() string

		// @asti method=GetName
		GetName() string

		// @asti method=GetMetadata
		GetMetadata() map[string]interface{}

		// @asti method=Validate
		Validate() error
	} `json:"complex_interface"`
}

// @asti type=FunctionStruct validation=strict
type FunctionStruct struct {
	// @asti field=SimpleFunc required=true
	SimpleFunc func() `json:"simple_func"`

	// @asti field=StringFunc required=true
	StringFunc func() string `json:"string_func"`

	// @asti field=ParamFunc required=true
	ParamFunc func(string) error `json:"param_func"`

	// @asti field=ComplexFunc required=true
	ComplexFunc func(context.Context, string, int) (result interface{}, err error) `json:"complex_func"`

	// @asti field=VariadicFunc required=true
	VariadicFunc func(string, ...interface{}) error `json:"variadic_func"`

	// @asti field=ReturnFunc required=true
	ReturnFunc func() func() error `json:"return_func"`
}

// @asti type=EmbeddedStruct validation=strict
type EmbeddedStruct struct {
	// @asti field=SingleFieldStruct embedded=true
	SingleFieldStruct

	// @asti field=PointerStruct embedded=true
	*PointerStruct

	// @asti field=CustomField required=true
	CustomField string `json:"custom_field"`
}

// @asti type=TaggedStruct validation=strict
type TaggedStruct struct {
	// @asti field=JSONField required=true
	JSONField string `json:"json_field" xml:"xml_field" yaml:"yaml_field"`

	// @asti field=ValidateField required=true
	ValidateField string `json:"validate_field" validate:"required,min=1,max=100"`

	// @asti field=DBField required=true
	DBField string `json:"db_field" db:"db_field_name" sql:"sql_field_name"`

	// @asti field=CustomTag required=true
	CustomTag string `json:"custom_tag" custom:"custom_value" another:"another_value"`
}

// @asti name="CornerCaseService" timeout=60
// @asti category=service model=corner
type CornerCaseService interface {
	// @asti method=ProcessEmptyStruct retry=1 timeout=5
	// @asti validation=empty authorization=public
	ProcessEmptyStruct(ctx context.Context, data EmptyStruct) (result EmptyStruct, err error)

	// @asti method=ProcessPointerStruct retry=3 timeout=15
	// @asti validation=pointer authorization=user
	ProcessPointerStruct(ctx context.Context, data *PointerStruct) (result *PointerStruct, err error)

	// @asti method=ProcessArrayStruct timeout=20
	// @asti validation=array authorization=user
	ProcessArrayStruct(ctx context.Context, data ArrayStruct) (result ArrayStruct, err error)

	// @asti method=ProcessSliceStruct retry=2 timeout=25
	// @asti validation=slice authorization=user
	ProcessSliceStruct(ctx context.Context, data SliceStruct) (result SliceStruct, err error)

	// @asti method=ProcessMapStruct timeout=30
	// @asti validation=map authorization=user
	ProcessMapStruct(ctx context.Context, data MapStruct) (result MapStruct, err error)

	// @asti method=ProcessChannelStruct retry=5 timeout=35
	// @asti validation=channel authorization=admin
	ProcessChannelStruct(ctx context.Context, data ChannelStruct) (result ChannelStruct, err error)

	// @asti method=ProcessInterfaceStruct timeout=40
	// @asti validation=interface authorization=user
	ProcessInterfaceStruct(ctx context.Context, data InterfaceStruct) (result InterfaceStruct, err error)

	// @asti method=ProcessFunctionStruct retry=3 timeout=45
	// @asti validation=function authorization=admin
	ProcessFunctionStruct(ctx context.Context, data FunctionStruct) (result FunctionStruct, err error)

	// @asti method=ProcessEmbeddedStruct timeout=50
	// @asti validation=embedded authorization=user
	ProcessEmbeddedStruct(ctx context.Context, data EmbeddedStruct) (result EmbeddedStruct, err error)

	// @asti method=ProcessTaggedStruct retry=2 timeout=55
	// @asti validation=tagged authorization=user
	ProcessTaggedStruct(ctx context.Context, data TaggedStruct) (result TaggedStruct, err error)
}
