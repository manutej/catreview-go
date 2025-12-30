# Coupling Heatmap Visualizations

Multi-dimensional representation of coupling metrics for quick architectural assessment.

---

## Problem Statement

The current representation shows coupling as a flat table:
```
| Component | Afferent | Efferent | Instability | Abstractness |
```

**Issues**:
- Difficult to spot patterns across 505 objects
- No visual clustering of related components
- Can't see the "Main Sequence" distance (A + I = 1)
- No quick identification of refactoring candidates

**Solution**: Multi-dimensional heatmaps with intelligent sorting and color coding.

---

## Heatmap Type 1: Coupling Matrix

### Concept

2D matrix showing afferent (rows) vs efferent (columns) with color-coded cells.

### ASCII Visualization

```
                        EFFERENT COUPLING →
                   0    5   10   15   20   25   30   40   50+
              ┌────┬────┬────┬────┬────┬────┬────┬────┬────┐
          0   │████│    │    │    │    │    │    │    │    │ ← imports
              ├────┼────┼────┼────┼────┼────┼────┼────┼────┤
          5   │░░░░│░░░░│    │    │    │    │    │    │    │
              ├────┼────┼────┼────┼────┼────┼────┼────┼────┤
    A    10   │    │    │    │    │    │    │    │    │    │
    F         ├────┼────┼────┼────┼────┼────┼────┼────┼────┤
    F    15   │    │    │    │    │    │    │    │    │    │
    E         ├────┼────┼────┼────┼────┼────┼────┼────┼────┤
    R    20   │    │    │    │    │    │    │    │    │    │
    E         ├────┼────┼────┼────┼────┼────┼────┼────┼────┤
    N    30   │████│    │    │    │    │    │    │    │    │ ← lipgloss
    T         ├────┼────┼────┼────┼────┼────┼────┼────┼────┤
         40   │████│    │    │    │    │    │    │    │    │ ← os
    ↓         ├────┼────┼────┼────┼────┼────┼────┼────┼────┤
         50+  │████│    │    │    │    │    │    │    │    │ ← fmt, bubbletea
              └────┴────┴────┴────┴────┴────┴────┴────┴────┘

Legend: ████ = High concentration (many objects)
        ░░░░ = Moderate concentration
        Empty = Few or no objects

QUADRANT INTERPRETATION:
┌─────────────────────────┬─────────────────────────┐
│ BOTTOM-RIGHT            │ TOP-RIGHT               │
│ (High Ca, High Ce)      │ (Low Ca, High Ce)       │
│                         │                         │
│ HUB COMPONENTS          │ UNSTABLE COMPONENTS     │
│ Risky: many deps both   │ Many outgoing deps      │
│ ways                    │ I = 1.0 (examples)      │
│                         │ ← tea.go, renderer.go   │
│                         │ ← screen.go             │
├─────────────────────────┼─────────────────────────┤
│ BOTTOM-LEFT             │ TOP-LEFT                │
│ (High Ca, Low Ce)       │ (Low Ca, Low Ce)        │
│                         │                         │
│ STABLE FOUNDATIONS      │ ISOLATED COMPONENTS     │
│ Many incoming, few out  │ Few dependencies either │
│ I = 0.0 (imports)       │ way. Leaf nodes.        │
│ ← fmt, os, time         │                         │
└─────────────────────────┴─────────────────────────┘
```

---

## Heatmap Type 2: Instability vs Abstractness (Main Sequence)

### Concept

The "Main Sequence" from Robert C. Martin's metrics:
- Ideal: `A + I = 1` (diagonal line)
- Zone of Pain: Low A, Low I (concrete but stable = hard to change)
- Zone of Uselessness: High A, High I (abstract but unstable = unused)

### ASCII Visualization

```
INSTABILITY-ABSTRACTNESS PLOT (Main Sequence)

  A
  B  1.0 ┬───────────────────────────────────────────────────┐
  S      │ ZONE OF USELESSNESS         ╲                     │
  T      │ (Abstract but unstable)       ╲                   │
  R  0.8 │                                 ╲                 │
  A      │                                   ╲ Main Sequence │
  C  0.6 │                                     ╲             │
  T      │                                       ╲           │
  N  0.4 │          △ *progressWriter.Write        ╲         │
  E      │          △ *mainModel.Next                ╲       │
  S  0.2 │    ● tea.go                                 ╲     │
  S      │    ● standard_renderer.go                     ╲   │
     0.0 ├─────○────○────○────○────────────────────────────○─┤
         │    │    │    │    │                             │ │
         0   0.2  0.4  0.6  0.8                          1.0 │
         │                                                   │
         │   ZONE OF PAIN                                    │
         │   (Concrete but stable = hard to change)          │
         └───────────────────────────────────────────────────┘
                           INSTABILITY →

Legend:
  ○ = File (I = 1.0, A = 0) - Most files cluster here
  ● = Core framework file (I = 1.0, A = 0) - Need refactoring
  △ = Interface/method (I varies, A > 0)
  ╲ = Main Sequence (ideal: A + I = 1)

BUBBLETEA OBSERVATION:
Most objects cluster at I=1.0, A=0 (bottom-right)
→ High instability, low abstraction
→ Examples are fine (terminal nodes)
→ Core files should be closer to Main Sequence
```

### Distance from Main Sequence

Calculate: `D = |A + I - 1|`

| Component | I | A | D | Assessment |
|-----------|---|---|---|------------|
| import:fmt | 0.0 | 0.0 | 1.0 | Too concrete for stability |
| tea.go | 1.0 | 0.0 | 0.0 | On sequence (unstable, concrete) |
| *progressWriter.Write | 0.5 | 0.5 | 0.0 | On sequence (balanced) |
| examples/* | 1.0 | 0.0 | 0.0 | Terminal nodes (acceptable) |

---

## Heatmap Type 3: Component-Level Coupling Matrix

### Concept

Show actual coupling between specific components, not just metrics.

### ASCII Visualization (Top 15 Components)

```
COUPLING MATRIX: Top 15 Components by Total Coupling

              │fmt │os  │bt  │lg  │tea │rndr│scrn│key │lip │opt │time│sync│strg│err │ctx │
──────────────┼────┼────┼────┼────┼────┼────┼────┼────┼────┼────┼────┼────┼────┼────┼────┤
fmt     (A:58)│ ·  │    │    │    │████│████│████│████│    │████│    │    │████│    │    │
os      (A:42)│    │ ·  │    │    │████│████│████│    │    │████│    │    │    │    │    │
bubbletea(52)│    │    │ ·  │    │    │    │    │    │    │    │    │    │    │    │    │
lipgloss(A30)│    │    │    │ ·  │    │████│    │    │    │    │    │    │    │    │    │
tea.go (E:53)│████│████│████│    │ ·  │████│    │    │    │████│████│████│    │████│████│
renderer(E55)│████│████│    │████│████│ ·  │████│    │████│    │████│████│    │    │    │
screen  (E32)│████│████│    │    │    │    │ ·  │    │████│    │████│████│    │    │    │
key     (E18)│████│    │    │    │    │    │    │ ·  │    │    │    │    │████│    │    │
options (E21)│████│    │████│    │████│    │    │    │ ·  │    │████│    │    │    │████│
time   (A:25)│    │    │    │    │████│████│████│    │████│ ·  │    │    │    │    │    │
sync    (A:4)│    │    │    │    │████│████│    │    │    │    │ ·  │    │    │    │    │
strings(A:21)│    │    │    │    │    │    │    │████│    │    │    │ ·  │    │    │    │
errors  (A:4)│    │    │    │    │████│    │    │    │    │    │    │    │ ·  │    │    │
context (A:5)│    │    │    │    │████│    │    │    │████│    │    │    │    │ ·  │    │
──────────────┴────┴────┴────┴────┴────┴────┴────┴────┴────┴────┴────┴────┴────┴────┴────┘

Legend: ████ = Direct dependency
        ·    = Self (diagonal)
        empty= No direct dependency

PATTERN ANALYSIS:
- Row density = Efferent coupling (many ████ = high Ce)
- Column density = Afferent coupling (many ████ = high Ca)
- Dense clusters = Tightly coupled subsystems
```

---

## Heatmap Type 4: File Coupling by Category

### Concept

Group files into categories and show inter-category coupling density.

```
INTER-CATEGORY COUPLING DENSITY

               │ Imports │ Core │ Support │ Examples │
───────────────┼─────────┼──────┼─────────┼──────────│
Imports        │   0%    │ 100% │   80%   │   60%    │ (stable foundation)
───────────────┼─────────┼──────┼─────────┼──────────│
Core (2 files) │  HIGH   │  10% │   20%   │    0%    │ (orchestration)
───────────────┼─────────┼──────┼─────────┼──────────│
Support (10)   │   MED   │  30% │   15%   │    0%    │ (utilities)
───────────────┼─────────┼──────┼─────────┼──────────│
Examples (50)  │   MED   │  95% │   20%   │    0%    │ (applications)
───────────────┴─────────┴──────┴─────────┴──────────┘

INTERPRETATION:
- Examples → Core: 95% coupling (all examples use tea.go)
- Core → Imports: 100% (all core files use standard imports)
- Inter-example coupling: 0% (complete isolation)
- Support → Core: 30% (some support files depend on core)
```

---

## Implementation: Color Scales

### For Terminal (ANSI Colors)

```python
def coupling_color(value, max_value):
    """Return ANSI color code based on coupling intensity"""
    ratio = value / max_value
    if ratio < 0.2:
        return '\033[92m'  # Green (low)
    elif ratio < 0.4:
        return '\033[93m'  # Yellow
    elif ratio < 0.6:
        return '\033[33m'  # Orange
    elif ratio < 0.8:
        return '\033[91m'  # Red
    else:
        return '\033[31m'  # Dark Red (high)

def instability_color(instability):
    """Color code for instability: green (stable) to red (unstable)"""
    if instability < 0.3:
        return '\033[92m'  # Green
    elif instability < 0.7:
        return '\033[93m'  # Yellow
    else:
        return '\033[91m'  # Red
```

### For Web (CSS/D3.js)

```javascript
const couplingScale = d3.scaleSequential(d3.interpolateYlOrRd)
    .domain([0, maxCoupling]);

const instabilityScale = d3.scaleSequential(d3.interpolateRdYlGn)
    .domain([1, 0]);  // Reversed: 0 is green (stable)

const mainSequenceDistance = d3.scaleSequential(d3.interpolateViridis)
    .domain([0, 1]);
```

---

## ASCII Heatmap Generator

```python
def generate_coupling_heatmap(coupling_metrics, width=60):
    """Generate ASCII coupling heatmap sorted by total coupling"""

    # Sort by total coupling
    sorted_items = sorted(
        coupling_metrics.items(),
        key=lambda x: x[1]['afferent_coupling'] + x[1]['efferent_coupling'],
        reverse=True
    )[:20]  # Top 20

    # Header
    output = []
    output.append('╔' + '═' * width + '╗')
    output.append('║' + 'COUPLING HEATMAP'.center(width) + '║')
    output.append('╠' + '═' * width + '╣')
    output.append('║' + 'Component'.ljust(25) + 'Ca'.center(6) + 'Ce'.center(6) + 'I'.center(6) + 'D'.center(6) + ' ▓▓▓▓▓▓▓ ║')
    output.append('╟' + '─' * width + '╢')

    max_coupling = max(
        m['afferent_coupling'] + m['efferent_coupling']
        for _, m in sorted_items
    )

    for obj_id, m in sorted_items:
        name = obj_id[:24].ljust(25)
        ca = str(m['afferent_coupling']).center(6)
        ce = str(m['efferent_coupling']).center(6)
        i = f"{m['instability']:.2f}".center(6)
        d = f"{abs(m['instability'] + m['abstractness'] - 1):.2f}".center(6)

        # Bar representation
        total = m['afferent_coupling'] + m['efferent_coupling']
        bar_len = int((total / max_coupling) * 10)
        bar = '█' * bar_len + '░' * (10 - bar_len)

        output.append(f'║{name}{ca}{ce}{i}{d} {bar} ║')

    output.append('╚' + '═' * width + '╝')

    return '\n'.join(output)
```

### Example Output

```
╔════════════════════════════════════════════════════════════╗
║                     COUPLING HEATMAP                       ║
╠════════════════════════════════════════════════════════════╣
║Component                 Ca    Ce     I      D    ▓▓▓▓▓▓▓  ║
╟────────────────────────────────────────────────────────────╢
║import:fmt               58     0   0.00   1.00  ██████░░░░ ║
║standard_renderer.go      0    55   1.00   0.00  █████░░░░░ ║
║tea.go                    0    53   1.00   0.00  █████░░░░░ ║
║import:bubbletea         52     0   0.00   1.00  █████░░░░░ ║
║import:os                42     0   0.00   1.00  ████░░░░░░ ║
║screen.go                 0    32   1.00   0.00  ███░░░░░░░ ║
║import:lipgloss          30     0   0.00   1.00  ███░░░░░░░ ║
║import:time              25     0   0.00   1.00  ██░░░░░░░░ ║
║nil_renderer.go           0    26   1.00   0.00  ██░░░░░░░░ ║
║import:strings           21     0   0.00   1.00  ██░░░░░░░░ ║
║options.go                0    21   1.00   0.00  ██░░░░░░░░ ║
║examples/views/main.go    0    19   1.00   0.00  █░░░░░░░░░ ║
║examples/autocomplete     0    18   1.00   0.00  █░░░░░░░░░ ║
║key.go                    0    18   1.00   0.00  █░░░░░░░░░ ║
║examples/cellbuffer       0    17   1.00   0.00  █░░░░░░░░░ ║
╚════════════════════════════════════════════════════════════╝

Legend: Ca = Afferent, Ce = Efferent, I = Instability, D = Distance from Main Sequence
        Bar: ████████ = Total Coupling (Ca + Ce)
```

---

## Key Insights from Heatmap Analysis

### Bubbletea Patterns Revealed

1. **Bimodal Distribution**: Two clusters
   - High afferent, zero efferent (imports)
   - Zero afferent, high efferent (framework + examples)

2. **No Middle Ground**: Very few components with balanced coupling
   - Most are either pure sources or pure sinks

3. **Main Sequence Violations**: Most components at D = 1.0
   - Suggests need for more interfaces/abstractions
   - Examples are acceptable (terminal nodes)

4. **Coupling Concentration**: Top 5 components = 60% of total morphisms
   - `fmt`, `bubbletea`, `tea.go`, `renderer`, `os`

---

## Next Steps

- `04-CATEGORICAL-FLOW.md`: Animate Elm architecture data flow
- `05-TOOLING-RECOMMENDATIONS.md`: Implementation tools and libraries
