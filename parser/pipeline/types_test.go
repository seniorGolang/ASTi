package pipeline

import (
	"context"
	"testing"

	"github.com/seniorGolang/asti/parser/models"
)

func TestTypeInfoPackageField(t *testing.T) {
	// Создаем тестовый парсер аннотаций
	annotationParser := &models.DefaultAnnotationParser{}

	// Создаем stage для сбора типов
	stage := NewStageTypeCollection(annotationParser)

	// Тестовые данные
	data := Data{
		Package: &models.Package{
			ModuleName:  "testmodule",
			PackagePath: "testpackage",
		},
	}

	// Выполняем обработку
	result, err := stage.Process(context.Background(), data)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Проверяем, что все типы имеют заполненное поле Package
	for typeName, typeInfo := range result.Types {
		if typeInfo.Package == "" {
			t.Errorf("Type %s has empty Package field", typeName)
		}
		t.Logf("Type: %s, Package: %s, Import: %s", typeName, typeInfo.Package, typeInfo.Import)
	}
}

func TestIsBasicType(t *testing.T) {
	annotationParser := &models.DefaultAnnotationParser{}
	stage := NewStageTypeCollection(annotationParser)

	// Тестируем базовые типы
	basicTypes := []string{
		"string", "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "bool", "byte", "rune", "error", "interface{}", "any",
	}

	for _, basicType := range basicTypes {
		if !stage.isBasicType(basicType) {
			t.Errorf("Expected %s to be a basic type", basicType)
		}
	}

	// Тестируем не базовые типы
	nonBasicTypes := []string{
		"mypackage.MyType",
		"strings.Builder",
		"time.Time",
		"CustomType",
	}

	for _, nonBasicType := range nonBasicTypes {
		if stage.isBasicType(nonBasicType) {
			t.Errorf("Expected %s to NOT be a basic type", nonBasicType)
		}
	}
}

func TestGetTypePackage(t *testing.T) {
	annotationParser := &models.DefaultAnnotationParser{}
	stage := NewStageTypeCollection(annotationParser)

	testCases := []struct {
		input    string
		expected string
	}{
		{"mypackage.MyType", "mypackage"},
		{"strings.Builder", "strings"},
		{"time.Time", "time"},
		{"CustomType", ""},
		{"", ""},
	}

	for _, tc := range testCases {
		result := stage.getTypePackage(tc.input)
		if result != tc.expected {
			t.Errorf("getTypePackage(%q) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

func TestGetTypeName(t *testing.T) {
	annotationParser := &models.DefaultAnnotationParser{}
	stage := NewStageTypeCollection(annotationParser)

	testCases := []struct {
		input    string
		expected string
	}{
		{"mypackage.MyType", "MyType"},
		{"strings.Builder", "Builder"},
		{"time.Time", "Time"},
		{"CustomType", "CustomType"},
		{"", ""},
	}

	for _, tc := range testCases {
		result := stage.getTypeName(tc.input)
		if result != tc.expected {
			t.Errorf("getTypeName(%q) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
} 