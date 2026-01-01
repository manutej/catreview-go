// graph.go provides DAG graph building from categorical models.
package viz

import (
	"encoding/json"
	"sort"

	"github.com/manu/catreview/pkg/analysis"
	"github.com/manu/catreview/pkg/category"
)

// Node represents a node in the visualization graph.
type Node struct {
	ID          string                 `json:"id"`
	Label       string                 `json:"label"`
	Type        string                 `json:"type"`    // file, import, function, etc.
	Layer       int                    `json:"layer"`   // Architectural layer (0-3)
	Afferent    int                    `json:"afferent"`
	Efferent    int                    `json:"efferent"`
	Instability float64                `json:"instability"`
	Abstractness float64               `json:"abstractness"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Edge represents an edge in the visualization graph.
type Edge struct {
	Source   string `json:"source"`
	Target   string `json:"target"`
	Type     string `json:"type"`   // dependency, import, call, cross_layer, inferred
	Weight   int    `json:"weight"` // Coupling strength
}

// Layer represents an architectural layer in the graph.
type Layer struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	NodeIDs     []string `json:"node_ids"`
}

// GraphStats contains summary statistics for the graph.
type GraphStats struct {
	TotalNodes    int     `json:"total_nodes"`
	TotalEdges    int     `json:"total_edges"`
	AvgCoupling   float64 `json:"avg_coupling"`
	MaxCoupling   int     `json:"max_coupling"`
	LayerCounts   []int   `json:"layer_counts"`
}

// Graph represents the complete visualization graph.
type Graph struct {
	Name      string           `json:"name"`
	Nodes     map[string]*Node `json:"nodes"`
	Edges     []*Edge          `json:"edges"`
	Layers    []*Layer         `json:"layers"`
	Stats     *GraphStats      `json:"stats"`
	IsDAG     bool             `json:"is_dag"`
}

// FilterOptions specifies filtering criteria for the graph.
type FilterOptions struct {
	MinAfferent      int
	MinEfferent      int
	MinTotalCoupling int
	MaxNodes         int
	LayerFilter      []int
	TypeFilter       []string
}

// GraphBuilder constructs visualization graphs from categorical data.
type GraphBuilder struct {
	category *category.Category
	coupling map[string]*analysis.CouplingMetrics
}

// NewGraphBuilder creates a new graph builder from a category.
func NewGraphBuilder(cat *category.Category) *GraphBuilder {
	return &GraphBuilder{
		category: cat,
		coupling: make(map[string]*analysis.CouplingMetrics),
	}
}

// Build constructs the visualization graph from the category.
func (b *GraphBuilder) Build() *Graph {
	graph := &Graph{
		Name:   b.category.Name,
		Nodes:  make(map[string]*Node),
		Edges:  make([]*Edge, 0),
		Layers: make([]*Layer, 4),
		Stats:  &GraphStats{LayerCounts: make([]int, 4)},
		IsDAG:  true,
	}

	// Initialize layers
	layerNames := []string{"Foundations", "Core", "Support", "Applications"}
	layerDescs := []string{
		"Stable imports and external dependencies (I=0)",
		"Framework core with high efferent coupling",
		"Supporting modules and utilities",
		"Application code (I=1.0)",
	}
	for i := 0; i < 4; i++ {
		graph.Layers[i] = &Layer{
			ID:          i,
			Name:        layerNames[i],
			Description: layerDescs[i],
			NodeIDs:     make([]string, 0),
		}
	}

	// Calculate coupling metrics
	b.calculateCoupling()

	// Build nodes from objects
	for _, obj := range b.category.Objects() {
		node := b.objectToNode(obj)
		graph.Nodes[node.ID] = node

		// Assign to layer
		graph.Layers[node.Layer].NodeIDs = append(graph.Layers[node.Layer].NodeIDs, node.ID)
		graph.Stats.LayerCounts[node.Layer]++
	}

	// Build edges from morphisms
	for _, morph := range b.category.Morphisms() {
		if morph.Source == morph.Target {
			continue // Skip identity morphisms
		}

		edge := &Edge{
			Source: morph.Source,
			Target: morph.Target,
			Type:   b.classifyEdge(morph, graph),
			Weight: 1,
		}

		// Check for cross-layer edges
		if srcNode, ok := graph.Nodes[morph.Source]; ok {
			if tgtNode, ok := graph.Nodes[morph.Target]; ok {
				if srcNode.Layer != tgtNode.Layer {
					edge.Type = "cross_layer"
				}
			}
		}

		graph.Edges = append(graph.Edges, edge)
	}

	// Sort layer nodes by coupling
	for _, layer := range graph.Layers {
		sort.Slice(layer.NodeIDs, func(i, j int) bool {
			ni, oki := graph.Nodes[layer.NodeIDs[i]]
			nj, okj := graph.Nodes[layer.NodeIDs[j]]
			if !oki || !okj {
				return false
			}
			ti := ni.Afferent + ni.Efferent
			tj := nj.Afferent + nj.Efferent
			return ti > tj
		})
	}

	// Calculate statistics
	graph.Stats.TotalNodes = len(graph.Nodes)
	graph.Stats.TotalEdges = len(graph.Edges)

	totalCoupling := 0
	for _, node := range graph.Nodes {
		coupling := node.Afferent + node.Efferent
		totalCoupling += coupling
		if coupling > graph.Stats.MaxCoupling {
			graph.Stats.MaxCoupling = coupling
		}
	}
	if graph.Stats.TotalNodes > 0 {
		graph.Stats.AvgCoupling = float64(totalCoupling) / float64(graph.Stats.TotalNodes)
	}

	// Check for cycles (simplified - just checks if we have back edges)
	graph.IsDAG = b.checkDAG(graph)

	return graph
}

// calculateCoupling computes coupling metrics for all objects.
func (b *GraphBuilder) calculateCoupling() {
	// Count afferent (incoming) and efferent (outgoing) for each object
	afferent := make(map[string]int)
	efferent := make(map[string]int)

	for _, morph := range b.category.Morphisms() {
		if morph.Source == morph.Target {
			continue
		}
		efferent[morph.Source]++
		afferent[morph.Target]++
	}

	for _, obj := range b.category.Objects() {
		ca := afferent[obj.ID]
		ce := efferent[obj.ID]
		instability := 0.0
		if ca+ce > 0 {
			instability = float64(ce) / float64(ca+ce)
		}

		b.coupling[obj.ID] = &analysis.CouplingMetrics{
			ObjectID:         obj.ID,
			AfferentCoupling: ca,
			EfferentCoupling: ce,
			Instability:      instability,
		}
	}
}

// objectToNode converts a category object to a graph node.
func (b *GraphBuilder) objectToNode(obj *category.Object) *Node {
	metrics := b.coupling[obj.ID]
	if metrics == nil {
		metrics = &analysis.CouplingMetrics{}
	}

	node := &Node{
		ID:          obj.ID,
		Label:       b.extractLabel(obj.ID),
		Type:        obj.Type,
		Afferent:    metrics.AfferentCoupling,
		Efferent:    metrics.EfferentCoupling,
		Instability: metrics.Instability,
		Metadata:    obj.Metadata,
	}

	// Assign layer based on instability and type
	node.Layer = b.assignLayer(node)

	return node
}

// extractLabel extracts a short label from an object ID.
func (b *GraphBuilder) extractLabel(id string) string {
	// For file paths, extract just the filename
	for i := len(id) - 1; i >= 0; i-- {
		if id[i] == '/' {
			return id[i+1:]
		}
	}
	return id
}

// assignLayer determines the architectural layer for a node.
func (b *GraphBuilder) assignLayer(node *Node) int {
	// Layer 0: Foundations - stable imports (I=0, high afferent)
	if node.Instability == 0 && node.Afferent > 0 {
		return 0
	}

	// Layer 3: Applications - unstable application code (I=1.0)
	if node.Instability == 1.0 && node.Efferent > 0 && node.Afferent == 0 {
		// Check if it's a main.go or example
		if node.Label == "main.go" || node.Type == "example" {
			return 3
		}
	}

	// Layer 1: Core - high efferent coupling (framework core)
	if node.Efferent > 40 {
		return 1
	}

	// Layer 2: Support - everything else
	return 2
}

// classifyEdge determines the type of an edge.
func (b *GraphBuilder) classifyEdge(morph *category.Morphism, graph *Graph) string {
	if morph.Type != "" {
		return morph.Type
	}
	return "dependency"
}

// checkDAG verifies the graph is a directed acyclic graph.
func (b *GraphBuilder) checkDAG(graph *Graph) bool {
	// Simple cycle detection using DFS
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var hasCycle func(nodeID string) bool
	hasCycle = func(nodeID string) bool {
		visited[nodeID] = true
		recStack[nodeID] = true

		// Check all outgoing edges
		for _, edge := range graph.Edges {
			if edge.Source == nodeID {
				if !visited[edge.Target] {
					if hasCycle(edge.Target) {
						return true
					}
				} else if recStack[edge.Target] {
					return true
				}
			}
		}

		recStack[nodeID] = false
		return false
	}

	for nodeID := range graph.Nodes {
		if !visited[nodeID] {
			if hasCycle(nodeID) {
				return false
			}
		}
	}

	return true
}

// Filter creates a new graph with only nodes matching the filter criteria.
func (g *Graph) Filter(opts *FilterOptions) *Graph {
	if opts == nil {
		return g
	}

	filtered := &Graph{
		Name:   g.Name,
		Nodes:  make(map[string]*Node),
		Edges:  make([]*Edge, 0),
		Layers: make([]*Layer, 4),
		Stats:  &GraphStats{LayerCounts: make([]int, 4)},
		IsDAG:  g.IsDAG,
	}

	// Initialize layers
	for i := 0; i < 4; i++ {
		filtered.Layers[i] = &Layer{
			ID:          g.Layers[i].ID,
			Name:        g.Layers[i].Name,
			Description: g.Layers[i].Description,
			NodeIDs:     make([]string, 0),
		}
	}

	// Collect nodes that pass filter, sorted by coupling
	type nodeCoupling struct {
		node     *Node
		coupling int
	}
	candidates := make([]nodeCoupling, 0)

	for _, node := range g.Nodes {
		totalCoupling := node.Afferent + node.Efferent

		// Apply filters
		if opts.MinAfferent > 0 && node.Afferent < opts.MinAfferent {
			continue
		}
		if opts.MinEfferent > 0 && node.Efferent < opts.MinEfferent {
			continue
		}
		if opts.MinTotalCoupling > 0 && totalCoupling < opts.MinTotalCoupling {
			continue
		}
		if len(opts.LayerFilter) > 0 {
			found := false
			for _, l := range opts.LayerFilter {
				if node.Layer == l {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if len(opts.TypeFilter) > 0 {
			found := false
			for _, t := range opts.TypeFilter {
				if node.Type == t {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		candidates = append(candidates, nodeCoupling{node: node, coupling: totalCoupling})
	}

	// Sort by coupling (highest first)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].coupling > candidates[j].coupling
	})

	// Apply max nodes limit
	if opts.MaxNodes > 0 && len(candidates) > opts.MaxNodes {
		candidates = candidates[:opts.MaxNodes]
	}

	// Add filtered nodes
	for _, nc := range candidates {
		filtered.Nodes[nc.node.ID] = nc.node
		filtered.Layers[nc.node.Layer].NodeIDs = append(filtered.Layers[nc.node.Layer].NodeIDs, nc.node.ID)
		filtered.Stats.LayerCounts[nc.node.Layer]++
	}

	// Add edges between filtered nodes
	for _, edge := range g.Edges {
		_, srcOK := filtered.Nodes[edge.Source]
		_, tgtOK := filtered.Nodes[edge.Target]
		if srcOK && tgtOK {
			filtered.Edges = append(filtered.Edges, edge)
		}
	}

	// Update statistics
	filtered.Stats.TotalNodes = len(filtered.Nodes)
	filtered.Stats.TotalEdges = len(filtered.Edges)

	totalCoupling := 0
	for _, node := range filtered.Nodes {
		coupling := node.Afferent + node.Efferent
		totalCoupling += coupling
		if coupling > filtered.Stats.MaxCoupling {
			filtered.Stats.MaxCoupling = coupling
		}
	}
	if filtered.Stats.TotalNodes > 0 {
		filtered.Stats.AvgCoupling = float64(totalCoupling) / float64(filtered.Stats.TotalNodes)
	}

	return filtered
}

// ExtractLayer creates a subgraph containing only nodes from the specified layer.
func (g *Graph) ExtractLayer(layerID int) *Graph {
	if layerID < 0 || layerID > 3 {
		return g
	}

	subgraph := &Graph{
		Name:   g.Name + "_layer" + string(rune('0'+layerID)),
		Nodes:  make(map[string]*Node),
		Edges:  make([]*Edge, 0),
		Layers: make([]*Layer, 4),
		Stats:  &GraphStats{LayerCounts: make([]int, 4)},
		IsDAG:  true,
	}

	// Initialize layers
	for i := 0; i < 4; i++ {
		subgraph.Layers[i] = &Layer{
			ID:          g.Layers[i].ID,
			Name:        g.Layers[i].Name,
			Description: g.Layers[i].Description,
			NodeIDs:     make([]string, 0),
		}
	}

	// Add nodes from specified layer
	for _, nodeID := range g.Layers[layerID].NodeIDs {
		if node, ok := g.Nodes[nodeID]; ok {
			subgraph.Nodes[nodeID] = node
			subgraph.Layers[layerID].NodeIDs = append(subgraph.Layers[layerID].NodeIDs, nodeID)
		}
	}

	// Add edges within layer and cross-layer edges
	for _, edge := range g.Edges {
		srcNode, srcOK := g.Nodes[edge.Source]
		tgtNode, tgtOK := g.Nodes[edge.Target]

		if !srcOK || !tgtOK {
			continue
		}

		// Include if source is in layer
		if srcNode.Layer == layerID {
			// Add target node if it's a cross-layer dependency
			if tgtNode.Layer != layerID {
				subgraph.Nodes[edge.Target] = tgtNode
				subgraph.Layers[tgtNode.Layer].NodeIDs = append(subgraph.Layers[tgtNode.Layer].NodeIDs, edge.Target)
			}
			subgraph.Edges = append(subgraph.Edges, edge)
		}
	}

	// Update statistics
	subgraph.Stats.TotalNodes = len(subgraph.Nodes)
	subgraph.Stats.TotalEdges = len(subgraph.Edges)

	return subgraph
}

// ToAdjacency converts the graph to an adjacency list representation.
func (g *Graph) ToAdjacency() map[string]map[string]map[string]interface{} {
	adj := make(map[string]map[string]map[string]interface{})

	for _, edge := range g.Edges {
		if adj[edge.Source] == nil {
			adj[edge.Source] = make(map[string]map[string]interface{})
		}
		adj[edge.Source][edge.Target] = map[string]interface{}{
			"type":   edge.Type,
			"weight": edge.Weight,
		}
	}

	return adj
}

// ToJSON serializes the graph to JSON.
func (g *Graph) ToJSON() ([]byte, error) {
	return json.MarshalIndent(g, "", "  ")
}
