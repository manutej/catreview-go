// Package functor implements functors for mapping between categories.
//
// Functors are structure-preserving maps between categories that:
// 1. Map objects to objects: F: Ob(C) → Ob(D)
// 2. Map morphisms to morphisms: F: Hom(A,B) → Hom(F(A),F(B))
// 3. Preserve composition: F(g ∘ f) = F(g) ∘ F(f)
// 4. Preserve identities: F(id_A) = id_{F(A)}
package functor

import (
	"fmt"

	"github.com/manu/catreview/pkg/category"
)

// Functor represents a structure-preserving mapping between categories.
//
// A functor F: C → D consists of:
// - An object mapping: Ob(C) → Ob(D)
// - A morphism mapping: Mor(C) → Mor(D)
//
// Functor laws:
// 1. F(g ∘ f) = F(g) ∘ F(f) for all composable morphisms f, g
// 2. F(id_A) = id_{F(A)} for all objects A
type Functor interface {
	// Name returns the functor's name
	Name() string

	// MapObject maps an object from source category to target category
	MapObject(obj *category.Object) (*category.Object, error)

	// MapMorphism maps a morphism from source category to target category
	MapMorphism(morph *category.Morphism) (*category.Morphism, error)

	// SourceCategory returns the source category
	SourceCategory() *category.Category

	// TargetCategory returns the target category
	TargetCategory() *category.Category

	// VerifyLaws verifies functor laws hold
	VerifyLaws() error
}

// BaseFunctor provides common functor functionality.
type BaseFunctor struct {
	name     string
	source   *category.Category
	target   *category.Category
	objMap   map[string]*category.Object
	morphMap map[string]*category.Morphism
}

// NewBaseFunctor creates a new base functor.
func NewBaseFunctor(name string, source, target *category.Category) *BaseFunctor {
	return &BaseFunctor{
		name:     name,
		source:   source,
		target:   target,
		objMap:   make(map[string]*category.Object),
		morphMap: make(map[string]*category.Morphism),
	}
}

// Name returns the functor's name.
func (f *BaseFunctor) Name() string {
	return f.name
}

// SourceCategory returns the source category.
func (f *BaseFunctor) SourceCategory() *category.Category {
	return f.source
}

// TargetCategory returns the target category.
func (f *BaseFunctor) TargetCategory() *category.Category {
	return f.target
}

// AddObjectMapping adds an explicit object mapping.
func (f *BaseFunctor) AddObjectMapping(sourceID string, target *category.Object) {
	f.objMap[sourceID] = target
}

// AddMorphismMapping adds an explicit morphism mapping.
func (f *BaseFunctor) AddMorphismMapping(sourceID string, target *category.Morphism) {
	f.morphMap[sourceID] = target
}

// GetObjectMapping retrieves a mapped object.
func (f *BaseFunctor) GetObjectMapping(sourceID string) (*category.Object, bool) {
	obj, exists := f.objMap[sourceID]
	return obj, exists
}

// GetMorphismMapping retrieves a mapped morphism.
func (f *BaseFunctor) GetMorphismMapping(sourceID string) (*category.Morphism, bool) {
	morph, exists := f.morphMap[sourceID]
	return morph, exists
}

// VerifyCompositionLaw verifies F(g ∘ f) = F(g) ∘ F(f).
func (f *BaseFunctor) VerifyCompositionLaw(mapper Functor) error {
	morphisms := f.source.Morphisms()

	// Check composition law for sample pairs
	checked := 0
	const maxChecks = 50

	for i := 0; i < len(morphisms) && checked < maxChecks; i++ {
		for j := 0; j < len(morphisms) && checked < maxChecks; j++ {
			fMorph := morphisms[i]
			gMorph := morphisms[j]

			// Skip if not composable
			if !fMorph.IsComposable(gMorph) {
				continue
			}

			// Compute g ∘ f in source category
			composed, err := f.source.Compose(fMorph, gMorph)
			if err != nil {
				continue
			}

			// Compute F(g ∘ f)
			fComposed, err1 := mapper.MapMorphism(composed)
			if err1 != nil {
				return fmt.Errorf("failed to map composed morphism: %v", err1)
			}

			// Compute F(f) and F(g)
			ff, err2 := mapper.MapMorphism(fMorph)
			if err2 != nil {
				return fmt.Errorf("failed to map f: %v", err2)
			}

			fg, err3 := mapper.MapMorphism(gMorph)
			if err3 != nil {
				return fmt.Errorf("failed to map g: %v", err3)
			}

			// Compute F(g) ∘ F(f) in target category
			fgComposed, err4 := f.target.Compose(ff, fg)
			if err4 != nil {
				return fmt.Errorf("failed to compose in target: %v", err4)
			}

			// Verify F(g ∘ f) = F(g) ∘ F(f)
			if fComposed.Source != fgComposed.Source || fComposed.Target != fgComposed.Target {
				return fmt.Errorf("composition law violated: F(%s∘%s) != F(%s)∘F(%s)",
					gMorph.ID, fMorph.ID, gMorph.ID, fMorph.ID)
			}

			checked++
		}
	}

	return nil
}

// VerifyIdentityLaw verifies F(id_A) = id_{F(A)}.
func (f *BaseFunctor) VerifyIdentityLaw(mapper Functor) error {
	objects := f.source.Objects()

	for _, obj := range objects {
		// Get identity morphism for this object
		idMorph := &category.Morphism{
			ID:     fmt.Sprintf("id_%s", obj.ID),
			Source: obj.ID,
			Target: obj.ID,
			Type:   "identity",
		}

		// Map the identity
		fIdMorph, err := mapper.MapMorphism(idMorph)
		if err != nil {
			return fmt.Errorf("failed to map identity for %s: %v", obj.ID, err)
		}

		// Map the object
		fObj, err := mapper.MapObject(obj)
		if err != nil {
			return fmt.Errorf("failed to map object %s: %v", obj.ID, err)
		}

		// Verify F(id_A) is identity for F(A)
		if fIdMorph.Source != fObj.ID || fIdMorph.Target != fObj.ID {
			return fmt.Errorf("identity law violated: F(id_%s) is not id_{F(%s)}",
				obj.ID, obj.ID)
		}
	}

	return nil
}

// PackageAbstractionFunctor maps file-level category to package-level category.
//
// This functor abstracts from individual files to packages:
// - Files → Packages
// - File dependencies → Package dependencies
//
// This is useful for viewing architecture at different granularities.
type PackageAbstractionFunctor struct {
	*BaseFunctor
}

// NewPackageAbstractionFunctor creates a functor from file-level to package-level.
func NewPackageAbstractionFunctor(source, target *category.Category) *PackageAbstractionFunctor {
	return &PackageAbstractionFunctor{
		BaseFunctor: NewBaseFunctor("FileToPackage", source, target),
	}
}

// MapObject maps a file object to its package object.
func (f *PackageAbstractionFunctor) MapObject(obj *category.Object) (*category.Object, error) {
	// Check cache first
	if cached, exists := f.GetObjectMapping(obj.ID); exists {
		return cached, nil
	}

	// Extract package from file metadata
	packageName, ok := obj.Metadata["package"].(string)
	if !ok {
		return nil, fmt.Errorf("object %s has no package metadata", obj.ID)
	}

	// Create or get package object
	pkgID := fmt.Sprintf("pkg:%s", packageName)
	pkgObj, exists := f.target.GetObject(pkgID)
	if !exists {
		pkgObj = category.NewObject(pkgID, "package", packageName, map[string]interface{}{
			"files": []string{obj.ID},
		})
		if err := f.target.AddObject(pkgObj); err != nil {
			return nil, err
		}
	} else {
		// Add this file to package's file list
		if files, ok := pkgObj.Metadata["files"].([]string); ok {
			pkgObj.Metadata["files"] = append(files, obj.ID)
		}
	}

	// Cache mapping
	f.AddObjectMapping(obj.ID, pkgObj)

	return pkgObj, nil
}

// MapMorphism maps a file dependency to a package dependency.
func (f *PackageAbstractionFunctor) MapMorphism(morph *category.Morphism) (*category.Morphism, error) {
	// Check cache first
	if cached, exists := f.GetMorphismMapping(morph.ID); exists {
		return cached, nil
	}

	// Get source and target objects from source category
	sourceObj, exists := f.source.GetObject(morph.Source)
	if !exists {
		return nil, fmt.Errorf("source object %s not found", morph.Source)
	}

	targetObj, exists := f.source.GetObject(morph.Target)
	if !exists {
		return nil, fmt.Errorf("target object %s not found", morph.Target)
	}

	// Map objects to packages
	sourcePkg, err := f.MapObject(sourceObj)
	if err != nil {
		return nil, err
	}

	targetPkg, err := f.MapObject(targetObj)
	if err != nil {
		return nil, err
	}

	// Don't create self-dependencies
	if sourcePkg.ID == targetPkg.ID {
		// Return identity morphism
		return &category.Morphism{
			ID:     fmt.Sprintf("id_%s", sourcePkg.ID),
			Source: sourcePkg.ID,
			Target: sourcePkg.ID,
			Type:   "identity",
		}, nil
	}

	// Create package dependency morphism
	pkgMorphID := fmt.Sprintf("dep:%s->%s", sourcePkg.ID, targetPkg.ID)
	pkgMorph, exists := f.target.GetMorphism(pkgMorphID)
	if !exists {
		pkgMorph = category.NewMorphism(
			pkgMorphID,
			sourcePkg.ID,
			targetPkg.ID,
			"dependency",
			map[string]interface{}{
				"source_files": []string{morph.Source},
				"target_files": []string{morph.Target},
			},
		)
		if err := f.target.AddMorphism(pkgMorph); err != nil {
			return nil, err
		}
	} else {
		// Add to existing dependency's file lists
		if sources, ok := pkgMorph.Metadata["source_files"].([]string); ok {
			pkgMorph.Metadata["source_files"] = append(sources, morph.Source)
		}
		if targets, ok := pkgMorph.Metadata["target_files"].([]string); ok {
			pkgMorph.Metadata["target_files"] = append(targets, morph.Target)
		}
	}

	// Cache mapping
	f.AddMorphismMapping(morph.ID, pkgMorph)

	return pkgMorph, nil
}

// VerifyLaws verifies functor laws for this specific functor.
func (f *PackageAbstractionFunctor) VerifyLaws() error {
	if err := f.VerifyIdentityLaw(f); err != nil {
		return err
	}
	if err := f.VerifyCompositionLaw(f); err != nil {
		return err
	}
	return nil
}
