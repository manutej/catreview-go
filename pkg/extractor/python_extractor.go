// Package extractor provides code analysis and categorical model extraction.
//
// PythonExtractor extracts categorical structures from Python source code using
// language-specific AST parsing while producing the same categorical model
// format as GoExtractor, ensuring uniformity across all language extractors.
package extractor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manu/catreview/pkg/category"
)

// PythonExtractor extracts categorical models from Python source code.
//
// Implementation Strategy:
//   - Use Python's ast module via subprocess or embedded Python interpreter
//   - Map Python constructs to categorical objects:
//     * Modules → Module objects
//     * Classes → Type objects
//     * Functions/Methods → Function objects
//   - Map Python dependencies to categorical morphisms:
//     * import/from statements → Import morphisms
//     * Function calls → Call morphisms
//     * Class inheritance → Inheritance morphisms
//     * Type hints → Type dependency morphisms
//
// This design ensures the same category.Category structure as GoExtractor,
// allowing all downstream analysis (complexity, functors, axioms) to work
// identically regardless of source language.
type PythonExtractor struct {
	category   *category.Category
	moduleMap  map[string]string // Maps file paths to module names
}

// NewPythonExtractor creates a new Python code extractor.
func NewPythonExtractor() *PythonExtractor {
	return &PythonExtractor{
		category:   category.NewCategory("python_codebase"),
		moduleMap:  make(map[string]string),
	}
}

// ExtractFromPath extracts categorical model from a Python project path.
//
// Implementation Plan (v1.1):
//   1. Walk directory tree finding .py files
//   2. For each .py file:
//      a. Parse using Python AST (ast.parse via subprocess or CGo)
//      b. Extract module docstring
//      c. Extract class definitions → Type objects
//      d. Extract function/method definitions → Function objects
//      e. Extract import statements → Import morphisms
//      f. Extract function calls → Call morphisms
//      g. Extract class inheritance → Inheritance morphisms
//      h. Extract type hints → Type dependency morphisms
//   3. Create identity morphisms for all objects
//   4. Verify category axioms (associativity, identity)
//
// Implements the Extractor interface.
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
			return fmt.Errorf("failed to extract from %s: %v", path, err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return e.category, nil
}

// extractFromFile extracts categorical structures from a single Python file.
//
// TODO (v1.1): Implement Python AST parsing.
// Options for Python parsing from Go:
//   1. Call python3 -m ast with JSON output and parse result
//   2. Use embedded Python via go-python/go-python3
//   3. Call Python script via exec.Command and parse structured output
//   4. Use tree-sitter-python parser (tree-sitter Go bindings)
//
// Recommended approach: exec.Command calling Python script that uses ast module
// and outputs JSON, providing clean separation and leveraging Python's built-in
// ast parsing capabilities.
//
// For now, this is a placeholder that creates minimal structure.
func (e *PythonExtractor) extractFromFile(filePath string) error {
	// TODO: Implement Python AST parsing
	// This skeleton creates a placeholder file object

	// Create file object
	fileObj := category.NewObject(
		filePath,
		"file",
		filepath.Base(filePath),
		map[string]interface{}{
			"language": "python",
			"path":     filePath,
			"status":   "skeleton_implementation",
		},
	)

	if err := e.category.AddObject(fileObj); err != nil {
		return err
	}

	// TODO: Extract module docstring
	// TODO: Extract import/from statements → Import morphisms
	// TODO: Extract class definitions → Type objects
	// TODO: Extract function/method definitions → Function objects
	// TODO: Extract function calls → Call morphisms
	// TODO: Extract class inheritance → Inheritance morphisms
	// TODO: Extract type hints → Type dependency morphisms

	return nil
}

// Language returns "python".
// Implements the Extractor interface.
func (e *PythonExtractor) Language() string {
	return "python"
}

// FileExtensions returns [".py"].
// Implements the Extractor interface.
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

func countUniqueModules(m map[string]string) int {
	unique := make(map[string]bool)
	for _, v := range m {
		unique[v] = true
	}
	return len(unique)
}

// Python-to-Category Mapping Reference
//
// Python Construct           → Categorical Object Type
// =======================================================
// # File: user.py            → Module object (ID: "user" or full import path)
// class User:                → Type object (ID: "user.User", kind: "class")
// def save(self):            → Function object (ID: "user.User.save")
// def process():             → Function object (ID: "user.process")
// self.name: str             → Field (metadata in class object)
//
// Python Dependency          → Categorical Morphism Type
// =======================================================
// import requests            → Import morphism (file → "requests")
// from typing import List    → Import morphism (file → "typing.List")
// requests.get(url)          → Call morphism (source → "requests.get")
// class User(Base):          → Inheritance morphism (User → Base, type: "inherits")
// def save(user: User):      → Type dependency morphism (save → User, type: "type_hint")
// self.repo.save()           → Call morphism (method → "repo.save")
//
// Special Python Considerations:
// - Dynamic typing: Track type hints when available, infer from usage otherwise
// - Duck typing: May not capture all runtime dependencies (document limitation)
// - __init__.py: Treat as package object with re-exports as morphisms
// - Decorators: Capture as metadata, create morphism to decorator function
// - Multiple inheritance: Create separate inheritance morphisms for each base
//
// This mapping ensures Python code produces identical categorical structure
// to Go/Java code, allowing unified analysis via the same complexity algorithms.
