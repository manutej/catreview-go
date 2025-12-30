# Visualization Tooling Recommendations

Implementation tools and libraries for generating clean dependency visualizations.

---

## Tool Categories

| Category | Use Case | Recommended Tools |
|----------|----------|-------------------|
| ASCII Diagrams | Terminal, README, docs | Custom generators, boxes.py |
| Static Graphs | GitHub, docs, reports | Mermaid.js, Graphviz/DOT |
| Interactive | Web dashboards | D3.js, Vis.js, Cytoscape.js |
| Heatmaps | Coupling analysis | Plotly, matplotlib, D3.js |
| Animation | Flow visualization | D3.js, anime.js, Framer Motion |

---

## 1. ASCII Diagram Generation

### Python Library: `rich` + Custom Box Drawing

```python
from rich.console import Console
from rich.table import Table
from rich.panel import Panel
from rich.tree import Tree

def generate_layer_tree(analysis_data):
    """Generate tree representation of architectural layers"""
    console = Console()

    tree = Tree("[bold blue]BUBBLETEA ARCHITECTURE[/]")

    # Layer 0: Imports
    imports = tree.add("[green]Layer 0: Stable Imports (I=0)[/]")
    for obj in analysis_data['top_coupled']:
        if obj['object_id'].startswith('import:'):
            imports.add(f"[dim]{obj['object_id']}[/] A:{obj['afferent_coupling']}")

    # Layer 1: Core
    core = tree.add("[yellow]Layer 1: Framework Core[/]")
    core.add("[bold]tea.go[/] E:53")
    core.add("[bold]standard_renderer.go[/] E:55")

    # Layer 2: Support
    support = tree.add("[cyan]Layer 2: Support Modules[/]")
    for obj_id, metrics in analysis_data['coupling_metrics'].items():
        if (not obj_id.startswith('import:') and
            not obj_id.startswith('examples/') and
            metrics['efferent_coupling'] > 10 and
            obj_id not in ['tea.go', 'standard_renderer.go']):
            support.add(f"{obj_id} E:{metrics['efferent_coupling']}")

    # Layer 3: Applications
    apps = tree.add("[magenta]Layer 3: Applications (I=1.0)[/]")
    for obj in analysis_data['top_unstable'][:5]:
        apps.add(f"[dim]{obj['object_id']}[/] E:{obj['efferent_coupling']}")

    console.print(tree)
```

### Output Example

```
BUBBLETEA ARCHITECTURE
├── Layer 0: Stable Imports (I=0)
│   ├── import:fmt A:58
│   ├── import:bubbletea A:52
│   ├── import:os A:42
│   └── import:lipgloss A:30
├── Layer 1: Framework Core
│   ├── tea.go E:53
│   └── standard_renderer.go E:55
├── Layer 2: Support Modules
│   ├── screen.go E:32
│   ├── key.go E:18
│   └── options.go E:21
└── Layer 3: Applications (I=1.0)
    ├── examples/views/main.go E:19
    ├── examples/autocomplete/main.go E:18
    └── examples/cellbuffer/main.go E:17
```

---

## 2. Mermaid.js Integration

### Advantages
- GitHub renders natively
- Simple text-based syntax
- Good for documentation

### Generator Script

```python
def generate_mermaid_graph(analysis_data, threshold=15):
    """Generate Mermaid.js graph definition"""

    lines = ["```mermaid", "flowchart BT"]

    # Subgraphs for layers
    lines.append('    subgraph L0["Stable Imports"]')
    for obj in analysis_data['top_coupled']:
        if obj['object_id'].startswith('import:') and obj['afferent_coupling'] >= 25:
            clean_id = obj['object_id'].replace('import:', '').replace('/', '_').replace('.', '_')
            lines.append(f'        {clean_id}["{obj["object_id"]}<br/>A:{obj["afferent_coupling"]}"]')
    lines.append('    end')

    lines.append('    subgraph L1["Framework Core"]')
    lines.append('        tea["tea.go<br/>E:53"]')
    lines.append('        renderer["standard_renderer.go<br/>E:55"]')
    lines.append('    end')

    lines.append('    subgraph L3["Applications"]')
    for obj in analysis_data['top_unstable'][:3]:
        clean_id = obj['object_id'].replace('/', '_').replace('.', '_')
        lines.append(f'        {clean_id}["{obj["object_id"]}<br/>E:{obj["efferent_coupling"]}"]')
    lines.append('    end')

    # Edges
    lines.append('    tea --> fmt')
    lines.append('    tea --> os')
    lines.append('    renderer --> fmt')
    lines.append('    renderer --> lipgloss')
    lines.append('    examples_views_main_go --> tea')

    # Styling
    lines.append('    style L0 fill:#e8f5e9')
    lines.append('    style L1 fill:#fff3e0')
    lines.append('    style L3 fill:#fce4ec')

    lines.append("```")

    return '\n'.join(lines)
```

---

## 3. Graphviz/DOT for Complex Graphs

### Advantages
- Handles large graphs well
- Multiple layout algorithms (dot, neato, fdp)
- High-quality output (SVG, PDF)

### Generator

```python
def generate_dot_graph(analysis_data):
    """Generate Graphviz DOT file"""

    lines = [
        "digraph BubbleteaDependencies {",
        "    rankdir=BT;",
        "    node [shape=box, style=filled];",
        "",
        "    // Rank groupings for layers",
        '    { rank=same; ',
    ]

    # Add imports to same rank
    import_nodes = []
    for obj in analysis_data['top_coupled']:
        if obj['object_id'].startswith('import:') and obj['afferent_coupling'] >= 20:
            node_id = obj['object_id'].replace('import:', 'imp_').replace('/', '_').replace('.', '_')
            import_nodes.append(node_id)
    lines.append('    { rank=same; ' + '; '.join(import_nodes) + '; }')

    # Core nodes
    lines.append('    { rank=same; tea; renderer; }')

    # Node definitions with attributes
    lines.append('')
    lines.append('    // Import nodes (green = stable)')
    for obj in analysis_data['top_coupled']:
        if obj['object_id'].startswith('import:') and obj['afferent_coupling'] >= 20:
            node_id = obj['object_id'].replace('import:', 'imp_').replace('/', '_').replace('.', '_')
            label = f"{obj['object_id']}\\nA:{obj['afferent_coupling']}"
            lines.append(f'    {node_id} [label="{label}", fillcolor="#c8e6c9"];')

    lines.append('')
    lines.append('    // Core nodes (orange = hub)')
    lines.append('    tea [label="tea.go\\nE:53", fillcolor="#ffe0b2"];')
    lines.append('    renderer [label="standard_renderer.go\\nE:55", fillcolor="#ffe0b2"];')

    lines.append('')
    lines.append('    // Application nodes (pink = unstable)')
    for obj in analysis_data['top_unstable'][:5]:
        node_id = obj['object_id'].replace('/', '_').replace('.', '_')
        label = f"{obj['object_id']}\\nE:{obj['efferent_coupling']}"
        lines.append(f'    {node_id} [label="{label}", fillcolor="#f8bbd0"];')

    # Edges
    lines.append('')
    lines.append('    // Dependencies')
    lines.append('    tea -> imp_fmt;')
    lines.append('    tea -> imp_os;')
    lines.append('    renderer -> imp_fmt;')
    lines.append('    renderer -> imp_lipgloss;')
    for obj in analysis_data['top_unstable'][:5]:
        node_id = obj['object_id'].replace('/', '_').replace('.', '_')
        lines.append(f'    {node_id} -> tea;')

    lines.append("}")

    return '\n'.join(lines)
```

### Generate SVG

```bash
# Generate SVG from DOT file
dot -Tsvg dependencies.dot -o dependencies.svg

# Or PNG
dot -Tpng dependencies.dot -o dependencies.png

# Different layout for force-directed
neato -Tsvg dependencies.dot -o dependencies-force.svg
```

---

## 4. D3.js for Interactive Visualization

### Force-Directed Graph

```javascript
// data: { nodes: [...], links: [...] }

const width = 960;
const height = 600;

const svg = d3.select("#graph")
    .append("svg")
    .attr("width", width)
    .attr("height", height);

const simulation = d3.forceSimulation(data.nodes)
    .force("link", d3.forceLink(data.links).id(d => d.id).distance(100))
    .force("charge", d3.forceManyBody().strength(-300))
    .force("center", d3.forceCenter(width / 2, height / 2))
    .force("y", d3.forceY().y(d => {
        // Stratify by layer
        const layerHeight = height / 4;
        return (3 - d.layer) * layerHeight;
    }).strength(0.5));

// Color scale based on instability
const colorScale = d3.scaleSequential(d3.interpolateRdYlGn)
    .domain([1, 0]);

// Node size based on total coupling
const sizeScale = d3.scaleSqrt()
    .domain([0, 60])
    .range([5, 25]);

// Draw links
const link = svg.append("g")
    .selectAll("line")
    .data(data.links)
    .join("line")
    .attr("stroke", "#999")
    .attr("stroke-opacity", 0.6)
    .attr("stroke-width", d => Math.sqrt(d.weight));

// Draw nodes
const node = svg.append("g")
    .selectAll("circle")
    .data(data.nodes)
    .join("circle")
    .attr("r", d => sizeScale(d.afferent + d.efferent))
    .attr("fill", d => colorScale(d.instability))
    .call(drag(simulation));

// Labels
const label = svg.append("g")
    .selectAll("text")
    .data(data.nodes)
    .join("text")
    .text(d => d.id.split('/').pop())
    .attr("font-size", "10px");

// Tooltip
node.on("mouseover", (event, d) => {
    tooltip.html(`
        <strong>${d.id}</strong><br/>
        Afferent: ${d.afferent}<br/>
        Efferent: ${d.efferent}<br/>
        Instability: ${d.instability.toFixed(2)}
    `);
});

simulation.on("tick", () => {
    link
        .attr("x1", d => d.source.x)
        .attr("y1", d => d.source.y)
        .attr("x2", d => d.target.x)
        .attr("y2", d => d.target.y);

    node
        .attr("cx", d => d.x)
        .attr("cy", d => d.y);

    label
        .attr("x", d => d.x + 10)
        .attr("y", d => d.y + 3);
});
```

---

## 5. Heatmap with Plotly (Python)

```python
import plotly.express as px
import plotly.graph_objects as go
import pandas as pd

def generate_coupling_heatmap(analysis_data):
    """Generate interactive coupling heatmap"""

    # Extract data
    records = []
    for obj_id, metrics in analysis_data['coupling_metrics'].items():
        if not obj_id.startswith('import:') and not obj_id.startswith('main.'):
            records.append({
                'Component': obj_id[:30],
                'Afferent': metrics['afferent_coupling'],
                'Efferent': metrics['efferent_coupling'],
                'Instability': metrics['instability'],
                'Abstractness': metrics['abstractness'],
                'Total': metrics['afferent_coupling'] + metrics['efferent_coupling']
            })

    df = pd.DataFrame(records)
    df = df.sort_values('Total', ascending=False).head(20)

    # Create heatmap
    fig = px.imshow(
        df[['Afferent', 'Efferent', 'Instability', 'Abstractness']].values,
        labels=dict(x="Metric", y="Component", color="Value"),
        x=['Afferent', 'Efferent', 'Instability', 'Abstractness'],
        y=df['Component'].tolist(),
        color_continuous_scale='RdYlGn_r',
        aspect='auto'
    )

    fig.update_layout(
        title='Coupling Metrics Heatmap',
        height=600
    )

    return fig
```

---

## 6. Cytoscape.js for Network Analysis

```javascript
const cy = cytoscape({
    container: document.getElementById('cy'),
    elements: {
        nodes: data.nodes.map(n => ({
            data: {
                id: n.id,
                label: n.id.split('/').pop(),
                afferent: n.afferent,
                efferent: n.efferent,
                instability: n.instability,
                layer: n.layer
            }
        })),
        edges: data.links.map(l => ({
            data: {
                source: l.source,
                target: l.target,
                weight: l.weight
            }
        }))
    },
    style: [
        {
            selector: 'node',
            style: {
                'background-color': d => {
                    // Color by instability: green (0) to red (1)
                    const i = d.data('instability');
                    return `rgb(${Math.round(255 * i)}, ${Math.round(255 * (1-i))}, 0)`;
                },
                'width': d => 10 + d.data('afferent') + d.data('efferent'),
                'height': d => 10 + d.data('afferent') + d.data('efferent'),
                'label': 'data(label)'
            }
        },
        {
            selector: 'edge',
            style: {
                'width': d => Math.sqrt(d.data('weight')),
                'line-color': '#ccc',
                'target-arrow-color': '#ccc',
                'target-arrow-shape': 'triangle'
            }
        }
    ],
    layout: {
        name: 'dagre',  // Good for hierarchical
        rankDir: 'BT',
        nodeSep: 50,
        rankSep: 100
    }
});

// Enable zoom/pan
cy.userZoomingEnabled(true);
cy.userPanningEnabled(true);

// Click handler for details
cy.on('tap', 'node', function(evt){
    const node = evt.target;
    showDetails(node.data());
});
```

---

## 7. Static SVG Generation (Python)

For embedding in markdown without JavaScript:

```python
import svgwrite

def generate_svg_diagram(analysis_data, width=800, height=600):
    """Generate static SVG dependency diagram"""

    dwg = svgwrite.Drawing(size=(width, height))

    # Background
    dwg.add(dwg.rect(insert=(0, 0), size=(width, height), fill='#fafafa'))

    # Layer heights
    layer_y = {0: height - 80, 1: height - 200, 2: height - 320, 3: 80}

    # Draw imports (layer 0)
    imports = [o for o in analysis_data['top_coupled']
               if o['object_id'].startswith('import:') and o['afferent_coupling'] >= 25]

    x_step = width / (len(imports) + 1)
    for i, imp in enumerate(imports):
        x = x_step * (i + 1)
        y = layer_y[0]

        # Node
        dwg.add(dwg.rect(
            insert=(x - 40, y - 15),
            size=(80, 30),
            fill='#c8e6c9',
            rx=5, ry=5
        ))
        dwg.add(dwg.text(
            imp['object_id'].replace('import:', ''),
            insert=(x, y),
            text_anchor='middle',
            font_size='10px'
        ))
        dwg.add(dwg.text(
            f"A:{imp['afferent_coupling']}",
            insert=(x, y + 12),
            text_anchor='middle',
            font_size='8px',
            fill='#666'
        ))

    # Draw core (layer 1)
    core_nodes = [
        {'id': 'tea.go', 'x': width/3, 'efferent': 53},
        {'id': 'renderer.go', 'x': 2*width/3, 'efferent': 55}
    ]

    for node in core_nodes:
        x, y = node['x'], layer_y[1]

        dwg.add(dwg.rect(
            insert=(x - 50, y - 20),
            size=(100, 40),
            fill='#ffe0b2',
            rx=5, ry=5
        ))
        dwg.add(dwg.text(
            node['id'],
            insert=(x, y),
            text_anchor='middle',
            font_size='11px',
            font_weight='bold'
        ))
        dwg.add(dwg.text(
            f"E:{node['efferent']}",
            insert=(x, y + 15),
            text_anchor='middle',
            font_size='9px',
            fill='#666'
        ))

    # Draw edges (simplified)
    # tea.go -> imports
    dwg.add(dwg.line(
        start=(width/3, layer_y[1] + 20),
        end=(x_step, layer_y[0] - 15),
        stroke='#999',
        stroke_width=1
    ))

    return dwg.tostring()
```

---

## Tool Selection Matrix

| Requirement | Recommended Tool | Alternative |
|-------------|------------------|-------------|
| README/docs (static) | Mermaid.js | Graphviz DOT |
| Terminal display | Rich (Python) | Custom ASCII |
| GitHub rendering | Mermaid.js | SVG embed |
| Interactive exploration | D3.js | Cytoscape.js |
| Large graphs (>500 nodes) | Cytoscape.js | Graphviz |
| Heatmaps | Plotly | D3.js |
| Animation | D3.js | Framer Motion |
| Offline/print | Graphviz (SVG/PDF) | Static SVG |

---

## Integration with catreview-go

### Proposed CLI Commands

```bash
# Generate ASCII tree for terminal
catreview viz --format=ascii --style=tree examples/bubbletea

# Generate Mermaid for markdown
catreview viz --format=mermaid --threshold=15 examples/bubbletea > diagram.md

# Generate DOT for Graphviz
catreview viz --format=dot examples/bubbletea | dot -Tsvg -o graph.svg

# Generate JSON for D3.js
catreview viz --format=json examples/bubbletea > data.json

# Generate SVG directly
catreview viz --format=svg --width=1200 --height=800 examples/bubbletea > graph.svg
```

### Go Implementation Sketch

```go
package viz

type Format string

const (
    FormatASCII   Format = "ascii"
    FormatMermaid Format = "mermaid"
    FormatDOT     Format = "dot"
    FormatJSON    Format = "json"
    FormatSVG     Format = "svg"
)

type Generator interface {
    Generate(analysis *Analysis, opts *Options) (string, error)
}

type Options struct {
    Format          Format
    Threshold       int     // Min coupling to include
    MaxNodes        int     // Max nodes to show
    ShowImports     bool
    ShowExamples    bool
    GroupByLayer    bool
    Width, Height   int     // For SVG
}

func NewGenerator(format Format) Generator {
    switch format {
    case FormatASCII:
        return &ASCIIGenerator{}
    case FormatMermaid:
        return &MermaidGenerator{}
    case FormatDOT:
        return &DOTGenerator{}
    case FormatJSON:
        return &JSONGenerator{}
    case FormatSVG:
        return &SVGGenerator{}
    default:
        return &ASCIIGenerator{}
    }
}
```

---

## Next Steps

- `06-OIS-INTEGRATION.md`: Leveraging OIS dependency-observer for enhanced analysis
