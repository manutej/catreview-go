package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/manu/catreview/pkg/category"
	"github.com/manu/catreview/pkg/extractor"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: analyze_project <path-to-python-project> [output-name]")
		fmt.Println("Example: analyze_project /Users/manu/Documents/LUXOR/PROJECTS/BARQUE barque")
		os.Exit(1)
	}

	projectPath := os.Args[1]
	outputName := "analysis"
	if len(os.Args) >= 3 {
		outputName = os.Args[2]
	}

	fmt.Printf("═══════════════════════════════════════════════════════════════\n")
	fmt.Printf("  Categorical Code Analysis Framework\n")
	fmt.Printf("═══════════════════════════════════════════════════════════════\n\n")

	// Verify project path exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "❌ Error: Project path does not exist: %s\n", projectPath)
		os.Exit(1)
	}

	fmt.Printf("📂 Project Path: %s\n", projectPath)
	fmt.Printf("📊 Output Name: %s\n\n", outputName)

	// Count Python files first
	pythonFiles, err := countPythonFiles(projectPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error counting files: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("📈 Python Files Found: %d\n\n", pythonFiles)

	if pythonFiles == 0 {
		fmt.Println("⚠️  No Python files found in project")
		os.Exit(0)
	}

	// Create Python extractor
	fmt.Println("🔧 Initializing Python extractor...")
	pythonExtractor := extractor.NewPythonExtractor()

	// Extract categorical model
	fmt.Println("🔍 Extracting categorical model...")
	startTime := time.Now()

	cat, err := pythonExtractor.ExtractFromPath(projectPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error extracting: %v\n", err)
		os.Exit(1)
	}

	extractionTime := time.Since(startTime)
	fmt.Printf("✅ Extraction completed in %v\n\n", extractionTime)

	// Print comprehensive statistics
	printStatistics(cat, pythonFiles)

	// Export to JSON
	outputPath := fmt.Sprintf("%s-analysis.json", outputName)
	fmt.Printf("\n💾 Exporting categorical model to: %s\n", outputPath)

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cat); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error encoding JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Export successful")

	// Generate summary report
	generateSummaryReport(cat, outputName, projectPath, pythonFiles, extractionTime)

	fmt.Printf("\n═══════════════════════════════════════════════════════════════\n")
	fmt.Printf("  Analysis Complete!\n")
	fmt.Printf("═══════════════════════════════════════════════════════════════\n\n")
	fmt.Printf("📄 Next Steps:\n")
	fmt.Printf("  1. Review: cat %s-summary.txt\n", outputName)
	fmt.Printf("  2. Visualize: go run visualize.go --input %s-analysis.json --output %s\n", outputName, outputName)
	fmt.Printf("  3. Explore: jq '.objects | length' %s-analysis.json\n\n", outputName)
}

func countPythonFiles(root string) (int, error) {
	count := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if !info.IsDir() && filepath.Ext(path) == ".py" {
			count++
		}
		return nil
	})
	return count, err
}

func printStatistics(cat *category.Category, totalFiles int) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("  Categorical Model Statistics")
	fmt.Println("═══════════════════════════════════════════════════════════════\n")

	objects := cat.Objects()
	morphisms := cat.Morphisms()

	fmt.Printf("📦 Total Objects: %d\n", len(objects))
	fmt.Printf("🔗 Total Morphisms: %d\n", len(morphisms))
	fmt.Printf("📐 Ratio: %.2f morphisms per object\n\n", float64(len(morphisms))/float64(len(objects)))

	// Categorize objects by type
	objectsByType := make(map[string]int)
	for _, obj := range objects {
		objectsByType[obj.Type]++
	}

	fmt.Println("📊 Objects by Type:")
	// Sort by count descending
	type kv struct {
		Type  string
		Count int
	}
	var sorted []kv
	for t, c := range objectsByType {
		sorted = append(sorted, kv{t, c})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Count > sorted[j].Count
	})

	for _, item := range sorted {
		percentage := float64(item.Count) / float64(len(objects)) * 100
		fmt.Printf("  • %-12s: %5d (%.1f%%)\n", item.Type, item.Count, percentage)
	}

	// Categorize morphisms by type
	morphismsByType := make(map[string]int)
	for _, morph := range morphisms {
		morphismsByType[morph.Type]++
	}

	fmt.Println("\n🔗 Morphisms by Type:")
	var sortedMorphs []kv
	for t, c := range morphismsByType {
		sortedMorphs = append(sortedMorphs, kv{t, c})
	}
	sort.Slice(sortedMorphs, func(i, j int) bool {
		return sortedMorphs[i].Count > sortedMorphs[j].Count
	})

	for _, item := range sortedMorphs {
		percentage := float64(item.Count) / float64(len(morphisms)) * 100
		fmt.Printf("  • %-15s: %5d (%.1f%%)\n", item.Type, item.Count, percentage)
	}

	// Verify category axioms
	fmt.Println("\n🔬 Category Axioms Verification:")
	identityCount := 0
	for _, morph := range morphisms {
		if morph.Type == "identity" {
			identityCount++
		}
	}

	if identityCount == len(objects) {
		fmt.Printf("  ✅ Identity: All %d objects have identity morphisms\n", len(objects))
	} else {
		fmt.Printf("  ⚠️  Identity: %d/%d objects have identity morphisms\n", identityCount, len(objects))
	}

	// Composition check (simplified - just verify no broken references)
	brokenMorphisms := 0
	for _, morph := range morphisms {
		sourceExists := false
		targetExists := false
		for _, obj := range objects {
			if obj.ID == morph.Source {
				sourceExists = true
			}
			if obj.ID == morph.Target {
				targetExists = true
			}
		}
		if !sourceExists || !targetExists {
			brokenMorphisms++
		}
	}

	if brokenMorphisms == 0 {
		fmt.Printf("  ✅ Composition: All morphisms reference valid objects\n")
	} else {
		fmt.Printf("  ⚠️  Composition: %d morphisms reference missing objects\n", brokenMorphisms)
	}

	// Calculate complexity metrics
	fmt.Println("\n📈 Complexity Metrics:")
	avgInDegree := float64(len(morphisms)) / float64(len(objects))
	fmt.Printf("  • Avg In-Degree: %.2f\n", avgInDegree)

	// Count inheritance depth (classes with bases)
	classesWithBases := 0
	for _, obj := range objects {
		if obj.Type == "class" {
			if bases, ok := obj.Metadata["bases"].([]interface{}); ok && len(bases) > 0 {
				classesWithBases++
			}
		}
	}
	fmt.Printf("  • Classes with Inheritance: %d/%d\n", classesWithBases, objectsByType["class"])

	// Count functions with calls
	functionsWithCalls := 0
	for _, morph := range morphisms {
		if morph.Type == "function_call" {
			functionsWithCalls++
		}
	}
	fmt.Printf("  • Function Calls: %d\n", functionsWithCalls)
}

func generateSummaryReport(cat *category.Category, outputName, projectPath string, totalFiles int, duration time.Duration) {
	reportPath := fmt.Sprintf("%s-summary.txt", outputName)
	file, err := os.Create(reportPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Warning: Could not create summary report: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "═══════════════════════════════════════════════════════════════\n")
	fmt.Fprintf(file, "  Categorical Code Analysis Summary\n")
	fmt.Fprintf(file, "═══════════════════════════════════════════════════════════════\n\n")

	fmt.Fprintf(file, "Project: %s\n", projectPath)
	fmt.Fprintf(file, "Analyzed: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "Duration: %v\n", duration)
	fmt.Fprintf(file, "Python Files: %d\n\n", totalFiles)

	objects := cat.Objects()
	morphisms := cat.Morphisms()

	fmt.Fprintf(file, "Objects: %d\n", len(objects))
	fmt.Fprintf(file, "Morphisms: %d\n", len(morphisms))
	fmt.Fprintf(file, "Ratio: %.2f morphisms per object\n\n", float64(len(morphisms))/float64(len(objects)))

	// Top modules by object count
	fmt.Fprintf(file, "═══════════════════════════════════════════════════════════════\n")
	fmt.Fprintf(file, "Top Modules (by object count)\n")
	fmt.Fprintf(file, "═══════════════════════════════════════════════════════════════\n\n")

	moduleObjects := make(map[string]int)
	for _, obj := range objects {
		// Extract module from object ID (first component before last dot)
		parts := splitObjectID(obj.ID)
		if len(parts) > 1 {
			module := parts[0]
			moduleObjects[module]++
		}
	}

	type moduleCount struct {
		Module string
		Count  int
	}
	var modules []moduleCount
	for mod, count := range moduleObjects {
		modules = append(modules, moduleCount{mod, count})
	}
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Count > modules[j].Count
	})

	for i, mod := range modules {
		if i >= 10 {
			break
		}
		fmt.Fprintf(file, "%2d. %-50s %5d objects\n", i+1, mod.Module, mod.Count)
	}

	fmt.Printf("✅ Summary report saved to: %s\n", reportPath)
}

func splitObjectID(id string) []string {
	var parts []string
	current := ""
	for _, ch := range id {
		if ch == '.' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
