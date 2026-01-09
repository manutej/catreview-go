// visualize_project_v3.go - Force TRUE vertical stacking using invisible edges
// Solves the horizontal sprawl problem for disconnected component graphs
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/manu/catreview/pkg/category"
)

func main() {
	inputFile := flag.String("input", "", "Input JSON file (e.g., barque-analysis.json)")
	outputPrefix := flag.String("output", "output", "Output prefix for .dot files")
	maxNodes := flag.Int("max-nodes", 100, "Maximum nodes per visualization (default: 100)")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Usage: visualize_project_v3 --input <analysis.json> [--output <prefix>] [--max-nodes <N>]")
		fmt.Println("Example: visualize_project_v3 --input barque-analysis.json --output barque --max-nodes 100")
		os.Exit(1)
	}

	fmt.Printf("═══════════════════════════════════════════════════════════════\n")
	fmt.Printf("  Categorical Visualization Generator V3 (FORCED Vertical)\n")
	fmt.Printf("═══════════════════════════════════════════════════════════════\n\n")

	fmt.Printf("📂 Input: %s\n", *inputFile)
	fmt.Printf("📊 Output Prefix: %s\n", *outputPrefix)
	fmt.Printf("📐 Max Nodes: %d\n\n", *maxNodes)

	// Load categorical model
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var cat category.Category
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cat); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error decoding JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Loaded categorical model\n")
	fmt.Printf("   Objects: %d\n", len(cat.Objects()))
	fmt.Printf("   Morphisms: %d\n\n", len(cat.Morphisms()))

	// Generate visualizations
	fmt.Println("Creating visualizations with forced vertical layout...")

	if err := generateCompositionGraphV3(&cat, *outputPrefix, *maxNodes); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error generating composition graph: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ %s-composition-v3.dot created\n\n", *outputPrefix)

	fmt.Println("To generate SVG, run:")
	fmt.Printf("  dot -Tsvg %s-composition-v3.dot -o %s-composition-v3.svg\n\n", *outputPrefix, *outputPrefix)
	fmt.Println("To generate PDF, run:")
	fmt.Printf("  dot -Tpdf %s-composition-v3.dot -o %s-composition-v3.pdf\n\n", *outputPrefix, *outputPrefix)
}

func generateCompositionGraphV3(cat *category.Category, prefix string, maxNodes int) error {
	filename := fmt.Sprintf("%s-composition-v3.dot", prefix)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Sample modules with most defines relationships
	type moduleStat struct {
		id      string
		objects []*category.Object
	}
	moduleStats := make(map[string]*moduleStat)

	for _, obj := range cat.Objects() {
		if obj.Type == "module" {
			moduleStats[obj.ID] = &moduleStat{id: obj.ID, objects: []*category.Object{}}
		}
	}

	for _, morph := range cat.Morphisms() {
		if morph.Type == "defines" {
			if stat, exists := moduleStats[morph.Source]; exists {
				obj, objExists := cat.GetObject(morph.Target)
				if objExists {
					stat.objects = append(stat.objects, obj)
				}
			}
		}
	}

	// Sort modules by object count
	var sorted []*moduleStat
	for _, stat := range moduleStats {
		if len(stat.objects) > 0 {
			sorted = append(sorted, stat)
		}
	}
	sort.Slice(sorted, func(i, j int) bool {
		return len(sorted[i].objects) > len(sorted[j].objects)
	})

	// Select top modules - limit to 20 for manageability
	limit := min(20, len(sorted))

	// SINGLE COLUMN LAYOUT for true vertical stacking
	modulesPerRow := 1
	numRows := limit

	// Generate DOT with FORCED vertical layering using invisible edges
	fmt.Fprintln(f, "digraph CompositionStructure {")
	fmt.Fprintln(f, "  // Graph layout settings - FORCED VERTICAL")
	fmt.Fprintln(f, "  rankdir=TB;  // Top-to-Bottom")
	fmt.Fprintln(f, "  ranksep=1.8;  // Large vertical spacing between rows")
	fmt.Fprintln(f, "  nodesep=1.0;  // Horizontal spacing within rows")
	fmt.Fprintln(f, "  compound=true;")
	fmt.Fprintln(f, "  newrank=true;  // Enable new ranking algorithm")
	fmt.Fprintln(f, "")
	fmt.Fprintf(f, "  label=\"Composition Structure (Top %d modules)\\nSingle Column Vertical Layout: Modules Stacked Top-to-Bottom\";\n", limit)
	fmt.Fprintln(f, "  fontsize=18;")
	fmt.Fprintln(f, "  fontname=\"Helvetica-Bold\";")
	fmt.Fprintln(f, "  labelloc=t;")
	fmt.Fprintln(f, "")

	// Legend box
	fmt.Fprintln(f, "  // Legend")
	fmt.Fprintln(f, "  subgraph cluster_legend {")
	fmt.Fprintln(f, "    label=\"Legend\";")
	fmt.Fprintln(f, "    style=filled;")
	fmt.Fprintln(f, "    fillcolor=lightyellow;")
	fmt.Fprintln(f, "    fontsize=12;")
	fmt.Fprintln(f, "    rank=source;")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "    legend_module [label=\"Module\", shape=box, style=\"rounded,filled\", fillcolor=lightgray];")
	fmt.Fprintln(f, "    legend_func [label=\"Function\", shape=ellipse, style=filled, fillcolor=lightgreen];")
	fmt.Fprintln(f, "    legend_class [label=\"Class\", shape=box, style=filled, fillcolor=lightblue];")
	fmt.Fprintln(f, "  }")
	fmt.Fprintln(f, "")

	selectedObjects := make(map[string]bool)
	invisibleEdgeTargets := []string{} // Track first node in each cluster for invisible edges

	// Generate modules in rows
	for i := 0; i < limit; i++ {
		stat := sorted[i]
		row := i / modulesPerRow
		col := i % modulesPerRow

		fmt.Fprintf(f, "  // Module %d (Row %d, Column %d)\n", i+1, row+1, col+1)
		fmt.Fprintf(f, "  subgraph cluster_%d {\n", i)
		fmt.Fprintf(f, "    label=\"📦 %s\";\n", getShortName(stat.id))
		fmt.Fprintln(f, "    style=\"rounded,filled\";")
		fmt.Fprintln(f, "    fillcolor=lightgray;")
		fmt.Fprintln(f, "    fontsize=13;")
		fmt.Fprintln(f, "")

		// Limit objects per module
		objLimit := min(10, len(stat.objects))
		var firstNodeID string
		for j := 0; j < objLimit; j++ {
			obj := stat.objects[j]
			selectedObjects[obj.ID] = true
			if j == 0 {
				firstNodeID = obj.ID
			}
			color := "lightgreen"
			shape := "ellipse"
			if obj.Type == "class" {
				color = "lightblue"
				shape = "box"
			}
			label := getShortName(obj.ID)
			fmt.Fprintf(f, "    \"%s\" [label=\"%s\", fillcolor=%s, shape=%s];\n", obj.ID, label, color, shape)
		}

		fmt.Fprintln(f, "  }")
		fmt.Fprintln(f, "")

		if firstNodeID != "" {
			invisibleEdgeTargets = append(invisibleEdgeTargets, firstNodeID)
		}
	}

	// Create vertical structure with invisible edges
	// This FORCES row-by-row vertical stacking
	fmt.Fprintln(f, "  // Invisible structural edges to force vertical layout")
	fmt.Fprintln(f, "  // These edges create dependency chains that enforce row ordering")
	for row := 0; row < numRows-1; row++ {
		startIdx := row * modulesPerRow
		nextRowStartIdx := (row + 1) * modulesPerRow

		if nextRowStartIdx < len(invisibleEdgeTargets) {
			// Connect first module of this row to first module of next row
			fmt.Fprintf(f, "  \"%s\" -> \"%s\" [style=invis, weight=10];\n",
				invisibleEdgeTargets[startIdx],
				invisibleEdgeTargets[nextRowStartIdx])
		}
	}
	fmt.Fprintln(f, "")

	// NO rank=same constraints - let invisible edges enforce vertical stacking
	fmt.Fprintln(f, "  // Vertical stacking enforced by invisible edge chain above")
	fmt.Fprintln(f, "")

	// Add call edges between selected objects
	fmt.Fprintln(f, "  // Inter-module function calls")
	for _, morph := range cat.Morphisms() {
		if morph.Type == "function_call" {
			if selectedObjects[morph.Source] && selectedObjects[morph.Target] {
				fmt.Fprintf(f, "  \"%s\" -> \"%s\" [color=green, style=dashed, penwidth=1.5];\n",
					morph.Source, morph.Target)
			}
		}
	}

	fmt.Fprintln(f, "}")
	return nil
}

func getShortName(fullID string) string {
	// Remove common prefixes for cleaner labels
	fullID = strings.TrimPrefix(fullID, "import:")
	fullID = strings.TrimPrefix(fullID, "venv.lib.python")
	fullID = strings.TrimPrefix(fullID, "Users.manu.Documents.LUXOR.")

	parts := strings.Split(fullID, ".")
	if len(parts) > 3 {
		// Take last 3 parts for deeply nested modules
		return strings.Join(parts[len(parts)-3:], ".")
	}
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return fullID
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
