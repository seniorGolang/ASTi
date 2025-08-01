package complextestcases

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/seniorGolang/asti/parser"
)

// TestPerformance тестирует производительность парсинга
func TestPerformance(t *testing.T) {
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
			var totalTime time.Duration
			iterations := 100

			for i := 0; i < iterations; i++ {
				start := time.Now()

				result, err := p.ParsePackage(ctx, packagePath)
				if err != nil {
					t.Fatalf("Failed to parse package %s on iteration %d: %v", packagePath, i, err)
				}

				duration := time.Since(start)
				totalTime += duration

				// Проверяем, что результат не пустой
				if len(result.Interfaces) == 0 {
					t.Errorf("No interfaces found in package %s on iteration %d", packagePath, i)
				}

				// Проверяем, что типы собраны
				if len(result.Types) == 0 {
					t.Errorf("No types collected in package %s on iteration %d", packagePath, i)
				}
			}

			avgTime := totalTime / time.Duration(iterations)
			t.Logf("Average parsing time for %s: %v", packagePath, avgTime)

			// Проверяем, что среднее время парсинга не превышает разумные пределы
			if avgTime > 100*time.Millisecond {
				t.Errorf("Average parsing time for %s is too high: %v", packagePath, avgTime)
			}
		})
	}
}

// TestStress тестирует парсер под нагрузкой
func TestStress(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	// Тестируем под нагрузкой - много одновременных запросов
	packagePath := "./complex_imports"
	concurrentRequests := 50

	t.Run("Stress_Concurrent", func(t *testing.T) {
		results := make(chan error, concurrentRequests)

		// Запускаем множество одновременных запросов
		for i := 0; i < concurrentRequests; i++ {
			go func(id int) {
				result, err := p.ParsePackage(ctx, packagePath)
				if err != nil {
					results <- fmt.Errorf("goroutine %d failed: %v", id, err)
					return
				}

				if len(result.Interfaces) == 0 {
					results <- fmt.Errorf("goroutine %d: no interfaces found", id)
					return
				}

				results <- nil
			}(i)
		}

		// Собираем результаты
		for i := 0; i < concurrentRequests; i++ {
			if err := <-results; err != nil {
				t.Errorf("Concurrent request failed: %v", err)
			}
		}
	})
}

// TestMemory тестирует использование памяти
func TestMemory(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	packagePath := "./complex_imports"
	iterations := 1000

	t.Run("Memory_Usage", func(t *testing.T) {
		// Запускаем много итераций для проверки утечек памяти
		for i := 0; i < iterations; i++ {
			result, err := p.ParsePackage(ctx, packagePath)
			if err != nil {
				t.Fatalf("Failed to parse package on iteration %d: %v", i, err)
			}

			// Проверяем, что результат валидный
			if len(result.Interfaces) == 0 {
				t.Errorf("No interfaces found on iteration %d", i)
			}

			// Проверяем, что типы собраны
			if len(result.Types) == 0 {
				t.Errorf("No types collected on iteration %d", i)
			}

			// Проверяем JSON сериализацию
			jsonData, err := p.ToJSON(result)
			if err != nil {
				t.Fatalf("Failed to serialize to JSON on iteration %d: %v", i, err)
			}

			if len(jsonData) == 0 {
				t.Errorf("Empty JSON data on iteration %d", i)
			}
		}
	})
}

// TestLargeData тестирует обработку больших объемов данных
func TestLargeData(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	// Тестируем самый сложный пакет
	packagePath := "./advanced_annotations"

	t.Run("Large_Data_Processing", func(t *testing.T) {
		// Запускаем парсинг много раз для имитации больших объемов данных
		for i := 0; i < 100; i++ {
			result, err := p.ParsePackage(ctx, packagePath)
			if err != nil {
				t.Fatalf("Failed to parse large data on iteration %d: %v", i, err)
			}

			// Проверяем сложность результата
			if len(result.Interfaces) == 0 {
				t.Errorf("No interfaces found in large data on iteration %d", i)
			}

			// Проверяем количество типов
			if len(result.Types) < 3 {
				t.Errorf("Expected at least 3 types in large data, got %d on iteration %d", len(result.Types), i)
			}

			// Проверяем сложность аннотаций
			totalAnnotations := 0
			for _, iface := range result.Interfaces {
				totalAnnotations += len(iface.Annotations)
				for _, method := range iface.Methods {
					totalAnnotations += len(method.Annotations)
				}
			}

			if totalAnnotations < 10 {
				t.Errorf("Expected at least 10 annotations in large data, got %d on iteration %d", totalAnnotations, i)
			}
		}
	})
}

// TestConcurrentParsers тестирует работу нескольких парсеров одновременно
func TestConcurrentParsers(t *testing.T) {
	ctx := context.Background()
	packagePath := "./complex_imports"
	parserCount := 10

	t.Run("Concurrent_Parsers", func(t *testing.T) {
		results := make(chan error, parserCount)

		// Создаем несколько парсеров и запускаем их одновременно
		for i := 0; i < parserCount; i++ {
			go func(id int) {
				p := parser.NewParser()

				result, err := p.ParsePackage(ctx, packagePath)
				if err != nil {
					results <- fmt.Errorf("parser %d failed: %v", id, err)
					return
				}

				if len(result.Interfaces) == 0 {
					results <- fmt.Errorf("parser %d: no interfaces found", id)
					return
				}

				// Проверяем JSON сериализацию
				jsonData, err := p.ToJSON(result)
				if err != nil {
					results <- fmt.Errorf("parser %d: JSON serialization failed: %v", id, err)
					return
				}

				if len(jsonData) == 0 {
					results <- fmt.Errorf("parser %d: empty JSON data", id)
					return
				}

				results <- nil
			}(i)
		}

		// Собираем результаты
		for i := 0; i < parserCount; i++ {
			if err := <-results; err != nil {
				t.Errorf("Concurrent parser failed: %v", err)
			}
		}
	})
}

// TestParserReuse тестирует повторное использование парсера
func TestParserReuse(t *testing.T) {
	ctx := context.Background()
	p := parser.NewParser()

	packages := []string{
		"./complex_imports",
		"./cyclic_structures",
		"./generic_types",
		"./nested_types",
		"./advanced_annotations",
		"./corner_cases",
		"./edge_cases",
	}

	t.Run("Parser_Reuse", func(t *testing.T) {
		// Используем один парсер для всех пакетов
		for _, packagePath := range packages {
			result, err := p.ParsePackage(ctx, packagePath)
			if err != nil {
				t.Fatalf("Failed to parse package %s: %v", packagePath, err)
			}

			if len(result.Interfaces) == 0 {
				t.Errorf("No interfaces found in package %s", packagePath)
			}

			// Проверяем JSON сериализацию
			jsonData, err := p.ToJSON(result)
			if err != nil {
				t.Fatalf("Failed to serialize package %s to JSON: %v", packagePath, err)
			}

			if len(jsonData) == 0 {
				t.Errorf("Empty JSON data for package %s", packagePath)
			}

			// Проверяем десериализацию
			deserialized, err := p.FromJSON(jsonData)
			if err != nil {
				t.Fatalf("Failed to deserialize package %s from JSON: %v", packagePath, err)
			}

			if len(deserialized.Interfaces) != len(result.Interfaces) {
				t.Errorf("Interface count mismatch after deserialization for %s", packagePath)
			}
		}
	})
}
