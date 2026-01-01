// catreview - Categorical Codebase Review Tool
//
// A tool for analyzing software architecture using category theory.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/manu/catreview/pkg/analysis"
	"github.com/manu/catreview/pkg/category"
	"github.com/manu/catreview/pkg/extractor"
	"github.com/manu/catreview/pkg/functor"
	"github.com/manu/catreview/pkg/viz"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "catreview",
		Short: "Categorical codebase review and analysis",
		Long: `catreview analyzes software architecture using category theory.

It extracts categorical models from codebases and provides insights into:
- Architectural complexity (Basu-Isik diagram complexity)
- Coupling and instability metrics
- Dependency cycles
- Abstraction levels via functors`,
	}

	extractCmd = &cobra.Command{
		Use:   "extract [path]",
		Short: "Extract categorical model from codebase",
		Args:  cobra.ExactArgs(1),
		RunE:  runExtract,
	}

	analyzeCmd = &cobra.Command{
		Use:   "analyze [model.json]",
		Short: "Analyze categorical model and generate report",
		Args:  cobra.ExactArgs(1),
		RunE:  runAnalyze,
	}

	verifyCmd = &cobra.Command{
		Use:   "verify [model.json]",
		Short: "Verify category axioms",
		Args:  cobra.ExactArgs(1),
		RunE:  runVerify,
	}

	abstractCmd = &cobra.Command{
		Use:   "abstract [model.json]",
		Short: "Create package-level abstraction via functor",
		Args:  cobra.ExactArgs(1),
		RunE:  runAbstract,
	}

	vizCmd = &cobra.Command{
		Use:   "viz [model.json]",
		Short: "Visualize categorical model as dependency graph",
		Long: `Visualize a categorical model as a dependency graph with layered architecture.

Outputs the graph in various formats:
  ascii   - ASCII art visualization (default)
  mermaid - Mermaid.js flowchart syntax
  dot     - Graphviz DOT format
  json    - JSON representation

Supports filtering by layer and coupling thresholds for cleaner views.`,
		Args: cobra.ExactArgs(1),
		RunE: runViz,
	}

	// Flags
	outputFile      string
	formatJSON      bool
	maxCycles       int
	failOnViolation bool

	// Viz flags
	vizFormat      string
	vizLayer       int
	vizMinCoupling int
	vizMaxNodes    int
	vizHeatmap     bool
	vizLayered     bool
)

func init() {
	// Extract command flags
	extractCmd.Flags().StringVarP(&outputFile, "output", "o", "model.json", "Output file for categorical model")
	extractCmd.Flags().BoolVar(&formatJSON, "pretty", false, "Pretty-print JSON output")

	// Analyze command flags
	analyzeCmd.Flags().StringVarP(&outputFile, "output", "o", "report.json", "Output file for analysis report")
	analyzeCmd.Flags().BoolVar(&formatJSON, "pretty", true, "Pretty-print JSON output")

	// Verify command flags
	verifyCmd.Flags().IntVar(&maxCycles, "max-cycles", -1, "Maximum allowed cycles (-1 = no limit)")
	verifyCmd.Flags().BoolVar(&failOnViolation, "fail-on-violation", false, "Exit with error on axiom violation")

	// Abstract command flags
	abstractCmd.Flags().StringVarP(&outputFile, "output", "o", "abstract.json", "Output file for abstracted model")
	abstractCmd.Flags().BoolVar(&formatJSON, "pretty", true, "Pretty-print JSON output")

	// Viz command flags
	vizCmd.Flags().StringVarP(&vizFormat, "format", "f", "ascii", "Output format: ascii, mermaid, dot, json")
	vizCmd.Flags().IntVarP(&vizLayer, "layer", "l", -1, "Extract specific layer (0-3), -1 for all")
	vizCmd.Flags().IntVar(&vizMinCoupling, "min-coupling", 0, "Minimum total coupling to include node")
	vizCmd.Flags().IntVar(&vizMaxNodes, "max-nodes", 0, "Maximum nodes to display (0 = unlimited)")
	vizCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
	vizCmd.Flags().BoolVar(&vizHeatmap, "heatmap", false, "Generate coupling heatmap instead of graph")
	vizCmd.Flags().BoolVar(&vizLayered, "layered", false, "Generate detailed layered ASCII view")

	rootCmd.AddCommand(extractCmd, analyzeCmd, verifyCmd, abstractCmd, vizCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runExtract(cmd *cobra.Command, args []string) error {
	path := args[0]

	fmt.Printf("Extracting categorical model from: %s\n", path)

	// Create extractor
	ext := extractor.NewGoExtractor()

	// Extract from path
	cat, err := ext.ExtractFromPath(path)
	if err != nil {
		return fmt.Errorf("extraction failed: %v", err)
	}

	// Print statistics
	stats := cat.Stats()
	fmt.Printf("Extracted:\n")
	fmt.Printf("  Objects:   %d\n", stats["objects"])
	fmt.Printf("  Morphisms: %d\n", stats["morphisms"])
	fmt.Printf("  Identities: %d\n", stats["identities"])

	// Save to file
	if err := saveCategory(cat, outputFile); err != nil {
		return fmt.Errorf("failed to save model: %v", err)
	}

	fmt.Printf("Model saved to: %s\n", outputFile)
	return nil
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	modelFile := args[0]

	fmt.Printf("Analyzing categorical model: %s\n", modelFile)

	// Load model
	cat, err := loadCategory(modelFile)
	if err != nil {
		return fmt.Errorf("failed to load model: %v", err)
	}

	// Generate report
	report, err := analysis.GenerateReport(cat)
	if err != nil {
		return fmt.Errorf("analysis failed: %v", err)
	}

	// Print summary
	fmt.Printf("\nCategorical Analysis Report\n")
	fmt.Printf("===========================\n\n")
	fmt.Printf("Category Statistics:\n")
	fmt.Printf("  Objects:    %d\n", report.CategoryStats["objects"])
	fmt.Printf("  Morphisms:  %d\n", report.CategoryStats["morphisms"])
	fmt.Printf("\nComplexity Metrics:\n")
	fmt.Printf("  Diagram Complexity:    %.2f\n", report.DiagramComplexity)
	fmt.Printf("  Kolmogorov Complexity: %d bytes\n", report.KolmogorovComplexity)
	fmt.Printf("\nDependency Analysis:\n")
	fmt.Printf("  Cycles Found: %d\n", len(report.Cycles))

	if len(report.Cycles) > 0 {
		fmt.Printf("\n  Top 5 Cycles:\n")
		for i, cycle := range report.Cycles {
			if i >= 5 {
				break
			}
			fmt.Printf("    Cycle %d (length %d): %v\n", i+1, cycle.Length, cycle.Objects)
		}
	}

	fmt.Printf("\nTop 5 Most Unstable Components:\n")
	for i, m := range report.TopUnstable {
		if i >= 5 {
			break
		}
		fmt.Printf("  %s: I=%.2f (Ce=%d, Ca=%d)\n",
			m.ObjectID, m.Instability, m.EfferentCoupling, m.AfferentCoupling)
	}

	fmt.Printf("\nTop 5 Most Coupled Components:\n")
	for i, m := range report.TopCoupled {
		if i >= 5 {
			break
		}
		total := m.AfferentCoupling + m.EfferentCoupling
		fmt.Printf("  %s: %d total (Ce=%d, Ca=%d)\n",
			m.ObjectID, total, m.EfferentCoupling, m.AfferentCoupling)
	}

	// Save full report
	if err := saveJSON(report, outputFile); err != nil {
		return fmt.Errorf("failed to save report: %v", err)
	}

	fmt.Printf("\nFull report saved to: %s\n", outputFile)
	return nil
}

func runVerify(cmd *cobra.Command, args []string) error {
	modelFile := args[0]

	fmt.Printf("Verifying category axioms: %s\n", modelFile)

	// Load model
	cat, err := loadCategory(modelFile)
	if err != nil {
		return fmt.Errorf("failed to load model: %v", err)
	}

	// Verify axioms
	fmt.Printf("Checking associativity and identity laws...\n")
	if err := cat.VerifyAxioms(); err != nil {
		fmt.Printf("❌ Axiom verification FAILED: %v\n", err)
		if failOnViolation {
			return err
		}
		return nil
	}
	fmt.Printf("✅ Category axioms verified successfully\n")

	// Check cycles if limit specified
	if maxCycles >= 0 {
		fmt.Printf("\nChecking for dependency cycles (max allowed: %d)...\n", maxCycles)
		cycleAnalyzer := analysis.NewCycleAnalyzer(cat)
		cycles := cycleAnalyzer.FindCycles()

		fmt.Printf("Found %d cycles\n", len(cycles))

		if len(cycles) > maxCycles {
			fmt.Printf("❌ Cycle limit exceeded: %d > %d\n", len(cycles), maxCycles)
			if failOnViolation {
				return fmt.Errorf("too many cycles: %d > %d", len(cycles), maxCycles)
			}
		} else {
			fmt.Printf("✅ Cycle count within limit\n")
		}
	}

	return nil
}

func runAbstract(cmd *cobra.Command, args []string) error {
	modelFile := args[0]

	fmt.Printf("Creating package-level abstraction from: %s\n", modelFile)

	// Load file-level model
	fileCat, err := loadCategory(modelFile)
	if err != nil {
		return fmt.Errorf("failed to load model: %v", err)
	}

	// Create target category for packages
	pkgCat := category.NewCategory("package_level")

	// Create abstraction functor
	f := functor.NewPackageAbstractionFunctor(fileCat, pkgCat)

	// Apply functor to all objects
	fmt.Printf("Mapping files to packages...\n")
	for _, obj := range fileCat.Objects() {
		if _, err := f.MapObject(obj); err != nil {
			fmt.Printf("Warning: failed to map %s: %v\n", obj.ID, err)
		}
	}

	// Apply functor to all morphisms
	fmt.Printf("Mapping dependencies...\n")
	for _, morph := range fileCat.Morphisms() {
		if _, err := f.MapMorphism(morph); err != nil {
			// Silently skip errors (many will fail due to missing objects)
			continue
		}
	}

	// Verify functor laws
	fmt.Printf("Verifying functor laws...\n")
	if err := f.VerifyLaws(); err != nil {
		fmt.Printf("Warning: functor law verification failed: %v\n", err)
	} else {
		fmt.Printf("✅ Functor laws verified\n")
	}

	// Print statistics
	stats := pkgCat.Stats()
	fmt.Printf("\nAbstracted Category:\n")
	fmt.Printf("  Packages:           %d\n", stats["objects"])
	fmt.Printf("  Package Dependencies: %d\n", stats["morphisms"])

	// Save abstracted model
	if err := saveCategory(pkgCat, outputFile); err != nil {
		return fmt.Errorf("failed to save abstracted model: %v", err)
	}

	fmt.Printf("\nAbstracted model saved to: %s\n", outputFile)
	return nil
}

func runViz(cmd *cobra.Command, args []string) error {
	modelFile := args[0]

	fmt.Fprintf(os.Stderr, "Loading categorical model: %s\n", modelFile)

	// Load model
	cat, err := loadCategory(modelFile)
	if err != nil {
		return fmt.Errorf("failed to load model: %v", err)
	}

	// Build visualization graph
	builder := viz.NewGraphBuilder(cat)
	graph := builder.Build()

	fmt.Fprintf(os.Stderr, "Built graph: %d nodes, %d edges, DAG: %v\n",
		graph.Stats.TotalNodes, graph.Stats.TotalEdges, graph.IsDAG)

	// Apply filters if specified
	if vizMinCoupling > 0 || vizMaxNodes > 0 {
		filterOpts := &viz.FilterOptions{
			MinTotalCoupling: vizMinCoupling,
			MaxNodes:         vizMaxNodes,
		}
		graph = graph.Filter(filterOpts)
		fmt.Fprintf(os.Stderr, "After filtering: %d nodes, %d edges\n",
			graph.Stats.TotalNodes, graph.Stats.TotalEdges)
	}

	// Extract specific layer if requested
	if vizLayer >= 0 && vizLayer <= 3 {
		graph = graph.ExtractLayer(vizLayer)
		fmt.Fprintf(os.Stderr, "Extracted layer %d: %d nodes, %d edges\n",
			vizLayer, graph.Stats.TotalNodes, graph.Stats.TotalEdges)
	}

	// Generate output
	var output string

	if vizHeatmap {
		// Generate heatmap view
		output = viz.GenerateHeatmap(graph)
	} else if vizLayered {
		// Generate detailed layered view
		output = viz.GenerateLayeredASCII(graph)
	} else {
		// Generate standard visualization
		format := viz.Format(vizFormat)
		generator := viz.NewGenerator(format)
		output, err = generator.Generate(graph)
		if err != nil {
			return fmt.Errorf("generation failed: %v", err)
		}
	}

	// Output to file or stdout
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write output: %v", err)
		}
		fmt.Fprintf(os.Stderr, "Output written to: %s\n", outputFile)
	} else {
		fmt.Print(output)
	}

	return nil
}

// Helper functions

func saveCategory(cat *category.Category, filename string) error {
	return saveJSON(cat, filename)
}

func loadCategory(filename string) (*category.Category, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cat category.Category
	if err := json.Unmarshal(data, &cat); err != nil {
		return nil, err
	}

	return &cat, nil
}

func saveJSON(v interface{}, filename string) error {
	var data []byte
	var err error

	if formatJSON {
		data, err = json.MarshalIndent(v, "", "  ")
	} else {
		data, err = json.Marshal(v)
	}

	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
