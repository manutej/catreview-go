# Data Structure Analysis

Understanding the catreview JSON structure to identify visualization opportunities.

---

## analysis.json Structure

```json
{
  "category_stats": {
    "identities": 505,    // Self-morphisms (objects)
    "morphisms": 769,     // All morphisms
    "objects": 505        // Unique objects
  },
  "diagram_complexity": 3075.12,    // Fan-in x fan-out summation
  "kolmogorov_complexity": 49163,   // Information density metric

  "coupling_metrics": {
    "<object_id>": {
      "object_id": "string",
      "afferent_coupling": number,   // Incoming dependencies
      "efferent_coupling": number,   // Outgoing dependencies
      "instability": number,         // I = Ce / (Ca + Ce), range [0,1]
      "abstractness": number         // Interface vs implementation
    }
  },

  "cycles": null,  // Array of cycle paths, or null if DAG

  "top_unstable": [
    // Objects with highest instability (I = 1.0)
    { "object_id": "...", "efferent_coupling": N, "instability": 1 }
  ],

  "top_coupled": [
    // Objects with highest total coupling
    { "object_id": "...", "afferent_coupling": N } // or efferent_coupling
  ]
}
```

---

## Object Types Identified

From the Bubbletea analysis:

### 1. Files (Source Code)
```
tea.go                    E:53 (core orchestrator)
standard_renderer.go      E:55 (rendering engine)
screen.go                 E:32 (screen management)
key.go                    E:18 (input handling)
examples/*/main.go        E:2-19, I:1.0 (isolated apps)
```

### 2. Imports (Dependencies)
```
import:fmt                A:58 (most used)
import:os                 A:42 (OS integration)
import:github.com/charmbracelet/bubbletea    A:52 (framework)
import:github.com/charmbracelet/lipgloss     A:30 (styling)
import:time               A:25 (timing)
import:strings            A:21 (string ops)
```

### 3. Functions/Methods
```
main.*model.updateInputs      A:1, E:1, I:0.5
main.SleepPrintln             A:1, E:0, I:0, abstractness:0.5
tea.Program.Run               (implicit from model)
```

### 4. Types (Structs/Interfaces)
```
main.cellbuffer               abstractness:0.1
main.*progressWriter.Write    abstractness:0.5 (interface method)
```

---

## Visualization Opportunities by Data Type

### A. Coupling Distribution Chart

**Data Required**:
- `coupling_metrics` values
- Filter by object type (file vs import vs function)

**Visualization**: Scatter plot
- X-axis: Efferent coupling (outgoing)
- Y-axis: Afferent coupling (incoming)
- Color: Instability (red = 1.0, green = 0.0)
- Size: Total coupling (Ca + Ce)

**Insight**: Identify hub components (high both) vs terminal nodes (high efferent only)

---

### B. Stratified Dependency Graph

**Data Required**:
- All objects with `efferent_coupling > threshold`
- Implicit edges: file -> import, file -> file (via shared imports)

**Visualization**: Layered DAG
```
Layer 0: Terminal Primitives (imports with A > 20, E = 0)
Layer 1: Core Framework (tea.go, standard_renderer.go)
Layer 2: Supporting Files (screen.go, key.go, options.go)
Layer 3: Examples (all apps with I = 1.0)
```

**Algorithm**:
1. Sort by `afferent_coupling` descending (stable foundations first)
2. Build topological order
3. Assign layers based on max dependency depth

---

### C. Instability Heatmap

**Data Required**:
- `top_unstable` array
- `coupling_metrics` for all files

**Visualization**: Matrix heatmap
- Rows: All non-import objects
- Columns: Metrics (Ca, Ce, I, A)
- Color intensity: Normalized value

**Insight**: Quick scan for refactoring candidates (high I, low A = bad)

---

### D. Import Dependency Sunburst

**Data Required**:
- All `import:*` objects
- Group by package namespace

**Visualization**: Hierarchical sunburst/treemap
```
github.com/
  charmbracelet/
    bubbletea (A:52)
    lipgloss (A:30)
    bubbles/
      key (A:8)
      spinner (A:7)
      progress (A:5)
```

**Insight**: External dependency footprint and concentration

---

## Key Metrics for Filtering

### High-Impact Thresholds

| Metric | Threshold | Count in Bubbletea | Purpose |
|--------|-----------|-------------------|---------|
| Efferent > 20 | 5 files | `tea.go`, `standard_renderer.go`, `screen.go`, etc. | Core hub identification |
| Afferent > 20 | 6 imports | `fmt`, `os`, `bubbletea`, `lipgloss`, `time`, `strings` | Foundation identification |
| Instability = 1.0 | ~50 objects | All examples + most files | Terminal node identification |
| Abstractness > 0.3 | ~20 functions | Interface implementations | Abstraction points |

### Recommended Filter Sets

**Clean Overview (< 30 nodes)**:
```
WHERE (efferent_coupling > 15) OR (afferent_coupling > 25)
```

**Framework Focus (< 50 nodes)**:
```
WHERE NOT object_id LIKE 'examples/%'
  AND NOT object_id LIKE 'import:%'
  AND efferent_coupling > 5
```

**Full Architecture (< 100 nodes)**:
```
WHERE (efferent_coupling > 10) OR (afferent_coupling > 10)
```

---

## Data Transformations Needed

### 1. Edge Generation

The analysis.json doesn't store explicit edges. Generate from:
- File imports (parse `import:*` afferent sources)
- Function calls (match function object_ids to their files)
- Shared imports (files using same import = implicit coupling)

### 2. Layer Assignment

```python
def assign_layer(obj):
    if obj['afferent_coupling'] > 0 and obj['efferent_coupling'] == 0:
        return 0  # Terminal import/interface
    elif obj['object_id'].startswith('import:'):
        return 0  # Import layer
    elif obj['instability'] == 0:
        return 1  # Stable core
    elif 'examples/' in obj['object_id']:
        return 3  # Application layer
    else:
        return 2  # Supporting layer
```

### 3. Normalization for Heatmaps

```python
max_ca = max(obj['afferent_coupling'] for obj in objects)
max_ce = max(obj['efferent_coupling'] for obj in objects)

def normalize(obj):
    return {
        'ca_norm': obj['afferent_coupling'] / max_ca,
        'ce_norm': obj['efferent_coupling'] / max_ce,
        'instability': obj['instability'],  # Already normalized
        'distance': abs(obj['instability'] + obj['abstractness'] - 1)  # Main sequence
    }
```

---

## Model.json Structure (Full Graph)

The model.json contains the complete categorical structure:

```json
{
  "objects": [
    { "id": "tea.go", "type": "file", ... },
    { "id": "tea.Program.Run", "type": "function", ... }
  ],
  "morphisms": [
    { "source": "tea.go", "target": "import:fmt", "type": "import" },
    { "source": "tea.Program.Run", "target": "tea.go", "type": "contains" }
  ]
}
```

**Key**: This provides explicit edges that analysis.json summarizes into coupling metrics.

---

## Recommended Processing Pipeline

```
model.json ──┬──> Filter by coupling threshold
             │
             ├──> Extract explicit morphisms
             │
             └──> Group objects by type (file/import/function)
                        │
                        v
              ┌─────────────────────┐
              │  Layered Graph Gen  │
              │  (topological sort) │
              └──────────┬──────────┘
                         │
          ┌──────────────┼──────────────┐
          v              v              v
    ASCII Diagram   Mermaid.js    D3.js/Vis.js
    (terminal)      (GitHub)      (interactive)
```

---

## Next Steps

1. **02-LAYERED-GRAPH-APPROACH.md**: Design stratified visualization
2. **03-COUPLING-HEATMAPS.md**: Build multi-dimensional coupling display
3. **04-CATEGORICAL-FLOW.md**: Animate Elm architecture data flow
