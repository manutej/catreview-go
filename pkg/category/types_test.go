package category

import (
	"testing"
)

func TestObjectCreation(t *testing.T) {
	obj := NewObject("obj1", "module", "MyModule", nil)

	if obj.ID != "obj1" {
		t.Errorf("Expected ID 'obj1', got '%s'", obj.ID)
	}
	if obj.Type != "module" {
		t.Errorf("Expected Type 'module', got '%s'", obj.Type)
	}
	if obj.Name != "MyModule" {
		t.Errorf("Expected Name 'MyModule', got '%s'", obj.Name)
	}
	if obj.Metadata == nil {
		t.Error("Expected metadata to be initialized")
	}
}

func TestMorphismComposability(t *testing.T) {
	// f: A → B
	f := NewMorphism("f", "A", "B", "dependency", nil)
	// g: B → C
	g := NewMorphism("g", "B", "C", "dependency", nil)
	// h: C → D
	h := NewMorphism("h", "C", "D", "dependency", nil)

	// f and g should be composable
	if !f.IsComposable(g) {
		t.Error("f and g should be composable")
	}

	// g and h should be composable
	if !g.IsComposable(h) {
		t.Error("g and h should be composable")
	}

	// f and h should NOT be composable
	if f.IsComposable(h) {
		t.Error("f and h should not be composable")
	}
}

func TestCategoryObjectAddition(t *testing.T) {
	cat := NewCategory("test")

	obj := NewObject("obj1", "module", "Test", nil)
	if err := cat.AddObject(obj); err != nil {
		t.Errorf("Failed to add object: %v", err)
	}

	// Verify object was added
	retrieved, exists := cat.GetObject("obj1")
	if !exists {
		t.Error("Object should exist after adding")
	}
	if retrieved.ID != "obj1" {
		t.Errorf("Retrieved wrong object: %s", retrieved.ID)
	}

	// Verify identity morphism was created
	identity, exists := cat.GetMorphism("id_obj1")
	if !exists {
		t.Error("Identity morphism should be created")
	}
	if identity.Source != "obj1" || identity.Target != "obj1" {
		t.Error("Identity morphism should be obj1 → obj1")
	}

	// Cannot add duplicate
	if err := cat.AddObject(obj); err == nil {
		t.Error("Should not allow duplicate object")
	}
}

func TestCategoryMorphismAddition(t *testing.T) {
	cat := NewCategory("test")

	// Add objects first
	objA := NewObject("A", "module", "ModuleA", nil)
	objB := NewObject("B", "module", "ModuleB", nil)
	cat.AddObject(objA)
	cat.AddObject(objB)

	// Add morphism
	morph := NewMorphism("f", "A", "B", "dependency", nil)
	if err := cat.AddMorphism(morph); err != nil {
		t.Errorf("Failed to add morphism: %v", err)
	}

	// Verify morphism was added
	retrieved, exists := cat.GetMorphism("f")
	if !exists {
		t.Error("Morphism should exist after adding")
	}
	if retrieved.ID != "f" {
		t.Errorf("Retrieved wrong morphism: %s", retrieved.ID)
	}

	// Cannot add morphism with nonexistent source
	invalidMorph := NewMorphism("g", "X", "A", "dependency", nil)
	if err := cat.AddMorphism(invalidMorph); err == nil {
		t.Error("Should not allow morphism with nonexistent source")
	}

	// Cannot add morphism with nonexistent target
	invalidMorph2 := NewMorphism("h", "A", "Y", "dependency", nil)
	if err := cat.AddMorphism(invalidMorph2); err == nil {
		t.Error("Should not allow morphism with nonexistent target")
	}
}

func TestMorphismComposition(t *testing.T) {
	cat := NewCategory("test")

	// Create objects: A → B → C
	objA := NewObject("A", "module", "ModuleA", nil)
	objB := NewObject("B", "module", "ModuleB", nil)
	objC := NewObject("C", "module", "ModuleC", nil)
	cat.AddObject(objA)
	cat.AddObject(objB)
	cat.AddObject(objC)

	// Create morphisms: f: A → B, g: B → C
	f := NewMorphism("f", "A", "B", "dependency", nil)
	g := NewMorphism("g", "B", "C", "dependency", nil)
	cat.AddMorphism(f)
	cat.AddMorphism(g)

	// Compose: g ∘ f: A → C
	composed, err := cat.Compose(f, g)
	if err != nil {
		t.Errorf("Composition failed: %v", err)
	}

	if composed.Source != "A" {
		t.Errorf("Composed source should be A, got %s", composed.Source)
	}
	if composed.Target != "C" {
		t.Errorf("Composed target should be C, got %s", composed.Target)
	}

	// Cannot compose incompatible morphisms
	h := NewMorphism("h", "C", "A", "dependency", nil)
	cat.AddMorphism(h)
	_, err = cat.Compose(f, h) // f: A→B, h: C→A (not composable)
	if err == nil {
		t.Error("Should not allow composition of incompatible morphisms")
	}
}

func TestCategoryAxiomVerification(t *testing.T) {
	cat := NewCategory("test")

	// Create a simple category: A → B → C
	objA := NewObject("A", "module", "ModuleA", nil)
	objB := NewObject("B", "module", "ModuleB", nil)
	objC := NewObject("C", "module", "ModuleC", nil)
	cat.AddObject(objA)
	cat.AddObject(objB)
	cat.AddObject(objC)

	f := NewMorphism("f", "A", "B", "dependency", nil)
	g := NewMorphism("g", "B", "C", "dependency", nil)
	cat.AddMorphism(f)
	cat.AddMorphism(g)

	// Verify axioms (identity and associativity)
	if err := cat.VerifyAxioms(); err != nil {
		t.Errorf("Axiom verification failed: %v", err)
	}
}

func TestCategoryStats(t *testing.T) {
	cat := NewCategory("test")

	// Initially empty (except identities will be added with objects)
	stats := cat.Stats()
	if stats["objects"] != 0 {
		t.Errorf("Expected 0 objects, got %d", stats["objects"])
	}

	// Add objects
	objA := NewObject("A", "module", "ModuleA", nil)
	objB := NewObject("B", "module", "ModuleB", nil)
	cat.AddObject(objA)
	cat.AddObject(objB)

	stats = cat.Stats()
	if stats["objects"] != 2 {
		t.Errorf("Expected 2 objects, got %d", stats["objects"])
	}
	if stats["identities"] != 2 {
		t.Errorf("Expected 2 identities, got %d", stats["identities"])
	}

	// Add morphism
	f := NewMorphism("f", "A", "B", "dependency", nil)
	cat.AddMorphism(f)

	stats = cat.Stats()
	if stats["morphisms"] != 1 { // Excludes identity morphisms
		t.Errorf("Expected 1 morphism (excluding identities), got %d", stats["morphisms"])
	}
}

func TestAssociativityLaw(t *testing.T) {
	cat := NewCategory("test")

	// Create objects: A → B → C → D
	objA := NewObject("A", "module", "ModuleA", nil)
	objB := NewObject("B", "module", "ModuleB", nil)
	objC := NewObject("C", "module", "ModuleC", nil)
	objD := NewObject("D", "module", "ModuleD", nil)
	cat.AddObject(objA)
	cat.AddObject(objB)
	cat.AddObject(objC)
	cat.AddObject(objD)

	// Create morphisms: f: A→B, g: B→C, h: C→D
	f := NewMorphism("f", "A", "B", "dependency", nil)
	g := NewMorphism("g", "B", "C", "dependency", nil)
	h := NewMorphism("h", "C", "D", "dependency", nil)
	cat.AddMorphism(f)
	cat.AddMorphism(g)
	cat.AddMorphism(h)

	// Compute (h ∘ g) ∘ f
	hg, err := cat.Compose(g, h)
	if err != nil {
		t.Fatalf("Failed to compose h∘g: %v", err)
	}
	left, err := cat.Compose(f, hg)
	if err != nil {
		t.Fatalf("Failed to compose (h∘g)∘f: %v", err)
	}

	// Compute h ∘ (g ∘ f)
	gf, err := cat.Compose(f, g)
	if err != nil {
		t.Fatalf("Failed to compose g∘f: %v", err)
	}
	right, err := cat.Compose(gf, h)
	if err != nil {
		t.Fatalf("Failed to compose h∘(g∘f): %v", err)
	}

	// Verify associativity: (h ∘ g) ∘ f = h ∘ (g ∘ f)
	if left.Source != right.Source || left.Target != right.Target {
		t.Error("Associativity law violated")
	}
	if left.Source != "A" || left.Target != "D" {
		t.Errorf("Composed morphism should be A→D, got %s→%s", left.Source, left.Target)
	}
}

func TestIdentityLaw(t *testing.T) {
	cat := NewCategory("test")

	// Create objects: A → B
	objA := NewObject("A", "module", "ModuleA", nil)
	objB := NewObject("B", "module", "ModuleB", nil)
	cat.AddObject(objA)
	cat.AddObject(objB)

	// Create morphism: f: A → B
	f := NewMorphism("f", "A", "B", "dependency", nil)
	cat.AddMorphism(f)

	// Get identity morphisms
	idA, _ := cat.GetMorphism("id_A")
	idB, _ := cat.GetMorphism("id_B")

	// Test left identity: id_B ∘ f = f
	left, err := cat.Compose(f, idB)
	if err != nil {
		t.Fatalf("Failed to compose f∘id_B: %v", err)
	}
	if left.Source != f.Source || left.Target != f.Target {
		t.Error("Left identity law violated")
	}

	// Test right identity: f ∘ id_A = f
	right, err := cat.Compose(idA, f)
	if err != nil {
		t.Fatalf("Failed to compose id_A∘f: %v", err)
	}
	if right.Source != f.Source || right.Target != f.Target {
		t.Error("Right identity law violated")
	}
}
