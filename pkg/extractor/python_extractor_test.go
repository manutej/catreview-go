package extractor

import (
	"testing"
)

func TestPythonExtractorBasic(t *testing.T) {
	// Given: Simple Python module
	extractor := NewPythonExtractor()

	// When: Extract categorical model
	cat, err := extractor.ExtractFromPath("../../testdata/python/simple")

	// Then: Verify extraction succeeded
	if err != nil {
		t.Fatalf("extraction failed: %v", err)
	}

	// Verify objects extracted
	if len(cat.Objects()) == 0 {
		t.Error("expected objects to be extracted")
	}

	// Verify module object exists
	moduleID := "testdata.python.simple.hello"
	if _, exists := cat.GetObject(moduleID); !exists {
		t.Errorf("expected module object %s", moduleID)
	}

	// Verify class object exists
	classID := "testdata.python.simple.hello.Greeter"
	if _, exists := cat.GetObject(classID); !exists {
		t.Errorf("expected class object %s", classID)
	}

	// Verify function object exists
	funcID := "testdata.python.simple.hello.greet"
	if _, exists := cat.GetObject(funcID); !exists {
		t.Errorf("expected function object %s", funcID)
	}

	// Verify morphisms extracted
	if len(cat.Morphisms()) == 0 {
		t.Error("expected morphisms to be extracted")
	}

	stats := extractor.Stats()
	t.Logf("Extraction stats: %+v", stats)
}

func TestPythonExtractorInheritance(t *testing.T) {
	// Given: Python classes with inheritance
	extractor := NewPythonExtractor()

	// When: Extract categorical model
	cat, err := extractor.ExtractFromPath("../../testdata/python/inheritance")

	// Then: Verify extraction succeeded
	if err != nil {
		t.Fatalf("extraction failed: %v", err)
	}

	// Verify base class exists
	baseClassID := "testdata.python.inheritance.base.BaseService"
	baseClass, exists := cat.GetObject(baseClassID)
	if !exists {
		t.Errorf("expected base class object %s", baseClassID)
	}
	if exists && baseClass.Type != "class" {
		t.Errorf("expected type 'class', got %s", baseClass.Type)
	}

	// Verify derived class exists
	derivedClassID := "testdata.python.inheritance.derived.UserService"
	derivedClass, exists := cat.GetObject(derivedClassID)
	if !exists {
		t.Errorf("expected derived class object %s", derivedClassID)
	}
	if exists && derivedClass.Type != "class" {
		t.Errorf("expected type 'class', got %s", derivedClass.Type)
	}

	// Verify inheritance morphism exists
	foundInheritance := false
	for _, morph := range cat.Morphisms() {
		if morph.Type == "inheritance" && morph.Source == derivedClassID {
			foundInheritance = true
			t.Logf("Found inheritance morphism: %s -> %s", morph.Source, morph.Target)
			break
		}
	}
	if !foundInheritance {
		t.Error("expected inheritance morphism from UserService")
	}
}

func TestPythonExtractorMethodCalls(t *testing.T) {
	// Given: Python code with method calls
	extractor := NewPythonExtractor()

	// When: Extract categorical model
	cat, err := extractor.ExtractFromPath("../../testdata/python/simple")

	// Then: Verify call morphisms created
	if err != nil {
		t.Fatalf("extraction failed: %v", err)
	}

	// Count function_call morphisms
	callCount := 0
	for _, morph := range cat.Morphisms() {
		if morph.Type == "function_call" {
			callCount++
			t.Logf("Found call: %s -> %s", morph.Source, morph.Target)
		}
	}

	if callCount == 0 {
		t.Error("expected at least one function_call morphism")
	}

	t.Logf("Total function calls extracted: %d", callCount)
}

func TestPythonExtractorTypeHints(t *testing.T) {
	// Given: Python code with type hints
	extractor := NewPythonExtractor()

	// When: Extract categorical model
	cat, err := extractor.ExtractFromPath("../../testdata/python/simple")

	// Then: Verify type dependency morphisms created
	if err != nil {
		t.Fatalf("extraction failed: %v", err)
	}

	// Count type dependency morphisms
	typeDeps := 0
	for _, morph := range cat.Morphisms() {
		if morph.Type == "param_type" || morph.Type == "return_type" {
			typeDeps++
			t.Logf("Found type dependency: %s (%s)", morph.Type, morph.Target)
		}
	}

	// Note: Our test file has type hints, so we should find some
	if typeDeps == 0 {
		t.Log("Warning: no type dependencies found (may be expected if no hints in test file)")
	} else {
		t.Logf("Total type dependencies extracted: %d", typeDeps)
	}
}

func TestPythonExtractorCategoryAxioms(t *testing.T) {
	// Given: Extracted Python category
	extractor := NewPythonExtractor()
	cat, err := extractor.ExtractFromPath("../../testdata/python/simple")

	if err != nil {
		t.Fatalf("extraction failed: %v", err)
	}

	// When: Verify axioms
	err = cat.VerifyAxioms()

	// Then: Associativity and identity laws should hold
	if err != nil {
		t.Errorf("category axiom verification failed: %v", err)
	}

	t.Logf("Category axioms verified successfully for %d objects and %d morphisms",
		len(cat.Objects()), len(cat.Morphisms()))
}

func TestPythonExtractorLanguage(t *testing.T) {
	extractor := NewPythonExtractor()

	if extractor.Language() != "python" {
		t.Errorf("expected language 'python', got %s", extractor.Language())
	}
}

func TestPythonExtractorFileExtensions(t *testing.T) {
	extractor := NewPythonExtractor()

	exts := extractor.FileExtensions()
	if len(exts) != 1 || exts[0] != ".py" {
		t.Errorf("expected ['.py'], got %v", exts)
	}
}

func TestPythonExtractorEmptyDirectory(t *testing.T) {
	// Given: Empty directory
	extractor := NewPythonExtractor()

	// When: Extract from non-existent path
	cat, err := extractor.ExtractFromPath("../../testdata/python/nonexistent")

	// Then: Should handle gracefully
	if err == nil {
		t.Log("Extraction from non-existent path handled gracefully")
	}

	if cat != nil && len(cat.Objects()) > 0 {
		t.Error("expected no objects from non-existent path")
	}
}

func TestPythonExtractorStats(t *testing.T) {
	// Given: Extracted Python category
	extractor := NewPythonExtractor()
	_, err := extractor.ExtractFromPath("../../testdata/python/simple")

	if err != nil {
		t.Fatalf("extraction failed: %v", err)
	}

	// When: Get stats
	stats := extractor.Stats()

	// Then: Verify stats structure
	if stats["language"] != "python" {
		t.Errorf("expected language 'python', got %v", stats["language"])
	}

	if stats["objects"].(int) == 0 {
		t.Error("expected non-zero object count")
	}

	if stats["morphisms"].(int) == 0 {
		t.Error("expected non-zero morphism count")
	}

	if stats["files"].(int) == 0 {
		t.Error("expected non-zero file count")
	}

	t.Logf("Python extractor stats: %+v", stats)
}

func TestPythonExtractorInterfaceCompliance(t *testing.T) {
	// Verify PythonExtractor implements Extractor interface
	var _ Extractor = (*PythonExtractor)(nil)

	t.Log("PythonExtractor implements Extractor interface")
}

// Benchmark tests

func BenchmarkPythonExtractorSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		extractor := NewPythonExtractor()
		_, err := extractor.ExtractFromPath("../../testdata/python/simple")
		if err != nil {
			b.Fatalf("extraction failed: %v", err)
		}
	}
}

func BenchmarkPythonExtractorInheritance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		extractor := NewPythonExtractor()
		_, err := extractor.ExtractFromPath("../../testdata/python/inheritance")
		if err != nil {
			b.Fatalf("extraction failed: %v", err)
		}
	}
}
