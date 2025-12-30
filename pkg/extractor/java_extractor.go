// Package extractor provides code analysis and categorical model extraction.
//
// JavaExtractor extracts categorical structures from Java source code using
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

// JavaExtractor extracts categorical models from Java source code.
//
// Implementation Strategy:
//   - Use external Java parser (javaparser, Eclipse JDT, or javac AST)
//   - Map Java constructs to categorical objects:
//     * Packages → Package objects
//     * Classes/Interfaces → Type objects
//     * Methods → Function objects
//   - Map Java dependencies to categorical morphisms:
//     * import statements → Import morphisms
//     * Method calls → Call morphisms
//     * Class inheritance/implements → Inheritance morphisms
//     * Field type references → Type dependency morphisms
//
// This design ensures the same category.Category structure as GoExtractor,
// allowing all downstream analysis (complexity, functors, axioms) to work
// identically regardless of source language.
type JavaExtractor struct {
	category   *category.Category
	packageMap map[string]string // Maps file paths to package names
}

// NewJavaExtractor creates a new Java code extractor.
func NewJavaExtractor() *JavaExtractor {
	return &JavaExtractor{
		category:   category.NewCategory("java_codebase"),
		packageMap: make(map[string]string),
	}
}

// ExtractFromPath extracts categorical model from a Java project path.
//
// Implementation Plan (v1.1):
//   1. Walk directory tree finding .java files
//   2. For each .java file:
//      a. Parse using Java AST parser
//      b. Extract package declaration
//      c. Extract class/interface declarations → Objects
//      d. Extract import statements → Import morphisms
//      e. Extract method declarations → Function objects
//      f. Extract method calls → Call morphisms
//      g. Extract inheritance/implements → Inheritance morphisms
//   3. Create identity morphisms for all objects
//   4. Verify category axioms (associativity, identity)
//
// Implements the Extractor interface.
func (e *JavaExtractor) ExtractFromPath(root string) (*category.Category, error) {
	// Walk the directory tree
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Java files and test files
		if info.IsDir() || !strings.HasSuffix(path, ".java") || strings.Contains(path, "/test/") {
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

// extractFromFile extracts categorical structures from a single Java file.
//
// TODO (v1.1): Implement Java AST parsing.
// Options for Java parsing from Go:
//   1. Call javac with -Xprint flag and parse output
//   2. Use javaparser via exec/subprocess
//   3. Call Eclipse JDT Core via JNI/CGo
//   4. Parse Java grammar using go parser generator (antlr4, participle)
//
// For now, this is a placeholder that creates minimal structure.
func (e *JavaExtractor) extractFromFile(filePath string) error {
	// TODO: Implement Java AST parsing
	// This skeleton creates a placeholder file object

	// Create file object
	fileObj := category.NewObject(
		filePath,
		"file",
		filepath.Base(filePath),
		map[string]interface{}{
			"language": "java",
			"path":     filePath,
			"status":   "skeleton_implementation",
		},
	)

	if err := e.category.AddObject(fileObj); err != nil {
		return err
	}

	// TODO: Extract package declaration
	// TODO: Extract import statements → Import morphisms
	// TODO: Extract class/interface declarations → Type objects
	// TODO: Extract method declarations → Function objects
	// TODO: Extract method calls → Call morphisms
	// TODO: Extract inheritance/implements → Inheritance morphisms

	return nil
}

// Language returns "java".
// Implements the Extractor interface.
func (e *JavaExtractor) Language() string {
	return "java"
}

// FileExtensions returns [".java"].
// Implements the Extractor interface.
func (e *JavaExtractor) FileExtensions() []string {
	return []string{".java"}
}

// Stats returns extraction statistics.
func (e *JavaExtractor) Stats() map[string]interface{} {
	stats := e.category.Stats()
	return map[string]interface{}{
		"objects":   stats["objects"],
		"morphisms": stats["morphisms"],
		"files":     len(e.packageMap),
		"packages":  countUniquePackages(e.packageMap),
		"language":  "java",
	}
}

func countUniquePackages(m map[string]string) int {
	unique := make(map[string]bool)
	for _, v := range m {
		unique[v] = true
	}
	return len(unique)
}

// Java-to-Category Mapping Reference
//
// Java Construct          → Categorical Object Type
// ====================================================
// package com.example     → Package object (ID: "com.example")
// class User              → Type object (ID: "com.example.User", kind: "class")
// interface Repository    → Type object (ID: "com.example.Repository", kind: "interface")
// public void save()      → Function object (ID: "com.example.User.save")
// private String name     → Field (metadata in class object)
//
// Java Dependency         → Categorical Morphism Type
// ====================================================
// import com.foo.Bar      → Import morphism (file → "com.foo.Bar")
// bar.doSomething()       → Call morphism (source → target function)
// class User extends Base → Inheritance morphism (User → Base, type: "extends")
// implements Serializable → Inheritance morphism (User → Serializable, type: "implements")
// Field field             → Type dependency morphism (User → Field, type: "field_type")
//
// This mapping ensures Java code produces identical categorical structure
// to Go code, allowing unified analysis via the same complexity algorithms.
