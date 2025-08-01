package pipeline

import (
	"context"

	"github.com/seniorGolang/asti/parser/models"
)

type StageFilter struct {
	rules []FilterRule
}

func NewStageFilter(rules ...FilterRule) (stage *StageFilter) {

	if len(rules) == 0 {
		rules = []FilterRule{
			&AnnotationFilterRule{},
			&ContextFilterRule{},
			&ErrorFilterRule{},
			&NamedParametersFilterRule{},
		}
	}
	stage = &StageFilter{rules: rules}
	return
}

// Process выполняет фильтрацию интерфейсов
func (s *StageFilter) Process(ctx context.Context, data Data) (result Data, err error) {

	var validInterfaces []models.Interface
	for _, iface := range data.Interfaces {
		isValid := true
		for _, rule := range s.rules {
			if !rule.ShouldInclude(iface) {
				isValid = false
				break
			}
		}
		if isValid {
			validInterfaces = append(validInterfaces, iface)
		}
	}
	data.Interfaces = validInterfaces
	result = data
	return
}

type FilterRule interface {
	ShouldInclude(iface models.Interface) (shouldInclude bool)
}

type AnnotationFilterRule struct{}

func (r *AnnotationFilterRule) ShouldInclude(iface models.Interface) (shouldInclude bool) {
	shouldInclude = len(iface.Annotations) > 0
	return
}

type ContextFilterRule struct{}

func (r *ContextFilterRule) ShouldInclude(iface models.Interface) (shouldInclude bool) {

	for _, method := range iface.Methods {
		if len(method.Parameters) == 0 {
			shouldInclude = false
			return
		}
		firstParam := method.Parameters[0]
		if firstParam.Type != "context.Context" || firstParam.Name != "ctx" {
			shouldInclude = false
			return
		}
	}
	shouldInclude = true
	return
}

type ErrorFilterRule struct{}

func (r *ErrorFilterRule) ShouldInclude(iface models.Interface) (shouldInclude bool) {

	for _, method := range iface.Methods {
		if len(method.Results) == 0 {
			shouldInclude = false
			return
		}

		lastResult := method.Results[len(method.Results)-1]
		if lastResult.Type != "error" || lastResult.Name != "err" {
			shouldInclude = false
			return
		}
	}
	shouldInclude = true
	return
}

type NamedParametersFilterRule struct{}

func (r *NamedParametersFilterRule) ShouldInclude(iface models.Interface) (shouldInclude bool) {

	for _, method := range iface.Methods {
		for _, param := range method.Parameters {
			if param.Name == "" {
				shouldInclude = false
				return
			}
		}
		for _, result := range method.Results {
			if result.Name == "" {
				shouldInclude = false
				return
			}
		}
	}
	shouldInclude = true
	return
}
