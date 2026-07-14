// Package extractor: Java support.
//
// java_extractor.go implements a dependency-free, structural extractor for Java
// source code. Unlike the Go extractor (which uses the standard library's
// go/ast), Go has no built-in Java parser, so this file provides a self-contained
// scanner that recognises Java's declaration-level structure:
//
//   - package declaration       -> categorical "package" of every object
//   - import statements         -> import morphisms (file -> imported_package)
//   - class/interface/enum/      -> type Objects (kind: class|interface|enum|
//     record/@interface             record|annotation), including nested types
//   - extends / implements      -> inheritance morphisms (type -> supertype)
//   - methods & constructors    -> function Objects
//   - method calls (same type)  -> function_call morphisms (best-effort)
//
// The extractor is intentionally structural rather than a full Java grammar: it
// resolves the constructs that matter for categorical architecture analysis
// (objects and their dependency morphisms) while remaining robust to code it
// does not fully understand. It produces exactly the same category.Category
// shape as GoExtractor, so all downstream analysis (complexity, functors,
// cycles, visualization) works unchanged.
package extractor

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/manu/catreview/pkg/category"
)

// JavaExtractor extracts categorical models from Java source code.
type JavaExtractor struct {
	category *category.Category
	fileSet  map[string]bool // Set of parsed file paths (for Stats)
	packages map[string]bool // Set of distinct package names (for Stats)

	// typeByName indexes defined type Objects by their simple name so that
	// extends/implements targets and same-name references can be resolved to a
	// concrete Object id when possible.
	typeByName map[string]string
}

// NewJavaExtractor creates a new Java code extractor.
func NewJavaExtractor() *JavaExtractor {
	return &JavaExtractor{
		category:   category.NewCategory("java_codebase"),
		fileSet:    make(map[string]bool),
		packages:   make(map[string]bool),
		typeByName: make(map[string]string),
	}
}

// ExtractFromPath extracts a categorical model from a Java project path.
// Implements the Extractor interface.
func (e *JavaExtractor) ExtractFromPath(root string) (*category.Category, error) {
	// Reset per-run state so a factory-held instance can be reused safely.
	e.category = category.NewCategory("java_codebase")
	e.fileSet = make(map[string]bool)
	e.packages = make(map[string]bool)
	e.typeByName = make(map[string]string)

	// Two passes: first collect declared types (so inheritance and call
	// morphisms can be resolved to concrete Objects), then extract relationships.
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Skip common build/output directories to avoid noise.
			base := info.Name()
			if base == "target" || base == "build" || base == "out" || base == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".java") {
			return nil
		}
		// Skip generated package-info / module-info marker files: they carry no
		// type declarations and only add noise.
		base := filepath.Base(path)
		if base == "module-info.java" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Parse each file into an intermediate representation.
	parsed := make([]*javaFile, 0, len(files))
	for _, path := range files {
		jf, perr := parseJavaFile(path)
		if perr != nil {
			// Best-effort: skip files we cannot read, but keep going.
			continue
		}
		parsed = append(parsed, jf)
	}

	// Pass 1: register every declared type Object so later morphisms resolve.
	for _, jf := range parsed {
		e.fileSet[jf.path] = true
		e.packages[jf.pkg] = true
		for _, t := range jf.types {
			e.addTypeObject(jf, t)
		}
	}

	// Pass 2: files, imports, methods, inheritance and call morphisms.
	for _, jf := range parsed {
		e.addFileAndRelations(jf)
	}

	return e.category, nil
}

// addTypeObject registers a type declaration as a categorical Object.
func (e *JavaExtractor) addTypeObject(jf *javaFile, t javaType) {
	obj := category.NewObject(
		t.id,
		t.kind,
		t.name,
		map[string]interface{}{
			"package":  jf.pkg,
			"file":     jf.path,
			"doc":      t.doc,
			"exported": t.isPublic,
		},
	)
	if err := e.category.AddObject(obj); err != nil {
		// Duplicate id (e.g. same nested name in two files) — keep the first.
		return
	}
	// Index by simple name for later resolution. First writer wins.
	if _, exists := e.typeByName[t.name]; !exists {
		e.typeByName[t.name] = t.id
	}
}

// addFileAndRelations creates the file Object and all morphisms rooted at it.
func (e *JavaExtractor) addFileAndRelations(jf *javaFile) {
	fileObj := category.NewObject(
		jf.path,
		"file",
		filepath.Base(jf.path),
		map[string]interface{}{
			"package": jf.pkg,
			"path":    jf.path,
			"imports": len(jf.imports),
		},
	)
	if err := e.category.AddObject(fileObj); err != nil {
		return
	}

	// Imports -> imported_package Objects + import morphisms.
	for _, imp := range jf.imports {
		e.addImport(jf.path, imp)
	}

	// Types: defines morphism (file -> type) + inheritance morphisms.
	for _, t := range jf.types {
		morphID := fmt.Sprintf("defines:%s->%s", jf.path, t.id)
		defMorph := category.NewMorphism(morphID, jf.path, t.id, "defines",
			map[string]interface{}{"kind": t.kind})
		e.category.AddMorphism(defMorph)

		for _, super := range t.supertypes {
			e.addInheritance(t.id, super)
		}
	}

	// Methods: function Objects + defines morphisms, then call morphisms.
	for _, m := range jf.methods {
		e.addMethod(jf, m)
	}
	for _, m := range jf.methods {
		e.addCalls(jf, m)
	}
}

// addImport creates the imported_package Object (once) and an import morphism.
func (e *JavaExtractor) addImport(sourceFile, importPath string) {
	targetID := fmt.Sprintf("import:%s", importPath)
	if _, exists := e.category.GetObject(targetID); !exists {
		impObj := category.NewObject(targetID, "imported_package", importPath,
			map[string]interface{}{"import_path": importPath})
		if err := e.category.AddObject(impObj); err != nil {
			return
		}
	}
	morphID := fmt.Sprintf("import:%s->%s", sourceFile, importPath)
	if _, exists := e.category.GetMorphism(morphID); !exists {
		morph := category.NewMorphism(morphID, sourceFile, targetID, "import",
			map[string]interface{}{"import_path": importPath})
		e.category.AddMorphism(morph)
	}
}

// addInheritance links a type to a supertype. The supertype is resolved to a
// declared Object when its simple name is known; otherwise an external_type
// placeholder Object is created (mirroring how imports become Objects), so the
// resulting category stays valid (every morphism has real endpoints).
func (e *JavaExtractor) addInheritance(typeID, superName string) {
	targetID, resolved := e.typeByName[superName]
	if !resolved {
		targetID = fmt.Sprintf("external_type:%s", superName)
		if _, exists := e.category.GetObject(targetID); !exists {
			ext := category.NewObject(targetID, "external_type", superName,
				map[string]interface{}{"resolved": false})
			if err := e.category.AddObject(ext); err != nil {
				return
			}
		}
	}
	morphID := fmt.Sprintf("inherits:%s->%s", typeID, targetID)
	if _, exists := e.category.GetMorphism(morphID); !exists {
		morph := category.NewMorphism(morphID, typeID, targetID, "inheritance",
			map[string]interface{}{"supertype": superName, "resolved": resolved})
		e.category.AddMorphism(morph)
	}
}

// addMethod creates a function Object for a method/constructor and the
// defines morphism from its file.
func (e *JavaExtractor) addMethod(jf *javaFile, m javaMethod) {
	kind := "method"
	if m.isConstructor {
		kind = "constructor"
	}
	obj := category.NewObject(
		m.id,
		"function",
		m.name,
		map[string]interface{}{
			"package":  jf.pkg,
			"file":     jf.path,
			"class":    m.typeID,
			"exported": m.isPublic,
			"kind":     kind,
			"doc":      m.doc,
		},
	)
	if err := e.category.AddObject(obj); err != nil {
		return
	}
	morphID := fmt.Sprintf("defines:%s->%s", jf.path, m.id)
	morph := category.NewMorphism(morphID, jf.path, m.id, "defines",
		map[string]interface{}{"kind": "function"})
	e.category.AddMorphism(morph)
}

// addCalls adds best-effort function_call morphisms for calls resolvable to
// another method on the same type (including this./super.). Unresolvable calls
// (to imported or JDK methods) are intentionally dropped, matching the Go
// extractor's behaviour of only recording edges with real endpoints.
func (e *JavaExtractor) addCalls(jf *javaFile, m javaMethod) {
	if m.body == "" {
		return
	}
	// Index sibling methods (same enclosing type) by simple name.
	siblings := make(map[string]string)
	for _, other := range jf.methods {
		if other.typeID == m.typeID {
			siblings[other.name] = other.id
		}
	}
	seen := make(map[string]bool)
	for _, call := range callRe.FindAllStringSubmatch(m.body, -1) {
		name := call[2]
		if javaKeywords[name] {
			continue
		}
		targetID, ok := siblings[name]
		if !ok || targetID == m.id {
			continue
		}
		if seen[targetID] {
			continue
		}
		seen[targetID] = true
		morphID := fmt.Sprintf("calls:%s->%s", m.id, targetID)
		if _, exists := e.category.GetMorphism(morphID); !exists {
			morph := category.NewMorphism(morphID, m.id, targetID, "function_call",
				map[string]interface{}{"target": name})
			e.category.AddMorphism(morph)
		}
	}
}

// Stats returns extraction statistics (parallels GoExtractor.Stats).
func (e *JavaExtractor) Stats() map[string]interface{} {
	stats := e.category.Stats()
	return map[string]interface{}{
		"objects":   stats["objects"],
		"morphisms": stats["morphisms"],
		"files":     len(e.fileSet),
		"packages":  len(e.packages),
	}
}

// Language returns the source language name. Implements the Extractor interface.
func (e *JavaExtractor) Language() string { return "java" }

// FileExtensions returns the handled file extensions. Implements Extractor.
func (e *JavaExtractor) FileExtensions() []string { return []string{".java"} }

// ---------------------------------------------------------------------------
// Parsing internals
// ---------------------------------------------------------------------------

// javaFile is the intermediate representation of a parsed Java source file.
type javaFile struct {
	path    string
	pkg     string
	imports []string
	types   []javaType
	methods []javaMethod
}

// javaType is a declared class/interface/enum/record/annotation.
type javaType struct {
	id         string // fully-qualified id: pkg.Outer.Inner
	name       string // simple name
	kind       string // class|interface|enum|record|annotation
	supertypes []string
	isPublic   bool
	doc        string
}

// javaMethod is a declared method or constructor.
type javaMethod struct {
	id            string // pkg.Type.name
	name          string
	typeID        string // enclosing type id
	isPublic      bool
	isConstructor bool
	doc           string
	bodyStart     int    // start offset of the method body in cleaned source
	body          string // stripped body text (for call extraction)
}

var (
	packageRe = regexp.MustCompile(`(?m)^\s*package\s+([\w.]+)\s*;`)
	importRe  = regexp.MustCompile(`(?m)^\s*import\s+(?:static\s+)?([\w.*]+)\s*;`)

	// annotationTypeRe matches an annotation type declaration: @interface Name
	annotationTypeRe = regexp.MustCompile(`@\s*interface\s+(\w+)`)
	// typeHeaderRe matches a class/interface/enum/record declaration name.
	typeHeaderRe = regexp.MustCompile(`\b(class|interface|enum|record)\s+(\w+)`)
	// methodHeaderRe matches a method/constructor header ending in a parameter
	// list (optionally followed by a throws clause). Group 1 is the name.
	methodHeaderRe = regexp.MustCompile(`(\w+)\s*\([^{;]*\)\s*(?:throws[\w\s,.]+)?$`)

	extendsRe    = regexp.MustCompile(`\bextends\s+([\w.]+(?:\s*,\s*[\w.]+)*)`)
	implementsRe = regexp.MustCompile(`\bimplements\s+([\w.]+(?:\s*,\s*[\w.]+)*)`)
	genericsRe   = regexp.MustCompile(`<[^<>]*>`)

	callRe = regexp.MustCompile(`(?:\b(this|super)\s*\.\s*)?(\w+)\s*\(`)
)

// javaKeywords are identifiers that look like calls/methods but are language
// keywords or control structures; they must never become methods or callees.
var javaKeywords = map[string]bool{
	"if": true, "for": true, "while": true, "switch": true, "catch": true,
	"synchronized": true, "return": true, "new": true, "else": true, "do": true,
	"try": true, "throw": true, "assert": true, "super": true, "this": true,
	"case": true, "instanceof": true, "yield": true,
}

// parseJavaFile reads and structurally parses a single Java file.
func parseJavaFile(path string) (*javaFile, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	src := string(raw)
	clean, javadocs := cleanSource(src)

	jf := &javaFile{path: path, pkg: "default"}

	if m := packageRe.FindStringSubmatch(clean); m != nil {
		jf.pkg = m[1]
	}
	for _, m := range importRe.FindAllStringSubmatch(clean, -1) {
		jf.imports = append(jf.imports, m[1])
	}

	scanDeclarations(clean, javadocs, jf)
	return jf, nil
}

// javadoc records the position and text of a /** ... */ documentation comment
// in the cleaned source coordinate space.
type javadoc struct {
	start int
	end   int
	text  string
}

// cleanSource replaces comment and string/char-literal content with spaces
// (preserving newlines and overall length so positions stay aligned), and
// returns any Javadoc (/** ... */) comments encountered. Neutralising literals
// and comments makes brace matching and keyword detection reliable.
func cleanSource(src string) (string, []javadoc) {
	out := []byte(src)
	var docs []javadoc
	n := len(src)
	blank := func(i int) {
		if src[i] != '\n' && src[i] != '\r' {
			out[i] = ' '
		}
	}
	for i := 0; i < n; {
		c := src[i]
		switch {
		case c == '/' && i+1 < n && src[i+1] == '/':
			out[i], out[i+1] = ' ', ' '
			i += 2
			for i < n && src[i] != '\n' {
				blank(i)
				i++
			}
		case c == '/' && i+1 < n && src[i+1] == '*':
			isDoc := i+2 < n && src[i+2] == '*'
			start := i
			out[i], out[i+1] = ' ', ' '
			i += 2
			for i < n && !(src[i] == '*' && i+1 < n && src[i+1] == '/') {
				blank(i)
				i++
			}
			if i < n {
				out[i] = ' '
				if i+1 < n {
					out[i+1] = ' '
				}
				i += 2
			}
			if isDoc {
				docs = append(docs, javadoc{start: start, end: i, text: cleanJavadocText(src[start:i])})
			}
		case c == '"':
			// Handle text blocks (""" ... """) and normal strings.
			if i+2 < n && src[i+1] == '"' && src[i+2] == '"' {
				i += 3
				for i+2 < n && !(src[i] == '"' && src[i+1] == '"' && src[i+2] == '"') {
					blank(i)
					i++
				}
				i += 3
			} else {
				out[i] = ' '
				i++
				for i < n && src[i] != '"' {
					if src[i] == '\\' && i+1 < n {
						blank(i)
						i++
					}
					blank(i)
					i++
				}
				if i < n {
					out[i] = ' '
					i++
				}
			}
		case c == '\'':
			out[i] = ' '
			i++
			for i < n && src[i] != '\'' {
				if src[i] == '\\' && i+1 < n {
					blank(i)
					i++
				}
				blank(i)
				i++
			}
			if i < n {
				out[i] = ' '
				i++
			}
		default:
			i++
		}
	}
	return string(out), docs
}

// cleanJavadocText reduces a raw javadoc block to a short first-line summary.
func cleanJavadocText(raw string) string {
	raw = strings.TrimPrefix(raw, "/**")
	raw = strings.TrimSuffix(raw, "*/")
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "@") {
			return line
		}
	}
	return ""
}

// scanDeclarations walks the cleaned source using brace matching to recover the
// nesting of types and methods. Every '{' is classified by the header text
// since the last statement break (';', '{' or '}'):
//   - a type header -> push a type scope (records a javaType)
//   - a method header inside a type scope -> push a method scope (records a
//     javaMethod, capturing its body span)
//   - anything else -> a plain block scope
func scanDeclarations(clean string, docs []javadoc, jf *javaFile) {
	type frame struct {
		kind      string // "root" | "type" | "method" | "block"
		typeID    string // enclosing type id (for methods/nested types)
		typeName  string // enclosing simple type name (constructor detection)
		methodIdx int    // index into jf.methods when kind=="method"
	}

	stack := []frame{{kind: "root"}}
	typePath := []string{} // simple names of open types, for id construction

	lastBreak := 0
	n := len(clean)

	findDoc := func(from, to int) string {
		for _, d := range docs {
			if d.start >= from && d.end <= to {
				return d.text
			}
		}
		return ""
	}

	for i := 0; i < n; i++ {
		switch clean[i] {
		case ';':
			lastBreak = i + 1
		case '}':
			top := stack[len(stack)-1]
			if top.kind == "method" {
				jf.methods[top.methodIdx].body = clean[jf.methods[top.methodIdx].bodyStart:i]
			}
			if top.kind == "type" && len(typePath) > 0 {
				typePath = typePath[:len(typePath)-1]
			}
			if len(stack) > 1 {
				stack = stack[:len(stack)-1]
			}
			lastBreak = i + 1
		case '{':
			header := strings.TrimSpace(clean[lastBreak:i])
			parent := stack[len(stack)-1]
			doc := findDoc(lastBreak, i)

			if kind, name, ok := matchTypeHeader(header); ok {
				fullName := name
				if len(typePath) > 0 {
					fullName = strings.Join(typePath, ".") + "." + name
				}
				typeID := jf.pkg + "." + fullName
				jf.types = append(jf.types, javaType{
					id:         typeID,
					name:       name,
					kind:       kind,
					supertypes: parseSupertypes(header),
					isPublic:   strings.Contains(header, "public"),
					doc:        doc,
				})
				typePath = append(typePath, name)
				stack = append(stack, frame{kind: "type", typeID: typeID, typeName: name})
			} else if parent.kind == "type" {
				if name, ok := matchMethodHeader(header); ok {
					methodID := parent.typeID + "." + name
					// Disambiguate overloads/duplicates by suffixing an index.
					methodID = uniqueMethodID(jf, methodID)
					jf.methods = append(jf.methods, javaMethod{
						id:            methodID,
						name:          name,
						typeID:        parent.typeID,
						isPublic:      strings.Contains(header, "public"),
						isConstructor: name == parent.typeName,
						doc:           doc,
						bodyStart:     i + 1,
					})
					stack = append(stack, frame{kind: "method", typeID: parent.typeID, methodIdx: len(jf.methods) - 1})
				} else {
					stack = append(stack, frame{kind: "block", typeID: parent.typeID})
				}
			} else {
				stack = append(stack, frame{kind: "block", typeID: parent.typeID})
			}
			lastBreak = i + 1
		}
	}
}

// matchTypeHeader determines whether a header declares a type, returning its
// categorical kind and simple name.
func matchTypeHeader(header string) (kind, name string, ok bool) {
	if m := annotationTypeRe.FindStringSubmatch(header); m != nil {
		return "annotation", m[1], true
	}
	if m := typeHeaderRe.FindStringSubmatch(header); m != nil {
		return m[1], m[2], true
	}
	return "", "", false
}

// matchMethodHeader determines whether a header declares a method/constructor,
// returning the method name. Control-flow constructs (if/for/…) are rejected via
// the keyword set.
func matchMethodHeader(header string) (name string, ok bool) {
	m := methodHeaderRe.FindStringSubmatch(header)
	if m == nil {
		return "", false
	}
	name = m[1]
	if javaKeywords[name] {
		return "", false
	}
	return name, true
}

// parseSupertypes returns the union of extends and implements targets, with
// generic parameters stripped and names reduced to their simple form.
func parseSupertypes(header string) []string {
	var out []string
	collect := func(re *regexp.Regexp) {
		if m := re.FindStringSubmatch(stripGenerics(header)); m != nil {
			for _, part := range strings.Split(m[1], ",") {
				part = strings.TrimSpace(part)
				if idx := strings.LastIndex(part, "."); idx >= 0 {
					part = part[idx+1:]
				}
				if part != "" {
					out = append(out, part)
				}
			}
		}
	}
	collect(extendsRe)
	collect(implementsRe)
	return out
}

// stripGenerics removes generic type parameters so that supertype lists parse
// cleanly (e.g. "implements Comparable<Foo>, Serializable").
func stripGenerics(s string) string {
	for genericsRe.MatchString(s) {
		s = genericsRe.ReplaceAllString(s, "")
	}
	return s
}

// uniqueMethodID ensures overloaded methods (same name) get distinct Object ids.
func uniqueMethodID(jf *javaFile, base string) string {
	id := base
	suffix := 1
	for methodIDExists(jf, id) {
		id = fmt.Sprintf("%s#%d", base, suffix)
		suffix++
	}
	return id
}

func methodIDExists(jf *javaFile, id string) bool {
	for _, m := range jf.methods {
		if m.id == id {
			return true
		}
	}
	return false
}
