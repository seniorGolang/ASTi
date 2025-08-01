package parser

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/seniorGolang/asti/parser/models"
	"github.com/seniorGolang/asti/parser/pipeline"
)

const defaultAnnotationPrefix = "@asti"

type Parser struct {
	annotationPrefix string
	annotationParser models.AnnotationParser
	pipeline         *pipeline.Pipeline
}

func NewParser(options ...Option) (parser *Parser) {

	parser = &Parser{
		annotationPrefix: defaultAnnotationPrefix,
	}
	for _, apply := range options {
		apply(parser)
	}
	annotationParser := models.NewAnnotationParser(parser.annotationPrefix)
	parser.annotationParser = annotationParser
	parser.pipeline = pipeline.NewPipeline(
		pipeline.NewStageModule(),
		pipeline.NewStageAST(annotationParser),
		pipeline.NewStageFilter(),
		pipeline.NewStageTypeCollection(annotationParser),
		pipeline.NewStageSerialization(),
	)
	return
}

// ParsePackage парсит пакет и возвращает информацию об интерфейсах
func (p *Parser) ParsePackage(ctx context.Context, packagePath string) (result *models.Package, err error) {

	if _, err = os.Stat(packagePath); os.IsNotExist(err) {
		err = fmt.Errorf("package path does not exist: %s", packagePath)
		return
	}
	var absPath string
	absPath, err = filepath.Abs(packagePath)
	if err != nil {
		err = fmt.Errorf("failed to get absolute path: %w", err)
		return
	}
	initialData := pipeline.Data{
		Package: &models.Package{
			PackagePath: absPath,
		},
	}

	var resultData pipeline.Data
	resultData, err = p.pipeline.Execute(ctx, initialData)
	if err != nil {
		err = fmt.Errorf("pipeline execution failed: %w", err)
		return
	}
	result = resultData.Package
	return
}

// ParsePackageToJSON парсит пакет и возвращает JSON
func (p *Parser) ParsePackageToJSON(ctx context.Context, packagePath string) (jsonData []byte, err error) {

	var pkg *models.Package
	pkg, err = p.ParsePackage(ctx, packagePath)
	if err != nil {
		return nil, err
	}
	jsonData, err = p.ToJSON(pkg)
	return
}

// ToJSON сериализует пакет в JSON
func (p *Parser) ToJSON(pkg *models.Package) (jsonData []byte, err error) {

	serializationStage := pipeline.NewStageSerialization()
	jsonData, err = serializationStage.ToJSON(pkg)
	return
}

// FromJSON десериализует пакет из JSON
func (p *Parser) FromJSON(jsonData []byte) (pkg *models.Package, err error) {

	serializationStage := pipeline.NewStageSerialization()
	pkg, err = serializationStage.FromJSON(jsonData)
	return
}

func (p *Parser) GetAnnotationPrefix() (prefix string) {

	prefix = p.annotationPrefix
	return
}

func (p *Parser) SetAnnotationPrefix(prefix string) {

	if prefix == "" {
		prefix = "@asti"
	}

	p.annotationPrefix = prefix
	p.annotationParser = models.NewAnnotationParser(prefix)

	p.pipeline = pipeline.NewPipeline(
		pipeline.NewStageAST(p.annotationParser),
		pipeline.NewStageFilter(),
		pipeline.NewStageTypeCollection(p.annotationParser),
		pipeline.NewStageSerialization(),
	)
}
