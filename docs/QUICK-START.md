# Catreview-Go Quick Start Guide

**Get started with categorical codebase analysis in 5 minutes.**

---

## TL;DR

```bash
# 1. Build the tool
git clone https://github.com/yourusername/catreview-go
cd catreview-go
go build -o catreview ./cmd/catreview

# 2. Analyze your Go codebase
./catreview extract ./path/to/your/go/project -o model.json
./catreview analyze model.json

# 3. Read the report
cat report.json | jq '.complexity, .top_unstable[0:5]'
```

**Done!** You now have categorical insights into your codebase.

---

## Installation

### Option 1: Build from Source (Recommended)

```bash
git clone https://github.com/yourusername/catreview-go
cd catreview-go
go build -o catreview ./cmd/catreview

# Optional: Install globally
sudo mv catreview /usr/local/bin/
```

### Option 2: Go Install (Future)

```bash
go install github.com/yourusername/catreview-go/cmd/catreview@latest
```

### Verify Installation

```bash
catreview --version
# Expected: catreview version 1.0.0
```

---

## Your First Analysis (3 Steps)

### Step 1: Extract Categorical Model

Point catreview at your Go project:

```bash
cd /path/to/your/go/project
catreview extract ./pkg -o model.json
```

**Output:**
```
Extracting categorical model from: ./pkg
Parsed 42 Go files
Extracted:
  Objects:   215
  Morphisms: 342
  Identities: 215
Model saved to: model.json
```

**What happened?**
- Catreview parsed your Go code using `go/ast`
- Extracted **Objects** (packages, types, functions)
- Extracted **Morphisms** (imports, function calls, dependencies)
- Created identity morphisms for each object
- Saved categorical model as JSON

---

### Step 2: Analyze Complexity

Run complexity analysis on the extracted model:

```bash
catreview analyze model.json -o report.json --pretty
```

**Output:**
```
Categorical Analysis Report
===========================

Category Statistics:
  Objects:    215
  Morphisms:  342

Complexity Metrics:
  Diagram Complexity:    1,425.67
  Kolmogorov Complexity: 18,432 bytes

Morphisms/Objects Ratio: 1.59 ✅ (healthy: 1.0-2.0)

Dependency Analysis:
  Cycles Found: 0 ✅

Top 5 Most Unstable Components:
  1. pkg/server/handler.go: I=1.00 (Ce=15, Ca=0)
  2. pkg/client/client.go:  I=1.00 (Ce=12, Ca=0)
  3. pkg/utils/helpers.go:  I=1.00 (Ce=8, Ca=0)
  ...

Top 5 Most Coupled Components:
  1. pkg/server/handler.go: 15 total (Ce=15, Ca=0)
  2. pkg/client/client.go:  12 total (Ce=12, Ca=0)
  ...

Full report saved to: report.json
```

**What happened?**
- Computed **Basu-Isik diagram complexity** (1,425.67)
- Estimated **Kolmogorov complexity** via gzip (18,432 bytes)
- Calculated **M/O ratio** (1.59 - healthy!)
- Detected **dependency cycles** (0 - good!)
- Identified **unstable components** (I=1.00 means leaf modules)
- Found **coupling hotspots** (handler.go has 15 dependencies)

---

### Step 3: Verify Quality

Check if your architecture satisfies category theory axioms:

```bash
catreview verify model.json --max-cycles 0 --fail-on-violation
```

**Output:**
```
Verifying category axioms: model.json

Checking associativity law: (h ∘ g) ∘ f = h ∘ (g ∘ f)
  Verified 1,284 composition chains
  ✅ Associativity holds

Checking identity law: f ∘ id_A = f and id_B ∘ f = f
  Verified 342 morphisms
  ✅ Identity holds

Category axioms verified successfully ✅

Checking for dependency cycles (max allowed: 0)...
  Found 0 cycles
  ✅ Cycle count within limit

Verification passed ✅
```

**What happened?**
- Verified **associativity** of composition chains
- Verified **identity** morphisms behave correctly
- Detected **dependency cycles** (none found)
- Exit code 0 = all checks passed

---

## Understanding the Output

### Complexity Metrics Explained

| Metric | What It Means | Good Range | Your Value |
|--------|---------------|------------|------------|
| **Diagram Complexity** | Sum of object + morphism + composition complexity | N/A (relative) | 1,425.67 |
| **Kolmogorov Complexity** | Compressed size (information content) | N/A (relative) | 18,432 bytes |
| **M/O Ratio** | Morphisms per Object (interaction density) | 1.0 - 2.0 | 1.59 ✅ |
| **Cycles** | Circular dependencies | 0 ideal | 0 ✅ |

### Coupling Metrics Explained

| Metric | Formula | Meaning | Range |
|--------|---------|---------|-------|
| **Afferent Coupling (Ca)** | Incoming deps | How many depend ON you | 0 - ∞ |
| **Efferent Coupling (Ce)** | Outgoing deps | How many you depend ON | 0 - ∞ |
| **Instability (I)** | Ce / (Ca + Ce) | 0 = stable, 1 = unstable | 0.0 - 1.0 |

**Interpretation**:
- **I = 1.00** (Ce > 0, Ca = 0): Leaf module - normal for implementation files
- **I = 0.00** (Ce = 0, Ca > 0): Root module - normal for interfaces
- **I ≈ 0.50**: Balanced - both depends on others AND is depended upon

### What to Look For

#### ✅ Good Signs
- **M/O ratio 1.0-2.0**: Healthy interaction density
- **0 cycles**: Clean dependency structure
- **Linear complexity scaling**: No "complexity debt"
- **Leaf modules I=1.00**: Proper separation of concerns

#### ⚠️ Warning Signs
- **M/O ratio > 2.0**: Over-coupled, consider refactoring
- **Cycles > 0**: Circular dependencies (unless domain-driven like Git trees)
- **Total coupling > 50**: Coordination bottleneck
- **I ≈ 0.5 everywhere**: Lack of architectural layers

---

## Common Use Cases

### Use Case 1: Pre-Commit Quality Gate

Add to `.git/hooks/pre-commit`:

```bash
#!/bin/bash
catreview extract ./pkg -o /tmp/model.json
catreview verify /tmp/model.json --max-cycles 0 --fail-on-violation || exit 1
```

**Blocks commits** if dependency cycles are introduced.

---

### Use Case 2: CI/CD Integration

Add to `.github/workflows/catreview.yml`:

```yaml
name: Categorical Review
on: [push, pull_request]

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Install catreview
        run: |
          git clone https://github.com/yourusername/catreview-go
          cd catreview-go && go build -o catreview ./cmd/catreview
          sudo mv catreview /usr/local/bin/

      - name: Analyze
        run: |
          catreview extract ./pkg -o model.json
          catreview analyze model.json -o report.json --pretty

      - name: Verify Quality Gates
        run: catreview verify model.json --max-cycles 0 --fail-on-violation

      - name: Upload Report
        uses: actions/upload-artifact@v3
        with:
          name: categorical-report
          path: report.json
```

---

### Use Case 3: Track Complexity Over Time

```bash
#!/bin/bash
# track-complexity.sh

DATE=$(date +%Y-%m-%d)
catreview extract ./pkg -o model.json
COMPLEXITY=$(catreview analyze model.json | grep "Diagram Complexity" | awk '{print $3}')

echo "$DATE,$COMPLEXITY" >> complexity-history.csv

# Plot with gnuplot
gnuplot -e "set datafile separator ','; plot 'complexity-history.csv' using 2 with lines"
```

**Result**: Trend line showing complexity growth over commits.

---

### Use Case 4: Find Refactoring Targets

```bash
catreview analyze model.json --pretty | jq '.top_coupled[0:3]'
```

**Output:**
```json
[
  {
    "component": "pkg/server/handler.go",
    "total_coupling": 42,
    "efferent": 42,
    "afferent": 0,
    "instability": 1.00
  },
  ...
]
```

**Action**: Refactor `handler.go` by extracting dependencies to separate modules.

---

## Advanced: Package-Level Abstraction

Create a higher-level view by mapping files to packages:

```bash
catreview abstract model.json -o abstract.json --pretty
```

**Output:**
```
Creating package-level abstraction from: model.json

Applying functor: File → Package
Mapping 215 files to 12 packages...
Mapping 342 file dependencies to package dependencies...

Verifying functor laws:
  ✅ Composition preservation: F(g ∘ f) = F(g) ∘ F(f)
  ✅ Identity preservation:    F(id_A) = id_{F(A)}

Functor laws verified ✅

Abstracted Category:
  Packages:             12
  Package Dependencies: 18

Abstracted model saved to: abstract.json
```

**What happened?**
- Applied a **functor** F: FileCat → PkgCat
- Mapped 215 files → 12 packages
- Mapped 342 file deps → 18 package deps
- Verified functor satisfies category theory laws

**Now analyze at package level:**

```bash
catreview analyze abstract.json --pretty
```

**Result**: Architectural view focusing on package interactions, not file details.

---

## Interpreting Real Examples

### Example 1: Glow (Simple CLI)

```
Objects:    215
Morphisms:  238
M/O Ratio:  1.11
Complexity: 1,169.34
Cycles:     0
```

**Interpretation**: Minimal interaction density (1.11), lowest complexity - characteristic of simple CLI tools with linear pipelines.

---

### Example 2: Bubbletea (TUI Framework)

```
Objects:    308
Morphisms:  567
M/O Ratio:  1.84
Complexity: 2,077.66
Cycles:     0
```

**Interpretation**: High M/O ratio (1.84) reflects Elm Architecture's message-driven design - many morphisms (messages) between objects (components).

---

### Example 3: Crush (AI CLI)

```
Objects:    2,372
Morphisms:  3,815
M/O Ratio:  1.61 ≈ φ (golden ratio!)
Complexity: 15,677.02
Cycles:     0
```

**Interpretation**: Golden ratio M/O suggests natural aesthetic balance in architecture. Linear complexity scaling confirms compositional consistency.

---

### Example 4: Soft-Serve (Git Server)

```
Objects:    1,191
Morphisms:  1,874
M/O Ratio:  1.57
Complexity: 7,787.00
Cycles:     12 (all self-loops in Git types)
```

**Interpretation**: 12 cycles are **expected** - Git's tree/commit/reference types are inherently recursive. Domain-driven design, not architectural smell.

---

## Troubleshooting

### Issue: "No objects extracted"

**Problem**: Model shows 0 objects after extraction.

**Solution**: Check you're pointing to Go source files:
```bash
# ❌ Wrong - no .go files
catreview extract ./bin -o model.json

# ✅ Correct - contains .go files
catreview extract ./pkg -o model.json
catreview extract ./src -o model.json
```

---

### Issue: "Category axiom violation"

**Problem**: Verification fails with associativity or identity errors.

**Solution**: This indicates a bug in the extractor. Report as issue with:
```bash
catreview extract ./pkg -o bug-model.json
# Attach bug-model.json to issue report
```

---

### Issue: "Too many cycles"

**Problem**: `verify` reports unexpected cycles.

**Solution**: Check if cycles are domain-driven (like Git trees):
```bash
catreview analyze model.json | jq '.cycles'
```

If cycles are in recursive data structures → **expected**
If cycles are in business logic → **refactoring target**

---

### Issue: "Model file too large"

**Problem**: `model.json` is 100+ MB for large codebases.

**Solution**: Use `--no-pretty` to reduce size:
```bash
catreview extract ./pkg -o model.json  # No --pretty flag
```

Or analyze specific packages:
```bash
catreview extract ./pkg/server -o server-model.json
catreview extract ./pkg/client -o client-model.json
```

---

## What's Next?

### Learn More
- **README.md** - Complete user guide
- **PRODUCTION-GUIDE.md** - Validation results, real-world examples
- **examples/** - Analysis results from 5 real repositories

### Contribute
- Add TypeScript extractor (v1.1 roadmap)
- Add Java extractor (v1.1 roadmap)
- Create D3.js visualization dashboard (v1.2 roadmap)

### Apply
- Run on your codebase
- Set up CI/CD quality gates
- Track complexity over time
- Identify refactoring targets

---

## Summary

**In 5 minutes, you learned to:**

1. ✅ **Extract** categorical models from Go code
2. ✅ **Analyze** complexity and coupling metrics
3. ✅ **Verify** category theory axioms
4. ✅ **Interpret** M/O ratios, instability, cycles
5. ✅ **Integrate** into CI/CD pipelines

**Key Takeaways:**

- **M/O ratio 1.0-2.0** = healthy interaction
- **0 cycles** = clean dependencies (unless domain-driven)
- **I=1.00 for leaves** = proper layering
- **Linear complexity** = no technical debt

**Go analyze your codebase now!** 🚀

```bash
catreview extract ./your-project -o model.json
catreview analyze model.json
```

---

**Questions?** Open an issue: https://github.com/yourusername/catreview-go/issues
