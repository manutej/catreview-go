// Package extractor provides code analysis and categorical model extraction.
//
// This package uses Go's AST parsing to extract categorical structures from code:
// - Packages, types, functions → Objects
// - Imports, calls, type references → Morphisms
package extractor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/manu/catreview/pkg/category"
)

// GoExtractor extracts categorical models from Go source code.
type GoExtractor struct {
	fset       *token.FileSet
	category   *category.Category
	packageMap map[string]string // Maps file paths to package names
}

// NewGoExtractor creates a new Go code extractor.
func NewGoExtractor() *GoExtractor {
	return &GoExtractor{
		fset:       token.NewFileSet(),
		category:   category.NewCategory("go_codebase"),
		packageMap: make(map[string]string),
	}
}

// ExtractFromPath extracts categorical model from a Go project path.
func (e *GoExtractor) ExtractFromPath(root string) (*category.Category, error) {
	// Walk the directory tree
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files and test files
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Parse the file
		if err := e.extractFromFile(path); err != nil {
			return fmt.Errorf("failed to extract from %s: %v", path, err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return e.category, nil
}

// extractFromFile extracts categorical structures from a single Go file.
func (e *GoExtractor) extractFromFile(filePath string) error {
	// Parse the file
	f, err := parser.ParseFile(e.fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	pkgName := f.Name.Name
	e.packageMap[filePath] = pkgName

	// Create file object
	fileObj := category.NewObject(
		filePath,
		"file",
		filepath.Base(filePath),
		map[string]interface{}{
			"package":  pkgName,
			"path":     filePath,
			"doc":      getDocComment(f.Doc),
			"imports":  len(f.Imports),
		},
	)
	if err := e.category.AddObject(fileObj); err != nil {
		return err
	}

	// Extract imports as morphisms
	for _, imp := range f.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		if err := e.extractImport(filePath, importPath); err != nil {
			// Continue on error - some imports might be external
			continue
		}
	}

	// Extract type declarations
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			e.extractGenDecl(filePath, pkgName, d)
		case *ast.FuncDecl:
			e.extractFuncDecl(filePath, pkgName, d)
		}
	}

	return nil
}

// extractImport creates a dependency morphism for an import.
func (e *GoExtractor) extractImport(sourceFile, importPath string) error {
	// For now, create an object for the imported package
	// In a full implementation, we'd resolve the import to actual files
	targetID := fmt.Sprintf("import:%s", importPath)

	// Create imported package object if it doesn't exist
	if _, exists := e.category.GetObject(targetID); !exists {
		impObj := category.NewObject(
			targetID,
			"imported_package",
			importPath,
			map[string]interface{}{
				"import_path": importPath,
			},
		)
		if err := e.category.AddObject(impObj); err != nil {
			return err
		}
	}

	// Create import dependency morphism
	morphID := fmt.Sprintf("import:%s->%s", sourceFile, importPath)
	if _, exists := e.category.GetMorphism(morphID); !exists {
		morph := category.NewMorphism(
			morphID,
			sourceFile,
			targetID,
			"import",
			map[string]interface{}{
				"import_path": importPath,
			},
		)
		if err := e.category.AddMorphism(morph); err != nil {
			return err
		}
	}

	return nil
}

// extractGenDecl extracts general declarations (types, consts, vars).
func (e *GoExtractor) extractGenDecl(filePath, pkgName string, decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			e.extractTypeSpec(filePath, pkgName, s, decl.Doc)
		}
	}
}

// extractTypeSpec extracts a type declaration.
func (e *GoExtractor) extractTypeSpec(filePath, pkgName string, spec *ast.TypeSpec, doc *ast.CommentGroup) {
	typeName := spec.Name.Name
	typeID := fmt.Sprintf("%s.%s", pkgName, typeName)

	// Determine type kind
	typeKind := "type"
	switch spec.Type.(type) {
	case *ast.StructType:
		typeKind = "struct"
	case *ast.InterfaceType:
		typeKind = "interface"
	}

	// Create type object
	typeObj := category.NewObject(
		typeID,
		typeKind,
		typeName,
		map[string]interface{}{
			"package":  pkgName,
			"file":     filePath,
			"doc":      getDocComment(doc),
		},
	)
	if err := e.category.AddObject(typeObj); err != nil {
		return
	}

	// Create "defines" morphism from file to type
	morphID := fmt.Sprintf("defines:%s->%s", filePath, typeID)
	morph := category.NewMorphism(
		morphID,
		filePath,
		typeID,
		"defines",
		map[string]interface{}{
			"kind": typeKind,
		},
	)
	e.category.AddMorphism(morph)

	// Extract struct field types as dependencies
	if structType, ok := spec.Type.(*ast.StructType); ok {
		e.extractStructFields(typeID, structType)
	}
}

// extractStructFields extracts field type dependencies from a struct.
func (e *GoExtractor) extractStructFields(structID string, structType *ast.StructType) {
	for _, field := range structType.Fields.List {
		// Extract type references
		e.extractTypeRef(structID, field.Type)
	}
}

// extractTypeRef extracts type references (dependencies on other types).
func (e *GoExtractor) extractTypeRef(sourceID string, expr ast.Expr) {
	switch t := expr.(type) {
	case *ast.Ident:
		// Reference to type in same package
		if t.Obj != nil {
			// Local type reference - would need to resolve
		}
	case *ast.SelectorExpr:
		// Reference to type in another package
		if ident, ok := t.X.(*ast.Ident); ok {
			targetType := fmt.Sprintf("%s.%s", ident.Name, t.Sel.Name)
			// Create dependency morphism
			morphID := fmt.Sprintf("uses:%s->%s", sourceID, targetType)
			if _, exists := e.category.GetMorphism(morphID); !exists {
				morph := category.NewMorphism(
					morphID,
					sourceID,
					targetType,
					"type_dependency",
					map[string]interface{}{
						"type": targetType,
					},
				)
				// Note: This might fail if target doesn't exist yet
				e.category.AddMorphism(morph)
			}
		}
	case *ast.StarExpr:
		// Pointer type
		e.extractTypeRef(sourceID, t.X)
	case *ast.ArrayType:
		// Array/slice type
		e.extractTypeRef(sourceID, t.Elt)
	case *ast.MapType:
		// Map type
		e.extractTypeRef(sourceID, t.Key)
		e.extractTypeRef(sourceID, t.Value)
	}
}

// extractFuncDecl extracts a function declaration.
func (e *GoExtractor) extractFuncDecl(filePath, pkgName string, decl *ast.FuncDecl) {
	funcName := decl.Name.Name

	// Determine if it's a method
	var funcID string
	if decl.Recv != nil && len(decl.Recv.List) > 0 {
		// It's a method - attach to receiver type
		recvType := exprToString(decl.Recv.List[0].Type)
		funcID = fmt.Sprintf("%s.%s.%s", pkgName, recvType, funcName)
	} else {
		// Regular function
		funcID = fmt.Sprintf("%s.%s", pkgName, funcName)
	}

	// Create function object
	funcObj := category.NewObject(
		funcID,
		"function",
		funcName,
		map[string]interface{}{
			"package": pkgName,
			"file":    filePath,
			"doc":     getDocComment(decl.Doc),
			"is_exported": ast.IsExported(funcName),
		},
	)
	if err := e.category.AddObject(funcObj); err != nil {
		return
	}

	// Create "defines" morphism from file to function
	morphID := fmt.Sprintf("defines:%s->%s", filePath, funcID)
	morph := category.NewMorphism(
		morphID,
		filePath,
		funcID,
		"defines",
		map[string]interface{}{
			"kind": "function",
		},
	)
	e.category.AddMorphism(morph)

	// Extract function calls from body
	if decl.Body != nil {
		ast.Inspect(decl.Body, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				e.extractFunctionCall(funcID, call)
			}
			return true
		})
	}
}

// extractFunctionCall extracts function call relationships.
func (e *GoExtractor) extractFunctionCall(sourceFunc string, call *ast.CallExpr) {
	var targetFunc string

	switch fun := call.Fun.(type) {
	case *ast.Ident:
		// Call to function in same package
		targetFunc = fun.Name
	case *ast.SelectorExpr:
		// Call to method or function in another package
		targetFunc = exprToString(fun)
	default:
		return
	}

	// Create call morphism
	morphID := fmt.Sprintf("calls:%s->%s", sourceFunc, targetFunc)
	if _, exists := e.category.GetMorphism(morphID); !exists {
		morph := category.NewMorphism(
			morphID,
			sourceFunc,
			targetFunc,
			"function_call",
			map[string]interface{}{
				"target": targetFunc,
			},
		)
		// Note: Might fail if target doesn't exist
		e.category.AddMorphism(morph)
	}
}

// Helper functions

func getDocComment(cg *ast.CommentGroup) string {
	if cg == nil {
		return ""
	}
	return cg.Text()
}

func exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return exprToString(e.X) + "." + e.Sel.Name
	case *ast.StarExpr:
		return "*" + exprToString(e.X)
	default:
		return ""
	}
}

// Stats returns extraction statistics.
func (e *GoExtractor) Stats() map[string]interface{} {
	stats := e.category.Stats()
	return map[string]interface{}{
		"objects":     stats["objects"],
		"morphisms":   stats["morphisms"],
		"files":       len(e.packageMap),
		"packages":    countUnique(e.packageMap),
	}
}

func countUnique(m map[string]string) int {
	unique := make(map[string]bool)
	for _, v := range m {
		unique[v] = true
	}
	return len(unique)
}
