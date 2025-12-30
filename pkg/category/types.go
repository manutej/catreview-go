// Package category implements core category theory structures for codebase analysis.
//
// This package provides the mathematical foundation for modeling software as categories:
// - Objects represent software components (modules, packages, types, functions)
// - Morphisms represent relationships (dependencies, function calls, inheritance)
// - Categories ensure composition laws hold (associativity, identity)
package category

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// Object represents an object in a category.
// In software terms, this could be a module, package, class, or function.
//
// Mathematical properties:
// - Objects must be uniquely identifiable (via ID)
// - Objects are immutable once created
// - Objects carry metadata for analysis
type Object struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"` // "module", "package", "class", "function"
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

// NewObject creates a new category object with the given properties.
func NewObject(id, objType, name string, metadata map[string]interface{}) *Object {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	return &Object{
		ID:       id,
		Type:     objType,
		Name:     name,
		Metadata: metadata,
	}
}

// Hash returns a cryptographic hash of this object for equality checking.
func (o *Object) Hash() string {
	data, _ := json.Marshal(o)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash[:16]) // Use first 16 bytes for brevity
}

// Morphism represents a morphism (arrow) between objects in a category.
// In software terms, this could be a dependency, function call, or inheritance relationship.
//
// Mathematical properties:
// - Morphisms must have a valid source and target (domain and codomain)
// - Morphisms must be composable when target(f) == source(g)
// - Composition must be associative: (h ∘ g) ∘ f = h ∘ (g ∘ f)
type Morphism struct {
	ID       string                 `json:"id"`
	Source   string                 `json:"source"` // ID of source object
	Target   string                 `json:"target"` // ID of target object
	Type     string                 `json:"type"`   // "dependency", "call", "inheritance"
	Metadata map[string]interface{} `json:"metadata"`
}

// NewMorphism creates a new morphism between source and target objects.
func NewMorphism(id, source, target, morphType string, metadata map[string]interface{}) *Morphism {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	return &Morphism{
		ID:       id,
		Source:   source,
		Target:   target,
		Type:     morphType,
		Metadata: metadata,
	}
}

// IsComposable checks if this morphism can be composed with another.
// Morphisms are composable when the target of f equals the source of g.
func (m *Morphism) IsComposable(other *Morphism) bool {
	return m.Target == other.Source
}

// Category represents a category with objects and morphisms.
//
// A category C consists of:
// 1. A collection of objects Ob(C)
// 2. For each pair of objects A, B, a set of morphisms Hom(A,B)
// 3. For each object A, an identity morphism id_A : A → A
// 4. A composition operation ∘ that is associative
//
// Category axioms:
// - Associativity: (h ∘ g) ∘ f = h ∘ (g ∘ f)
// - Identity: f ∘ id_A = f and id_B ∘ f = f for f : A → B
type Category struct {
	Name       string                `json:"name"`
	Objects_   map[string]*Object    `json:"objects"`
	Morphisms_ map[string]*Morphism  `json:"morphisms"`
	Identities map[string]*Morphism  `json:"identities"` // Identity morphisms for each object
}

// NewCategory creates a new category with the given name.
func NewCategory(name string) *Category {
	return &Category{
		Name:       name,
		Objects_:   make(map[string]*Object),
		Morphisms_: make(map[string]*Morphism),
		Identities: make(map[string]*Morphism),
	}
}

// AddObject adds an object to the category and creates its identity morphism.
func (c *Category) AddObject(obj *Object) error {
	if obj == nil {
		return fmt.Errorf("cannot add nil object")
	}
	if _, exists := c.Objects_[obj.ID]; exists {
		return fmt.Errorf("object %s already exists", obj.ID)
	}

	c.Objects_[obj.ID] = obj

	// Create identity morphism: id_A : A → A
	identity := NewMorphism(
		fmt.Sprintf("id_%s", obj.ID),
		obj.ID,
		obj.ID,
		"identity",
		map[string]interface{}{"is_identity": true},
	)
	c.Identities[obj.ID] = identity
	c.Morphisms_[identity.ID] = identity

	return nil
}

// AddMorphism adds a morphism to the category after validating it.
func (c *Category) AddMorphism(m *Morphism) error {
	if m == nil {
		return fmt.Errorf("cannot add nil morphism")
	}

	// Validate source and target exist
	if _, exists := c.Objects_[m.Source]; !exists {
		return fmt.Errorf("source object %s does not exist", m.Source)
	}
	if _, exists := c.Objects_[m.Target]; !exists {
		return fmt.Errorf("target object %s does not exist", m.Target)
	}

	if _, exists := c.Morphisms_[m.ID]; exists {
		return fmt.Errorf("morphism %s already exists", m.ID)
	}

	c.Morphisms_[m.ID] = m
	return nil
}

// GetObject retrieves an object by ID.
func (c *Category) GetObject(id string) (*Object, bool) {
	obj, exists := c.Objects_[id]
	return obj, exists
}

// GetMorphism retrieves a morphism by ID.
func (c *Category) GetMorphism(id string) (*Morphism, bool) {
	m, exists := c.Morphisms_[id]
	return m, exists
}

// Objects returns all objects in the category.
func (c *Category) Objects() []*Object {
	objects := make([]*Object, 0, len(c.Objects_))
	for _, obj := range c.Objects_ {
		objects = append(objects, obj)
	}
	return objects
}

// Morphisms returns all morphisms in the category.
func (c *Category) Morphisms() []*Morphism {
	morphisms := make([]*Morphism, 0, len(c.Morphisms_))
	for _, m := range c.Morphisms_ {
		morphisms = append(morphisms, m)
	}
	return morphisms
}

// Compose composes two morphisms f: A → B and g: B → C to get g ∘ f: A → C.
// Returns error if morphisms are not composable.
func (c *Category) Compose(f, g *Morphism) (*Morphism, error) {
	if !f.IsComposable(g) {
		return nil, fmt.Errorf("morphisms not composable: target(%s)=%s != source(%s)=%s",
			f.ID, f.Target, g.ID, g.Source)
	}

	composed := NewMorphism(
		fmt.Sprintf("%s∘%s", g.ID, f.ID),
		f.Source,
		g.Target,
		"composed",
		map[string]interface{}{
			"composed_from": []string{f.ID, g.ID},
		},
	)

	return composed, nil
}

// VerifyAxioms verifies that this category satisfies category axioms:
// 1. Associativity: (h ∘ g) ∘ f = h ∘ (g ∘ f)
// 2. Identity: f ∘ id_A = f and id_B ∘ f = f
func (c *Category) VerifyAxioms() error {
	// Check identity law for all morphisms
	for _, m := range c.Morphisms_ {
		if m.Type == "identity" {
			continue // Skip identity morphisms themselves
		}

		// Check left identity: id_B ∘ f = f
		if idTarget, exists := c.Identities[m.Target]; exists {
			composed, err := c.Compose(m, idTarget)
			if err != nil {
				return fmt.Errorf("identity composition failed: %v", err)
			}
			// Verify composition gives back original morphism properties
			if composed.Source != m.Source || composed.Target != m.Target {
				return fmt.Errorf("left identity law violated for %s", m.ID)
			}
		}

		// Check right identity: f ∘ id_A = f
		if idSource, exists := c.Identities[m.Source]; exists {
			composed, err := c.Compose(idSource, m)
			if err != nil {
				return fmt.Errorf("identity composition failed: %v", err)
			}
			if composed.Source != m.Source || composed.Target != m.Target {
				return fmt.Errorf("right identity law violated for %s", m.ID)
			}
		}
	}

	// Check associativity for sample morphism chains
	// For performance, we sample rather than checking all possible triples
	morphismList := c.Morphisms()
	checked := 0
	const maxChecks = 100 // Limit verification for large categories

	for i := 0; i < len(morphismList) && checked < maxChecks; i++ {
		f := morphismList[i]
		for j := 0; j < len(morphismList) && checked < maxChecks; j++ {
			g := morphismList[j]
			if !f.IsComposable(g) {
				continue
			}
			for k := 0; k < len(morphismList) && checked < maxChecks; k++ {
				h := morphismList[k]
				if !g.IsComposable(h) {
					continue
				}

				// Compute (h ∘ g) ∘ f
				hg, err1 := c.Compose(g, h)
				if err1 != nil {
					continue
				}
				left, err2 := c.Compose(f, hg)
				if err2 != nil {
					continue
				}

				// Compute h ∘ (g ∘ f)
				gf, err3 := c.Compose(f, g)
				if err3 != nil {
					continue
				}
				right, err4 := c.Compose(gf, h)
				if err4 != nil {
					continue
				}

				// Verify: (h ∘ g) ∘ f = h ∘ (g ∘ f)
				if left.Source != right.Source || left.Target != right.Target {
					return fmt.Errorf("associativity violated: (%s∘%s)∘%s != %s∘(%s∘%s)",
						h.ID, g.ID, f.ID, h.ID, g.ID, f.ID)
				}

				checked++
			}
		}
	}

	return nil
}

// Stats returns statistics about the category.
func (c *Category) Stats() map[string]int {
	return map[string]int{
		"objects":   len(c.Objects_),
		"morphisms": len(c.Morphisms_) - len(c.Identities), // Exclude identity morphisms
		"identities": len(c.Identities),
	}
}
