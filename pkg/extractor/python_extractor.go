// Package extractor provides code analysis and categorical model extraction.
//
// PythonExtractor extracts categorical structures from Python source code using
// language-specific AST parsing while producing the same categorical model
// format as GoExtractor, ensuring uniformity across all language extractors.
package extractor

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/manu/catreview/pkg/category"
)

// PythonExtractor extracts categorical models from Python source code.
type PythonExtractor struct {
	category     *category.Category
	moduleMap    map[string]string // Maps file paths to module names
	parserScript string            // Path to Python AST parser script
}

// PythonParsedData represents the JSON output from Python AST parser.
type PythonParsedData struct {
	Module    string                   `json:"module"`
	File      string                   `json:"file"`
	Imports   []PythonImport           `json:"imports"`
	Classes   []PythonClass            `json:"classes"`
	Functions []PythonFunction         `json:"functions"`
	Calls     []PythonCall             `json:"calls"`
	Error     string                   `json:"error,omitempty"`
}

type PythonImport struct {
	Type   string  `json:"type"`   // "import" or "from_import"
	Name   string  `json:"name"`
	Module string  `json:"module,omitempty"` // For from_import
	Asname *string `json:"asname"`
	Line   int     `json:"line"`
}

type PythonClass struct {
	Name       string            `json:"name"`
	Bases      []string          `json:"bases"`
	Decorators []string          `json:"decorators"`
	Methods    []PythonFunction  `json:"methods"`
	Line       int               `json:"line"`
	Docstring  *string           `json:"docstring"`
}

type PythonFunction struct {
	Name       string              `json:"name"`
	Parameters []PythonParameter   `json:"parameters"`
	ReturnType *string             `json:"return_type"`
	Decorators []string            `json:"decorators"`
	Line       int                 `json:"line"`
	IsMethod   bool                `json:"is_method"`
	Async      bool                `json:"async"`
	Docstring  *string             `json:"docstring"`
}

type PythonParameter struct {
	Name string  `json:"name"`
	Type *string `json:"type"`
}

type PythonCall struct {
	Func            string  `json:"func"`
	Line            int     `json:"line"`
	ContextClass    *string `json:"context_class"`
	ContextFunction *string `json:"context_function"`
}

// NewPythonExtractor creates a new Python code extractor.
func NewPythonExtractor() *PythonExtractor {
	return &PythonExtractor{
		category:     category.NewCategory("python_codebase"),
		moduleMap:    make(map[string]string),
		parserScript: "tools/python-parser/parse_python.py",
	}
}

// ExtractFromPath extracts categorical model from a Python project path.
func (e *PythonExtractor) ExtractFromPath(root string) (*category.Category, error) {
	// Walk the directory tree
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Python files, test files, and __pycache__
		if info.IsDir() || !strings.HasSuffix(path, ".py") ||
			strings.Contains(path, "test_") || strings.Contains(path, "__pycache__") {
			return nil
		}

		// Parse the file
		if err := e.extractFromFile(path); err != nil {
			// Log error but continue - partial extraction better than failure
			fmt.Fprintf(os.Stderr, "Warning: failed to extract from %s: %v\n", path, err)
			return nil
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return e.category, nil
}

// extractFromFile extracts categorical structures from a single Python file.
func (e *PythonExtractor) extractFromFile(filePath string) error {
	// Call Python AST parser subprocess
	parsed, err := e.parsePythonFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	// Handle parse errors
	if parsed.Error != "" {
		return fmt.Errorf("python syntax error: %s", parsed.Error)
	}

	moduleName := parsed.Module
	e.moduleMap[filePath] = moduleName

	// Create module object
	moduleObj := category.NewObject(
		moduleName,
		"module",
		filepath.Base(filePath),
		map[string]interface{}{
			"language": "python",
			"file":     filePath,
		},
	)
	if err := e.category.AddObject(moduleObj); err != nil {
		return err
	}

	// Extract imports as morphisms
	for _, imp := range parsed.Imports {
		if err := e.extractImport(moduleName, imp); err != nil {
			// Continue on error - some imports might be external
			continue
		}
	}

	// Extract classes
	for _, class := range parsed.Classes {
		if err := e.extractClass(moduleName, filePath, class); err != nil {
			return err
		}
	}

	// Extract functions
	for _, function := range parsed.Functions {
		if err := e.extractFunction(moduleName, filePath, function); err != nil {
			return err
		}
	}

	// Extract calls (with context)
	for _, call := range parsed.Calls {
		e.extractCall(moduleName, call)
	}

	return nil
}

// parsePythonFile calls the Python AST parser script via subprocess.
func (e *PythonExtractor) parsePythonFile(filePath string) (*PythonParsedData, error) {
	cmd := exec.Command("python3", e.parserScript, filePath)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("python parser failed: %w", err)
	}

	var parsed PythonParsedData
	if err := json.Unmarshal(output, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %w", err)
	}

	return &parsed, nil
}

// extractImport creates a dependency morphism for an import.
func (e *PythonExtractor) extractImport(sourceModule string, imp PythonImport) error {
	var targetID string

	if imp.Type == "from_import" {
		// from module import name
		if imp.Module != "" {
			targetID = fmt.Sprintf("import:%s.%s", imp.Module, imp.Name)
		} else {
			targetID = fmt.Sprintf("import:%s", imp.Name)
		}
	} else {
		// import module
		targetID = fmt.Sprintf("import:%s", imp.Name)
	}

	// Create imported module object if it doesn't exist
	if _, exists := e.category.GetObject(targetID); !exists {
		impObj := category.NewObject(
			targetID,
			"imported_module",
			imp.Name,
			map[string]interface{}{
				"import_type": imp.Type,
				"module":      imp.Module,
			},
		)
		if err := e.category.AddObject(impObj); err != nil {
			return err
		}
	}

	// Create import dependency morphism
	morphID := fmt.Sprintf("import:%s->%s", sourceModule, targetID)
	if _, exists := e.category.GetMorphism(morphID); !exists {
		morph := category.NewMorphism(
			morphID,
			sourceModule,
			targetID,
			"import",
			map[string]interface{}{
				"import_type": imp.Type,
				"line":        imp.Line,
			},
		)
		if err := e.category.AddMorphism(morph); err != nil {
			return err
		}
	}

	return nil
}

// extractClass creates objects and morphisms for a Python class.
func (e *PythonExtractor) extractClass(moduleName, filePath string, class PythonClass) error {
	classID := fmt.Sprintf("%s.%s", moduleName, class.Name)

	// Create class object
	classObj := category.NewObject(
		classID,
		"class",
		class.Name,
		map[string]interface{}{
			"language":   "python",
			"module":     moduleName,
			"file":       filePath,
			"line":       class.Line,
			"bases":      class.Bases,
			"decorators": class.Decorators,
			"docstring":  getStringPtrValue(class.Docstring),
		},
	)
	if err := e.category.AddObject(classObj); err != nil {
		return err
	}

	// Create "defines" morphism from module to class
	morphID := fmt.Sprintf("defines:%s->%s", moduleName, classID)
	morph := category.NewMorphism(
		morphID,
		moduleName,
		classID,
		"defines",
		map[string]interface{}{
			"kind": "class",
		},
	)
	e.category.AddMorphism(morph)

	// Extract inheritance morphisms
	for _, base := range class.Bases {
		if base != "" {
			inheritMorphID := fmt.Sprintf("inherits:%s->%s", classID, base)
			if _, exists := e.category.GetMorphism(inheritMorphID); !exists {
				inheritMorph := category.NewMorphism(
					inheritMorphID,
					classID,
					base,
					"inheritance",
					map[string]interface{}{},
				)
				// Note: Might fail if base class doesn't exist yet
				e.category.AddMorphism(inheritMorph)
			}
		}
	}

	// Extract decorator morphisms
	for _, decorator := range class.Decorators {
		if decorator != "" {
			decorMorphID := fmt.Sprintf("decorator:%s->%s", classID, decorator)
			if _, exists := e.category.GetMorphism(decorMorphID); !exists {
				decorMorph := category.NewMorphism(
					decorMorphID,
					classID,
					decorator,
					"decorator",
					map[string]interface{}{},
				)
				e.category.AddMorphism(decorMorph)
			}
		}
	}

	// Extract methods
	for _, method := range class.Methods {
		methodID := fmt.Sprintf("%s.%s", classID, method.Name)
		methodObj := category.NewObject(
			methodID,
			"function",
			method.Name,
			map[string]interface{}{
				"language":    "python",
				"class":       classID,
				"file":        filePath,
				"line":        method.Line,
				"parameters":  method.Parameters,
				"return_type": getStringPtrValue(method.ReturnType),
				"decorators":  method.Decorators,
				"async":       method.Async,
				"docstring":   getStringPtrValue(method.Docstring),
			},
		)
		if err := e.category.AddObject(methodObj); err != nil {
			return err
		}

		// Create "defines" morphism from class to method
		methodMorphID := fmt.Sprintf("defines:%s->%s", classID, methodID)
		methodMorph := category.NewMorphism(
			methodMorphID,
			classID,
			methodID,
			"defines",
			map[string]interface{}{
				"kind": "method",
			},
		)
		e.category.AddMorphism(methodMorph)

		// Extract type hint dependencies from parameters
		for _, param := range method.Parameters {
			if param.Type != nil && *param.Type != "" {
				e.extractTypeHint(methodID, *param.Type, "param_type")
			}
		}

		// Extract return type dependency
		if method.ReturnType != nil && *method.ReturnType != "" {
			e.extractTypeHint(methodID, *method.ReturnType, "return_type")
		}
	}

	return nil
}

// extractFunction creates objects for a top-level Python function.
func (e *PythonExtractor) extractFunction(moduleName, filePath string, function PythonFunction) error {
	funcID := fmt.Sprintf("%s.%s", moduleName, function.Name)

	// Create function object
	funcObj := category.NewObject(
		funcID,
		"function",
		function.Name,
		map[string]interface{}{
			"language":    "python",
			"module":      moduleName,
			"file":        filePath,
			"line":        function.Line,
			"parameters":  function.Parameters,
			"return_type": getStringPtrValue(function.ReturnType),
			"decorators":  function.Decorators,
			"async":       function.Async,
			"docstring":   getStringPtrValue(function.Docstring),
		},
	)
	if err := e.category.AddObject(funcObj); err != nil {
		return err
	}

	// Create "defines" morphism from module to function
	morphID := fmt.Sprintf("defines:%s->%s", moduleName, funcID)
	morph := category.NewMorphism(
		morphID,
		moduleName,
		funcID,
		"defines",
		map[string]interface{}{
			"kind": "function",
		},
	)
	e.category.AddMorphism(morph)

	// Extract type hint dependencies from parameters
	for _, param := range function.Parameters {
		if param.Type != nil && *param.Type != "" {
			e.extractTypeHint(funcID, *param.Type, "param_type")
		}
	}

	// Extract return type dependency
	if function.ReturnType != nil && *function.ReturnType != "" {
		e.extractTypeHint(funcID, *function.ReturnType, "return_type")
	}

	return nil
}

// extractTypeHint creates a type dependency morphism.
func (e *PythonExtractor) extractTypeHint(sourceID, typeName, morphType string) {
	morphID := fmt.Sprintf("type_dep:%s->%s", sourceID, typeName)
	if _, exists := e.category.GetMorphism(morphID); !exists {
		morph := category.NewMorphism(
			morphID,
			sourceID,
			typeName,
			morphType,
			map[string]interface{}{
				"type": typeName,
			},
		)
		// Note: Might fail if type doesn't exist
		e.category.AddMorphism(morph)
	}
}

// extractCall creates a function call morphism.
func (e *PythonExtractor) extractCall(moduleName string, call PythonCall) {
	// Determine source context
	var sourceID string
	if call.ContextClass != nil && *call.ContextClass != "" {
		if call.ContextFunction != nil && *call.ContextFunction != "" {
			// Method call context
			sourceID = fmt.Sprintf("%s.%s.%s", moduleName, *call.ContextClass, *call.ContextFunction)
		} else {
			// Class context (unlikely)
			sourceID = fmt.Sprintf("%s.%s", moduleName, *call.ContextClass)
		}
	} else if call.ContextFunction != nil && *call.ContextFunction != "" {
		// Top-level function context
		sourceID = fmt.Sprintf("%s.%s", moduleName, *call.ContextFunction)
	} else {
		// Module-level code
		sourceID = moduleName
	}

	// Create call morphism
	targetFunc := call.Func
	morphID := fmt.Sprintf("calls:%s->%s:%d", sourceID, targetFunc, call.Line)
	if _, exists := e.category.GetMorphism(morphID); !exists {
		morph := category.NewMorphism(
			morphID,
			sourceID,
			targetFunc,
			"function_call",
			map[string]interface{}{
				"target": targetFunc,
				"line":   call.Line,
			},
		)
		// Note: Might fail if target doesn't exist
		e.category.AddMorphism(morph)
	}
}

// Language returns "python".
func (e *PythonExtractor) Language() string {
	return "python"
}

// FileExtensions returns [".py"].
func (e *PythonExtractor) FileExtensions() []string {
	return []string{".py"}
}

// Stats returns extraction statistics.
func (e *PythonExtractor) Stats() map[string]interface{} {
	stats := e.category.Stats()
	return map[string]interface{}{
		"objects":   stats["objects"],
		"morphisms": stats["morphisms"],
		"files":     len(e.moduleMap),
		"modules":   countUniqueModules(e.moduleMap),
		"language":  "python",
	}
}

// Helper functions

func countUniqueModules(m map[string]string) int {
	unique := make(map[string]bool)
	for _, v := range m {
		unique[v] = true
	}
	return len(unique)
}

func getStringPtrValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
