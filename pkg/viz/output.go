// output.go provides multiple output formats for graph visualization.
package viz

import (
	"fmt"
	"sort"
	"strings"
)

// Format represents an output format for visualization.
type Format string

const (
	FormatASCII   Format = "ascii"
	FormatMermaid Format = "mermaid"
	FormatDOT     Format = "dot"
	FormatJSON    Format = "json"
)

// Generator generates visualization output in a specific format.
type Generator interface {
	Generate(g *Graph) (string, error)
}

// NewGenerator creates a generator for the specified format.
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
	default:
		return &ASCIIGenerator{}
	}
}

// ASCIIGenerator produces ASCII art visualizations.
type ASCIIGenerator struct{}

// Generate creates an ASCII visualization of the graph.
func (g *ASCIIGenerator) Generate(graph *Graph) (string, error) {
	var sb strings.Builder
	width := 78

	// Title
	sb.WriteString(boxTop(width))
	sb.WriteString(boxLine(centerString(graph.Name+" DEPENDENCY GRAPH", width-4), width))
	sb.WriteString(boxMid(width))

	// Stats
	dagStatus := "YES"
	if !graph.IsDAG {
		dagStatus = "NO (cycles detected)"
	}
	sb.WriteString(boxLine(fmt.Sprintf("Nodes: %d  Edges: %d  DAG: %s",
		graph.Stats.TotalNodes, graph.Stats.TotalEdges, dagStatus), width))
	sb.WriteString(boxMid(width))

	// Layers (top to bottom: applications -> foundations)
	for i := 3; i >= 0; i-- {
		if i < len(graph.Layers) {
			layer := graph.Layers[i]
			g.generateLayer(&sb, graph, layer, width, i)

			if i > 0 {
				// Draw connector between layers
				sb.WriteString(boxLine(centerString("|", width-4), width))
				sb.WriteString(boxLine(centerString("v", width-4), width))
			}
		}
	}

	sb.WriteString(boxBottom(width))

	// Legend
	sb.WriteString("\nLegend: A=Afferent, E=Efferent, I=Instability\n")
	sb.WriteString("Layers: 0=Foundations, 1=Core, 2=Support, 3=Applications\n")

	return sb.String(), nil
}

// generateLayer renders a single layer to the string builder.
func (g *ASCIIGenerator) generateLayer(sb *strings.Builder, graph *Graph, layer *Layer, width int, layerIdx int) {
	sb.WriteString(boxLine("", width))
	sb.WriteString(boxLine(fmt.Sprintf("LAYER %d: %s", layer.ID, strings.ToUpper(layer.Name)), width))

	// Get top nodes by coupling
	maxNodes := 5
	nodeIDs := layer.NodeIDs
	if len(nodeIDs) > maxNodes {
		nodeIDs = nodeIDs[:maxNodes]
	}

	if len(nodeIDs) == 0 {
		sb.WriteString(boxLine("  (no nodes)", width))
	} else {
		for _, nodeID := range nodeIDs {
			if node, ok := graph.Nodes[nodeID]; ok {
				nodeStr := fmt.Sprintf("  [%s] A:%d E:%d I:%.2f",
					truncate(node.Label, 20), node.Afferent, node.Efferent, node.Instability)
				sb.WriteString(boxLine(nodeStr, width))
			}
		}
		if len(layer.NodeIDs) > maxNodes {
			sb.WriteString(boxLine(fmt.Sprintf("  ... and %d more", len(layer.NodeIDs)-maxNodes), width))
		}
	}
}

// MermaidGenerator produces Mermaid.js syntax.
type MermaidGenerator struct{}

// Generate creates a Mermaid diagram definition.
func (g *MermaidGenerator) Generate(graph *Graph) (string, error) {
	var sb strings.Builder

	sb.WriteString("```mermaid\n")
	sb.WriteString("flowchart BT\n")

	// Define subgraphs for each layer
	layerNames := []string{"Foundations", "Core", "Support", "Applications"}
	layerStyles := []string{"fill:#e8f5e9", "fill:#fff3e0", "fill:#e3f2fd", "fill:#fce4ec"}

	for i := 0; i < len(graph.Layers); i++ {
		layer := graph.Layers[i]
		if len(layer.NodeIDs) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("    subgraph L%d[\"%s\"]\n", i, layerNames[i]))

		// Add nodes (limit to top 10 per layer)
		maxNodes := 10
		nodeIDs := layer.NodeIDs
		if len(nodeIDs) > maxNodes {
			nodeIDs = nodeIDs[:maxNodes]
		}

		for _, nodeID := range nodeIDs {
			if node, ok := graph.Nodes[nodeID]; ok {
				safeID := sanitizeMermaidID(nodeID)
				label := fmt.Sprintf("%s<br/>A:%d E:%d",
					escapeHTML(truncate(node.Label, 15)), node.Afferent, node.Efferent)
				sb.WriteString(fmt.Sprintf("        %s[\"%s\"]\n", safeID, label))
			}
		}
		sb.WriteString("    end\n")
	}

	sb.WriteString("\n")

	// Add edges (limit to prevent clutter)
	edgeCount := 0
	maxEdges := 50
	for _, edge := range graph.Edges {
		if edgeCount >= maxEdges {
			break
		}

		// Only include edges between top-coupled nodes
		sourceNode, sourceOK := graph.Nodes[edge.Source]
		targetNode, targetOK := graph.Nodes[edge.Target]
		if !sourceOK || !targetOK {
			continue
		}

		// Filter to significant edges
		if (sourceNode.Afferent+sourceNode.Efferent) < 10 &&
		   (targetNode.Afferent+targetNode.Efferent) < 10 {
			continue
		}

		sourceID := sanitizeMermaidID(edge.Source)
		targetID := sanitizeMermaidID(edge.Target)
		sb.WriteString(fmt.Sprintf("    %s --> %s\n", sourceID, targetID))
		edgeCount++
	}

	sb.WriteString("\n")

	// Add styles
	for i, style := range layerStyles {
		sb.WriteString(fmt.Sprintf("    style L%d %s\n", i, style))
	}

	sb.WriteString("```\n")

	return sb.String(), nil
}

// DOTGenerator produces Graphviz DOT syntax.
type DOTGenerator struct{}

// Generate creates a DOT graph definition.
func (g *DOTGenerator) Generate(graph *Graph) (string, error) {
	var sb strings.Builder

	sb.WriteString("digraph G {\n")
	sb.WriteString("    rankdir=BT;\n")
	sb.WriteString("    node [shape=box, style=filled];\n")
	sb.WriteString("    edge [arrowsize=0.7];\n\n")

	// Color schemes for each layer
	layerColors := []string{"#c8e6c9", "#ffe0b2", "#bbdefb", "#f8bbd0"}

	// Group nodes by layer
	for i := 0; i < len(graph.Layers); i++ {
		layer := graph.Layers[i]
		if len(layer.NodeIDs) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("    // Layer %d: %s\n", i, layer.Name))
		sb.WriteString(fmt.Sprintf("    subgraph cluster_%d {\n", i))
		sb.WriteString(fmt.Sprintf("        label=\"%s\";\n", layer.Name))
		sb.WriteString(fmt.Sprintf("        style=filled;\n"))
		sb.WriteString(fmt.Sprintf("        fillcolor=\"%s\";\n", layerColors[i]))

		// Limit nodes per layer
		maxNodes := 15
		nodeIDs := layer.NodeIDs
		if len(nodeIDs) > maxNodes {
			nodeIDs = nodeIDs[:maxNodes]
		}

		for _, nodeID := range nodeIDs {
			if node, ok := graph.Nodes[nodeID]; ok {
				safeID := sanitizeDOTID(nodeID)
				label := fmt.Sprintf("%s\\nA:%d E:%d",
					escapeDOT(truncate(node.Label, 20)), node.Afferent, node.Efferent)

				// Intensity based on coupling
				intensity := min(255, 128 + (node.Afferent+node.Efferent)*3)
				nodeColor := fmt.Sprintf("#%02x%02x%02x", intensity, intensity, intensity)

				sb.WriteString(fmt.Sprintf("        %s [label=\"%s\", fillcolor=\"%s\"];\n",
					safeID, label, nodeColor))
			}
		}
		sb.WriteString("    }\n\n")
	}

	// Add edges
	sb.WriteString("    // Edges\n")
	edgeCount := 0
	maxEdges := 100
	for _, edge := range graph.Edges {
		if edgeCount >= maxEdges {
			break
		}

		sourceID := sanitizeDOTID(edge.Source)
		targetID := sanitizeDOTID(edge.Target)

		penWidth := fmt.Sprintf("%.1f", 0.5+float64(edge.Weight)*0.2)
		color := "#666666"
		if edge.Type == "cross_layer" {
			color = "#cc0000"
		}

		sb.WriteString(fmt.Sprintf("    %s -> %s [penwidth=%s, color=\"%s\"];\n",
			sourceID, targetID, penWidth, color))
		edgeCount++
	}

	sb.WriteString("}\n")

	return sb.String(), nil
}

// JSONGenerator produces JSON output.
type JSONGenerator struct{}

// Generate creates JSON representation of the graph.
func (g *JSONGenerator) Generate(graph *Graph) (string, error) {
	data, err := graph.ToJSON()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Helper functions

func boxTop(width int) string {
	return fmt.Sprintf("╔%s╗\n", strings.Repeat("═", width-2))
}

func boxBottom(width int) string {
	return fmt.Sprintf("╚%s╝\n", strings.Repeat("═", width-2))
}

func boxMid(width int) string {
	return fmt.Sprintf("╠%s╣\n", strings.Repeat("═", width-2))
}

func boxLine(content string, width int) string {
	// Pad content to fit box width
	padded := content
	if len(padded) > width-4 {
		padded = padded[:width-7] + "..."
	}
	padding := width - 4 - len(padded)
	return fmt.Sprintf("║ %s%s ║\n", padded, strings.Repeat(" ", max(0, padding)))
}

func centerString(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	leftPad := (width - len(s)) / 2
	rightPad := width - len(s) - leftPad
	return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func sanitizeMermaidID(id string) string {
	// Mermaid IDs can't contain certain characters
	result := strings.ReplaceAll(id, "/", "_")
	result = strings.ReplaceAll(result, ".", "_")
	result = strings.ReplaceAll(result, ":", "_")
	result = strings.ReplaceAll(result, "-", "_")
	result = strings.ReplaceAll(result, "*", "_")
	return result
}

func sanitizeDOTID(id string) string {
	// DOT IDs need to be quoted if they contain special chars
	result := strings.ReplaceAll(id, "/", "_")
	result = strings.ReplaceAll(result, ".", "_")
	result = strings.ReplaceAll(result, ":", "_")
	result = strings.ReplaceAll(result, "-", "_")
	result = strings.ReplaceAll(result, "*", "_")
	return "n_" + result
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func escapeDOT(s string) string {
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// GenerateLayeredASCII creates a detailed ASCII visualization with all layers.
func GenerateLayeredASCII(graph *Graph) string {
	var sb strings.Builder
	width := 80

	sb.WriteString(strings.Repeat("═", width) + "\n")
	sb.WriteString(centerString(graph.Name+" - LAYERED DEPENDENCY VISUALIZATION", width) + "\n")
	sb.WriteString(strings.Repeat("═", width) + "\n\n")

	dagStatus := "✓ DAG Verified"
	if !graph.IsDAG {
		dagStatus = "✗ Cycles Detected"
	}

	sb.WriteString(fmt.Sprintf("Stats: %d nodes, %d edges, %s\n\n",
		graph.Stats.TotalNodes, graph.Stats.TotalEdges, dagStatus))

	// Render each layer from top to bottom
	for i := 3; i >= 0; i-- {
		if i >= len(graph.Layers) {
			continue
		}
		layer := graph.Layers[i]

		sb.WriteString(strings.Repeat("─", width) + "\n")
		sb.WriteString(fmt.Sprintf("LAYER %d: %s (%s)\n",
			layer.ID, strings.ToUpper(layer.Name), layer.Description))
		sb.WriteString(strings.Repeat("─", width) + "\n")

		if len(layer.NodeIDs) == 0 {
			sb.WriteString("  (empty)\n")
			continue
		}

		// Get nodes sorted by total coupling
		nodes := make([]*Node, 0, len(layer.NodeIDs))
		for _, id := range layer.NodeIDs {
			if node, ok := graph.Nodes[id]; ok {
				nodes = append(nodes, node)
			}
		}
		sort.Slice(nodes, func(a, b int) bool {
			return (nodes[a].Afferent + nodes[a].Efferent) > (nodes[b].Afferent + nodes[b].Efferent)
		})

		// Display top nodes
		maxDisplay := 8
		for j, node := range nodes {
			if j >= maxDisplay {
				sb.WriteString(fmt.Sprintf("  ... and %d more components\n", len(nodes)-maxDisplay))
				break
			}

			// Format: [label] A:## E:## I:0.## (type)
			line := fmt.Sprintf("  ┌ %-25s A:%-3d E:%-3d I:%.2f",
				truncate(node.Label, 25), node.Afferent, node.Efferent, node.Instability)
			sb.WriteString(line + "\n")
		}

		sb.WriteString("\n")

		// Show dependencies to next layer down
		if i > 0 {
			sb.WriteString(centerString("│", width) + "\n")
			sb.WriteString(centerString("▼ depends on", width) + "\n")
		}
	}

	sb.WriteString(strings.Repeat("═", width) + "\n")

	return sb.String()
}

// GenerateHeatmap creates an ASCII coupling heatmap.
func GenerateHeatmap(graph *Graph) string {
	var sb strings.Builder

	sb.WriteString("COUPLING HEATMAP (sorted by total coupling)\n")
	sb.WriteString("═══════════════════════════════════════════════════════════\n")
	sb.WriteString(fmt.Sprintf("%-30s %6s %6s %6s %s\n",
		"Component", "Ca", "Ce", "I", "Coupling"))
	sb.WriteString("───────────────────────────────────────────────────────────\n")

	// Collect and sort nodes
	nodes := make([]*Node, 0, len(graph.Nodes))
	for _, node := range graph.Nodes {
		nodes = append(nodes, node)
	}
	sort.Slice(nodes, func(i, j int) bool {
		ti := nodes[i].Afferent + nodes[i].Efferent
		tj := nodes[j].Afferent + nodes[j].Efferent
		return ti > tj
	})

	// Find max coupling for bar scaling
	maxCoupling := 1
	for _, node := range nodes {
		total := node.Afferent + node.Efferent
		if total > maxCoupling {
			maxCoupling = total
		}
	}

	// Display top 20
	for i, node := range nodes {
		if i >= 20 {
			break
		}

		total := node.Afferent + node.Efferent
		barLen := (total * 15) / maxCoupling
		bar := strings.Repeat("█", barLen) + strings.Repeat("░", 15-barLen)

		sb.WriteString(fmt.Sprintf("%-30s %6d %6d %6.2f %s\n",
			truncate(node.Label, 30),
			node.Afferent,
			node.Efferent,
			node.Instability,
			bar))
	}

	sb.WriteString("───────────────────────────────────────────────────────────\n")
	sb.WriteString("Ca=Afferent (incoming), Ce=Efferent (outgoing), I=Instability\n")

	return sb.String()
}
