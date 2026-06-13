// Package catreview_test guards the C01 fix: the module path declared in
// go.mod must match the repository's VCS path (github.com/manutej/catreview-go),
// and every internal import must use that path. A mismatch breaks
// `go install github.com/manutej/catreview-go/...@version` and any downstream
// consumer with:
//
//	module declares its path as: <declared>
//	        but was required as: github.com/manutej/catreview-go
package catreview_test

import (
	"bufio"
	"bytes"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const wantModulePath = "github.com/manutej/catreview-go"

// staleModulePrefix is the pre-fix module path; no import may reference it.
const staleModulePrefix = "github.com/manu/catreview"

func declaredModulePath(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile("go.mod")
	if err != nil {
		t.Fatalf("reading go.mod: %v", err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if rest, ok := strings.CutPrefix(line, "module "); ok {
			return strings.Trim(strings.TrimSpace(rest), `"`)
		}
	}
	t.Fatal("go.mod: no module declaration found")
	return ""
}

func TestModulePathMatchesRepository(t *testing.T) {
	if got := declaredModulePath(t); got != wantModulePath {
		t.Fatalf("go.mod declares module %q, want %q (mismatch breaks `go install %s/...`)",
			got, wantModulePath, wantModulePath)
	}
}

func TestInternalImportsUseDeclaredModulePath(t *testing.T) {
	modPath := declaredModulePath(t)
	fset := token.NewFileSet()

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if name != "." && (strings.HasPrefix(name, ".") || name == "vendor" || name == "testdata") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			t.Errorf("%s: parse error: %v", path, err)
			return nil
		}
		for _, imp := range file.Imports {
			val := strings.Trim(imp.Path.Value, `"`)
			if val == staleModulePrefix || strings.HasPrefix(val, staleModulePrefix+"/") {
				t.Errorf("%s imports stale pre-fix module path %q; use %q", path, val, modPath)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking repository: %v", err)
	}
}
