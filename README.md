# ASTi - Go Interface Parser

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

**ASTi** - это библиотека для парсинга Go интерфейсов с поддержкой кастомных аннотаций. Библиотека позволяет извлекать
метаданные из Go кода, анализировать структуры интерфейсов и генерировать документацию или код на основе аннотаций.

## 🚀 Основные возможности

- **Парсинг Go интерфейсов** с полной поддержкой всех типов данных
- **Кастомные аннотации** с гибким синтаксисом
- **Pipeline архитектура** для расширяемой обработки
- **JSON сериализация** результатов парсинга
- **Поддержка дженериков** и сложных типов
- **Циклические структуры** и рекурсивные типы
- **Валидация и фильтрация** данных
- **OpenTelemetry** интеграция для observability

## 📦 Установка

```bash
go get github.com/seniorGolang/asti
```

## 🎯 Быстрый старт

### Базовое использование

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/seniorGolang/asti/parser"
)

func main() {
    // Создаем парсер с дефолтным префиксом аннотаций @asti
    p := parser.NewParser()

    // Парсим пакет
    ctx := context.Background()
    pkg, err := p.ParsePackage(ctx, "./examples/sample")
    if err != nil {
        log.Fatalf("Failed to parse package: %v", err)
    }

    // Выводим информацию об интерфейсах
    fmt.Printf("Found %d interfaces\n", len(pkg.Interfaces))
    for _, iface := range pkg.Interfaces {
        fmt.Printf("- %s: %d methods\n", iface.Name, len(iface.Methods))
    }
}
```

### Пример Go кода с аннотациями

```go
// @asti version=1.0 author="Team"
package examples

import "context"

// @asti name="UserService" timeout=30 retry=3
type UserService interface {
    // @asti method=CreateUser timeout=10
    CreateUser(ctx context.Context, name string, email string) (userID string, err error)

    // @asti method=GetUser timeout=5
    GetUser(ctx context.Context, userID string) (user User, err error)
}

// @asti type=UserModel validation=strict
type User struct {
    // @asti field=Name required=true maxLength=100
    Name string `json:"name"`
}
```

## 📚 Документация

### Структура проекта

```
ASTi/
├── parser/                  # Основной пакет парсера
│   ├── models/              # Модели данных
│   │   ├── annotation.go    # Аннотации и их парсинг
│   │   ├── interface.go     # Интерфейсы и методы
│   │   ├── position.go      # Позиции в коде
│   │   ├── result.go        # Результаты парсинга
│   │   └── type.go          # Типы данных
│   ├── pipeline/            # Pipeline обработки
│   │   ├── ast.go           # AST анализ
│   │   ├── filter.go        # Фильтрация
│   │   ├── pipeline.go      # Pipeline логика
│   │   ├── serialization.go # Сериализация
│   │   ├── types.go         # Типы pipeline
│   │   └── validation.go    # Валидация
│   ├── options.go           # Опции конфигурации
│   └── parser.go            # Основной API
├── examples/                # Примеры использования
│   ├── sample/              # Базовые примеры
│   ├── complex/             # Сложные сценарии
│   └── usage/               # Примеры использования
└── memory_bank/             # Документация проекта
```

### API Reference

#### Parser

```go
// NewParser создает новый парсер с опциональными настройками
func NewParser(options ...Option) *Parser

// ParsePackage парсит пакет и возвращает информацию об интерфейсах
func (p *Parser) ParsePackage(ctx context.Context, packagePath string) (*models.Package, error)

// ParsePackageToJSON парсит пакет и возвращает JSON
func (p *Parser) ParsePackageToJSON(ctx context.Context, packagePath string) ([]byte, error)

// ToJSON сериализует пакет в JSON
func (p *Parser) ToJSON(pkg *models.Package) ([]byte, error)

// FromJSON десериализует пакет из JSON
func (p *Parser) FromJSON(jsonData []byte) (*models.Package, error)
```

#### Опции

```go
// WithAnnotationPrefix устанавливает кастомный префикс аннотаций
func WithAnnotationPrefix(prefix string) Option
```

#### Модели данных

```go
type Package struct {
    ModulePath  string              `json:"modulePath"`
    ModuleName  string              `json:"moduleName"`
    PackagePath string              `json:"packagePath"`
    Annotations Annotations         `json:"annotations"`
    Interfaces  []Interface         `json:"interfaces"`
    Types       map[string]TypeInfo `json:"types"`
}

type Interface struct {
    Name        string      `json:"name"`
    Package     string      `json:"package"`
    Methods     []Method    `json:"methods"`
    Annotations Annotations `json:"annotations,omitempty"`
    Position    Position    `json:"position"`
}

type Method struct {
    Name         string      `json:"name"`
    Parameters   []Variable  `json:"parameters"`
    Results      []Variable  `json:"results"`
    Annotations  Annotations `json:"annotations,omitempty"`
    Position     Position    `json:"position"`
    Serializable bool        `json:"serializable"`
}
```

### Аннотации

Библиотека поддерживает кастомные аннотации с гибким синтаксисом:

```go
// Базовый синтаксис
// @asti key=value

// Строковые значения в кавычках
// @asti name="UserService" version="1.0"

// Множественные аннотации
// @asti timeout=30 retry=3 validation=strict

// Аннотации на разных уровнях
// @asti package=main version=1.0
type MyInterface interface {
    // @asti method=Create timeout=10
    Create(ctx context.Context, data string) (result string, err error)
}
```

## 🔧 Конфигурация

### Настройка префикса аннотаций

```go
// Использование дефолтного префикса @asti
parser := parser.NewParser()

// Использование кастомного префикса
parser := parser.NewParser(parser.WithAnnotationPrefix("@custom"))
```

### Pipeline конфигурация

```go
// Создание кастомного pipeline
pipeline := pipeline.NewPipeline(
    pipeline.NewStageAST(annotationParser),
    pipeline.NewStageFilter(),
    pipeline.NewStageTypeCollection(annotationParser),
    pipeline.NewStageSerialization(),
)
```

## 🤝 Вклад в проект

Вклад в развитие проекта приветствуется!

### Требования к коду

- Следуйте [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Используйте `go fmt` для форматирования
- Добавляйте тесты для новой функциональности
- Обновляйте документацию при изменении API

### Процесс разработки

1. Форкните репозиторий
2. Создайте ветку для новой функции
3. Внесите изменения с тестами
4. Убедитесь, что все тесты проходят
5. Создайте Pull Request

## 📄 Лицензия

Этот проект лицензирован под Apache License 2.0 - см. файл [LICENSE](LICENSE) для деталей.

## 🆘 Поддержка

- **Issues**: [GitHub Issues](https://github.com/seniorGolang/asti/issues)
- **Discussions**: [GitHub Discussions](https://github.com/seniorGolang/asti/discussions)
- **Documentation**: [pkg.go.dev](https://pkg.go.dev/github.com/seniorGolang/asti)
