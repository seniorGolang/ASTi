package pipeline

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/seniorGolang/asti/parser/models"
)

type StageSerialization struct{}

func NewStageSerialization() (stage *StageSerialization) {

	stage = &StageSerialization{}
	return
}

// Process выполняет сериализацию результата
func (s *StageSerialization) Process(ctx context.Context, data Data) (result Data, err error) {

	if data.Package == nil {
		return data, fmt.Errorf("package data is required for serialization")
	}
	data.Package.Interfaces = data.Interfaces
	data.Package.Types = data.Types
	if err = s.validatePackage(data.Package); err != nil {
		err = fmt.Errorf("package validation failed: %w", err)
		return
	}
	result = data
	return
}

func (s *StageSerialization) validatePackage(pkg *models.Package) (err error) {

	if pkg.PackagePath == "" {
		err = fmt.Errorf("package path is required")
		return
	}
	interfaceIDs := make(map[string]bool)
	for _, iface := range pkg.Interfaces {
		if iface.ID == "" {
			err = fmt.Errorf("interface %s has empty ID", iface.Name)
			return
		}
		if interfaceIDs[iface.ID] {
			err = fmt.Errorf("duplicate interface ID: %s", iface.ID)
			return
		}
		interfaceIDs[iface.ID] = true
	}
	for _, iface := range pkg.Interfaces {
		methodIDs := make(map[string]bool)
		for _, method := range iface.Methods {
			if method.ID == "" {
				err = fmt.Errorf("method %s in interface %s has empty ID", method.Name, iface.Name)
				return
			}
			if methodIDs[method.ID] {
				err = fmt.Errorf("duplicate method ID %s in interface %s", method.ID, iface.Name)
				return
			}
			methodIDs[method.ID] = true
		}
	}
	return
}

// ToJSON сериализует пакет в JSON
func (s *StageSerialization) ToJSON(pkg *models.Package) (jsonData []byte, err error) {

	jsonData, err = json.MarshalIndent(pkg, "", "  ")
	return
}

// FromJSON десериализует пакет из JSON
func (s *StageSerialization) FromJSON(jsonData []byte) (pkg *models.Package, err error) {

	pkg = &models.Package{}
	if err = json.Unmarshal(jsonData, pkg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return pkg, nil
}
