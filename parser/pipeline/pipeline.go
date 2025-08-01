package pipeline

import (
	"context"

	"github.com/seniorGolang/asti/parser/models"
)

type Stage interface {
	Process(ctx context.Context, data Data) (result Data, err error)
}

type ValidationRule interface {
	Validate(iface models.Interface) (err error)
}

type TypeCollector interface {
	CollectTypes(ctx context.Context, interfaces []models.Interface) (types map[string]models.TypeInfo, err error)
}

type Data struct {
	Package     *models.Package
	Interfaces  []models.Interface
	Types       map[string]models.TypeInfo
	Annotations map[string]models.Annotations
	Errors      []error
}

type Pipeline struct {
	stages []Stage
}

func NewPipeline(stages ...Stage) (pipeline *Pipeline) {

	pipeline = &Pipeline{stages: stages}
	return
}

func (p *Pipeline) AddStage(stage Stage) {

	p.stages = append(p.stages, stage)
}

// Execute выполняет все этапы pipeline
func (p *Pipeline) Execute(ctx context.Context, initialData Data) (result Data, err error) {

	data := initialData
	for _, stage := range p.stages {
		if data, err = stage.Process(ctx, data); err != nil {
			result = data
			return
		}
	}
	result = data
	return
}
