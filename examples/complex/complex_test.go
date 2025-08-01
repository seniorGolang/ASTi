package complextestcases

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/seniorGolang/asti/parser"
)

// TestComplexImports тестирует парсинг сложных импортов с множественными алиасами
func TestComplexImports(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	packagePath := "./complex_imports"
	result, err := p.ParsePackage(ctx, packagePath)
	if err != nil {
		t.Fatalf("Failed to parse complex imports package: %v", err)
	}

	// Проверяем, что найден интерфейс
	if len(result.Interfaces) == 0 {
		t.Fatal("No interfaces found in complex imports package")
	}

	// Проверяем аннотации пакета
	if result.Annotations["name"] != "ComplexImports" {
		t.Errorf("Expected package name 'ComplexImports', got '%s'", result.Annotations["name"])
	}

	// Проверяем интерфейс
	service := result.Interfaces[0]
	if service.Name != "ComplexImportService" {
		t.Errorf("Expected interface name 'ComplexImportService', got '%s'", service.Name)
	}

	// Проверяем методы с алиасами
	expectedMethods := []string{
		"ProcessWithSQL",
		"ProcessWithAliasSQL",
		"ProcessWithDBAlias",
		"ProcessJSON",
		"ProcessJSONAlias",
		"ProcessHTTP2",
		"ProcessModels",
		"ProcessTime",
	}

	for _, methodName := range expectedMethods {
		found := false
		for _, method := range service.Methods {
			if method.Name == methodName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Method '%s' not found in interface", methodName)
		}
	}

	// Проверяем типы с алиасами
	if len(result.Types) == 0 {
		t.Error("No types collected for complex imports")
	}

	// Проверяем наличие типов с алиасами
	expectedTypes := []string{
		"sql.Result",
		"json.RawMessage",
		"time.Time",
		"models.Annotations",
	}

	for _, typeName := range expectedTypes {
		found := false
		for typeKey := range result.Types {
			if typeKey == typeName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Type '%s' not found in collected types", typeName)
		}
	}
}

// TestCyclicStructures тестирует парсинг циклических структур
func TestCyclicStructures(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	packagePath := "./cyclic_structures"
	result, err := p.ParsePackage(ctx, packagePath)
	if err != nil {
		t.Fatalf("Failed to parse cyclic structures package: %v", err)
	}

	// Проверяем интерфейс
	if len(result.Interfaces) == 0 {
		t.Fatal("No interfaces found in cyclic structures package")
	}

	service := result.Interfaces[0]
	if service.Name != "CyclicStructureService" {
		t.Errorf("Expected interface name 'CyclicStructureService', got '%s'", service.Name)
	}

	// Проверяем методы с циклическими типами
	expectedMethods := []string{
		"CreateNode",
		"BuildGraph",
		"TraverseGraph",
		"CreateRecursive",
		"ProcessMatrix",
		"UpdateAdjacencyMap",
	}

	for _, methodName := range expectedMethods {
		found := false
		for _, method := range service.Methods {
			if method.Name == methodName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Method '%s' not found in interface", methodName)
		}
	}

	// Проверяем типы с циклическими ссылками (базовые типы)
	expectedTypes := []string{
		"cyclicstructures.Node",
		"cyclicstructures.ComplexGraph",
		"cyclicstructures.RecursiveType",
	}

	for _, typeName := range expectedTypes {
		found := false
		for typeKey := range result.Types {
			if typeKey == typeName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Type '%s' not found in collected types", typeName)
		}
	}
}

// TestGenericTypes тестирует парсинг дженериков
func TestGenericTypes(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	packagePath := "./generic_types"
	result, err := p.ParsePackage(ctx, packagePath)
	if err != nil {
		t.Fatalf("Failed to parse generic types package: %v", err)
	}

	// Проверяем интерфейс
	if len(result.Interfaces) == 0 {
		t.Fatal("No interfaces found in generic types package")
	}

	service := result.Interfaces[0]
	if service.Name != "GenericTypeService" {
		t.Errorf("Expected interface name 'GenericTypeService', got '%s'", service.Name)
	}

	// Проверяем методы с дженериками
	expectedMethods := []string{
		"CreateContainer",
		"ProcessComplexGeneric",
		"BuildGenericTree",
		"ProcessMultiGeneric",
		"GenericMapOperation",
		"GenericSliceOperation",
	}

	for _, methodName := range expectedMethods {
		found := false
		for _, method := range service.Methods {
			if method.Name == methodName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Method '%s' not found in interface", methodName)
		}
	}

	// Проверяем дженерик типы (базовые типы)
	expectedTypes := []string{
		"generictypes.GenericContainer",
		"generictypes.ComplexGeneric",
		"generictypes.GenericNode",
		"generictypes.MultiGeneric",
	}



	for _, typeName := range expectedTypes {
		found := false
		for typeKey := range result.Types {
			if typeKey == typeName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Type '%s' not found in collected types", typeName)
		}
	}
}

// TestNestedTypes тестирует парсинг глубоко вложенных типов
func TestNestedTypes(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	packagePath := "./nested_types"
	result, err := p.ParsePackage(ctx, packagePath)
	if err != nil {
		t.Fatalf("Failed to parse nested types package: %v", err)
	}

	// Проверяем интерфейс
	if len(result.Interfaces) == 0 {
		t.Fatal("No interfaces found in nested types package")
	}

	service := result.Interfaces[0]
	if service.Name != "NestedTypeService" {
		t.Errorf("Expected interface name 'NestedTypeService', got '%s'", service.Name)
	}

	// Проверяем методы с вложенными типами
	expectedMethods := []string{
		"ProcessDeepNested",
		"ProcessComplexSlice",
		"ProcessComplexMap",
		"ProcessInterfaceEmbedding",
		"ProcessFunctionType",
		"ProcessChannelStruct",
	}

	for _, methodName := range expectedMethods {
		found := false
		for _, method := range service.Methods {
			if method.Name == methodName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Method '%s' not found in interface", methodName)
		}
	}

	// Проверяем сложные типы
	expectedTypes := []string{
		"nestedtypes.DeepNestedStruct",
		"nestedtypes.ComplexSliceStruct",
		"nestedtypes.ComplexMapStruct",
		"nestedtypes.InterfaceEmbedding",
		"nestedtypes.FunctionType",
		"nestedtypes.ChannelStruct",
	}

	for _, typeName := range expectedTypes {
		found := false
		for typeKey := range result.Types {
			if typeKey == typeName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Type '%s' not found in collected types", typeName)
		}
	}
}

// TestAdvancedAnnotations тестирует парсинг продвинутых аннотаций
func TestAdvancedAnnotations(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	packagePath := "./advanced_annotations"
	result, err := p.ParsePackage(ctx, packagePath)
	if err != nil {
		t.Fatalf("Failed to parse advanced annotations package: %v", err)
	}

	// Проверяем аннотации пакета
	expectedPackageAnnotations := map[string]string{
		"name":    "AdvancedAnnotations",
		"version": "3.0",
		"author":  "Senior Developer",
		"team":    "Backend",
		"tags":    "complex,advanced,testing",
	}

	for key, expectedValue := range expectedPackageAnnotations {
		if result.Annotations[key] != expectedValue {
			t.Errorf("Package annotation '%s' expected '%s', got '%s'", key, expectedValue, result.Annotations[key])
		}
	}

	// Проверяем интерфейс
	if len(result.Interfaces) == 0 {
		t.Fatal("No interfaces found in advanced annotations package")
	}

	service := result.Interfaces[0]
	if service.Name != "AdvancedAnnotationService" {
		t.Errorf("Expected interface name 'AdvancedAnnotationService', got '%s'", service.Name)
	}

	// Проверяем аннотации интерфейса
	expectedInterfaceAnnotations := map[string]string{
		"timeout":  "180",
		"category": "service",
		"model":    "business",
	}

	for key, expectedValue := range expectedInterfaceAnnotations {
		if service.Annotations[key] != expectedValue {
			t.Errorf("Interface annotation '%s' expected '%s', got '%s'", key, expectedValue, service.Annotations[key])
		}
	}

	// Проверяем методы с продвинутыми аннотациями
	expectedMethods := []string{
		"CreateAdvancedStruct",
		"UpdateAdvancedStruct",
		"ValidateWithRules",
		"ProcessEnum",
		"BulkCreate",
		"ComplexValidation",
	}

	for _, methodName := range expectedMethods {
		found := false
		for _, method := range service.Methods {
			if method.Name == methodName {
				found = true
				// Проверяем наличие аннотаций у метода
				if len(method.Annotations) == 0 {
					t.Errorf("Method '%s' has no annotations", methodName)
				}
				break
			}
		}
		if !found {
			t.Errorf("Method '%s' not found in interface", methodName)
		}
	}

	// Проверяем сложные типы с аннотациями
	expectedTypes := []string{
		"advancedannotations.AdvancedStruct",
		"advancedannotations.ComplexEnum",
		"advancedannotations.ValidationRules",
	}

	for _, typeName := range expectedTypes {
		found := false
		for typeKey, typeInfo := range result.Types {
			if typeKey == typeName {
				found = true
				// Проверяем наличие аннотаций у типа
				if len(typeInfo.Annotations) == 0 {
					t.Errorf("Type '%s' has no annotations", typeName)
				}
				break
			}
		}
		if !found {
			t.Errorf("Type '%s' not found in collected types", typeName)
		}
	}
}

// TestComplexJSONSerialization тестирует JSON сериализацию сложных структур
func TestComplexJSONSerialization(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	// Тестируем все пакеты
	packages := []string{
		"./complex_imports",
		"./cyclic_structures",
		"./generic_types",
		"./nested_types",
		"./advanced_annotations",
		"./corner_cases",
		"./edge_cases",
	}

	for _, packagePath := range packages {
		t.Run(fmt.Sprintf("JSON_%s", filepath.Base(packagePath)), func(t *testing.T) {
			result, err := p.ParsePackage(ctx, packagePath)
			if err != nil {
				t.Fatalf("Failed to parse package %s: %v", packagePath, err)
			}

			// Сериализуем в JSON
			jsonData, err := p.ToJSON(result)
			if err != nil {
				t.Fatalf("Failed to serialize package %s to JSON: %v", packagePath, err)
			}

			// Проверяем, что JSON не пустой
			if len(jsonData) == 0 {
				t.Errorf("JSON data is empty for package %s", packagePath)
			}

			// Проверяем, что JSON валидный
			var parsed map[string]interface{}
			if err := json.Unmarshal(jsonData, &parsed); err != nil {
				t.Errorf("Invalid JSON for package %s: %v", packagePath, err)
			}

			// Проверяем наличие обязательных полей
			if _, ok := parsed["interfaces"]; !ok {
				t.Errorf("Missing 'interfaces' field in JSON for package %s", packagePath)
			}

			if _, ok := parsed["types"]; !ok {
				t.Errorf("Missing 'types' field in JSON for package %s", packagePath)
			}

			// Десериализуем обратно
			deserialized, err := p.FromJSON(jsonData)
			if err != nil {
				t.Fatalf("Failed to deserialize package %s from JSON: %v", packagePath, err)
			}

			// Проверяем, что десериализованный результат совпадает с исходным
			if deserialized.PackagePath != result.PackagePath {
				t.Errorf("Package path mismatch after deserialization for %s", packagePath)
			}

			if len(deserialized.Interfaces) != len(result.Interfaces) {
				t.Errorf("Interface count mismatch after deserialization for %s", packagePath)
			}
		})
	}
}

// TestComplexValidation тестирует валидацию сложных интерфейсов
func TestComplexValidation(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	// Тестируем все пакеты на соответствие правилам валидации
	packages := []string{
		"./complex_imports",
		"./cyclic_structures",
		"./generic_types",
		"./nested_types",
		"./advanced_annotations",
		"./corner_cases",
		"./edge_cases",
	}

	for _, packagePath := range packages {
		t.Run(fmt.Sprintf("Validation_%s", filepath.Base(packagePath)), func(t *testing.T) {
			result, err := p.ParsePackage(ctx, packagePath)
			if err != nil {
				t.Fatalf("Failed to parse package %s: %v", packagePath, err)
			}

			// Проверяем, что все интерфейсы имеют аннотации
			for _, iface := range result.Interfaces {
				if len(iface.Annotations) == 0 {
					t.Errorf("Interface %s has no annotations in package %s", iface.Name, packagePath)
				}

				// Проверяем, что все методы имеют context.Context как первый параметр
				for _, method := range iface.Methods {
					if len(method.Parameters) == 0 {
						t.Errorf("Method %s has no parameters in interface %s", method.Name, iface.Name)
						continue
					}

					firstParam := method.Parameters[0]
					if firstParam.Type != "context.Context" {
						t.Errorf("Method %s first parameter should be context.Context, got %s", method.Name, firstParam.Type)
					}

					if firstParam.Name != "ctx" {
						t.Errorf("Method %s first parameter should be named 'ctx', got %s", method.Name, firstParam.Name)
					}

					// Проверяем, что последний результат - error
					if len(method.Results) == 0 {
						t.Errorf("Method %s has no results in interface %s", method.Name, iface.Name)
						continue
					}

					lastResult := method.Results[len(method.Results)-1]
					if lastResult.Type != "error" {
						t.Errorf("Method %s last result should be error, got %s", method.Name, lastResult.Type)
					}

					if lastResult.Name != "err" {
						t.Errorf("Method %s last result should be named 'err', got %s", method.Name, lastResult.Name)
					}

					// Проверяем, что все параметры и результаты именованы
					for _, param := range method.Parameters {
						if param.Name == "" {
							t.Errorf("Method %s has unnamed parameter of type %s", method.Name, param.Type)
						}
					}

					for _, result := range method.Results {
						if result.Name == "" {
							t.Errorf("Method %s has unnamed result of type %s", method.Name, result.Type)
						}
					}
				}
			}
		})
	}
}

// TestComplexErrorHandling тестирует обработку ошибок в сложных случаях
func TestComplexErrorHandling(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	// Тестируем несуществующий пакет
	_, err := p.ParsePackage(ctx, "./non_existent_package")
	if err == nil {
		t.Error("Expected error for non-existent package, got nil")
	}

	// Тестируем пустую директорию
	emptyDir := "./empty_test_dir"
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatalf("Failed to create empty test directory: %v", err)
	}
	defer os.RemoveAll(emptyDir)

	result, err := p.ParsePackage(ctx, emptyDir)
	if err != nil {
		t.Fatalf("Failed to parse empty directory: %v", err)
	}

	// Пустая директория должна вернуть пустой результат
	if len(result.Interfaces) != 0 {
		t.Errorf("Expected 0 interfaces in empty directory, got %d", len(result.Interfaces))
	}

	if len(result.Types) != 0 {
		t.Errorf("Expected 0 types in empty directory, got %d", len(result.Types))
	}
}

// TestCornerCases тестирует парсинг корнер кейсов
func TestCornerCases(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	packagePath := "./corner_cases"
	result, err := p.ParsePackage(ctx, packagePath)
	if err != nil {
		t.Fatalf("Failed to parse corner cases package: %v", err)
	}

	// Проверяем интерфейс
	if len(result.Interfaces) == 0 {
		t.Fatal("No interfaces found in corner cases package")
	}

	service := result.Interfaces[0]
	if service.Name != "CornerCaseService" {
		t.Errorf("Expected interface name 'CornerCaseService', got '%s'", service.Name)
	}

	// Проверяем методы с корнер кейсами
	expectedMethods := []string{
		"ProcessEmptyStruct",
		"ProcessPointerStruct",
		"ProcessArrayStruct",
		"ProcessSliceStruct",
		"ProcessMapStruct",
		"ProcessChannelStruct",
		"ProcessInterfaceStruct",
		"ProcessFunctionStruct",
		"ProcessEmbeddedStruct",
		"ProcessTaggedStruct",
	}

	for _, methodName := range expectedMethods {
		found := false
		for _, method := range service.Methods {
			if method.Name == methodName {
				found = true
				// Проверяем наличие аннотаций у метода
				if len(method.Annotations) == 0 {
					t.Errorf("Method '%s' has no annotations", methodName)
				}
				break
			}
		}
		if !found {
			t.Errorf("Method '%s' not found in interface", methodName)
		}
	}

	// Проверяем сложные типы с корнер кейсами (базовые типы)
	expectedTypes := []string{
		"cornercases.EmptyStruct",
		"cornercases.PointerStruct",
		"cornercases.ArrayStruct",
		"cornercases.SliceStruct",
		"cornercases.MapStruct",
		"cornercases.ChannelStruct",
		"cornercases.InterfaceStruct",
		"cornercases.FunctionStruct",
		"cornercases.EmbeddedStruct",
		"cornercases.TaggedStruct",
	}

	for _, typeName := range expectedTypes {
		found := false
		for typeKey, typeInfo := range result.Types {
			if typeKey == typeName {
				found = true
				// Проверяем наличие аннотаций у типа
				if len(typeInfo.Annotations) == 0 {
					t.Errorf("Type '%s' has no annotations", typeName)
				}
				break
			}
		}
		if !found {
			t.Errorf("Type '%s' not found in collected types", typeName)
		}
	}

	// Проверяем аннотации пакета
	expectedPackageAnnotations := map[string]string{
		"name":     "CornerCases",
		"version":  "1.0",
		"category": "testing",
		"model":    "corner",
	}

	for key, expectedValue := range expectedPackageAnnotations {
		if result.Annotations[key] != expectedValue {
			t.Errorf("Package annotation '%s' expected '%s', got '%s'", key, expectedValue, result.Annotations[key])
		}
	}
}

// TestEdgeCases тестирует парсинг граничных случаев
func TestEdgeCases(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	packagePath := "./edge_cases"
	result, err := p.ParsePackage(ctx, packagePath)
	if err != nil {
		t.Fatalf("Failed to parse edge cases package: %v", err)
	}

	// Проверяем, что найдены только валидные интерфейсы
	validInterfaces := []string{
		"ValidInterface",
		"VeryLongInterfaceNameThatExceedsNormalLimitsAndTestsHowTheParserHandlesExtremelyLongIdentifiers",
		"MultipleAnnotations",
		"ComplexTypesInterface",
		"SpecialChars",
		"ManyMethodsInterface",
	}

	// Проверяем, что невалидные интерфейсы пропущены
	invalidInterfaces := []string{
		"UnannotatedInterface",
		"InvalidContextInterface",
		"InvalidErrorInterface",
		"UnnamedParamsInterface",
		"UnnamedResultsInterface",
	}

	// Проверяем наличие валидных интерфейсов
	for _, interfaceName := range validInterfaces {
		found := false
		for _, iface := range result.Interfaces {
			if iface.Name == interfaceName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Valid interface '%s' not found in parsed result", interfaceName)
		}
	}

	// Проверяем отсутствие невалидных интерфейсов
	for _, interfaceName := range invalidInterfaces {
		found := false
		for _, iface := range result.Interfaces {
			if iface.Name == interfaceName {
				found = true
				break
			}
		}
		if found {
			t.Errorf("Invalid interface '%s' should not be found in parsed result", interfaceName)
		}
	}

	// Проверяем аннотации пакета
	if result.Annotations["name"] != "EdgeCases" {
		t.Errorf("Expected package name 'EdgeCases', got '%s'", result.Annotations["name"])
	}

	// Проверяем интерфейс с множественными аннотациями
	for _, iface := range result.Interfaces {
		if iface.Name == "MultipleAnnotations" {
			expectedAnnotations := map[string]string{
				"name":     "MultipleAnnotations",
				"version":  "2.0",
				"timeout":  "45",
				"category": "service",
				"model":    "business",
				"author":   "Senior Developer",
				"team":     "Backend",
			}

			for key, expectedValue := range expectedAnnotations {
				if iface.Annotations[key] != expectedValue {
					t.Errorf("Interface annotation '%s' expected '%s', got '%s'", key, expectedValue, iface.Annotations[key])
				}
			}
		}
	}

	// Проверяем интерфейс с очень длинными именами
	for _, iface := range result.Interfaces {
		if iface.Name == "VeryLongInterfaceNameThatExceedsNormalLimitsAndTestsHowTheParserHandlesExtremelyLongIdentifiers" {
			if len(iface.Methods) == 0 {
				t.Error("Very long interface should have methods")
			}

			// Проверяем очень длинный метод
			for _, method := range iface.Methods {
				if method.Name == "VeryLongMethodNameThatExceedsNormalLimitsAndTestsHowTheParserHandlesExtremelyLongIdentifiers" {
					if len(method.Parameters) == 0 {
						t.Error("Very long method should have parameters")
					}
					break
				}
			}
		}
	}

	// Проверяем интерфейс с очень большим количеством методов
	for _, iface := range result.Interfaces {
		if iface.Name == "ManyMethodsInterface" {
			if len(iface.Methods) != 10 {
				t.Errorf("ManyMethodsInterface should have 10 methods, got %d", len(iface.Methods))
			}
		}
	}
}

// TestComplexPerformance тестирует производительность парсинга сложных структур
func TestComplexPerformance(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	// Тестируем производительность парсинга всех пакетов
	packages := []string{
		"./complex_imports",
		"./cyclic_structures",
		"./generic_types",
		"./nested_types",
		"./advanced_annotations",
		"./corner_cases",
		"./edge_cases",
	}

	for _, packagePath := range packages {
		t.Run(fmt.Sprintf("Performance_%s", filepath.Base(packagePath)), func(t *testing.T) {
			// Запускаем парсинг несколько раз для измерения производительности
			for i := 0; i < 10; i++ {
				result, err := p.ParsePackage(ctx, packagePath)
				if err != nil {
					t.Fatalf("Failed to parse package %s on iteration %d: %v", packagePath, i, err)
				}

				// Проверяем, что результат не пустой
				if len(result.Interfaces) == 0 {
					t.Errorf("No interfaces found in package %s on iteration %d", packagePath, i)
				}

				// Проверяем, что типы собраны
				if len(result.Types) == 0 {
					t.Errorf("No types collected in package %s on iteration %d", packagePath, i)
				}
			}
		})
	}
}
