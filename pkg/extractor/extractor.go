package extractor

import (
	"catreview/pkg/category"
)

// Extractor is the language-agnostic interface for extracting categorical models.
// Each language implementation (Go, Java, Python, TypeScript, etc.) must satisfy this interface.
//
// The extractor's responsibility is to:
//  1. Parse source code using language-specific AST parsers
//  2. Map language constructs to categorical Objects (modules, classes, functions)
//  3. Map dependencies to categorical Morphisms (imports, calls, inheritance)
//  4. Create identity morphisms for each object
//  5. Return a complete Category that satisfies category axioms
//
// This design ensures the core Go implementation (analysis, functors, complexity)
// remains unchanged regardless of source language.
type Extractor interface {
	// ExtractFromPath extracts a categorical model from a source directory.
	// Returns a Category with Objects, Morphisms, and Identities populated.
	//
	// Parameters:
	//   root: Absolute or relative path to source directory
	//
	// Returns:
	//   *category.Category: Complete categorical model
	//   error: Parse errors, file access issues, or extraction failures
	ExtractFromPath(root string) (*category.Category, error)

	// Language returns the name of the source language (e.g., "go", "java", "python").
	// Used for logging and language detection.
	Language() string

	// FileExtensions returns the file extensions this extractor handles (e.g., [".go"], [".java"]).
	// Used for automatic language detection.
	FileExtensions() []string
}

// ExtractorFactory creates language-specific extractors based on detected language.
type ExtractorFactory struct {
	extractors map[string]Extractor
}

// NewExtractorFactory creates a factory with all available extractors registered.
func NewExtractorFactory() *ExtractorFactory {
	factory := &ExtractorFactory{
		extractors: make(map[string]Extractor),
	}

	// Register Go extractor (always available)
	factory.Register(&GoExtractor{})

	// Future extractors will be registered here:
	// factory.Register(&JavaExtractor{})    // v1.1
	// factory.Register(&PythonExtractor{})  // v1.1
	// factory.Register(&TypeScriptExtractor{}) // v1.2

	return factory
}

// Register adds an extractor to the factory.
func (f *ExtractorFactory) Register(e Extractor) {
	f.extractors[e.Language()] = e
}

// GetExtractor retrieves an extractor by language name.
// Returns nil if language is not supported.
func (f *ExtractorFactory) GetExtractor(language string) Extractor {
	return f.extractors[language]
}

// DetectLanguage attempts to detect the language of a source directory
// by examining file extensions.
//
// Returns the detected language name or empty string if unknown.
func (f *ExtractorFactory) DetectLanguage(root string) string {
	// Implementation would scan directory for file extensions
	// and match against registered extractors.
	// For now, this is a placeholder.
	return ""
}

// SupportedLanguages returns a list of all supported languages.
func (f *ExtractorFactory) SupportedLanguages() []string {
	languages := make([]string, 0, len(f.extractors))
	for lang := range f.extractors {
		languages = append(languages, lang)
	}
	return languages
}
