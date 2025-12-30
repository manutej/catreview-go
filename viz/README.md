# Categorical Visualization Research

**Purpose**: Research cleaner, more expressive visualizations for mapping complex categorical dependencies from catreview analysis data.

**Data Sources**:
- `examples/bubbletea/analysis.json` (769 morphisms, 505 objects)
- `examples/bubbletea/model.json` (full categorical model)
- `docs/BUBBLETEA-ARCHITECTURE.md` (existing visualizations)

**Key Challenges**:
1. High morphism count (769) creates visual noise in full graphs
2. Coupling metrics need multi-dimensional representation (afferent/efferent/instability)
3. Categorical structures (functors, products, coproducts) need intuitive mapping
4. Clean separation between framework core and isolated examples

---

## Research Documents

| File | Focus | Status |
|------|-------|--------|
| `01-DATA-STRUCTURE-ANALYSIS.md` | JSON schema understanding | Complete |
| `02-LAYERED-GRAPH-APPROACH.md` | Hierarchical visualization | Complete |
| `03-COUPLING-HEATMAPS.md` | Multi-dimensional coupling viz | Complete |
| `04-CATEGORICAL-FLOW.md` | Elm architecture animation | Complete |
| `05-TOOLING-RECOMMENDATIONS.md` | Implementation tools | Complete |
| `06-OIS-INTEGRATION.md` | Leveraging OIS dependency-observer | Complete |

---

## Quick Reference: Visualization Approaches

### 1. Filtered Layered Graphs
Show only high-impact dependencies (efferent > 10) with stratified layers.

### 2. Coupling Heatmaps
2D matrix with color-coded coupling strength, sorted by instability.

### 3. Radial Dependency Trees
Core components at center, applications at periphery.

### 4. Sankey Flow Diagrams
Data flow from framework core to terminal applications.

### 5. Interactive Force-Directed
D3.js/Vis.js for exploration with zoom/filter.

---

## Key Insight from OIS

The **dependency-observer** agent from OIS provides a structured approach:
```
System -> Observable<DependencyGraph>
```

With built-in support for:
- Coupling analysis (afferent/efferent/instability calculation)
- Layer extraction and topological sorting
- Cycle detection
- Data flow mapping

This aligns perfectly with catreview's categorical model.
