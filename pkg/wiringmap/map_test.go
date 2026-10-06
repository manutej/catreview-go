package wiringmap

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/manu/catreview/pkg/extractor"
)

func testPkgDir(t *testing.T) string {
	root, err := ModuleRoot(".")
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(root, "pkg")
}

func TestFromCategoryDeterministic(t *testing.T) {
	pkgDir := testPkgDir(t)
	ext := extractor.NewGoExtractor()
	cat, err := ext.ExtractFromPath(pkgDir)
	if err != nil {
		t.Fatal(err)
	}
	modRoot, err := ModuleRoot(".")
	if err != nil {
		t.Fatal(err)
	}
	opts := BuildOptions{Repo: "catreview-go", ModuleRoot: modRoot, WalkRoot: pkgDir}
	doc1, err := FromCategory(cat, opts)
	if err != nil {
		t.Fatal(err)
	}
	doc2, err := FromCategory(cat, opts)
	if err != nil {
		t.Fatal(err)
	}
	b1, _ := MarshalJSON(doc1, false)
	b2, _ := MarshalJSON(doc2, false)
	if !bytes.Equal(b1, b2) {
		t.Fatal("non-deterministic output")
	}
}

func TestFromCategoryAbsolutePathSameAsRelative(t *testing.T) {
	modRoot, err := ModuleRoot(".")
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(modRoot)
	pkgDir := filepath.Join(modRoot, "pkg")
	relWalk := "pkg"
	extRel := extractor.NewGoExtractor()
	catRel, err := extRel.ExtractFromPath(relWalk)
	if err != nil {
		t.Fatal(err)
	}
	extAbs := extractor.NewGoExtractor()
	catAbs, err := extAbs.ExtractFromPath(pkgDir)
	if err != nil {
		t.Fatal(err)
	}
	optsRel := BuildOptions{Repo: "catreview-go", ModuleRoot: modRoot, WalkRoot: relWalk}
	optsAbs := BuildOptions{Repo: "catreview-go", ModuleRoot: modRoot, WalkRoot: pkgDir}
	docRel, err := FromCategory(catRel, optsRel)
	if err != nil {
		t.Fatal(err)
	}
	docAbs, err := FromCategory(catAbs, optsAbs)
	if err != nil {
		t.Fatal(err)
	}
	bRel, _ := MarshalJSON(docRel, false)
	bAbs, _ := MarshalJSON(docAbs, false)
	if !bytes.Equal(bRel, bAbs) {
		t.Fatal("relative vs absolute path output differs")
	}
}

func TestEdgeDedupe(t *testing.T) {
	ext := extractor.NewGoExtractor()
	pkgDir := testPkgDir(t)
	cat, err := ext.ExtractFromPath(pkgDir)
	if err != nil {
		t.Fatal(err)
	}
	modRoot, _ := ModuleRoot(".")
	doc, err := FromCategory(cat, BuildOptions{Repo: "catreview-go", ModuleRoot: modRoot, WalkRoot: pkgDir})
	if err != nil {
		t.Fatal(err)
	}
	seen := make(map[string]struct{})
	for _, e := range doc.Edges {
		key := e.From + "|" + e.To + "|" + e.EdgeKind
		if _, ok := seen[key]; ok {
			t.Fatalf("duplicate edge %s", key)
		}
		seen[key] = struct{}{}
	}
}
