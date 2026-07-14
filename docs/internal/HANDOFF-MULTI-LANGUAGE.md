# Handoff — Multi-Language Support (Java Extractor)

**Branch**: `claude/multi-language-java-support-hwat62`
**Base**: `master` @ `ea7e11f`
**Status**: ✅ Java extractor implemented, wired, tested, and verified end-to-end
**Next**: Parallel-agent deep review → then secondary language set (Python, TypeScript, Rust)

---

## 1. What this session delivered

| # | Deliverable | Location |
|---|-------------|----------|
| 1 | **Full Java extractor** — dependency-free structural parser | `pkg/extractor/java_extractor.go` (~680 LOC) |
| 2 | **Java extractor tests** — 4 suites, all passing | `pkg/extractor/java_extractor_test.go` |
| 3 | **Factory + language detection wired into CLI** | `pkg/extractor/extractor.go`, `cmd/catreview/main.go` |
| 4 | **Dev/analysis workflow doc updated** (RMP/Ralph loop) | `docs/CATEGORICAL-REPO-ANALYSIS.md` |
| 5 | **README + INDEX reconciled** with reality | `README.md`, `docs/INDEX.md` |
| 6 | **This handoff** | `docs/internal/HANDOFF-MULTI-LANGUAGE.md` |

Both extractors now reset per-run state at the start of `ExtractFromPath`, so a
factory-held instance is safe to reuse (see `go_extractor.go` / `java_extractor.go`).

---

## 2. How the Java extractor works

Go has no standard-library Java parser, so rather than take a CGO/`javaparser`
dependency, the Java extractor is a **self-contained structural scanner**:

1. **`cleanSource`** — a character state machine that neutralises comments and
   string/char/text-block literals (replacing their contents with spaces while
   preserving newlines and total length). This makes brace matching and keyword
   detection reliable and captures Javadoc summaries.
2. **`scanDeclarations`** — walks the cleaned source with a brace-depth stack.
   Every `{` is classified by the header text since the last `;`/`{`/`}`:
   - type header (`class`/`interface`/`enum`/`record`/`@interface`) → type scope
   - method header (ends in a parameter list) **inside a type scope** → method scope
   - anything else → a plain block scope
   Because methods are only recognised inside a *type* scope, control-flow blocks
   (`if`/`for`/`while`/…) inside method bodies are never mistaken for methods.
3. **Category assembly** — two passes: pass 1 registers all declared type Objects
   (so inheritance/call targets can resolve); pass 2 adds files, imports, methods,
   inheritance, and same-type call morphisms.

### Categorical mapping (identical shape to the Go extractor)

| Java construct | Object type | Notes |
|----------------|-------------|-------|
| `.java` file | `file` | id = file path |
| `class` / `interface` / `enum` / `record` / `@interface` | `class`/`interface`/`enum`/`record`/`annotation` | id = `pkg.Outer.Inner` |
| method / constructor | `function` | id = `pkg.Type.name` (overloads suffixed `#n`) |
| `import x.y.Z` | `imported_package` | id = `import:x.y.Z` |
| unresolved supertype | `external_type` | id = `external_type:Name` |

| Java dependency | Morphism type |
|-----------------|---------------|
| file defines type/method | `defines` |
| `import` | `import` |
| `extends` / `implements` | `inheritance` (resolved to a declared type when the simple name is known, else an `external_type` placeholder) |
| same-type method call | `function_call` (best-effort) |

All downstream analysis (complexity, functors, cycles, viz) is unchanged.

### Known limitations (candidates for the deep review / next phase)

- **Interface/abstract method signatures are not Objects** (they have no body).
  This matches the Go extractor's treatment of interface methods, but capturing
  the interface contract could enrich the model — decide intentionally.
- **Call resolution is same-type only.** Cross-class/instance calls
  (`repo.findAll()`, static `Foo.bar()`) are dropped, mirroring the Go
  extractor's "only edges with real endpoints" behaviour. Import-aware
  resolution is a future enhancement.
- **Supertypes resolve by simple name.** Two types with the same simple name in
  different packages could collide; first-writer-wins in `typeByName`.
- **Not a full Java grammar.** Exotic constructs (annotations with `{}` array
  values, sealed/permits clauses, generic bounds) parse structurally but are not
  semantically modelled beyond names.

---

## 3. How to verify (reviewer quick start)

```bash
# Build + static checks
go build ./pkg/... ./cmd/...
go vet ./pkg/... ./cmd/...

# Tests (category axioms + Java extractor suites)
go test ./pkg/... -count=1

# End-to-end on a Java tree (auto-detect), then analyse + verify
go build -o /tmp/catreview ./cmd/catreview
/tmp/catreview extract <java-project> -o model.json --pretty   # prints "Language: java"
/tmp/catreview analyze model.json
/tmp/catreview verify model.json --max-cycles 0 --fail-on-violation

# Force a language on a polyglot repo
/tmp/catreview extract <path> --lang java -o model.json
```

**Verified this session**: a 2-file Java service extracted 11 objects / 12
morphisms, inheritance resolved `UserService → Repository`, `findById → lookup`
call edge captured, cycles = 0, category axioms ✅. Go self-analysis still works
via the factory path (173 objects on `./pkg`).

> Note: `go build ./...` fails in `examples/python/` (pre-existing: multiple
> `main` funcs in one package). That directory is example scripts, not part of
> the core module — build `./pkg/... ./cmd/...` to avoid it. Cleaning up
> `examples/python` is an optional item for the next phase.

---

## 4. Deep-review checklist (for the parallel agents)

Suggested independent review dimensions:

1. **Correctness / parser robustness** — adversarial Java inputs: annotations
   with `{}` args, nested generics, text blocks, `enum` constants with bodies,
   anonymous/local classes, `record` compact constructors, unicode identifiers.
2. **Categorical validity** — every extracted model must pass `VerifyAxioms()`;
   confirm no dangling morphism endpoints; confirm identity morphisms per object.
3. **Parity with GoExtractor** — object/morphism id conventions, metadata keys,
   `Stats()` shape; anything a downstream consumer relies on.
4. **CLI/UX** — `--lang` validation, auto-detect on polyglot repos, error
   messages, exit codes.
5. **Docs accuracy** — README architecture section, workflow doc, INDEX FAQ all
   reflect Go + Java; no remaining "Go only" claims.
6. **Test coverage gaps** — abstract methods, call resolution, external_type
   placeholders, large real-world repo (e.g. a Spring module).

---

## 5. Adding the secondary language set (next phase)

The architecture is drop-in. To add a language `L`:

1. Create `pkg/extractor/L_extractor.go` implementing `Extractor`
   (`ExtractFromPath`, `Language`, `FileExtensions`) and producing the same
   `category.Category` shape (see the mapping table above).
2. Reset per-run state at the top of `ExtractFromPath` (factory reuse safety).
3. Register it in `NewExtractorFactory()` via `factory.Register(NewLExtractor())`.
   Language auto-detection picks it up automatically from `FileExtensions()`.
4. Add `pkg/extractor/L_extractor_test.go` mirroring the Java suite (kinds,
   nesting, comments/strings ignored, inheritance, calls, axiom validity).
5. Update `README.md` architecture diagram, `docs/CATEGORICAL-REPO-ANALYSIS.md`
   supported-languages notes, and `docs/INDEX.md` FAQ.

Priority order discussed: **Java (done) → Python → TypeScript → Rust.**

---

## 6. Open questions for the maintainer

- Should interface/abstract method **signatures** become `function` Objects
  (richer contract modelling) or stay out for Go parity?
- Do you want **import-aware cross-type call resolution** now, or keep it
  same-type until the multi-language analysis needs it?
- Should we add a **GitHub Actions CI workflow** (build/vet/test) as part of the
  enterprise-readiness pass before your friend reviews? (Not done this session —
  "the workflow" here referred to the dev/analysis RMP doc.)
