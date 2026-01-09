// visualize_project.go - Create sampled Graphviz visualizations from project analysis
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
		fmt.Println("Usage: visualize_project --input <analysis.json> [--output <prefix>] [--max-nodes <N>]")
		fmt.Println("Example: visualize_project --input barque-analysis.json --output barque --max-nodes 100")
		os.Exit(1)
	}

	fmt.Printf("═══════════════════════════════════════════════════════════════\n")
	fmt.Printf("  Categorical Visualization Generator\n")
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

	if err := generateInheritanceGraph(&cat, *outputPrefix, *maxNodes); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error generating inheritance graph: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ %s-inheritance.dot created\n", *outputPrefix)

	if err := generateCallGraph(&cat, *outputPrefix, *maxNodes); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error generating call graph: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ %s-calls.dot created\n", *outputPrefix)

	if err := generateModuleDependencyGraph(&cat, *outputPrefix, *maxNodes); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error generating module graph: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ %s-modules.dot created\n", *outputPrefix)

	if err := generateCompositionGraph(&cat, *outputPrefix, *maxNodes); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error generating composition graph: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ %s-composition.dot created\n\n", *outputPrefix)

	fmt.Println("To generate SVG files, run:")
	fmt.Printf("  for f in %s-*.dot; do dot -Tsvg \"$f\" -o \"${f%.dot}.svg\"; done\n\n", *outputPrefix)
}

func generateInheritanceGraph(cat *category.Category, prefix string, maxNodes int) error {
	filename := fmt.Sprintf("%s-inheritance.dot", prefix)
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

	// Generate DOT
	fmt.Fprintln(f, "digraph InheritanceHierarchy {")
	fmt.Fprintln(f, "  rankdir=BT;")
	fmt.Fprintln(f, "  node [shape=box, style=filled, fillcolor=lightblue];")
	fmt.Fprintln(f, "  edge [color=darkblue, penwidth=2];")
	fmt.Fprintf(f, "  label=\"Inheritance Hierarchy (Top %d classes by connectivity)\";\n", limit)
	fmt.Fprintln(f, "  fontsize=16;")
	fmt.Fprintln(f, "")

	for classID := range selectedClasses {
		label := getShortName(classID)
		fmt.Fprintf(f, "  \"%s\" [label=\"%s\"];\n", classID, label)
	}
	fmt.Fprintln(f, "")

	for _, morph := range cat.Morphisms() {
		if morph.Type == "inheritance" {
			if selectedClasses[morph.Source] && selectedClasses[morph.Target] {
				fmt.Fprintf(f, "  \"%s\" -> \"%s\" [label=\"inherits\"];\n", morph.Source, morph.Target)
			}
		}
	}

	fmt.Fprintln(f, "}")
	return nil
}

func generateCallGraph(cat *category.Category, prefix string, maxNodes int) error {
	filename := fmt.Sprintf("%s-calls.dot", prefix)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Sample functions with most calls
	type funcStat struct {
		id     string
		degree int
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

	// Generate DOT
	fmt.Fprintln(f, "digraph CallGraph {")
	fmt.Fprintln(f, "  rankdir=LR;")
	fmt.Fprintln(f, "  node [shape=ellipse, style=filled];")
	fmt.Fprintln(f, "  edge [color=darkgreen];")
	fmt.Fprintf(f, "  label=\"Function Call Graph (Top %d functions by call frequency)\";\n", limit)
	fmt.Fprintln(f, "  fontsize=16;")
	fmt.Fprintln(f, "")

	for funcID := range selectedFuncs {
		obj, exists := cat.GetObject(funcID)
		color := "lightgreen"
		if exists && obj.Type == "class" {
			color = "lightblue"
		}
		label := getShortName(funcID)
		fmt.Fprintf(f, "  \"%s\" [label=\"%s\", fillcolor=%s];\n", funcID, label, color)
	}
	fmt.Fprintln(f, "")

	for _, morph := range cat.Morphisms() {
		if morph.Type == "function_call" {
			if selectedFuncs[morph.Source] && selectedFuncs[morph.Target] {
				fmt.Fprintf(f, "  \"%s\" -> \"%s\";\n", morph.Source, morph.Target)
			}
		}
	}

	fmt.Fprintln(f, "}")
	return nil
}

func generateModuleDependencyGraph(cat *category.Category, prefix string, maxNodes int) error {
	filename := fmt.Sprintf("%s-modules.dot", prefix)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Sample modules with most imports
	type moduleStat struct {
		id     string
		degree int
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

	// Generate DOT
	fmt.Fprintln(f, "digraph ModuleDependencies {")
	fmt.Fprintln(f, "  rankdir=TB;")
	fmt.Fprintln(f, "  node [shape=folder, style=filled, fillcolor=wheat];")
	fmt.Fprintln(f, "  edge [color=purple];")
	fmt.Fprintf(f, "  label=\"Module Dependencies (Top %d modules by import count)\";\n", limit)
	fmt.Fprintln(f, "  fontsize=16;")
	fmt.Fprintln(f, "")

	for modID := range selectedModules {
		obj, exists := cat.GetObject(modID)
		color := "wheat"
		shape := "folder"
		label := modID
		if exists {
			label = getShortName(obj.ID)
			if obj.Type == "imported_module" {
				color = "lightcyan"
				shape = "component"
			}
		}
		fmt.Fprintf(f, "  \"%s\" [label=\"%s\", fillcolor=%s, shape=%s];\n", modID, label, color, shape)
	}
	fmt.Fprintln(f, "")

	for _, morph := range cat.Morphisms() {
		if morph.Type == "import" {
			if selectedModules[morph.Source] && selectedModules[morph.Target] {
				fmt.Fprintf(f, "  \"%s\" -> \"%s\";\n", morph.Source, morph.Target)
			}
		}
	}

	fmt.Fprintln(f, "}")
	return nil
}

func generateCompositionGraph(cat *category.Category, prefix string, maxNodes int) error {
	filename := fmt.Sprintf("%s-composition.dot", prefix)
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

	// Generate DOT
	fmt.Fprintln(f, "digraph CompositionStructure {")
	fmt.Fprintln(f, "  rankdir=TB;")
	fmt.Fprintln(f, "  node [shape=box, style=\"rounded,filled\"];")
	fmt.Fprintln(f, "  edge [penwidth=1.5];")
	fmt.Fprintf(f, "  label=\"Composition Structure (Top %d modules by object count)\";\n", limit)
	fmt.Fprintln(f, "  fontsize=16;")
	fmt.Fprintln(f, "")

	selectedObjects := make(map[string]bool)
	for i := 0; i < limit; i++ {
		stat := sorted[i]
		fmt.Fprintf(f, "  subgraph cluster_%d {\n", i)
		fmt.Fprintf(f, "    label=\"%s\";\n", getShortName(stat.id))
		fmt.Fprintln(f, "    style=filled;")
		fmt.Fprintln(f, "    fillcolor=lightgray;")

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
	}
	fmt.Fprintln(f, "")

	// Add call edges between selected objects
	for _, morph := range cat.Morphisms() {
		if morph.Type == "function_call" {
			if selectedObjects[morph.Source] && selectedObjects[morph.Target] {
				fmt.Fprintf(f, "  \"%s\" -> \"%s\" [color=green, style=dashed];\n", morph.Source, morph.Target)
			}
		}
	}

	fmt.Fprintln(f, "}")
	return nil
}

func getShortName(fullID string) string {
	parts := strings.Split(fullID, ".")
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
