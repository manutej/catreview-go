// visualize_project_v2.go - Create hierarchical vertical layered visualizations
// Addresses horizontal sprawl with proper vertical stacking and abstraction levels
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
		fmt.Println("Usage: visualize_project_v2 --input <analysis.json> [--output <prefix>] [--max-nodes <N>]")
		fmt.Println("Example: visualize_project_v2 --input barque-analysis.json --output barque --max-nodes 100")
		os.Exit(1)
	}

	fmt.Printf("═══════════════════════════════════════════════════════════════\n")
	fmt.Printf("  Categorical Visualization Generator V2 (Vertical Layered)\n")
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
	fmt.Println("Creating visualizations...")

	if err := generateInheritanceGraphV2(&cat, *outputPrefix, *maxNodes); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error generating inheritance graph: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ %s-inheritance-v2.dot created\n", *outputPrefix)

	if err := generateCallGraphV2(&cat, *outputPrefix, *maxNodes); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error generating call graph: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ %s-calls-v2.dot created\n", *outputPrefix)

	if err := generateModuleDependencyGraphV2(&cat, *outputPrefix, *maxNodes); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error generating module graph: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ %s-modules-v2.dot created\n", *outputPrefix)

	if err := generateCompositionGraphV2(&cat, *outputPrefix, *maxNodes); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error generating composition graph: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ %s-composition-v2.dot created\n\n", *outputPrefix)

	fmt.Println("To generate SVG files, run:")
	fmt.Printf("  for f in %s-*-v2.dot; do dot -Tsvg \"$f\" -o \"${f%%.dot}.svg\"; done\n\n", *outputPrefix)
	fmt.Println("To generate PDFs, run:")
	fmt.Printf("  for f in %s-*-v2.dot; do dot -Tpdf \"$f\" -o \"${f%%.dot}.pdf\"; done\n\n", *outputPrefix)
}

func generateInheritanceGraphV2(cat *category.Category, prefix string, maxNodes int) error {
	filename := fmt.Sprintf("%s-inheritance-v2.dot", prefix)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Sample classes with most inheritance relationships
	type classStat struct {
		obj    *category.Object
		degree int
	}
	classStats := make(map[string]*classStat)

	for _, obj := range cat.Objects() {
		if obj.Type == "class" {
			classStats[obj.ID] = &classStat{obj: obj, degree: 0}
		}
	}

	for _, morph := range cat.Morphisms() {
		if morph.Type == "inheritance" {
			if stat, exists := classStats[morph.Source]; exists {
				stat.degree++
			}
			if stat, exists := classStats[morph.Target]; exists {
				stat.degree++
			}
		}
	}

	// Sort by degree and take top N
	var sorted []*classStat
	for _, stat := range classStats {
		sorted = append(sorted, stat)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].degree > sorted[j].degree
	})

	selectedClasses := make(map[string]bool)
	limit := min(maxNodes, len(sorted))
	for i := 0; i < limit; i++ {
		selectedClasses[sorted[i].obj.ID] = true
	}

	// Generate DOT with vertical layering
	fmt.Fprintln(f, "digraph InheritanceHierarchy {")
	fmt.Fprintln(f, "  // Graph layout settings")
	fmt.Fprintln(f, "  rankdir=BT;  // Bottom-to-Top (child → parent)")
	fmt.Fprintln(f, "  ranksep=0.8;  // Vertical spacing between ranks")
	fmt.Fprintln(f, "  nodesep=0.5;  // Horizontal spacing between nodes")
	fmt.Fprintln(f, "  splines=ortho;  // Orthogonal edges for clarity")
	fmt.Fprintln(f, "  ")
	fmt.Fprintf(f, "  label=\"Inheritance Hierarchy (Top %d classes)\\nVertical Layout: Base Classes (Top) → Derived Classes (Bottom)\";\n", limit)
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
	fmt.Fprintln(f, "    rank=source;  // Place legend at top")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "    legend_class [label=\"Class\", shape=box, style=filled, fillcolor=lightblue];")
	fmt.Fprintln(f, "    legend_base [label=\"Base Class\", shape=box, style=\"filled,bold\", fillcolor=lightblue, penwidth=2];")
	fmt.Fprintln(f, "    legend_edge [label=\"inherits from →\", shape=plaintext];")
	fmt.Fprintln(f, "  }")
	fmt.Fprintln(f, "")

	// Node definitions with abstraction level grouping
	fmt.Fprintln(f, "  // Classes grouped by abstraction level")
	for classID := range selectedClasses {
		label := getShortName(classID)
		// Bold border for base classes (degree > 3)
		style := "filled"
		penwidth := "1"
		if stat := classStats[classID]; stat != nil && stat.degree > 3 {
			style = "\"filled,bold\""
			penwidth = "2"
		}
		fmt.Fprintf(f, "  \"%s\" [label=\"%s\", shape=box, style=%s, fillcolor=lightblue, penwidth=%s];\n",
			classID, label, style, penwidth)
	}
	fmt.Fprintln(f, "")

	// Edges with labels
	fmt.Fprintln(f, "  // Inheritance edges")
	for _, morph := range cat.Morphisms() {
		if morph.Type == "inheritance" {
			if selectedClasses[morph.Source] && selectedClasses[morph.Target] {
				fmt.Fprintf(f, "  \"%s\" -> \"%s\" [color=darkblue, penwidth=2, arrowsize=0.8];\n",
					morph.Source, morph.Target)
			}
		}
	}

	fmt.Fprintln(f, "}")
	return nil
}

func generateCallGraphV2(cat *category.Category, prefix string, maxNodes int) error {
	filename := fmt.Sprintf("%s-calls-v2.dot", prefix)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Sample functions with most calls
	type funcStat struct {
		id     string
		degree int
		isClass bool
	}
	funcStats := make(map[string]*funcStat)

	for _, morph := range cat.Morphisms() {
		if morph.Type == "function_call" {
			if _, exists := funcStats[morph.Source]; !exists {
				funcStats[morph.Source] = &funcStat{id: morph.Source, degree: 0}
			}
			if _, exists := funcStats[morph.Target]; !exists {
				funcStats[morph.Target] = &funcStat{id: morph.Target, degree: 0}
			}
			funcStats[morph.Source].degree++
			funcStats[morph.Target].degree++
		}
	}

	// Determine if function is actually a class
	for id := range funcStats {
		if obj, exists := cat.GetObject(id); exists && obj.Type == "class" {
			funcStats[id].isClass = true
		}
	}

	var sorted []*funcStat
	for _, stat := range funcStats {
		sorted = append(sorted, stat)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].degree > sorted[j].degree
	})

	selectedFuncs := make(map[string]bool)
	limit := min(maxNodes, len(sorted))
	for i := 0; i < limit; i++ {
		selectedFuncs[sorted[i].id] = true
	}

	// Generate DOT with vertical layering
	fmt.Fprintln(f, "digraph CallGraph {")
	fmt.Fprintln(f, "  // Graph layout settings")
	fmt.Fprintln(f, "  rankdir=TB;  // Top-to-Bottom (caller → callee)")
	fmt.Fprintln(f, "  ranksep=1.0;")
	fmt.Fprintln(f, "  nodesep=0.6;")
	fmt.Fprintln(f, "  splines=true;")
	fmt.Fprintln(f, "")
	fmt.Fprintf(f, "  label=\"Function Call Graph (Top %d functions)\\nVertical Layout: Callers (Top) → Callees (Bottom)\";\n", limit)
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
	fmt.Fprintln(f, "    legend_func [label=\"Function\", shape=ellipse, style=filled, fillcolor=lightgreen];")
	fmt.Fprintln(f, "    legend_class [label=\"Class\", shape=ellipse, style=filled, fillcolor=lightblue];")
	fmt.Fprintln(f, "    legend_edge [label=\"calls →\", shape=plaintext];")
	fmt.Fprintln(f, "  }")
	fmt.Fprintln(f, "")

	// Node definitions
	fmt.Fprintln(f, "  // Functions and classes")
	for funcID := range selectedFuncs {
		stat := funcStats[funcID]
		color := "lightgreen"
		if stat.isClass {
			color = "lightblue"
		}
		label := getShortName(funcID)
		fmt.Fprintf(f, "  \"%s\" [label=\"%s\", shape=ellipse, style=filled, fillcolor=%s];\n",
			funcID, label, color)
	}
	fmt.Fprintln(f, "")

	// Edges
	fmt.Fprintln(f, "  // Call edges")
	for _, morph := range cat.Morphisms() {
		if morph.Type == "function_call" {
			if selectedFuncs[morph.Source] && selectedFuncs[morph.Target] {
				fmt.Fprintf(f, "  \"%s\" -> \"%s\" [color=darkgreen];\n", morph.Source, morph.Target)
			}
		}
	}

	fmt.Fprintln(f, "}")
	return nil
}

func generateModuleDependencyGraphV2(cat *category.Category, prefix string, maxNodes int) error {
	filename := fmt.Sprintf("%s-modules-v2.dot", prefix)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Sample modules with most imports
	type moduleStat struct {
		id       string
		degree   int
		isImport bool
	}
	moduleStats := make(map[string]*moduleStat)

	for _, morph := range cat.Morphisms() {
		if morph.Type == "import" {
			if _, exists := moduleStats[morph.Source]; !exists {
				moduleStats[morph.Source] = &moduleStat{id: morph.Source, degree: 0}
			}
			if _, exists := moduleStats[morph.Target]; !exists {
				moduleStats[morph.Target] = &moduleStat{id: morph.Target, degree: 0}
			}
			moduleStats[morph.Source].degree++
			moduleStats[morph.Target].degree++

			// Mark target as imported module
			if obj, exists := cat.GetObject(morph.Target); exists && obj.Type == "imported_module" {
				moduleStats[morph.Target].isImport = true
			}
		}
	}

	var sorted []*moduleStat
	for _, stat := range moduleStats {
		sorted = append(sorted, stat)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].degree > sorted[j].degree
	})

	selectedModules := make(map[string]bool)
	limit := min(maxNodes, len(sorted))
	for i := 0; i < limit; i++ {
		selectedModules[sorted[i].id] = true
	}

	// Group modules into abstraction layers
	layers := map[string][]string{
		"external": []string{},  // External dependencies
		"internal": []string{},  // Internal modules
	}

	for modID := range selectedModules {
		stat := moduleStats[modID]
		if stat.isImport {
			layers["external"] = append(layers["external"], modID)
		} else {
			layers["internal"] = append(layers["internal"], modID)
		}
	}

	// Generate DOT with vertical layering
	fmt.Fprintln(f, "digraph ModuleDependencies {")
	fmt.Fprintln(f, "  // Graph layout settings")
	fmt.Fprintln(f, "  rankdir=TB;  // Top-to-Bottom (dependencies → dependents)")
	fmt.Fprintln(f, "  ranksep=1.2;")
	fmt.Fprintln(f, "  nodesep=0.8;")
	fmt.Fprintln(f, "  splines=true;")
	fmt.Fprintln(f, "  compound=true;")
	fmt.Fprintln(f, "")
	fmt.Fprintf(f, "  label=\"Module Dependencies (Top %d modules)\\nVertical Layers: External Dependencies (Top) → Internal Modules (Bottom)\";\n", limit)
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
	fmt.Fprintln(f, "    legend_external [label=\"External Dependency\", shape=component, style=filled, fillcolor=lightcyan];")
	fmt.Fprintln(f, "    legend_internal [label=\"Internal Module\", shape=folder, style=filled, fillcolor=wheat];")
	fmt.Fprintln(f, "    legend_edge [label=\"imports →\", shape=plaintext];")
	fmt.Fprintln(f, "  }")
	fmt.Fprintln(f, "")

	// External dependencies layer
	if len(layers["external"]) > 0 {
		fmt.Fprintln(f, "  // Layer 1: External Dependencies")
		fmt.Fprintln(f, "  subgraph cluster_external {")
		fmt.Fprintln(f, "    label=\"External Dependencies\";")
		fmt.Fprintln(f, "    style=filled;")
		fmt.Fprintln(f, "    fillcolor=\"#E8F5E9\";")
		fmt.Fprintln(f, "    fontsize=14;")
		fmt.Fprintln(f, "    rank=same;")
		fmt.Fprintln(f, "")

		for _, modID := range layers["external"] {
			label := getShortName(modID)
			fmt.Fprintf(f, "    \"%s\" [label=\"%s\", shape=component, style=filled, fillcolor=lightcyan];\n",
				modID, label)
		}
		fmt.Fprintln(f, "  }")
		fmt.Fprintln(f, "")
	}

	// Internal modules layer
	if len(layers["internal"]) > 0 {
		fmt.Fprintln(f, "  // Layer 2: Internal Modules")
		fmt.Fprintln(f, "  subgraph cluster_internal {")
		fmt.Fprintln(f, "    label=\"Internal Modules\";")
		fmt.Fprintln(f, "    style=filled;")
		fmt.Fprintln(f, "    fillcolor=\"#FFF3E0\";")
		fmt.Fprintln(f, "    fontsize=14;")
		fmt.Fprintln(f, "    rank=same;")
		fmt.Fprintln(f, "")

		for _, modID := range layers["internal"] {
			label := getShortName(modID)
			fmt.Fprintf(f, "    \"%s\" [label=\"%s\", shape=folder, style=filled, fillcolor=wheat];\n",
				modID, label)
		}
		fmt.Fprintln(f, "  }")
		fmt.Fprintln(f, "")
	}

	// Edges
	fmt.Fprintln(f, "  // Import edges")
	for _, morph := range cat.Morphisms() {
		if morph.Type == "import" {
			if selectedModules[morph.Source] && selectedModules[morph.Target] {
				fmt.Fprintf(f, "  \"%s\" -> \"%s\" [color=purple];\n", morph.Source, morph.Target)
			}
		}
	}

	fmt.Fprintln(f, "}")
	return nil
}

func generateCompositionGraphV2(cat *category.Category, prefix string, maxNodes int) error {
	filename := fmt.Sprintf("%s-composition-v2.dot", prefix)
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

	// Select top modules
	limit := min(20, len(sorted)) // Max 20 modules for composition

	// Generate DOT with vertical layering
	fmt.Fprintln(f, "digraph CompositionStructure {")
	fmt.Fprintln(f, "  // Graph layout settings")
	fmt.Fprintln(f, "  rankdir=TB;  // Top-to-Bottom")
	fmt.Fprintln(f, "  ranksep=1.5;")
	fmt.Fprintln(f, "  nodesep=0.7;")
	fmt.Fprintln(f, "  compound=true;")
	fmt.Fprintln(f, "")
	fmt.Fprintf(f, "  label=\"Composition Structure (Top %d modules)\\nVertical Layout: Modules contain Functions and Classes\";\n", limit)
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
	for i := 0; i < limit; i++ {
		stat := sorted[i]
		fmt.Fprintf(f, "  // Module %d\n", i+1)
		fmt.Fprintf(f, "  subgraph cluster_%d {\n", i)
		fmt.Fprintf(f, "    label=\"📦 %s\";\n", getShortName(stat.id))
		fmt.Fprintln(f, "    style=\"rounded,filled\";")
		fmt.Fprintln(f, "    fillcolor=lightgray;")
		fmt.Fprintln(f, "    fontsize=13;")
		fmt.Fprintln(f, "")

		// Limit objects per module
		objLimit := min(10, len(stat.objects))
		for j := 0; j < objLimit; j++ {
			obj := stat.objects[j]
			selectedObjects[obj.ID] = true
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
	}

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
