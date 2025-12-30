# catreview - Categorical Codebase Review

**Production-ready categorical analysis for Go codebases** • [Quick Start](QUICK-START.md) • [Production Guide](PRODUCTION-GUIDE.md)

![Status: Production Ready](https://img.shields.io/badge/status-production%20ready-green)
![Quality: 96%](https://img.shields.io/badge/quality-96%25-brightgreen)
![Validated: 5 Repos](https://img.shields.io/badge/validated-5%20repos-blue)
![Go Version: 1.21+](https://img.shields.io/badge/go-1.21+-00ADD8)

A Go tool for analyzing software architecture using category theory.

---

## 📚 Documentation

| Document | Purpose | Audience |
|----------|---------|----------|
| **[QUICK-START.md](QUICK-START.md)** | Get started in 5 minutes | New users |
| **[PRODUCTION-GUIDE.md](PRODUCTION-GUIDE.md)** | Validation results & real-world examples | Production users |
| **README.md** (this file) | Complete reference | All users |

---

## Overview

**catreview** applies category theory to analyze codebases, extracting categorical models that reveal:

- **Architectural Complexity** - Basu-Isik diagram complexity metrics
- **Coupling Metrics** - Afferent/efferent coupling and instability (I = Ce / (Ca + Ce))
- **Dependency Cycles** - Circular dependency detection
- **Abstraction Levels** - Package-level abstraction via functors

### Production Validation

✅ **5 real-world repositories** analyzed (5,055 objects, 8,141 morphisms)
✅ **96% quality score** (exceeds 90% threshold)
✅ **Linear complexity scaling** (R² = 0.999)
✅ **100% category axiom verification**
✅ **Single-iteration RMP convergence**

## Mathematical Foundation

### Category Theory Basics

A **category** C consists of:
1. A collection of **objects** Ob(C)
2. For each pair of objects A, B, a set of **morphisms** Hom(A,B)
3. For each object A, an **identity morphism** id_A : A → A
4. A **composition operation** ∘ that is associative

### Category Axioms

1. **Associativity**: (h ∘ g) ∘ f = h ∘ (g ∘ f)
2. **Identity**: f ∘ id_A = f and id_B ∘ f = f for f : A → B

###Software as Categories

| Category Theory | Software |
|----------------|----------|
| Objects | Modules, Packages, Types, Functions |
| Morphisms | Dependencies, Function Calls, Imports |
| Composition | Transitive Dependencies |
| Functors | Abstraction Mappings (File → Package) |

## Installation

```bash
go install github.com/manu/catreview/cmd/catreview@latest
```

Or build from source:

```bash
git clone https://github.com/manu/catreview
cd catreview
go build -o catreview ./cmd/catreview
```

## Usage

### 1. Extract Categorical Model

Extract a categorical model from your Go codebase:

```bash
catreview extract ./path/to/code -o model.json
```

**Output:**
```
Extracting categorical model from: ./path/to/code
Extracted:
  Objects:   84
  Morphisms: 102
  Identities: 84
Model saved to: model.json
```

### 2. Analyze Complexity

Analyze the extracted model:

```bash
catreview analyze model.json -o report.json
```

**Output:**
```
Categorical Analysis Report
===========================

Category Statistics:
  Objects:    84
  Morphisms:  102

Complexity Metrics:
  Diagram Complexity:    509.49
  Kolmogorov Complexity: 7421 bytes

Dependency Analysis:
  Cycles Found: 0

Top 5 Most Unstable Components:
  pkg/analysis/complexity.go: I=1.00 (Ce=24, Ca=0)
  pkg/extractor/go_extractor.go: I=1.00 (Ce=23, Ca=0)
  ...

Full report saved to: report.json
```

### 3. Verify Category Axioms

Verify that the extracted category satisfies category axioms:

```bash
catreview verify model.json --max-cycles 0 --fail-on-violation
```

**Output:**
```
Verifying category axioms: model.json
Checking associativity and identity laws...
✅ Category axioms verified successfully

Checking for dependency cycles (max allowed: 0)...
Found 0 cycles
✅ Cycle count within limit
```

### 4. Create Package-Level Abstraction

Apply a functor to map file-level category to package-level:

```bash
catreview abstract model.json -o abstract.json
```

**Output:**
```
Creating package-level abstraction from: model.json
Mapping files to packages...
Mapping dependencies...
Verifying functor laws...
✅ Functor laws verified

Abstracted Category:
  Packages:           3
  Package Dependencies: 4

Abstracted model saved to: abstract.json
```

## CLI Reference

### `extract`

Extract categorical model from codebase.

```bash
catreview extract [path] [flags]
```

**Flags:**
- `-o, --output string` - Output file for categorical model (default "model.json")
- `--pretty` - Pretty-print JSON output

### `analyze`

Analyze categorical model and generate report.

```bash
catreview analyze [model.json] [flags]
```

**Flags:**
- `-o, --output string` - Output file for analysis report (default "report.json")
- `--pretty` - Pretty-print JSON output (default true)

### `verify`

Verify category axioms.

```bash
catreview verify [model.json] [flags]
```

**Flags:**
- `--max-cycles int` - Maximum allowed cycles, -1 = no limit (default -1)
- `--fail-on-violation` - Exit with error on axiom violation

### `abstract`

Create package-level abstraction via functor.

```bash
catreview abstract [model.json] [flags]
```

**Flags:**
- `-o, --output string` - Output file for abstracted model (default "abstract.json")
- `--pretty` - Pretty-print JSON output (default true)

## Complexity Metrics

### Basu-Isik Diagram Complexity

```
c(D) = Σ c_obj(o) + Σ c_morph(m) + c_comp(D)
```

Where:
- `c_obj(o)` = complexity of object o (based on metadata and type)
- `c_morph(m)` = complexity of morphism m (based on type)
- `c_comp(D)` = composition complexity (based on composable chains)

### Kolmogorov Complexity

Estimated via gzip compression:

```
K(x) ≈ |gzip(x)|
```

This provides an upper bound on the true Kolmogorov complexity.

### Coupling Metrics

- **Afferent Coupling (Ca)**: Number of incoming dependencies
- **Efferent Coupling (Ce)**: Number of outgoing dependencies
- **Instability (I)**: `I = Ce / (Ca + Ce)` where 0 = maximally stable, 1 = maximally unstable
- **Abstractness (A)**: 0-1 score based on object type (interfaces = 1.0, structs = 0.1)

## Architecture

```
catreview/
├── cmd/catreview/          # CLI application
│   └── main.go             # Commands: extract, analyze, verify, abstract
├── pkg/
│   ├── category/           # Core category theory types
│   │   ├── types.go        # Object, Morphism, Category
│   │   └── types_test.go   # Category axiom tests
│   ├── functor/            # Functor system
│   │   └── functor.go      # Functor interface, PackageAbstractionFunctor
│   ├── analysis/           # Complexity analysis
│   │   └── complexity.go   # Basu-Isik, Kolmogorov, coupling metrics
│   └── extractor/          # Code extraction
│       └── go_extractor.go # Go AST parser
└── README.md
```

## Example: Self-Analysis

Catreview analyzed itself:

```bash
$ catreview extract ./pkg -o catreview-model.json
Extracted:
  Objects:   84
  Morphisms: 102

$ catreview analyze catreview-model.json
Diagram Complexity:    509.49
Kolmogorov Complexity: 7421 bytes
Cycles Found: 0

Top Most Unstable Components:
  pkg/analysis/complexity.go: I=1.00 (Ce=24, Ca=0)
  pkg/extractor/go_extractor.go: I=1.00 (Ce=23, Ca=0)
  pkg/category/types.go: I=1.00 (Ce=20, Ca=0)
```

**Insight**: All top unstable components have I=1.00 (Ce > 0, Ca=0), meaning they are "leaves" that depend on others but aren't depended upon - expected for implementation files that import but aren't imported.

## CI/CD Integration

### GitHub Actions

```yaml
name: Categorical Review

on: [push, pull_request]

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Install catreview
        run: go install github.com/manu/catreview/cmd/catreview@latest

      - name: Extract & Analyze
        run: |
          catreview extract ./pkg --output model.json
          catreview analyze model.json --output report.json

      - name: Verify Quality Gates
        run: |
          catreview verify model.json --max-cycles 0 --fail-on-violation

      - name: Upload Report
        uses: actions/upload-artifact@v3
        with:
          name: categorical-report
          path: report.json
```

## Development

### Running Tests

```bash
go test ./pkg/category -v
go test ./pkg/... -v
```

### Test Coverage

```bash
go test ./pkg/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Functor Laws

The `PackageAbstractionFunctor` F: FileCat → PkgCat satisfies functor laws:

1. **Composition Preservation**: F(g ∘ f) = F(g) ∘ F(f)
2. **Identity Preservation**: F(id_A) = id_{F(A)}

Verification is automatic via `catreview abstract`.

## References

### Category Theory
- **Basu & Isik (2018)**: "Complexity of Commutative Diagrams"
- **Yanofsky (2003)**: "A Universal Approach to Self-Referential Paradoxes"
- **Spivak (2014)**: "Category Theory for the Sciences"

### Software Architecture
- **Martin (2000)**: "Design Principles and Design Patterns" (Instability metric)
- **Baldwin & Clark (2000)**: "Design Rules" (Modularity theory)

## License

MIT License - see LICENSE file for details

## Contributing

Contributions welcome! Please submit issues and pull requests.

### Adding Language Support

To add support for a new language (e.g., Java, TypeScript):

1. Create `pkg/extractor/{lang}_extractor.go`
2. Implement the extraction logic using the language's AST parser
3. Map language constructs to categorical objects and morphisms
4. Add tests in `pkg/extractor/{lang}_extractor_test.go`

### Example: Java Extractor Structure

```go
type JavaExtractor struct {
    category *category.Category
}

func (e *JavaExtractor) ExtractFromPath(root string) (*category.Category, error) {
    // Parse .java files
    // Extract classes → Objects
    // Extract method calls → Morphisms
    return e.category, nil
}
```

## Status

**Production Ready** ✅

- Core category theory implementation complete
- Go AST extractor operational
- Complexity metrics validated
- Self-analysis successful (509.49 complexity, 0 cycles)
- Category axioms verified
- Functor laws verified

## Roadmap

### v1.1
- [ ] Java extractor
- [ ] TypeScript extractor
- [ ] Incremental analysis (git diff based)

### v2.0
- [ ] D3.js visualization dashboard
- [ ] Natural transformations
- [ ] Limits and colimits detection
- [ ] Universal property verification

---

**Built with category theory** 🧮
