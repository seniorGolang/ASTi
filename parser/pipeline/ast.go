package pipeline

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/seniorGolang/asti/parser/models"
)

type StageAST struct {
	annotationParser models.AnnotationParser
	currentTypes     map[string]bool // типы, определенные в текущем файле
	packageName      string          // имя текущего пакета
}

func NewStageAST(annotationParser models.AnnotationParser) (stage *StageAST) {

	stage = &StageAST{annotationParser: annotationParser}
	return
}

// Process выполняет парсинг AST
func (s *StageAST) Process(ctx context.Context, data Data) (result Data, err error) {

	if data.Package == nil {
		err = fmt.Errorf("package data is required for AST parsing")
		return
	}
	// Используем абсолютный путь для поиска файлов, если он доступен
	packagePath := data.Package.PackagePath
	if data.Annotations != nil {
		if absPathData, exists := data.Annotations["_absolutePackagePath"]; exists {
			if absPath, ok := absPathData["path"]; ok {
				packagePath = absPath
			}
		}
	}

	var files []string
	files, err = filepath.Glob(filepath.Join(packagePath, "*.go"))
	if err != nil {
		err = fmt.Errorf("failed to find Go files: %w", err)
		return
	}
	
	// Если нет Go файлов, возвращаем пустой результат
	if len(files) == 0 {
		// Для пустых директорий возвращаем пустой результат без ошибки
		data.Interfaces = []models.Interface{}
		data.Package.Annotations = make(models.Annotations)
		result = data
		return
	}

	// Сначала собираем все типы из всех файлов пакета
	s.currentTypes = make(map[string]bool)
	fset := token.NewFileSet()
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		var astFile *ast.File
		astFile, err = parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			err = fmt.Errorf("failed to parse file %s: %w", file, err)
			return
		}
		// Собираем типы из текущего файла
		for _, decl := range astFile.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok {
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						s.currentTypes[typeSpec.Name.Name] = true
					}
				}
			}
		}
		// Устанавливаем имя пакета из первого файла
		if s.packageName == "" && astFile.Name != nil {
			s.packageName = astFile.Name.Name
		}
	}

	var interfaces []models.Interface
	var packageAnnotations models.Annotations
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		var astFile *ast.File
		astFile, err = parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			err = fmt.Errorf("failed to parse file %s: %w", file, err)
			return
		}
		var fileInterfaces []models.Interface
		var filePackageAnnotations models.Annotations
		// Получаем имя модуля для формирования полного пути импорта
		moduleName := ""
		if data.Package != nil {
			moduleName = data.Package.ModuleName
		}
		fileInterfaces, filePackageAnnotations, err = s.extractInterfaces(ctx, astFile, fset, file, packagePath, moduleName)
		if err != nil {
			err = fmt.Errorf("failed to extract interfaces from %s: %w", file, err)
			return
		}
		interfaces = append(interfaces, fileInterfaces...)
		if filePackageAnnotations != nil {
			if packageAnnotations == nil {
				packageAnnotations = make(models.Annotations)
			}
			for k, v := range filePackageAnnotations {
				packageAnnotations[k] = v
			}
		}
	}
	data.Interfaces = interfaces
	data.Package.Annotations = packageAnnotations
	result = data
	return
}

func (s *StageAST) hasAnnotation(commentText string) bool {

	commentText = strings.TrimSpace(commentText)
	if !strings.HasPrefix(commentText, "//") {
		return false
	}
	content := strings.TrimSpace(strings.TrimPrefix(commentText, "//"))
	if !strings.HasPrefix(content, "@") {
		return false
	}
	prefix := strings.TrimPrefix(s.annotationParser.(*models.DefaultAnnotationParser).GetPrefix(), "@")
	return strings.HasPrefix(content[1:], prefix)
}

func (s *StageAST) extractInterfaces(ctx context.Context, astFile *ast.File, fset *token.FileSet, filename string, packagePath string, moduleName string) (interfaces []models.Interface, packageAnnotations models.Annotations, err error) {

	packageAnnotations = make(models.Annotations)
	if astFile.Doc != nil {
		for _, comment := range astFile.Doc.List {
			if s.hasAnnotation(comment.Text) {
				var annotations models.Annotations
				annotations, err = s.annotationParser.Parse(ctx, comment.Text)
				if err != nil {
					continue
				}
				for k, v := range annotations {
					packageAnnotations[k] = v
				}
			}
		}
	}
	for _, decl := range astFile.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
						var interfaceAnnotations models.Annotations
						if genDecl.Doc != nil {
							for _, comment := range genDecl.Doc.List {
								if s.hasAnnotation(comment.Text) {
									var annotations models.Annotations
									annotations, err = s.annotationParser.Parse(ctx, comment.Text)
									if err == nil {
										if interfaceAnnotations == nil {
											interfaceAnnotations = make(models.Annotations)
										}
										for k, v := range annotations {
											interfaceAnnotations[k] = v
										}
									}
								}
							}
						}
						if len(interfaceAnnotations) > 0 {
							pos := fset.Position(typeSpec.Pos())
							var relativePath string
							relativePath, err = filepath.Rel(packagePath, filename)
							if err != nil {
								relativePath = filename
							}

							// Определяем путь импорта для текущего пакета
							importPath := packagePath
							// Если это абсолютный путь, извлекаем относительный путь от корня модуля
							if filepath.IsAbs(importPath) {
								// Используем общую функцию для получения информации о модуле
								if modPath, err := FindModuleRoot(importPath); err == nil {
									if relPath, err := filepath.Rel(modPath, importPath); err == nil {
										importPath = relPath
									}
								}
							}

							// Формируем полный путь импорта с именем модуля
							fullImportPath := importPath
							if moduleName != "" {
								fullImportPath = moduleName + "/" + importPath
							}

							iface := models.Interface{
								Name:        typeSpec.Name.Name,
								Package:     astFile.Name.Name,
								Import:      fullImportPath,
								Annotations: interfaceAnnotations,
								Position: models.Position{
									File:   relativePath,
									Line:   pos.Line,
									Column: pos.Column,
								},
								ID: fmt.Sprintf("%s.%s", astFile.Name.Name, typeSpec.Name.Name),
							}

							var methods []models.Method
							methods, err = s.extractMethods(ctx, interfaceType, fset, filename, packagePath)
							if err != nil {
								err = fmt.Errorf("failed to extract methods: %w", err)
								return
							}
							iface.Methods = methods

							interfaces = append(interfaces, iface)
						}
					}
				}
			}
		}
	}
	return
}

func (s *StageAST) extractMethods(ctx context.Context, interfaceType *ast.InterfaceType, fset *token.FileSet, filename string, packagePath string) (methods []models.Method, err error) {

	methodAnnotations := make(map[string]models.Annotations)
	for _, field := range interfaceType.Methods.List {
		if _, ok := field.Type.(*ast.FuncType); ok {
			if len(field.Names) > 0 {
				methodName := field.Names[0].Name

				if field.Doc != nil {
					for _, comment := range field.Doc.List {
						if s.hasAnnotation(comment.Text) {
							var annotations models.Annotations
							annotations, err = s.annotationParser.Parse(ctx, comment.Text)
							if err == nil {
								if methodAnnotations[methodName] == nil {
									methodAnnotations[methodName] = make(models.Annotations)
								}
								for k, v := range annotations {
									methodAnnotations[methodName][k] = v
								}
							}
						}
					}
				}
			}
		}
	}
	for _, field := range interfaceType.Methods.List {
		if funcType, ok := field.Type.(*ast.FuncType); ok {
			if len(field.Names) > 0 {
				methodName := field.Names[0].Name
				pos := fset.Position(field.Pos())

				var relativePath string
				relativePath, err = filepath.Rel(packagePath, filename)
				if err != nil {
					relativePath = filename
				}

				method := models.Method{
					MethodInfo: models.MethodInfo{
						Name: methodName,
						Position: models.Position{
							File:   relativePath,
							Line:   pos.Line,
							Column: pos.Column,
						},
						Annotations: methodAnnotations[methodName],
					},
					ID: methodName,
				}

				var parameters []models.Variable
				parameters, err = s.extractVariables(funcType.Params)
				if err != nil {
					err = fmt.Errorf("failed to extract parameters: %w", err)
					return
				}
				method.Parameters = parameters

				var results []models.Variable
				results, err = s.extractVariables(funcType.Results)
				if err != nil {
					err = fmt.Errorf("failed to extract results: %w", err)
					return
				}
				method.Results = results

				methods = append(methods, method)
			}
		}
	}
	return
}

func (s *StageAST) extractVariables(fieldList *ast.FieldList) (variables []models.Variable, err error) {

	if fieldList == nil {
		return
	}
	for _, field := range fieldList.List {
		typeStr := s.typeToString(field.Type)
		if len(field.Names) > 0 {
			for _, name := range field.Names {
				variable := models.Variable{
					Name: name.Name,
					Type: typeStr,
				}
				s.analyzeTypeCharacteristics(field.Type, &variable)
				variables = append(variables, variable)
			}
		} else {
			variable := models.Variable{
				Name: "",
				Type: typeStr,
			}
			s.analyzeTypeCharacteristics(field.Type, &variable)
			variables = append(variables, variable)
		}
	}
	return
}

func (s *StageAST) typeToString(expr ast.Expr) (typeStr string) {

	switch t := expr.(type) {
	case *ast.Ident:
		// Если это тип из текущего файла, добавляем префикс пакета
		if s.currentTypes[t.Name] && s.packageName != "" {
			typeStr = s.packageName + "." + t.Name
		} else {
			typeStr = t.Name
		}
	case *ast.StarExpr:
		// Для указателей возвращаем базовый тип, указатель будет в analyzeTypeCharacteristics
		typeStr = s.typeToString(t.X)
	case *ast.ArrayType:
		// Для массивов возвращаем базовый тип, характеристики будут в analyzeTypeCharacteristics
		typeStr = s.typeToString(t.Elt)
	case *ast.SelectorExpr:
		typeStr = s.typeToString(t.X) + "." + t.Sel.Name
	case *ast.InterfaceType:
		typeStr = "interface{}"
	case *ast.MapType:
		// Для карт возвращаем базовый тип, характеристики будут в analyzeTypeCharacteristics
		typeStr = s.typeToString(t.Value)
	case *ast.ChanType:
		// Для каналов возвращаем базовый тип, характеристики будут в analyzeTypeCharacteristics
		typeStr = s.typeToString(t.Value)
	case *ast.Ellipsis:
		// Для variadic возвращаем базовый тип, характеристики будут в analyzeTypeCharacteristics
		typeStr = s.typeToString(t.Elt)
	case *ast.IndexExpr:
		// Для дженерик типов возвращаем базовый тип, параметры будут в analyzeTypeCharacteristics
		typeStr = s.typeToString(t.X)
	case *ast.IndexListExpr:
		// Для дженерик типов возвращаем базовый тип, параметры будут в analyzeTypeCharacteristics
		typeStr = s.typeToString(t.X)
	default:
		typeStr = "unknown"
	}
	return
}

func (s *StageAST) analyzeTypeCharacteristics(expr ast.Expr, variable *models.Variable) {

	switch t := expr.(type) {
	case *ast.StarExpr:
		variable.Pointer = true
		s.analyzeTypeCharacteristics(t.X, variable)
	case *ast.ArrayType:
		if t.Len == nil {
			variable.Slice = true
		} else {
			variable.Array = true
		}
		s.analyzeTypeCharacteristics(t.Elt, variable)
	case *ast.MapType:
		variable.Map = true
		s.analyzeTypeCharacteristics(t.Value, variable)
	case *ast.ChanType:
		variable.Channel = true
		s.analyzeTypeCharacteristics(t.Value, variable)
	case *ast.Ellipsis:
		variable.Variadic = true
		s.analyzeTypeCharacteristics(t.Elt, variable)
	case *ast.IndexExpr:
		variable.Generic = true
		s.analyzeTypeCharacteristics(t.X, variable)
	case *ast.IndexListExpr:
		variable.Generic = true
		s.analyzeTypeCharacteristics(t.X, variable)
	}
}
