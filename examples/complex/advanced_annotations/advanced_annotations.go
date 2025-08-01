// @asti name="AdvancedAnnotations" version=3.0
// @asti author="Senior Developer" team="Backend"
// @asti tags=complex,advanced,testing
package advancedannotations

import (
	"context"
	"time"
)

// @asti type=AdvancedStruct validation=strict
// @asti category=data model=entity
type AdvancedStruct struct {
	// @asti field=ID required=true validation=uuid
	ID string `json:"id" validate:"required,uuid"`

	// @asti field=Name required=true validation=string maxLength=255
	Name string `json:"name" validate:"required,max=255"`

	// @asti field=Email validation=email required=true
	Email string `json:"email" validate:"required,email"`

	// @asti field=Age validation=number min=0 max=150
	Age int `json:"age" validate:"min=0,max=150"`

	// @asti field=IsActive validation=boolean default=true
	IsActive bool `json:"is_active" validate:"boolean"`

	// @asti field=CreatedAt validation=timestamp auto=true
	CreatedAt time.Time `json:"created_at"`

	// @asti field=UpdatedAt validation=timestamp auto=true
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	// @asti field=Metadata validation=json optional=true
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// @asti field=Tags validation=array maxItems=10
	Tags []string `json:"tags" validate:"max=10"`

	// @asti field=Settings validation=object
	Settings struct {
		// @asti field=Theme validation=string enum=light,dark,auto
		Theme string `json:"theme" validate:"oneof=light dark auto"`

		// @asti field=Language validation=string enum=en,ru,es,fr
		Language string `json:"language" validate:"oneof=en ru es fr"`

		// @asti field=Notifications validation=boolean default=true
		Notifications bool `json:"notifications"`
	} `json:"settings"`
}

// @asti type=ComplexEnum validation=strict
// @asti category=enumeration model=constant
type ComplexEnum string

const (
	// @asti value=StatusActive description="Active status"
	StatusActive ComplexEnum = "active"

	// @asti value=StatusInactive description="Inactive status"
	StatusInactive ComplexEnum = "inactive"

	// @asti value=StatusPending description="Pending status"
	StatusPending ComplexEnum = "pending"

	// @asti value=StatusSuspended description="Suspended status"
	StatusSuspended ComplexEnum = "suspended"
)

// @asti type=ValidationRules validation=strict
// @asti category=validation model=rules
type ValidationRules struct {
	// @asti field=MinLength validation=number min=1 max=1000
	MinLength int `json:"min_length" validate:"min=1,max=1000"`

	// @asti field=MaxLength validation=number min=1 max=10000
	MaxLength int `json:"max_length" validate:"min=1,max=10000"`

	// @asti field=Pattern validation=regex
	Pattern string `json:"pattern" validate:"regexp"`

	// @asti field=AllowedValues validation=array
	AllowedValues []string `json:"allowed_values" validate:"required"`

	// @asti field=CustomValidator validation=function
	CustomValidator func(interface{}) error `json:"custom_validator"`
}

// @asti name="AdvancedAnnotationService" timeout=180
// @asti category=service model=business
type AdvancedAnnotationService interface {
	// @asti method=CreateAdvancedStruct retry=5 timeout=30
	// @asti validation=strict authorization=required
	CreateAdvancedStruct(ctx context.Context, data AdvancedStruct) (result *AdvancedStruct, err error)

	// @asti method=UpdateAdvancedStruct retry=3 timeout=45
	// @asti validation=partial authorization=required
	UpdateAdvancedStruct(ctx context.Context, id string, updates map[string]interface{}) (result *AdvancedStruct, err error)

	// @asti method=ValidateWithRules retry=2 timeout=60
	// @asti validation=custom authorization=optional
	ValidateWithRules(ctx context.Context, data interface{}, rules ValidationRules) (isValid bool, errors []string, err error)

	// @asti method=ProcessEnum timeout=15
	// @asti validation=enum authorization=public
	ProcessEnum(ctx context.Context, status ComplexEnum) (result string, err error)

	// @asti method=BulkCreate retry=10 timeout=120
	// @asti validation=bulk authorization=admin
	BulkCreate(ctx context.Context, items []AdvancedStruct) (results []*AdvancedStruct, errors []error, err error)

	// @asti method=ComplexValidation timeout=90
	// @asti validation=complex authorization=admin
	ComplexValidation(ctx context.Context, data map[string]interface{}) (isValid bool, details map[string]interface{}, err error)
}
