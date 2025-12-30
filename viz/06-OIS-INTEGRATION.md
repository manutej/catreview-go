# OIS Integration for Enhanced Visualization

Leveraging the OIS (Orchestration Intelligence System) framework to map complex dependencies.

---

## OIS Overview

The **OIS framework** at `/Users/manu/Documents/LUXOR/OIS/` provides structured agents following the 7-function pattern:

```
OBSERVE → REASON → GENERATE → REFINE → OPTIMIZE → INTEGRATE → REFLECT
```

Relevant agents for dependency visualization:

| Agent | Foundation | Function | Purpose |
|-------|------------|----------|---------|
| `dependency-observer` | systems-thinking | OBSERVE | Map dependency graphs |
| `emergence-analyzer` | systems-thinking | REASON | Analyze emergent patterns |
| `architecture-generator` | systems-thinking | GENERATE | Generate improved architectures |

---

## Mapping catreview to OIS

### Current catreview Output

```json
{
  "category_stats": { "objects": 505, "morphisms": 769 },
  "coupling_metrics": { ... },
  "top_coupled": [ ... ],
  "top_unstable": [ ... ],
  "cycles": null
}
```

### OIS DependencyObserver Output Schema

```typescript
{
  data: {
    nodes: ComponentNode[];
    edges: DependencyEdge[];
    layers: ArchitecturalLayer[];
    cycles: CircularDependency[];
    flows: DataFlow[];
  };
  metadata: {
    totalComponents: number;
    totalDependencies: number;
    maxDepth: number;
    cyclomaticComplexity: number;
  };
  quality: {
    completeness: number;
    confidence: number;
  };
}
```

### Transformation: catreview → OIS Format

```python
def catreview_to_ois(catreview_data):
    """Transform catreview analysis to OIS DependencyGraph format"""

    # Extract nodes from coupling_metrics
    nodes = []
    for obj_id, metrics in catreview_data['coupling_metrics'].items():
        node_type = classify_node_type(obj_id)
        layer = assign_layer(obj_id, metrics)

        nodes.append({
            'id': obj_id,
            'type': node_type,
            'layer': layer,
            'cohesion': 1.0 - metrics['instability'],  # Inverse approximation
            'afferent': metrics['afferent_coupling'],
            'efferent': metrics['efferent_coupling'],
            'instability': metrics['instability'],
            'abstractness': metrics['abstractness']
        })

    # Generate edges (from top_coupled relationships)
    edges = generate_edges(catreview_data)

    # Assign layers
    layers = [
        {'name': 'foundations', 'components': [], 'depth': 0},
        {'name': 'core', 'components': [], 'depth': 1},
        {'name': 'support', 'components': [], 'depth': 2},
        {'name': 'applications', 'components': [], 'depth': 3}
    ]

    for node in nodes:
        layers[node['layer']]['components'].append(node['id'])

    return {
        'data': {
            'nodes': nodes,
            'edges': edges,
            'layers': layers,
            'cycles': catreview_data.get('cycles', []) or [],
            'flows': []  # Would need model.json for detailed flows
        },
        'metadata': {
            'totalComponents': catreview_data['category_stats']['objects'],
            'totalDependencies': catreview_data['category_stats']['morphisms'],
            'maxDepth': 4,
            'cyclomaticComplexity': int(catreview_data['diagram_complexity'])
        },
        'quality': {
            'completeness': 0.95,
            'confidence': 0.90
        }
    }

def classify_node_type(obj_id):
    """Classify node type from object ID"""
    if obj_id.startswith('import:'):
        return 'import'
    elif obj_id.startswith('examples/'):
        return 'application'
    elif '.' in obj_id and not obj_id.endswith('.go'):
        return 'function'
    else:
        return 'file'

def assign_layer(obj_id, metrics):
    """Assign architectural layer based on coupling patterns"""
    if obj_id.startswith('import:'):
        return 0  # Foundations
    elif metrics['efferent_coupling'] >= 40:
        return 1  # Core
    elif obj_id.startswith('examples/'):
        return 3  # Applications
    else:
        return 2  # Support
```

---

## OIS Agent Integration

### Using dependency-observer with catreview Data

```python
# Invoke OIS dependency-observer via Task tool
task_prompt = """
Analyze the following catreview analysis data using the dependency-observer pattern.

Input data (catreview format):
{catreview_json}

Required analysis:
1. Map all dependency relationships
2. Identify architectural layers
3. Detect coupling patterns (tight, loose, decoupled)
4. Calculate instability metrics for components
5. Identify potential architectural violations

Output the analysis in OIS Observable<DependencyGraph> format.
"""

# The agent would:
# 1. Parse the catreview data
# 2. Transform to OIS schema
# 3. Apply systems-thinking patterns
# 4. Generate enhanced visualization data
```

### Using emergence-analyzer for Pattern Detection

```python
emergence_prompt = """
Given the following dependency graph (OIS format), analyze emergent patterns:

{ois_dependency_graph}

Identify:
1. Hub components (high centrality)
2. Bottleneck points
3. Cohesion clusters
4. Architectural layer violations
5. Emergent subsystem boundaries

Apply the EmergenceAnalyzer pattern to produce Reasoning<EmergentPatterns>.
"""
```

---

## Enhanced Visualization Pipeline

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                    OIS-ENHANCED VISUALIZATION PIPELINE                        │
└──────────────────────────────────────────────────────────────────────────────┘

┌─────────────┐      ┌─────────────────┐      ┌────────────────────┐
│  catreview  │─────→│ dependency-     │─────→│ emergence-         │
│  analysis   │      │ observer        │      │ analyzer           │
│  .json      │      │ (OBSERVE)       │      │ (REASON)           │
└─────────────┘      └────────┬────────┘      └──────────┬─────────┘
                              │                          │
                              ▼                          ▼
                    ┌─────────────────┐      ┌────────────────────┐
                    │ DependencyGraph │      │ EmergentPatterns   │
                    │                 │      │                    │
                    │ • nodes         │      │ • hubs             │
                    │ • edges         │      │ • bottlenecks      │
                    │ • layers        │      │ • violations       │
                    │ • cycles        │      │ • clusters         │
                    └────────┬────────┘      └──────────┬─────────┘
                             │                          │
                             └──────────┬───────────────┘
                                        │
                                        ▼
                              ┌────────────────────┐
                              │ Visualization      │
                              │ Generator          │
                              │                    │
                              │ • Mermaid          │
                              │ • D3.js            │
                              │ • ASCII            │
                              │ • Graphviz         │
                              └────────────────────┘
                                        │
                                        ▼
                              ┌────────────────────┐
                              │ Enhanced Output    │
                              │                    │
                              │ • Layer diagrams   │
                              │ • Hub highlights   │
                              │ • Pattern overlays │
                              │ • Violation markers│
                              └────────────────────┘
```

---

## OIS-Enhanced Visualization Features

### 1. Hub Detection Overlay

OIS emergence-analyzer identifies hubs using centrality metrics:

```
╔═══════════════════════════════════════════════════════════════════════════════╗
║                    HUB DETECTION OVERLAY                                       ║
╠═══════════════════════════════════════════════════════════════════════════════╣
║                                                                               ║
║  Layer 3: Applications                                                        ║
║  ┌─────────────────────────────────────────────────────────────────────────┐ ║
║  │  ○ views/main.go  ○ autocomplete  ○ cellbuffer  ○ chat  ○ ...          │ ║
║  └───────────────────────────────────────────────────┬─────────────────────┘ ║
║                                                       │                       ║
║  Layer 2: Support                                     │                       ║
║  ┌──────────────────────────────────────────────────┐│                       ║
║  │  ○ screen.go      ○ key.go        ○ options.go  ││                       ║
║  └────────────────────────────────────────────────┬─┘│                       ║
║                                                   │   │                       ║
║  Layer 1: Core (HUBS DETECTED)                    │   │                       ║
║  ┌────────────────────────────────────────────────┼───┼───────────────────┐ ║
║  │                                                │   │                   │ ║
║  │  ╔═══════════════════════╗    ╔═══════════════════════════════════╗   │ ║
║  │  ║ ◉ tea.go              ║    ║ ◉ standard_renderer.go            ║   │ ║
║  │  ║   HUB SCORE: 0.85     ║    ║   HUB SCORE: 0.82                 ║   │ ║
║  │  ║   PageRank: 0.12      ║    ║   PageRank: 0.11                  ║   │ ║
║  │  ║   Betweenness: 0.45   ║    ║   Betweenness: 0.42               ║   │ ║
║  │  ╚═══════════════════════╝    ╚═══════════════════════════════════╝   │ ║
║  │                                                                       │ ║
║  └───────────────────────────────────────────────────────────────────────┘ ║
║                                                                               ║
║  Layer 0: Foundations                                                         ║
║  ┌─────────────────────────────────────────────────────────────────────────┐ ║
║  │  ◉ fmt (A:58)   ◉ bubbletea (A:52)   ○ os (A:42)   ○ lipgloss (A:30)   │ ║
║  │  HUB: 0.78      HUB: 0.72                                               │ ║
║  └─────────────────────────────────────────────────────────────────────────┘ ║
║                                                                               ║
║  Legend: ◉ = Hub (centrality > 0.5)   ○ = Non-hub                            ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
```

### 2. Architectural Violation Markers

```
╔═══════════════════════════════════════════════════════════════════════════════╗
║                    ARCHITECTURAL VIOLATION DETECTION                          ║
╠═══════════════════════════════════════════════════════════════════════════════╣
║                                                                               ║
║  CLEAN ARCHITECTURE:                                                          ║
║  ───────────────────                                                         ║
║                                                                               ║
║  Layer 3 → Layer 2 → Layer 1 → Layer 0                                       ║
║  (Apps)    (Support)  (Core)    (Imports)                                    ║
║                                                                               ║
║  ✓ No violations detected in bubbletea                                       ║
║  ✓ All dependencies flow downward (DAG verified)                             ║
║  ✓ No layer-skipping dependencies                                            ║
║                                                                               ║
║  HYPOTHETICAL VIOLATION EXAMPLE:                                              ║
║  ─────────────────────────────────                                           ║
║                                                                               ║
║  If screen.go directly imported from examples/:                               ║
║                                                                               ║
║  Layer 3: Applications                                                        ║
║  ┌───────────────────┐                                                       ║
║  │ examples/foo.go   │                                                       ║
║  └─────────┬─────────┘                                                       ║
║            │                                                                  ║
║            │ ⚠️ VIOLATION: Layer skip                                        ║
║            │    (L2 depends on L3)                                           ║
║            │                                                                  ║
║  Layer 2:  │                                                                  ║
║  ┌─────────▼─────────┐                                                       ║
║  │ screen.go         │                                                       ║
║  │ [VIOLATION]       │                                                       ║
║  └───────────────────┘                                                       ║
║                                                                               ║
║  Recommendation: Extract shared code to Layer 1 abstraction                  ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
```

### 3. Cohesion Cluster Visualization

```
╔═══════════════════════════════════════════════════════════════════════════════╗
║                    COHESION CLUSTER ANALYSIS                                   ║
╠═══════════════════════════════════════════════════════════════════════════════╣
║                                                                               ║
║  Clusters identified by EmergenceAnalyzer:                                    ║
║                                                                               ║
║  ┌─────────────────────────────────────────────────────────────────────────┐ ║
║  │ CLUSTER A: Input Handling                                               │ ║
║  │ ╭───────────────────────────────────────────────────────────╮           │ ║
║  │ │  key.go ←→ key_other.go ←→ key_windows.go ←→ key_seq.go  │           │ ║
║  │ │                                                           │           │ ║
║  │ │  Cohesion Score: 0.89 (HIGH)                              │           │ ║
║  │ │  Internal Morphisms: 24                                   │           │ ║
║  │ │  External Dependencies: 8                                 │           │ ║
║  │ ╰───────────────────────────────────────────────────────────╯           │ ║
║  └─────────────────────────────────────────────────────────────────────────┘ ║
║                                                                               ║
║  ┌─────────────────────────────────────────────────────────────────────────┐ ║
║  │ CLUSTER B: Rendering Pipeline                                           │ ║
║  │ ╭───────────────────────────────────────────────────────────╮           │ ║
║  │ │  standard_renderer.go ←→ screen.go ←→ nil_renderer.go    │           │ ║
║  │ │                                                           │           │ ║
║  │ │  Cohesion Score: 0.76 (MODERATE)                          │           │ ║
║  │ │  Internal Morphisms: 45                                   │           │ ║
║  │ │  External Dependencies: 32                                │           │ ║
║  │ ╰───────────────────────────────────────────────────────────╯           │ ║
║  └─────────────────────────────────────────────────────────────────────────┘ ║
║                                                                               ║
║  ┌─────────────────────────────────────────────────────────────────────────┐ ║
║  │ CLUSTER C: Orchestration Core                                           │ ║
║  │ ╭───────────────────────────────────────────────────────────╮           │ ║
║  │ │  tea.go ←→ options.go ←→ commands.go ←→ exec.go          │           │ ║
║  │ │                                                           │           │ ║
║  │ │  Cohesion Score: 0.81 (HIGH)                              │           │ ║
║  │ │  Internal Morphisms: 58                                   │           │ ║
║  │ │  External Dependencies: 41                                │           │ ║
║  │ ╰───────────────────────────────────────────────────────────╯           │ ║
║  └─────────────────────────────────────────────────────────────────────────┘ ║
║                                                                               ║
║  Inter-Cluster Coupling:                                                      ║
║  ─────────────────────────                                                   ║
║                                                                               ║
║  A ←──(12 morphisms)──→ C                                                    ║
║  B ←──(28 morphisms)──→ C                                                    ║
║  A ←──(4 morphisms)───→ B                                                    ║
║                                                                               ║
║  Observation: Cluster C (Orchestration) acts as bridge between A and B       ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
```

---

## Integration Code

### catreview CLI Extension

```go
// cmd/catreview/viz.go

package main

import (
    "github.com/catreview-go/pkg/ois"
    "github.com/catreview-go/pkg/viz"
)

type VizCommand struct {
    Format      string `help:"Output format (ascii, mermaid, dot, json, svg)" default:"ascii"`
    OISEnhance  bool   `help:"Use OIS agents for enhanced analysis" default:"false"`
    Threshold   int    `help:"Min coupling threshold" default:"10"`
}

func (c *VizCommand) Run(analysis *Analysis) error {
    var graph *viz.DependencyGraph

    if c.OISEnhance {
        // Transform to OIS format
        oisData := ois.TransformCatreview(analysis)

        // Run OIS agents
        observed := ois.RunDependencyObserver(oisData)
        patterns := ois.RunEmergenceAnalyzer(observed)

        // Enhance graph with OIS insights
        graph = viz.EnhanceWithOIS(observed, patterns)
    } else {
        graph = viz.FromCatreview(analysis)
    }

    // Generate output
    generator := viz.NewGenerator(c.Format)
    output, err := generator.Generate(graph, &viz.Options{
        Threshold: c.Threshold,
        ShowHubs:  c.OISEnhance,
        ShowClusters: c.OISEnhance,
    })

    if err != nil {
        return err
    }

    fmt.Println(output)
    return nil
}
```

### OIS Agent Invocation (via Task Tool)

```python
# For Claude Code integration
def invoke_ois_enhancement(catreview_json):
    """Invoke OIS agents via Task tool for enhanced visualization"""

    # Step 1: dependency-observer
    observer_result = Task(
        subagent_type="dependency-observer",
        prompt=f"""
        Analyze this catreview data and produce OIS DependencyGraph:

        {catreview_json}

        Apply dependency-observer patterns:
        1. Extract all dependency relationships
        2. Identify architectural layers
        3. Calculate coupling metrics
        4. Detect any circular dependencies
        """
    )

    # Step 2: emergence-analyzer
    emergence_result = Task(
        subagent_type="emergence-analyzer",
        prompt=f"""
        Given this DependencyGraph, identify emergent patterns:

        {observer_result}

        Apply emergence-analyzer patterns:
        1. Identify hub components (high centrality)
        2. Detect cohesion clusters
        3. Find architectural violations
        4. Discover implicit boundaries
        """
    )

    return merge_results(observer_result, emergence_result)
```

---

## Summary

OIS integration enhances catreview visualization by:

1. **Standardized Schema**: Transform catreview output to OIS DependencyGraph format
2. **Pattern Detection**: Use emergence-analyzer to identify hubs, clusters, violations
3. **Enhanced Overlays**: Add hub markers, cluster boundaries, violation warnings
4. **Quality Metrics**: Include confidence scores from OIS agents
5. **Extensibility**: Future integration with architecture-generator for refactoring suggestions

The OIS framework's 7-function pattern (OBSERVE → REASON → GENERATE → ...) provides a systematic approach to analyzing and visualizing complex dependencies.

---

## Next Steps (Implementation)

1. Create `pkg/ois/transform.go` for catreview → OIS conversion
2. Add `--ois-enhance` flag to catreview viz command
3. Implement hub/cluster overlay rendering
4. Write tests with bubbletea dataset
5. Document OIS agent integration patterns
