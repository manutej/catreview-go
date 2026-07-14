package extractor

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/manu/catreview/pkg/category"
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

	// Register built-in extractors. Each is constructed via its New*
	// constructor so internal state (categories, index maps) is initialised.
	factory.Register(NewGoExtractor())
	factory.Register(NewJavaExtractor())

	// Future extractors will be registered here:
	// factory.Register(NewPythonExtractor())     // v1.2
	// factory.Register(NewTypeScriptExtractor()) // v1.2

	return factory
}

// extensionToLanguage maps a file extension to the language handled by a
// registered extractor. Built lazily from the registered extractors.
func (f *ExtractorFactory) extensionToLanguage() map[string]string {
	m := make(map[string]string)
	for lang, e := range f.extractors {
		for _, ext := range e.FileExtensions() {
			m[ext] = lang
		}
	}
	return m
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

// DetectLanguage detects the primary language of a source directory by counting
// files per registered extension and returning the language with the most files.
// Common build/output and VCS directories are skipped so vendored or generated
// files do not skew detection.
//
// Returns the detected language name, or an empty string if no registered
// language's files are found.
func (f *ExtractorFactory) DetectLanguage(root string) string {
	ext2lang := f.extensionToLanguage()
	counts := make(map[string]int)

	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Ignore unreadable entries; detection is best-effort.
		}
		if info.IsDir() {
			switch info.Name() {
			case ".git", "node_modules", "vendor", "target", "build", "out", "dist":
				return filepath.SkipDir
			}
			return nil
		}
		if lang, ok := ext2lang[strings.ToLower(filepath.Ext(path))]; ok {
			counts[lang]++
		}
		return nil
	})

	best, bestCount := "", 0
	for lang, c := range counts {
		if c > bestCount {
			best, bestCount = lang, c
		}
	}
	return best
}

// ExtractorForPath detects the language of root and returns the matching
// extractor, or nil and the detected (possibly empty) language if unsupported.
func (f *ExtractorFactory) ExtractorForPath(root string) (Extractor, string) {
	lang := f.DetectLanguage(root)
	if lang == "" {
		return nil, ""
	}
	return f.GetExtractor(lang), lang
}

// SupportedLanguages returns a list of all supported languages.
func (f *ExtractorFactory) SupportedLanguages() []string {
	languages := make([]string, 0, len(f.extractors))
	for lang := range f.extractors {
		languages = append(languages, lang)
	}
	return languages
}
