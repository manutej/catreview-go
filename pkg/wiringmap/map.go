package wiringmap

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/manu/catreview/pkg/category"
)

// BuildOptions configures wiringmap generation from an extracted category.
type BuildOptions struct {
	Repo       string
	ModuleRoot string
	WalkRoot   string
}

// FromCategory maps a categorical extract model to wiringmap.v0.
func FromCategory(cat *category.Category, opts BuildOptions) (*Document, error) {
	moduleRoot, err := filepath.Abs(opts.ModuleRoot)
	if err != nil {
		return nil, err
	}
	walkRoot, err := filepath.Abs(opts.WalkRoot)
	if err != nil {
		return nil, err
	}
	ref, err := RelToModule(moduleRoot, walkRoot)
	if err != nil {
		return nil, err
	}

	relFile := make(map[string]string) // file object id -> relative path
	sanitizeKeys := make(map[string]string)

	for _, obj := range cat.Objects() {
		if obj.Type != "file" {
			continue
		}
		rel, err := normalizeFilePath(moduleRoot, obj.ID, obj.Metadata)
		if err != nil {
			return nil, err
		}
		relFile[obj.ID] = rel
		sanitizeKeys[rel] = SanitizeID(rel)
	}
	if err := CheckSanitizeCollisions(sanitizeKeys); err != nil {
		return nil, err
	}

	unitsByPath := make(map[string]*Unit)
	unitOrder := make([]string, 0)

	for fileID, rel := range relFile {
		obj, _ := cat.GetObject(fileID)
		pkg := ""
		if obj != nil && obj.Metadata != nil {
			if v, ok := obj.Metadata["package"].(string); ok {
				pkg = v
			}
		}
		u := &Unit{
			ID:    "unit:" + SanitizeID(rel),
			Path:  rel,
			Role:  "package " + pkg,
			Ports: nil,
		}
		unitsByPath[rel] = u
		unitOrder = append(unitOrder, rel)
	}
	sort.Strings(unitOrder)

	objectToPort := make(map[string]string)

	for _, obj := range cat.Objects() {
		switch obj.Type {
		case "function", "struct", "interface", "type":
		default:
			continue
		}
		if strings.Contains(obj.ID, "#") {
			return nil, fmt.Errorf("object id contains '#': %q", obj.ID)
		}
		fileRaw, ok := metaString(obj.Metadata, "file")
		if !ok {
			continue
		}
		rel, err := normalizeFilePath(moduleRoot, fileRaw, nil)
		if err != nil {
			return nil, err
		}
		u, ok := unitsByPath[rel]
		if !ok {
			continue
		}
		portID := "port:" + SanitizeID(rel) + "#" + obj.ID
		kind := portKind(obj.Type)
		exported := exportedFlag(obj)
		u.Ports = append(u.Ports, Port{
			ID:       portID,
			Symbol:   obj.Name,
			Kind:     kind,
			Exported: exported,
		})
		objectToPort[obj.ID] = portID
	}

	for _, rel := range unitOrder {
		u := unitsByPath[rel]
		sort.Slice(u.Ports, func(i, j int) bool {
			return u.Ports[i].ID < u.Ports[j].ID
		})
	}

	junctionByImport := make(map[string]Junction)
	for _, obj := range cat.Objects() {
		if obj.Type != "imported_package" {
			continue
		}
		importPath := obj.Name
		if v, ok := metaString(obj.Metadata, "import_path"); ok {
			importPath = v
		}
		jid := "junction:" + SanitizeID(importPath)
		junctionByImport[obj.ID] = Junction{
			ID:           jid,
			Label:        importPath,
			JunctionKind: "other",
		}
	}
	junctionList := make([]Junction, 0, len(junctionByImport))
	for _, j := range junctionByImport {
		junctionList = append(junctionList, j)
	}
	sort.Slice(junctionList, func(i, j int) bool {
		return junctionList[i].ID < junctionList[j].ID
	})

	importPathToJunction := make(map[string]string)
	for impObjID, j := range junctionByImport {
		importPathToJunction[impObjID] = j.ID
	}

	type edgeKey struct {
		from, to, kind string
	}
	edgeSet := make(map[edgeKey]struct{})
	dropped := 0

	for _, m := range cat.Morphisms() {
		switch m.Type {
		case "function_call":
			fromPort, okFrom := objectToPort[m.Source]
			toPort, okTo := objectToPort[m.Target]
			if !okFrom || !okTo {
				dropped++
				continue
			}
			edgeSet[edgeKey{fromPort, toPort, "call"}] = struct{}{}
		case "type_dependency":
			fromPort, okFrom := objectToPort[m.Source]
			toPort, okTo := objectToPort[m.Target]
			if !okFrom || !okTo {
				dropped++
				continue
			}
			edgeSet[edgeKey{fromPort, toPort, "ref"}] = struct{}{}
		case "import":
			rel, ok := relFile[m.Source]
			if !ok {
				dropped++
				continue
			}
			jid, ok := importPathToJunction[m.Target]
			if !ok {
				dropped++
				continue
			}
			fromUnit := "unit:" + SanitizeID(rel)
			edgeSet[edgeKey{fromUnit, jid, "import"}] = struct{}{}
		case "identity", "defines":
			continue
		default:
			dropped++
		}
	}

	edgeKeys := make([]edgeKey, 0, len(edgeSet))
	for k := range edgeSet {
		edgeKeys = append(edgeKeys, k)
	}
	sort.Slice(edgeKeys, func(i, j int) bool {
		a, b := edgeKeys[i], edgeKeys[j]
		if a.kind != b.kind {
			return a.kind < b.kind
		}
		if a.from != b.from {
			return a.from < b.from
		}
		return a.to < b.to
	})

	edges := make([]Edge, 0, len(edgeKeys))
	for i, k := range edgeKeys {
		edges = append(edges, Edge{
			ID:       fmt.Sprintf("edge:%s-%d", k.kind, i+1),
			From:     k.from,
			To:       k.to,
			EdgeKind: k.kind,
		})
	}

	units := make([]Unit, 0, len(unitOrder))
	for _, rel := range unitOrder {
		units = append(units, *unitsByPath[rel])
	}

	fileCount := len(units)
	sliceLabel := fmt.Sprintf("catreview wiringmap %s (%d files)", ref, fileCount)
	notes := fmt.Sprintf("%d external morphisms dropped", dropped)

	doc := &Document{
		SchemaVersion: "wiringmap.v0",
		Meta: Meta{
			Repo:  opts.Repo,
			Ref:   ref,
			Slice: sliceLabel,
			Notes: notes,
		},
		Junctions: junctionList,
		Units:     units,
		Edges:     edges,
	}
	return doc, nil
}

func normalizeFilePath(moduleRoot, fileID string, metadata map[string]interface{}) (string, error) {
	p := fileID
	if metadata != nil {
		if v, ok := metaString(metadata, "path"); ok && v != "" {
			p = v
		}
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(moduleRoot, abs)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(rel), nil
}

func metaString(md map[string]interface{}, key string) (string, bool) {
	if md == nil {
		return "", false
	}
	v, ok := md[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok && s != ""
}

func portKind(objType string) string {
	switch objType {
	case "function":
		return "method"
	default:
		return "type"
	}
}

func exportedFlag(obj *category.Object) bool {
	if obj.Metadata != nil {
		if v, ok := obj.Metadata["is_exported"].(bool); ok {
			return v
		}
	}
	if obj.Name == "" {
		return false
	}
	return unicode.IsUpper(rune(obj.Name[0]))
}
