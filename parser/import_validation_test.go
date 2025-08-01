package parser

import (
	"context"
	"strings"
	"testing"
)

// TestInterfaceImportValidation проверяет, что поле Import в Interface содержит полноценный импорт
func TestInterfaceImportValidation(t *testing.T) {
	// Создаем парсер
	p := NewParser()

	// Парсим пакет с интерфейсами
	result, err := p.ParsePackage(context.Background(), "../examples/sample")
	if err != nil {
		t.Fatalf("Failed to parse package: %v", err)
	}

	// Проверяем, что интерфейсы найдены
	if len(result.Interfaces) == 0 {
		t.Fatal("No interfaces found")
	}

	// Проверяем каждый интерфейс
	for _, iface := range result.Interfaces {
		t.Run("Interface_"+iface.Name, func(t *testing.T) {
			validateFullImport(t, iface.Import, "Interface "+iface.Name)
		})
	}
}

// TestTypeInfoImportValidation проверяет, что поле Import в TypeInfo содержит полноценный импорт
func TestTypeInfoImportValidation(t *testing.T) {
	// Создаем парсер
	p := NewParser()

	// Парсим пакет с типами
	result, err := p.ParsePackage(context.Background(), "../examples/sample")
	if err != nil {
		t.Fatalf("Failed to parse package: %v", err)
	}

	// Проверяем, что типы найдены
	if len(result.Types) == 0 {
		t.Fatal("No types found")
	}

	// Проверяем каждый тип
	for typeName, typeInfo := range result.Types {
		t.Run("Type_"+typeName, func(t *testing.T) {
			// Пропускаем встроенные типы Go, которые не должны иметь Import
			if isBuiltinType(typeName) {
				if typeInfo.Import != "" {
					t.Errorf("Builtin type %s should not have Import field, got: %s", typeName, typeInfo.Import)
				}
				return
			}

			// Для типов из внешних пакетов проверяем полноценный импорт
			if typeInfo.Package != "examples" && typeInfo.Import != "" {
				validateFullImport(t, typeInfo.Import, "Type "+typeName)
			}
		})
	}
}

// TestComplexPackageImportValidation проверяет импорты в сложном пакете
func TestComplexPackageImportValidation(t *testing.T) {
	// Создаем парсер
	p := NewParser()

	// Парсим сложный пакет
	result, err := p.ParsePackage(context.Background(), "../examples/complex/complex_imports")
	if err != nil {
		t.Fatalf("Failed to parse complex package: %v", err)
	}

	// Проверяем интерфейсы
	for _, iface := range result.Interfaces {
		t.Run("ComplexInterface_"+iface.Name, func(t *testing.T) {
			validateFullImport(t, iface.Import, "Complex Interface "+iface.Name)
		})
	}

			// Проверяем типы
		for typeName, typeInfo := range result.Types {
			t.Run("ComplexType_"+typeName, func(t *testing.T) {
				if isBuiltinType(typeName) {
					return
				}

				// Пропускаем стандартные библиотеки Go
				if typeInfo.Import != "" && isStandardLibraryImport(typeInfo.Import) {
					return
				}

				if typeInfo.Package != "compleximports" && typeInfo.Import != "" {
					validateFullImport(t, typeInfo.Import, "Complex Type "+typeName)
				}
			})
		}
}

// TestModuleNameInImport проверяет, что импорты содержат имя модуля
func TestModuleNameInImport(t *testing.T) {
	// Создаем парсер
	p := NewParser()

	// Парсим пакет
	result, err := p.ParsePackage(context.Background(), "../examples/sample")
	if err != nil {
		t.Fatalf("Failed to parse package: %v", err)
	}

	expectedModuleName := "github.com/seniorGolang/asti"

	// Проверяем интерфейсы
	for _, iface := range result.Interfaces {
		t.Run("InterfaceModule_"+iface.Name, func(t *testing.T) {
			if !strings.Contains(iface.Import, expectedModuleName) {
				t.Errorf("Interface %s Import '%s' does not contain expected module name '%s'", 
					iface.Name, iface.Import, expectedModuleName)
			}
		})
	}

	// Проверяем типы
	for typeName, typeInfo := range result.Types {
		t.Run("TypeModule_"+typeName, func(t *testing.T) {
			if isBuiltinType(typeName) {
				return
			}

			// Пропускаем стандартные библиотеки Go
			if typeInfo.Import != "" && isStandardLibraryImport(typeInfo.Import) {
				return
			}

			if typeInfo.Import != "" && !strings.Contains(typeInfo.Import, expectedModuleName) {
				t.Errorf("Type %s Import '%s' does not contain expected module name '%s'", 
					typeName, typeInfo.Import, expectedModuleName)
			}
		})
	}
}

// TestImportFormatValidation проверяет формат импортов
func TestImportFormatValidation(t *testing.T) {
	// Создаем парсер
	p := NewParser()

	// Парсим пакет
	result, err := p.ParsePackage(context.Background(), "../examples/sample")
	if err != nil {
		t.Fatalf("Failed to parse package: %v", err)
	}

	// Проверяем интерфейсы
	for _, iface := range result.Interfaces {
		t.Run("InterfaceFormat_"+iface.Name, func(t *testing.T) {
			validateImportFormat(t, iface.Import, "Interface "+iface.Name)
		})
	}

			// Проверяем типы
		for typeName, typeInfo := range result.Types {
			t.Run("TypeFormat_"+typeName, func(t *testing.T) {
				if isBuiltinType(typeName) || typeInfo.Import == "" {
					return
				}

				// Пропускаем стандартные библиотеки Go
				if isStandardLibraryImport(typeInfo.Import) {
					return
				}

				validateImportFormat(t, typeInfo.Import, "Type "+typeName)
			})
		}
}

// validateFullImport проверяет, что импорт является полноценным
func validateFullImport(t *testing.T, importPath, context string) {
	if importPath == "" {
		t.Errorf("%s: Import field is empty", context)
		return
	}

	// Проверяем, что это не относительный путь
	if importPath == "." || importPath == "./" || importPath == "/" {
		t.Errorf("%s: Import '%s' is a relative path, expected full import path", context, importPath)
		return
	}

	// Для встроенных типов Go (context, time, math/big и т.д.) не требуем доменное имя
	if isStandardLibraryImport(importPath) {
		// Проверяем только, что путь не начинается с ./ или /
		if strings.HasPrefix(importPath, "./") || strings.HasPrefix(importPath, "/") {
			t.Errorf("%s: Import '%s' starts with relative path indicator", context, importPath)
			return
		}
		t.Logf("%s: Valid standard library import path '%s'", context, importPath)
		return
	}

	// Проверяем, что путь содержит доменное имя (обычно github.com или другое)
	if !strings.Contains(importPath, ".") {
		t.Errorf("%s: Import '%s' does not contain domain name", context, importPath)
		return
	}

	// Проверяем, что путь не начинается с ./ или /
	if strings.HasPrefix(importPath, "./") || strings.HasPrefix(importPath, "/") {
		t.Errorf("%s: Import '%s' starts with relative path indicator", context, importPath)
		return
	}

	t.Logf("%s: Valid import path '%s'", context, importPath)
}

// validateImportFormat проверяет формат импорта
func validateImportFormat(t *testing.T, importPath, context string) {
	if importPath == "" {
		return // Пустые импорты допустимы для встроенных типов
	}

	// Проверяем, что путь не содержит пробелов
	if strings.Contains(importPath, " ") {
		t.Errorf("%s: Import '%s' contains spaces", context, importPath)
		return
	}

	// Проверяем, что путь не содержит недопустимых символов
	invalidChars := []string{"\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(importPath, char) {
			t.Errorf("%s: Import '%s' contains invalid character '%s'", context, importPath, char)
			return
		}
	}

	// Для стандартных библиотек Go проверяем только базовый формат
	if isStandardLibraryImport(importPath) {
		// Проверяем, что путь не начинается с ./ или /
		if strings.HasPrefix(importPath, "./") || strings.HasPrefix(importPath, "/") {
			t.Errorf("%s: Import '%s' starts with relative path indicator", context, importPath)
			return
		}
		t.Logf("%s: Valid standard library import format '%s'", context, importPath)
		return
	}

	// Проверяем, что путь имеет правильную структуру (домен/путь)
	parts := strings.Split(importPath, "/")
	if len(parts) < 2 {
		t.Errorf("%s: Import '%s' does not have proper structure (domain/path)", context, importPath)
		return
	}

	// Проверяем, что доменная часть содержит точку
	if !strings.Contains(parts[0], ".") {
		t.Errorf("%s: Import '%s' domain part '%s' does not contain dot", context, importPath, parts[0])
		return
	}

	t.Logf("%s: Valid import format '%s'", context, importPath)
}

// isStandardLibraryImport проверяет, является ли импорт стандартной библиотекой Go
func isStandardLibraryImport(importPath string) bool {
	standardLibs := map[string]bool{
		"context":        true,
		"time":           true,
		"math/big":       true,
		"encoding/json":  true,
		"database/sql":   true,
		"fmt":            true,
		"strings":        true,
		"strconv":        true,
		"os":             true,
		"io":             true,
		"bufio":          true,
		"bytes":          true,
		"crypto":         true,
		"crypto/md5":     true,
		"crypto/sha1":    true,
		"crypto/sha256":  true,
		"crypto/sha512":  true,
		"encoding/base64": true,
		"encoding/hex":   true,
		"encoding/xml":   true,
		"encoding/csv":   true,
		"encoding/gob":   true,
		"encoding/binary": true,
		"compress/gzip":  true,
		"compress/zlib":  true,
		"archive/zip":    true,
		"archive/tar":    true,
		"path":           true,
		"path/filepath":  true,
		"regexp":         true,
		"sort":           true,
		"container/heap": true,
		"container/list": true,
		"container/ring": true,
		"hash":           true,
		"hash/adler32":   true,
		"hash/crc32":     true,
		"hash/crc64":     true,
		"hash/fnv":       true,
		"hash/md5":       true,
		"hash/sha1":      true,
		"hash/sha256":    true,
		"hash/sha512":    true,
		"image":          true,
		"image/color":    true,
		"image/draw":     true,
		"image/gif":      true,
		"image/jpeg":     true,
		"image/png":      true,
		"log":            true,
		"log/syslog":     true,
		"math":           true,
		"math/cmplx":     true,
		"math/rand":      true,
		"mime":           true,
		"mime/multipart": true,
		"net":            true,
		"net/http":       true,
		"net/http/cgi":   true,
		"net/http/cookiejar": true,
		"net/http/fcgi":  true,
		"net/http/httptest": true,
		"net/http/httptrace": true,
		"net/http/httputil": true,
		"net/http/pprof": true,
		"net/mail":       true,
		"net/rpc":        true,
		"net/smtp":       true,
		"net/textproto":  true,
		"net/url":        true,
		"reflect":        true,
		"runtime":        true,
		"runtime/debug":  true,
		"runtime/pprof":  true,
		"runtime/trace":  true,
		"sync":           true,
		"sync/atomic":    true,
		"syscall":        true,
		"testing":        true,
		"testing/iotest": true,
		"testing/quick":  true,
		"text":           true,
		"text/scanner":   true,
		"text/tabwriter": true,
		"text/template":  true,
		"text/template/parse": true,
		"unicode":        true,
		"unicode/utf16":  true,
		"unicode/utf8":   true,
	}

	return standardLibs[importPath]
}

// isBuiltinType проверяет, является ли тип встроенным типом Go
func isBuiltinType(typeName string) bool {
	builtinTypes := map[string]bool{
		"string":     true,
		"int":        true,
		"int8":       true,
		"int16":      true,
		"int32":      true,
		"int64":      true,
		"uint":       true,
		"uint8":      true,
		"uint16":     true,
		"uint32":     true,
		"uint64":     true,
		"float32":    true,
		"float64":    true,
		"bool":       true,
		"byte":       true,
		"rune":       true,
		"error":      true,
		"interface{}": true,
		"any":        true,
		"complex64":  true,
		"complex128": true,
	}

	// Убираем указатели и слайсы для проверки базового типа
	baseType := strings.TrimPrefix(typeName, "*")
	baseType = strings.TrimPrefix(baseType, "[]")
	baseType = strings.TrimPrefix(baseType, "map[string]")
	baseType = strings.TrimPrefix(baseType, "map[int]")

	return builtinTypes[baseType]
} 