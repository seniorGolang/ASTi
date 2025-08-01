package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/seniorGolang/asti/parser/models"
)

func TestStageModule_Process(t *testing.T) {
	// Создаем временную директорию для теста
	tempDir, err := os.MkdirTemp("", "module_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Создаем go.mod файл в корне временной директории
	goModContent := `module github.com/test/module

go 1.24

require (
	github.com/some/dependency v1.0.0
)
`
	goModPath := filepath.Join(tempDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goModContent), 0600); err != nil {
		t.Fatalf("Failed to write go.mod file: %v", err)
	}

	// Создаем поддиректорию для пакета
	packageDir := filepath.Join(tempDir, "pkg", "service")
	if err := os.MkdirAll(packageDir, 0755); err != nil {
		t.Fatalf("Failed to create package directory: %v", err)
	}

	// Создаем тестовый Go файл
	goFileContent := `package service

// @asti service:user
type UserService interface {
	GetUser(id string) (*User, error)
}

type User struct {
	ID   string ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`
	goFilePath := filepath.Join(packageDir, "service.go")
	if err := os.WriteFile(goFilePath, []byte(goFileContent), 0600); err != nil {
		t.Fatalf("Failed to write Go file: %v", err)
	}

	// Тестируем StageModule
	data := Data{
		Package: &models.Package{
			PackagePath: packageDir,
		},
	}

	stage := NewStageModule()
	result, err := stage.Process(context.Background(), data)
	if err != nil {
		t.Fatalf("StageModule.Process failed: %v", err)
	}

	// Проверяем результаты
	if result.Package.ModuleName != "github.com/test/module" {
		t.Errorf("Expected module name 'github.com/test/module', got '%s'", result.Package.ModuleName)
	}

	expectedPackagePath := "pkg/service"
	if result.Package.PackagePath != expectedPackagePath {
		t.Errorf("Expected package path '%s', got '%s'", expectedPackagePath, result.Package.PackagePath)
	}
}

func TestStageModule_findModuleInfo(t *testing.T) {
	stage := NewStageModule()

	// Тест с несуществующей директорией
	_, _, err := stage.findModuleInfo("/non/existent/path")
	if err == nil {
		t.Error("Expected error for non-existent path, got nil")
	}
}

func TestStageModule_parseGoMod(t *testing.T) {

	// Создаем временный go.mod файл
	tempDir, err := os.MkdirTemp("", "gomod_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	goModPath := filepath.Join(tempDir, "go.mod")
	goModContent := `module github.com/test/example

go 1.24

require (
	github.com/some/dependency v1.0.0
)
`
	if err := os.WriteFile(goModPath, []byte(goModContent), 0600); err != nil {
		t.Fatalf("Failed to write go.mod file: %v", err)
	}

	moduleName, err := ParseGoMod(goModPath)
	if err != nil {
		t.Fatalf("ParseGoMod failed: %v", err)
	}

	if moduleName != "github.com/test/example" {
		t.Errorf("Expected module name 'github.com/test/example', got '%s'", moduleName)
	}
}

func TestStageModule_parseGoModWithComments(t *testing.T) {

	// Создаем временный go.mod файл с комментариями
	tempDir, err := os.MkdirTemp("", "gomod_comment_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	goModPath := filepath.Join(tempDir, "go.mod")
	goModContent := `module github.com/test/example // This is a comment

go 1.24
`
	if err := os.WriteFile(goModPath, []byte(goModContent), 0600); err != nil {
		t.Fatalf("Failed to write go.mod file: %v", err)
	}

	moduleName, err := ParseGoMod(goModPath)
	if err != nil {
		t.Fatalf("ParseGoMod failed: %v", err)
	}

	if moduleName != "github.com/test/example" {
		t.Errorf("Expected module name 'github.com/test/example', got '%s'", moduleName)
	}
}

func TestStageModule_parseGoModWithoutModule(t *testing.T) {

	// Создаем временный go.mod файл без объявления модуля
	tempDir, err := os.MkdirTemp("", "gomod_no_module_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	goModPath := filepath.Join(tempDir, "go.mod")
	goModContent := `go 1.24

require (
	github.com/some/dependency v1.0.0
)
`
	if err := os.WriteFile(goModPath, []byte(goModContent), 0600); err != nil {
		t.Fatalf("Failed to write go.mod file: %v", err)
	}

	_, err = ParseGoMod(goModPath)
	if err == nil {
		t.Error("Expected error for go.mod without module declaration, got nil")
	}
}
