package pipeline

import (
	"context"
	"strings"
	"testing"

	"github.com/seniorGolang/asti/parser/models"
)

func TestTypeInfoPackageFieldValidation(t *testing.T) {
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

func TestTypeInfoPackageFieldForExternalTypes(t *testing.T) {
	// Создаем тестовый парсер аннотаций
	annotationParser := &models.DefaultAnnotationParser{}

	// Создаем stage для сбора типов
	stage := NewStageTypeCollection(annotationParser)

	// Тестовые данные с импортами
	data := Data{
		Package: &models.Package{
			ModuleName:  "testmodule",
			PackagePath: "testpackage",
		},
	}

	// Симулируем импорты
	stage.imports = map[string]string{
		"time":    "time",
		"strings": "strings",
		"json":    "encoding/json",
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
		
		// Для типов из внешних пакетов проверяем, что есть Import
		if typeInfo.Package != "testpackage" && typeInfo.Import == "" {
			t.Errorf("Type %s from external package %s has empty Import field", typeName, typeInfo.Package)
		}
		
		t.Logf("Type: %s, Package: %s, Import: %s", typeName, typeInfo.Package, typeInfo.Import)
	}
}

// TestTypeInfoImportFieldValidation проверяет, что поле Import в TypeInfo содержит корректные импорты
func TestTypeInfoImportFieldValidation(t *testing.T) {
	// Создаем тестовый парсер аннотаций
	annotationParser := &models.DefaultAnnotationParser{}

	// Создаем stage для сбора типов
	stage := NewStageTypeCollection(annotationParser)

	// Тестовые данные
	data := Data{
		Package: &models.Package{
			ModuleName:  "github.com/test/module",
			PackagePath: "testpackage",
		},
	}

	// Симулируем импорты с полными путями
	stage.imports = map[string]string{
		"time":    "time",
		"strings": "strings",
		"json":    "encoding/json",
		"external": "github.com/external/package",
	}

	// Выполняем обработку
	result, err := stage.Process(context.Background(), data)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Проверяем каждый тип
	for typeName, typeInfo := range result.Types {
		t.Run("Type_"+typeName, func(t *testing.T) {
			// Проверяем, что поле Import не пустое для внешних типов
			if typeInfo.Package != "testpackage" && typeInfo.Import == "" {
				t.Errorf("Type %s from external package %s has empty Import field", typeName, typeInfo.Package)
			}

			// Проверяем, что импорт не является относительным путем
			if typeInfo.Import != "" {
				if typeInfo.Import == "." || typeInfo.Import == "./" || typeInfo.Import == "/" {
					t.Errorf("Type %s Import '%s' is a relative path, expected full import path", typeName, typeInfo.Import)
				}

				// Проверяем, что путь не начинается с ./ или /
				if strings.HasPrefix(typeInfo.Import, "./") || strings.HasPrefix(typeInfo.Import, "/") {
					t.Errorf("Type %s Import '%s' starts with relative path indicator", typeName, typeInfo.Import)
				}
			}

			t.Logf("Type: %s, Package: %s, Import: %s", typeName, typeInfo.Package, typeInfo.Import)
		})
	}
}

 