package pipeline

import (
	"context"
	"fmt"

	"github.com/seniorGolang/asti/parser/models"
)

type StageValidation struct {
	rules []ValidationRule
}

func NewStageValidation(rules ...ValidationRule) (stage *StageValidation) {

	if len(rules) == 0 {
		rules = []ValidationRule{
			&AnnotationRule{},
			&ContextRule{},
			&ErrorRule{},
			&NamedParametersRule{},
		}
	}
	stage = &StageValidation{rules: rules}
	return
}

// Process выполняет валидацию интерфейсов
func (s *StageValidation) Process(ctx context.Context, data Data) (result Data, err error) {

	var errors []error
	var validInterfaces []models.Interface
	for _, iface := range data.Interfaces {
		isValid := true
		for _, rule := range s.rules {
			if err = rule.Validate(iface); err != nil {
				errors = append(errors, fmt.Errorf("validation failed for interface %s: %w", iface.Name, err))
				isValid = false
				break
			}
		}
		if isValid {
			validInterfaces = append(validInterfaces, iface)
		}
	}
	data.Interfaces = validInterfaces
	data.Errors = append(data.Errors, errors...)
	return
}

type AnnotationRule struct{}

func (r *AnnotationRule) Validate(iface models.Interface) (err error) {

	if len(iface.Annotations) == 0 {
		err = fmt.Errorf("interface must have annotations")
		return
	}
	return
}

type ContextRule struct{}

func (r *ContextRule) Validate(iface models.Interface) (err error) {

	for _, method := range iface.Methods {
		if len(method.Parameters) == 0 {
			err = fmt.Errorf("method %s must have at least one parameter (context.Context)", method.Name)
			return
		}
		firstParam := method.Parameters[0]
		if firstParam.Type != "context.Context" {
			err = fmt.Errorf("method %s first parameter must be context.Context, got %s", method.Name, firstParam.Type)
			return
		}
		if firstParam.Name != "ctx" {
			err = fmt.Errorf("method %s first parameter must be named 'ctx', got %s", method.Name, firstParam.Name)
			return
		}
	}
	return
}

type ErrorRule struct{}

func (r *ErrorRule) Validate(iface models.Interface) (err error) {

	for _, method := range iface.Methods {
		if len(method.Results) == 0 {
			err = fmt.Errorf("method %s must return at least one value (error)", method.Name)
			return
		}
		lastResult := method.Results[len(method.Results)-1]
		if lastResult.Type != "error" {
			err = fmt.Errorf("method %s last result must be error, got %s", method.Name, lastResult.Type)
			return
		}
		if lastResult.Name != "err" {
			err = fmt.Errorf("method %s last result must be named 'err', got %s", method.Name, lastResult.Name)
			return
		}
	}
	return
}

type NamedParametersRule struct{}

func (r *NamedParametersRule) Validate(iface models.Interface) (err error) {

	for _, method := range iface.Methods {
		for i, param := range method.Parameters {
			if param.Name == "" {
				err = fmt.Errorf("method %s parameter %d must be named", method.Name, i+1)
				return
			}
		}
		for i, result := range method.Results {
			if result.Name == "" {
				err = fmt.Errorf("method %s result %d must be named", method.Name, i+1)
				return
			}
		}
	}
	return
}
