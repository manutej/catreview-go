# Phase 1 Projects - Visualization Guide
## hekat, hyperglyph, nanobanana-repo

**Generated**: 2025-12-30
**Projects**: 3 largest LUXOR projects (6,069 files, 117,728 objects)
**Visualizations**: 12 SVG files (4 per project)

---

## 📊 Overview

We've generated 4 categorical visualizations for each of the Phase 1 projects:

1. **Inheritance Graphs** - Class hierarchy relationships (Bottom-to-Top layout)
2. **Module Dependencies** - Import relationships between modules (Top-to-Bottom layout)
3. **Composition Graphs** - How modules organize classes and functions (Clustered layout)
4. **Call Graphs** - Function call relationships (Left-to-Right layout)

Each visualization shows the **top 100 most-connected nodes** to reveal architectural patterns without overwhelming complexity.

---

## 🎯 hekat Visualizations (50,049 objects, 2.38 density)

### hekat-inheritance.svg (62KB)

**What It Shows**: Top 100 classes by inheritance connectivity

**Key Insights to Look For**:
- **Inheritance Depth**: Count how many levels deep the hierarchy goes
  - Healthy: ≤3 levels (composition-friendly)
  - Warning: 4-5 levels (increasing coupling)
  - Problem: >5 levels (rigid hierarchy)

- **Fan-Out**: How many classes inherit from a single base
  - Healthy: <5 children per base (focused abstractions)
  - Warning: 5-10 children (possible over-abstraction)
  - Problem: >10 children (god class)

**Expected Patterns** (based on 3.9% inheritance):
- Shallow hierarchies (2-3 levels max)
- Multiple small trees rather than one large tree
- Many isolated classes (no inheritance)

**What This Reveals**: hekat's inheritance strategy - likely uses classes as containers with minimal hierarchy.

---

### hekat-modules.svg (201KB)

**What It Shows**: Top 100 modules by import connectivity

**Key Insights to Look For**:
- **Import Stars**: Modules imported by many others
  - These are your "core utilities" or "shared infrastructure"
  - Should be stable, well-tested, minimal dependencies

- **Import Hubs**: Modules that import many others
  - These are "orchestrators" or "facades"
  - High risk if they change (many dependents)

- **Circular Imports**: Modules that import each other (cycles)
  - Problem: Indicates tight coupling
  - Should be refactored into separate layers

- **Import Layers**: Can you see distinct layers? (utilities → domain → application)
  - Healthy: Clear layering (arrows flow one direction)
  - Warning: Mixed dependencies (arrows in all directions)

**Expected Patterns** (based on 17% import morphisms):
- Controlled import graph (not too many edges)
- Likely some clear utility modules at the bottom
- Application modules at the top

**What This Reveals**: hekat's dependency architecture - how modules are organized.

---

### hekat-composition.svg (115KB)

**What It Shows**: Top 20 modules with their classes and functions

**Key Insights to Look For**:
- **Module Size**: How many objects per module?
  - Healthy: 10-50 objects per module (focused)
  - Warning: 50-100 objects (getting large)
  - Problem: >100 objects (should split)

- **Function-to-Class Ratio**: Within each module
  - hekat's 75.9% functions suggests many pure functional modules
  - Look for modules with mostly functions (functional style)
  - vs modules with mostly classes (OOP style)

- **Internal Calls**: Dashed green edges show function calls within modules
  - High internal calls = cohesive module (good)
  - Many cross-module calls = potential coupling (review)

**Expected Patterns** (based on 75.9% functions):
- Many modules will be function-heavy
- Classes are likely data containers or simple wrappers
- Functional composition within modules

**What This Reveals**: hekat's internal organization - how code is structured within modules.

---

### hekat-calls.svg (775B - Empty)

**What It Shows**: Function call relationships

**Why It's Empty**: Our current Python extractor captures module-level relationships (imports, inheritance, definitions) but doesn't analyze function bodies to extract call relationships. This would require deeper AST traversal.

**Future Enhancement**: Add call graph extraction by visiting `ast.Call` nodes within function bodies. This would reveal:
- Most-called utility functions
- Call chain depth (nested calls)
- Cross-module function coupling

**Workaround**: Use module dependency graph to infer likely call patterns (imports suggest calls).

---

## 🎨 hyperglyph Visualizations (48,770 objects, 2.25 density)

### hyperglyph-inheritance.svg (64KB)

**What It Shows**: Top 100 classes by inheritance connectivity

**What Makes hyperglyph Special**:
- **81.4% functions** - most functional project
- **Only 3.1% inheritance** - lowest inheritance ratio
- **10.6% classes** - minimal OOP

**Expected Patterns**:
- Very shallow inheritance (1-2 levels max)
- Many isolated classes (data containers)
- Likely uses composition patterns instead

**Key Questions**:
1. Are there ANY deep hierarchies? (probably not)
2. Do classes exist mainly as data holders? (probably yes)
3. Is inheritance used at all? (sparingly)

**What This Reveals**: hyperglyph's extreme functional approach - classes exist but aren't hierarchical.

---

### hyperglyph-modules.svg (252KB - Largest!)

**What It Shows**: Top 100 modules by import connectivity

**Why It's Large**: hyperglyph has the most module-level organization (252KB vs 201KB hekat). This suggests:
- More fine-grained module decomposition
- More import relationships
- Possibly better separation of concerns

**Key Insights to Look For**:
- **Module Count**: More modules = better decomposition?
- **Import Density**: Are imports localized or global?
- **Layer Separation**: Can you identify:
  - Data transformation modules (inputs → outputs)
  - Utility modules (pure functions)
  - Orchestration modules (compose transformations)

**Expected Patterns** (based on 11.4% imports - lowest!):
- Clean module boundaries (few imports)
- Self-contained modules (low coupling)
- Functional pipeline structure

**What This Reveals**: hyperglyph's modular architecture - likely a well-decomposed functional system.

---

### hyperglyph-composition.svg (107KB)

**What It Shows**: Top 20 modules with their contents

**What Makes hyperglyph Special**:
- **81.4% functions** means modules should be VERY function-heavy
- Look for modules with 20-50+ functions and 0-5 classes
- Classes likely simple data structures (namedtuples, dataclasses)

**Key Insights to Look For**:
- **Pure Functional Modules**: Modules with only functions
  - These are your transformation/computation modules
  - Should have no side effects (pure functions)

- **Data Modules**: Modules with mostly classes
  - These define your data structures
  - Should be simple (few methods)

- **Mixed Modules**: Some functions + some classes
  - These are "convenience" modules
  - Classes hold state, functions operate on them

**Expected Pattern**: Clear separation between data (classes) and behavior (functions).

**What This Reveals**: hyperglyph's functional organization - pure functions operating on immutable data.

---

## 🍌 nanobanana-repo Visualizations (18,909 objects, 2.47 density)

### nanobanana-repo-inheritance.svg (67KB)

**What It Shows**: Top 100 classes by inheritance connectivity

**What Makes nanobanana Special**:
- **5.1% inheritance** - HIGHEST of all projects
- **16.4% classes** - high OOP usage
- **70.9% functions** - still functional-first

**Expected Patterns**:
- Deeper hierarchies than hekat/hyperglyph (3-4 levels)
- More inheritance chains visible
- Likely AI model abstractions (BaseModel → ConcreteModel)

**Key Questions**:
1. What are the base classes? (AI models? Data structures?)
2. How deep are the hierarchies? (3-4 levels expected)
3. Are there multiple inheritance trees or one dominant tree?

**What This Reveals**: nanobanana's OOP strategy - likely uses inheritance for AI model abstraction.

---

### nanobanana-repo-modules.svg (315KB - Largest!)

**What It Shows**: Top 100 modules by import connectivity

**Why It's Largest**: nanobanana has the most complex module dependencies (315KB). This suggests:
- Most imports (19.2% of morphisms)
- Complex dependency graph
- Possibly AI framework dependencies (PyTorch, Transformers, etc.)

**Key Insights to Look For**:
- **External Dependencies**: Look for import nodes (different color)
  - How many external libraries are used?
  - Are they concentrated or spread out?

- **Central Modules**: Which modules are import hubs?
  - These are likely "core" or "utils" modules
  - Or AI model definitions (imported everywhere)

- **Dependency Layers**: Can you identify:
  - Framework layer (PyTorch, etc.)
  - Model layer (AI models)
  - Application layer (orchestration)

**Expected Patterns** (based on 19.2% imports - highest):
- Heavy external dependencies (AI frameworks)
- Central model definition modules
- Complex import graph

**What This Reveals**: nanobanana's AI architecture - likely framework-heavy with complex dependencies.

---

### nanobanana-repo-composition.svg (116KB)

**What It Shows**: Top 20 modules with their contents

**What Makes nanobanana Special**:
- **16.4% classes** - expect to see more OOP modules
- **70.9% functions** - still function-heavy overall
- **5.1% inheritance** - classes use inheritance

**Key Insights to Look For**:
- **Model Modules**: Modules with many classes + inheritance
  - These define your AI models
  - Should show inheritance chains (BaseModel → Submodels)

- **Utility Modules**: Modules with mostly functions
  - Data preprocessing, postprocessing
  - Pure functional transformations

- **Orchestration Modules**: Mixed functions + classes
  - Pipeline definitions
  - Training loops

**Expected Pattern**: Clear separation between model definitions (classes) and data processing (functions).

**What This Reveals**: nanobanana's AI system structure - models as classes, processing as functions.

---

## 🔍 Comparative Insights

### What to Compare Across Projects

**1. Inheritance Patterns**
- **hekat**: 3.9% inheritance → Expect shallow hierarchies
- **hyperglyph**: 3.1% inheritance → Expect almost no hierarchies
- **nanobanana**: 5.1% inheritance → Expect deeper hierarchies

**Visual Check**: Which project has the "tallest" inheritance graph?

---

**2. Module Organization**
- **hekat**: 201KB modules graph → Moderate complexity
- **hyperglyph**: 252KB modules graph → High granularity
- **nanobanana**: 315KB modules graph → Complex dependencies

**Visual Check**: Which project has the most "tangled" import graph?

---

**3. Functional vs OOP Style**
- **hekat**: 75.9% functions → Functional-first
- **hyperglyph**: 81.4% functions → Pure functional
- **nanobanana**: 70.9% functions → Balanced

**Visual Check**: Which composition graphs show most functions vs classes per module?

---

## 📈 How to Read the Visualizations

### Node Colors & Shapes

**Classes**:
- Shape: Box (rectangle)
- Color: Light blue
- Label: Class name only (e.g., "UserService")

**Functions**:
- Shape: Ellipse (oval)
- Color: Light green
- Label: Function name only (e.g., "process_data")

**Modules**:
- Shape: Folder
- Color: Wheat (tan)
- Label: Module path (e.g., "core.models")

**Imported Modules**:
- Shape: Component (box with lines)
- Color: Light cyan
- Label: Module name (e.g., "import:pandas")

### Edge Types

**Inheritance** (in inheritance graphs):
- Color: Dark blue
- Thickness: 2 (thick)
- Direction: Child → Parent (bottom to top)
- Label: "inherits"

**Import** (in module graphs):
- Color: Purple
- Thickness: 1 (normal)
- Direction: Importer → Imported
- Label: None

**Defines** (in composition graphs):
- Color: Blue
- Style: Solid
- Direction: Module → Object
- Label: None

**Calls** (in composition graphs, if present):
- Color: Green
- Style: Dashed
- Direction: Caller → Callee
- Label: None

---

## 🎯 Action Items from Visualizations

### After Reviewing hekat Visualizations

1. **Document Best Practices**
   - Identify what makes hekat's architecture excellent
   - Extract patterns (inheritance depth, module organization)
   - Create "hekat pattern library"

2. **Find Reusable Components**
   - Look for central modules in module dependency graph
   - These are candidates for extraction to LUXOR/core/

3. **Measure Against Other Projects**
   - Use hekat as the "gold standard"
   - Compare inheritance depth, module organization

### After Reviewing hyperglyph Visualizations

1. **Study Functional Patterns**
   - How does hyperglyph achieve 81.4% functions?
   - What functional patterns are used?
   - Can other projects adopt these?

2. **Analyze Module Granularity**
   - Why does hyperglyph have such fine-grained modules?
   - Is this better than larger modules?
   - Measure module size distribution

3. **Extract Functional Library**
   - Identify pure functional utility modules
   - Extract to LUXOR/core/functional/

### After Reviewing nanobanana Visualizations

1. **Understand OOP Usage**
   - Why does nanobanana use more inheritance?
   - Is it justified by the domain (AI models)?
   - Can inheritance be reduced?

2. **Analyze Dependency Complexity**
   - Why are there so many imports (19.2%)?
   - Are all dependencies necessary?
   - Can external dependencies be reduced?

3. **Refactoring Candidates**
   - Identify deep inheritance chains (>3 levels)
   - Find highly-coupled modules (many imports)
   - Plan refactoring to reduce complexity

---

## 📊 Visualization File Index

```
Phase 1 Visualizations (12 files, ~1.9MB total)
════════════════════════════════════════════════

hekat (50,049 objects, 2.38 density):
  ✅ hekat-inheritance.svg      62KB   Top 100 classes by connectivity
  ✅ hekat-modules.svg          201KB   Top 100 modules by imports
  ✅ hekat-composition.svg      115KB   Top 20 modules with contents
  ✅ hekat-calls.svg            775B    (empty - see note)

hyperglyph (48,770 objects, 2.25 density):
  ✅ hyperglyph-inheritance.svg  64KB   Top 100 classes by connectivity
  ✅ hyperglyph-modules.svg      252KB   Top 100 modules by imports
  ✅ hyperglyph-composition.svg  107KB   Top 20 modules with contents
  ✅ hyperglyph-calls.svg        775B    (empty - see note)

nanobanana-repo (18,909 objects, 2.47 density):
  ✅ nanobanana-repo-inheritance.svg  67KB   Top 100 classes by connectivity
  ✅ nanobanana-repo-modules.svg      315KB   Top 100 modules by imports
  ✅ nanobanana-repo-composition.svg  116KB   Top 20 modules with contents
  ✅ nanobanana-repo-calls.svg        775B    (empty - see note)
```

**Note**: Call graphs are empty because current extractor doesn't analyze function bodies. This is a known limitation that can be enhanced in future versions.

---

## 🚀 Next Steps

1. **Review Opened Visualizations**
   - hekat-inheritance.svg (class hierarchies)
   - hekat-modules.svg (module dependencies)
   - hyperglyph-composition.svg (functional organization)
   - nanobanana-repo-inheritance.svg (OOP usage)

2. **Open Remaining Visualizations as Needed**
   ```bash
   # Open all hekat visualizations
   open hekat-*.svg

   # Open all hyperglyph visualizations
   open hyperglyph-*.svg

   # Open all nanobanana visualizations
   open nanobanana-repo-*.svg
   ```

3. **Extract Insights**
   - Take screenshots of interesting patterns
   - Document architectural decisions revealed
   - Identify refactoring opportunities

4. **Share with Team**
   - Use visualizations in architectural reviews
   - Present patterns to other projects
   - Create "LUXOR Architectural Patterns" guide

---

**All visualizations are located in**: `/Users/manu/Documents/LUXOR/catreview-go/`

**To regenerate with different parameters**:
```bash
# More nodes (more detail, larger files)
go run examples/python/visualize_project.go \
    --input hekat-analysis.json \
    --output hekat-detailed \
    --max-nodes 200

# Fewer nodes (simpler visualization)
go run examples/python/visualize_project.go \
    --input hekat-analysis.json \
    --output hekat-simple \
    --max-nodes 50
```

---

**Status**: ✅ **12 Visualizations Complete - Ready for Review**
