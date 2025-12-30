// Package analysis provides categorical complexity analysis for codebases.
//
// This package implements:
// - Basu-Isik diagram complexity: c(D) = Σc_objects + Σc_morphisms + c_composition
// - Kolmogorov complexity estimation via compression
// - Coupling metrics (afferent/efferent coupling, instability)
// - Cycle detection in dependency graphs
package analysis

import (
	"compress/gzip"
	"bytes"
	"encoding/json"
	"fmt"
	"math"

	"github.com/manu/catreview/pkg/category"
)

// ComplexityAnalyzer computes categorical complexity metrics.
type ComplexityAnalyzer struct {
	cat *category.Category
}

// NewComplexityAnalyzer creates a new complexity analyzer.
func NewComplexityAnalyzer(cat *category.Category) *ComplexityAnalyzer {
	return &ComplexityAnalyzer{cat: cat}
}

// DiagramComplexity computes Basu-Isik diagram complexity.
//
// Formula: c(D) = Σ c_obj(o) + Σ c_morph(m) + c_comp(D)
// where:
// - c_obj(o) = complexity of object o (based on metadata size)
// - c_morph(m) = complexity of morphism m (based on type and properties)
// - c_comp(D) = composition complexity (based on composable chains)
func (a *ComplexityAnalyzer) DiagramComplexity() float64 {
	objectComplexity := a.objectComplexity()
	morphismComplexity := a.morphismComplexity()
	compositionComplexity := a.compositionComplexity()

	return objectComplexity + morphismComplexity + compositionComplexity
}

// objectComplexity computes Σ c_obj(o) for all objects.
func (a *ComplexityAnalyzer) objectComplexity() float64 {
	total := 0.0
	for _, obj := range a.cat.Objects() {
		// Complexity based on metadata richness
		metadataSize := float64(len(obj.Metadata))
		typeComplexity := 1.0

		// Different object types have different base complexity
		switch obj.Type {
		case "struct", "interface":
			typeComplexity = 2.0
		case "function":
			typeComplexity = 1.5
		case "file", "package":
			typeComplexity = 0.5
		}

		total += typeComplexity * (1.0 + math.Log2(1.0+metadataSize))
	}
	return total
}

// morphismComplexity computes Σ c_morph(m) for all morphisms.
func (a *ComplexityAnalyzer) morphismComplexity() float64 {
	total := 0.0
	for _, morph := range a.cat.Morphisms() {
		if morph.Type == "identity" {
			continue // Identity morphisms have minimal complexity
		}

		// Base complexity by morphism type
		typeComplexity := 1.0
		switch morph.Type {
		case "import", "dependency":
			typeComplexity = 1.5
		case "function_call":
			typeComplexity = 2.0
		case "type_dependency":
			typeComplexity = 1.8
		case "composed":
			typeComplexity = 2.5
		}

		total += typeComplexity
	}
	return total
}

// compositionComplexity computes c_comp(D) based on composable chains.
func (a *ComplexityAnalyzer) compositionComplexity() float64 {
	// Build adjacency map for composition chains
	adjacency := make(map[string][]string)
	for _, morph := range a.cat.Morphisms() {
		if morph.Type == "identity" {
			continue
		}
		adjacency[morph.Target] = append(adjacency[morph.Target], morph.Source)
	}

	// Count composable chains of length 2+
	chains := 0.0
	for _, morph := range a.cat.Morphisms() {
		if morph.Type == "identity" {
			continue
		}
		// Count how many morphisms can compose with this one
		if sources, exists := adjacency[morph.Source]; exists {
			chains += float64(len(sources))
		}
	}

	// Complexity grows logarithmically with chain count
	return math.Log2(1.0 + chains)
}

// KolmogorovComplexity estimates Kolmogorov complexity via gzip compression.
//
// K(x) ≈ |gzip(x)|
//
// This is an upper bound on the true Kolmogorov complexity.
func (a *ComplexityAnalyzer) KolmogorovComplexity() (int, error) {
	// Serialize category to JSON
	data, err := json.Marshal(a.cat)
	if err != nil {
		return 0, err
	}

	// Compress with gzip
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(data); err != nil {
		return 0, err
	}
	if err := gz.Close(); err != nil {
		return 0, err
	}

	return buf.Len(), nil
}

// CouplingMetrics computes coupling metrics for each object.
type CouplingMetrics struct {
	ObjectID         string  `json:"object_id"`
	AfferentCoupling int     `json:"afferent_coupling"`  // Ca: incoming dependencies
	EfferentCoupling int     `json:"efferent_coupling"`  // Ce: outgoing dependencies
	Instability      float64 `json:"instability"`        // I = Ce / (Ca + Ce)
	Abstractness     float64 `json:"abstractness"`       // A (0-1, based on type)
}

// ComputeCoupling computes coupling metrics for all objects.
func (a *ComplexityAnalyzer) ComputeCoupling() map[string]*CouplingMetrics {
	metrics := make(map[string]*CouplingMetrics)

	// Initialize metrics for all objects
	for _, obj := range a.cat.Objects() {
		metrics[obj.ID] = &CouplingMetrics{
			ObjectID:         obj.ID,
			AfferentCoupling: 0,
			EfferentCoupling: 0,
			Abstractness:     a.computeAbstractness(obj),
		}
	}

	// Count incoming and outgoing dependencies
	for _, morph := range a.cat.Morphisms() {
		if morph.Type == "identity" {
			continue
		}

		// Efferent (outgoing) from source
		if m, exists := metrics[morph.Source]; exists {
			m.EfferentCoupling++
		}

		// Afferent (incoming) to target
		if m, exists := metrics[morph.Target]; exists {
			m.AfferentCoupling++
		}
	}

	// Compute instability: I = Ce / (Ca + Ce)
	for _, m := range metrics {
		total := m.AfferentCoupling + m.EfferentCoupling
		if total > 0 {
			m.Instability = float64(m.EfferentCoupling) / float64(total)
		}
	}

	return metrics
}

// computeAbstractness estimates abstractness of an object (0-1).
func (a *ComplexityAnalyzer) computeAbstractness(obj *category.Object) float64 {
	switch obj.Type {
	case "interface":
		return 1.0 // Fully abstract
	case "struct":
		return 0.1 // Mostly concrete
	case "function":
		if exported, ok := obj.Metadata["is_exported"].(bool); ok && exported {
			return 0.5 // Exported functions are semi-abstract (API)
		}
		return 0.2
	case "package":
		return 0.3
	default:
		return 0.0
	}
}

// CycleAnalyzer detects cycles in the dependency graph.
type CycleAnalyzer struct {
	cat *category.Category
}

// NewCycleAnalyzer creates a new cycle analyzer.
func NewCycleAnalyzer(cat *category.Category) *CycleAnalyzer {
	return &CycleAnalyzer{cat: cat}
}

// Cycle represents a detected dependency cycle.
type Cycle struct {
	Objects []string `json:"objects"`
	Length  int      `json:"length"`
}

// FindCycles detects all cycles in the dependency graph using DFS.
func (c *CycleAnalyzer) FindCycles() []*Cycle {
	// Build adjacency list
	adjacency := make(map[string][]string)
	for _, morph := range c.cat.Morphisms() {
		if morph.Type == "identity" {
			continue
		}
		adjacency[morph.Source] = append(adjacency[morph.Source], morph.Target)
	}

	// Track visited nodes and current path
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	var cycles []*Cycle

	// DFS from each node
	for _, obj := range c.cat.Objects() {
		if !visited[obj.ID] {
			path := []string{}
			c.dfs(obj.ID, adjacency, visited, recStack, path, &cycles)
		}
	}

	return cycles
}

// dfs performs depth-first search to detect cycles.
func (c *CycleAnalyzer) dfs(
	node string,
	adjacency map[string][]string,
	visited map[string]bool,
	recStack map[string]bool,
	path []string,
	cycles *[]*Cycle,
) {
	visited[node] = true
	recStack[node] = true
	path = append(path, node)

	// Visit neighbors
	for _, neighbor := range adjacency[node] {
		if !visited[neighbor] {
			c.dfs(neighbor, adjacency, visited, recStack, path, cycles)
		} else if recStack[neighbor] {
			// Found a cycle - extract it from path
			cycleStart := -1
			for i, n := range path {
				if n == neighbor {
					cycleStart = i
					break
				}
			}
			if cycleStart >= 0 {
				cycle := &Cycle{
					Objects: append([]string{}, path[cycleStart:]...),
					Length:  len(path) - cycleStart,
				}
				*cycles = append(*cycles, cycle)
			}
		}
	}

	recStack[node] = false
}

// Report generates a comprehensive analysis report.
type Report struct {
	CategoryStats    map[string]int           `json:"category_stats"`
	DiagramComplexity float64                 `json:"diagram_complexity"`
	KolmogorovComplexity int                  `json:"kolmogorov_complexity"`
	CouplingMetrics  map[string]*CouplingMetrics `json:"coupling_metrics"`
	Cycles           []*Cycle                 `json:"cycles"`
	TopUnstable      []*CouplingMetrics       `json:"top_unstable"`
	TopCoupled       []*CouplingMetrics       `json:"top_coupled"`
}

// GenerateReport creates a comprehensive analysis report.
func GenerateReport(cat *category.Category) (*Report, error) {
	complexityAnalyzer := NewComplexityAnalyzer(cat)
	cycleAnalyzer := NewCycleAnalyzer(cat)

	// Compute metrics
	diagramComplexity := complexityAnalyzer.DiagramComplexity()
	kolmogorov, err := complexityAnalyzer.KolmogorovComplexity()
	if err != nil {
		return nil, fmt.Errorf("failed to compute Kolmogorov complexity: %v", err)
	}

	couplingMetrics := complexityAnalyzer.ComputeCoupling()
	cycles := cycleAnalyzer.FindCycles()

	// Find top unstable and coupled components
	topUnstable := findTopN(couplingMetrics, 10, func(m *CouplingMetrics) float64 {
		return m.Instability
	})
	topCoupled := findTopN(couplingMetrics, 10, func(m *CouplingMetrics) float64 {
		return float64(m.AfferentCoupling + m.EfferentCoupling)
	})

	return &Report{
		CategoryStats:        cat.Stats(),
		DiagramComplexity:    diagramComplexity,
		KolmogorovComplexity: kolmogorov,
		CouplingMetrics:      couplingMetrics,
		Cycles:               cycles,
		TopUnstable:          topUnstable,
		TopCoupled:           topCoupled,
	}, nil
}

// findTopN finds top N components by a metric function.
func findTopN(metrics map[string]*CouplingMetrics, n int, scorer func(*CouplingMetrics) float64) []*CouplingMetrics {
	// Convert to slice
	slice := make([]*CouplingMetrics, 0, len(metrics))
	for _, m := range metrics {
		slice = append(slice, m)
	}

	// Simple selection sort for top N
	for i := 0; i < n && i < len(slice); i++ {
		maxIdx := i
		maxScore := scorer(slice[i])
		for j := i + 1; j < len(slice); j++ {
			score := scorer(slice[j])
			if score > maxScore {
				maxIdx = j
				maxScore = score
			}
		}
		slice[i], slice[maxIdx] = slice[maxIdx], slice[i]
	}

	if len(slice) > n {
		return slice[:n]
	}
	return slice
}
