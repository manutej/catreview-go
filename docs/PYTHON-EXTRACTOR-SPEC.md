# Python Extractor Specification

**Version**: 1.1.0
**Status**: Implementation Ready
**Based on**: GoExtractor v1.0 production patterns

---

## 1. Purpose

The Python Extractor implements the `Extractor` interface to extract categorical models from Python source code, producing identical categorical structures as GoExtractor to ensure uniform analysis across languages.

## 2. Design Principles

### 2.1 Uniformity with GoExtractor

The Python Extractor must produce **identical categorical structure** to enable:
- Same complexity analysis algorithms (Basu-Isik, Kolmogorov, Coupling)
- Same functor transformations (file → package abstraction)
- Same axiom verification (associativity, identity)
- Same cycle detection algorithms

### 2.2 Language-Native AST Parsing

- **Primary Approach**: Use Python's built-in `ast` module via subprocess
- **Alternative**: tree-sitter-python (if ast insufficient)
- **Rationale**: Leverage Python's excellent built-in AST while keeping all analysis in Go

### 2.3 Category Axiom Preservation

All extracted categories MUST satisfy:
1. **Associativity**: (h ∘ g) ∘ f = h ∘ (g ∘ f)
2. **Identity**: f ∘ id_A = f and id_B ∘ f = f

---

## 3. Python-to-Category Mapping

### 3.1 Objects (Ob(C))

| Python Construct | Object Type | Object ID Format | Example |
|------------------|-------------|------------------|---------|
| Module | `"module"` | Module path | `"app.services.user_service"` |
| Package | `"package"` | Package path | `"app.services"` |
| Class | `"class"` | Module.ClassName | `"app.services.user_service.UserService"` |
| Function (top-level) | `"function"` | Module.function_name | `"app.utils.helpers.format_date"` |
| Method | `"function"` | Class.method_name | `"UserService.save"` |

### 3.2 Morphisms (Hom(A,B))

| Python Dependency | Morphism Type | Source | Target | Example |
|-------------------|---------------|--------|--------|---------|
| import module | `"import"` | File/Module | Imported module | `user_service → user_model` |
| from X import Y | `"import"` | File/Module | Imported name | `user_service → User` |
| Function call | `"function_call"` | Caller | Called function | `save() → validate()` |
| Method call | `"function_call"` | Caller | Called method | `service.save() → save()` |
| Class inheritance | `"inheritance"` | Subclass | Base class | `AdminUser → User` |
| Type hint | `"type_dependency"` | Function/Method | Type | `save(user: User) → User` |
| Decorator | `"decorator"` | Decorated | Decorator | `login_required → require_auth` |

### 3.3 Identity Morphisms

For each object O, create identity morphism:
- **ID**: `"id_" + object.ID`
- **Type**: `"identity"`
- **Source**: O
- **Target**: O

---

## 4. Implementation Architecture

### 4.1 Subprocess Execution Pattern

```
┌─────────────────────────────────────────┐
│  PythonExtractor (Go)                   │
│  - Walks .py files                      │
│  - Calls subprocess per file            │
│  - Parses JSON output                   │
│  - Builds category.Category             │
└─────────────────────────────────────────┘
              │
              ▼ exec.Command
┌─────────────────────────────────────────┐
│  Python AST Parser (Python script)      │
│  - Uses ast.parse() to build AST        │
│  - NodeVisitor extracts declarations    │
│  - Outputs JSON to stdout               │
└─────────────────────────────────────────┘
              │
              ▼ JSON output
┌─────────────────────────────────────────┐
│  Structured Extraction Data             │
│  {                                      │
│    "module": "user_service",            │
│    "classes": [...],                    │
│    "functions": [...],                  │
│    "imports": [...],                    │
│    "calls": [...]                       │
│  }                                      │
└─────────────────────────────────────────┘
```

### 4.2 Python AST Parser Script

**Location**: `tools/python-parser/parse_python.py`

**Usage**:
```bash
python3 tools/python-parser/parse_python.py user_service.py
```

**Implementation** (excerpt):
```python
import ast
import json
import sys

class PythonExtractor(ast.NodeVisitor):
    def __init__(self):
        self.module_name = ""
        self.imports = []
        self.classes = []
        self.functions = []
        self.calls = []

    def visit_Import(self, node):
        for alias in node.names:
            self.imports.append({
                "type": "import",
                "name": alias.name,
                "asname": alias.asname,
                "line": node.lineno
            })
        self.generic_visit(node)

    def visit_ImportFrom(self, node):
        for alias in node.names:
            self.imports.append({
                "type": "from_import",
                "module": node.module,
                "name": alias.name,
                "asname": alias.asname,
                "line": node.lineno
            })
        self.generic_visit(node)

    def visit_ClassDef(self, node):
        class_info = {
            "name": node.name,
            "bases": [self.get_name(base) for base in node.bases],
            "methods": [],
            "line": node.lineno,
            "decorators": [self.get_name(d) for d in node.decorator_list]
        }

        # Extract methods
        for item in node.body:
            if isinstance(item, ast.FunctionDef):
                class_info["methods"].append({
                    "name": item.name,
                    "parameters": [arg.arg for arg in item.args.args],
                    "decorators": [self.get_name(d) for d in item.decorator_list],
                    "line": item.lineno
                })

        self.classes.append(class_info)
        self.generic_visit(node)

    def visit_FunctionDef(self, node):
        # Only top-level functions (methods handled in visit_ClassDef)
        if not isinstance(node.scope, ast.ClassDef):
            func_info = {
                "name": node.name,
                "parameters": [arg.arg for arg in node.args.args],
                "decorators": [self.get_name(d) for d in node.decorator_list],
                "line": node.lineno
            }
            self.functions.append(func_info)
        self.generic_visit(node)

    def visit_Call(self, node):
        call_info = {
            "func": self.get_name(node.func),
            "line": node.lineno
        }
        self.calls.append(call_info)
        self.generic_visit(node)

def main():
    if len(sys.argv) < 2:
        print("Usage: parse_python.py <file.py>", file=sys.stderr)
        sys.exit(1)

    with open(sys.argv[1]) as f:
        source = f.read()

    tree = ast.parse(source, filename=sys.argv[1])
    extractor = PythonExtractor()
    extractor.visit(tree)

    result = {
        "module": sys.argv[1].replace('.py', '').replace('/', '.'),
        "imports": extractor.imports,
        "classes": extractor.classes,
        "functions": extractor.functions,
        "calls": extractor.calls
    }

    print(json.dumps(result, indent=2))

if __name__ == "__main__":
    main()
```

**Output Format**:
```json
{
  "module": "app.services.user_service",
  "imports": [
    {"type": "from_import", "module": "app.models", "name": "User", "asname": null, "line": 1},
    {"type": "from_import", "module": "app.repositories", "name": "UserRepository", "asname": null, "line": 2}
  ],
  "classes": [
    {
      "name": "UserService",
      "bases": ["BaseService"],
      "decorators": [],
      "line": 5,
      "methods": [
        {
          "name": "save",
          "parameters": ["self", "user"],
          "decorators": [],
          "line": 10
        }
      ]
    }
  ],
  "functions": [
    {"name": "helper_function", "parameters": ["x", "y"], "decorators": [], "line": 25}
  ],
  "calls": [
    {"func": "repository.save", "line": 15},
    {"func": "user.validate", "line": 14}
  ]
}
```

---

## 5. Extraction Algorithm

### 5.1 High-Level Flow

```
1. ExtractFromPath(root string)
   ├─> Walk directory tree
   ├─> For each .py file:
   │   ├─> Parse file via Python ast subprocess
   │   ├─> Extract module → Create module Object
   │   ├─> Extract classes → Create class Objects
   │   ├─> Extract functions/methods → Create function Objects
   │   ├─> Extract imports → Create import Morphisms
   │   ├─> Extract calls → Create call Morphisms
   │   ├─> Extract inheritance → Create inheritance Morphisms
   │   └─> Extract type hints → Create type dependency Morphisms
   ├─> Create identity Morphisms for all Objects
   └─> Return category.Category
```

### 5.2 Detailed Extraction Steps

#### Step 1: Module Extraction
```go
moduleObj := category.NewObject(
    moduleName,            // ID: "app.services.user_service"
    "module",             // Type
    moduleName,           // Name
    map[string]interface{}{
        "language": "python",
        "file":     filePath,
    },
)
category.AddObject(moduleObj)
```

#### Step 2: Class Extraction
```go
for _, classDecl := range parsed.Classes {
    fqn := moduleName + "." + classDecl.Name
    classObj := category.NewObject(
        fqn,                     // ID: "app.services.user_service.UserService"
        "class",                 // Type
        classDecl.Name,          // "UserService"
        map[string]interface{}{
            "language":    "python",
            "module":      moduleName,
            "file":        filePath,
            "bases":       classDecl.Bases,
            "decorators":  classDecl.Decorators,
            "line":        classDecl.Line,
        },
    )
    category.AddObject(classObj)
}
```

#### Step 3: Function/Method Extraction
```go
// Top-level functions
for _, funcDecl := range parsed.Functions {
    funcID := moduleName + "." + funcDecl.Name
    funcObj := category.NewObject(
        funcID,                  // ID: "app.utils.format_date"
        "function",
        funcDecl.Name,           // "format_date"
        map[string]interface{}{
            "language":    "python",
            "module":      moduleName,
            "parameters":  funcDecl.Parameters,
            "decorators":  funcDecl.Decorators,
            "line":        funcDecl.Line,
        },
    )
    category.AddObject(funcObj)
}

// Methods
for _, methodDecl := range classDecl.Methods {
    methodID := fqn + "." + methodDecl.Name
    methodObj := category.NewObject(
        methodID,                // ID: "UserService.save"
        "function",
        methodDecl.Name,         // "save"
        map[string]interface{}{
            "language":    "python",
            "class":       fqn,
            "parameters":  methodDecl.Parameters,
            "decorators":  methodDecl.Decorators,
            "line":        methodDecl.Line,
        },
    )
    category.AddObject(methodObj)
}
```

#### Step 4: Import Morphism Creation
```go
for _, imp := range parsed.Imports {
    var target string
    if imp.Type == "from_import" {
        target = imp.Module + "." + imp.Name
    } else {
        target = imp.Name
    }

    importMorphism := category.NewMorphism(
        generateMorphismID(),
        moduleName,              // Source: importing module
        target,                  // Target: imported module/name
        "import",
        map[string]interface{}{
            "type":   imp.Type,
            "asname": imp.Asname,
            "line":   imp.Line,
        },
    )
    category.AddMorphism(importMorphism)
}
```

#### Step 5: Function Call Morphism Creation
```go
for _, call := range parsed.Calls {
    callMorphism := category.NewMorphism(
        generateMorphismID(),
        currentContext,          // Source: current function/method
        resolveCallTarget(call.Func), // Target: called function
        "function_call",
        map[string]interface{}{
            "line": call.Line,
            "file": filePath,
        },
    )
    category.AddMorphism(callMorphism)
}
```

#### Step 6: Inheritance Morphism Creation
```go
for _, base := range classDecl.Bases {
    inheritanceMorphism := category.NewMorphism(
        generateMorphismID(),
        fqn,                     // Source: derived class
        resolveType(base),       // Target: base class
        "inheritance",
        map[string]interface{}{},
    )
    category.AddMorphism(inheritanceMorphism)
}
```

#### Step 7: Decorator Morphisms
```go
for _, decorator := range classDecl.Decorators {
    decoratorMorphism := category.NewMorphism(
        generateMorphismID(),
        fqn,                     // Source: decorated class/function
        resolveDecorator(decorator), // Target: decorator function
        "decorator",
        map[string]interface{}{},
    )
    category.AddMorphism(decoratorMorphism)
}
```

---

## 6. Python-Specific Considerations

### 6.1 Dynamic Typing

**Challenge**: Python's duck typing means not all dependencies are statically analyzable.

**Solution**:
- Extract type hints when available (PEP 484)
- Track runtime calls via ast.Call nodes
- Document limitation: dynamic getattr/setattr calls may be missed

### 6.2 Multiple Inheritance

**Challenge**: Python supports multiple base classes.

**Solution**:
- Create separate inheritance morphism for each base
- Preserve MRO (Method Resolution Order) in metadata

### 6.3 __init__.py Packages

**Challenge**: __init__.py creates package namespaces with re-exports.

**Solution**:
- Treat __init__.py as package object
- Create import morphisms for re-exported names

### 6.4 Decorators

**Challenge**: Decorators modify function/class behavior.

**Solution**:
- Create decorator morphism from decorated → decorator
- Store decorator chain in metadata

---

## 7. Testing Requirements

### 7.1 Unit Tests

**File**: `pkg/extractor/python_extractor_test.go`

```go
func TestPythonExtractorBasic(t *testing.T) {
    // Given: Simple Python module
    // When: Extract categorical model
    // Then: Verify objects and morphisms extracted correctly
}

func TestPythonExtractorInheritance(t *testing.T) {
    // Given: Python class with inheritance
    // When: Extract categorical model
    // Then: Verify inheritance morphisms created
}

func TestPythonExtractorDecorators(t *testing.T) {
    // Given: Python functions with decorators
    // When: Extract categorical model
    // Then: Verify decorator morphisms created
}

func TestPythonExtractorTypeHints(t *testing.T) {
    // Given: Python code with type hints
    // When: Extract categorical model
    // Then: Verify type dependency morphisms created
}

func TestPythonExtractorCategoryAxioms(t *testing.T) {
    // Given: Extracted Python category
    // When: Verify axioms
    // Then: Associativity and identity laws hold
}
```

### 7.2 Integration Tests

**File**: `pkg/extractor/python_integration_test.go`

```go
func TestPythonExtractorRealProject(t *testing.T) {
    // Given: Small real Python project (Flask app)
    // When: Extract categorical model
    // Then: Verify realistic object/morphism counts
}

func TestPythonVsGoConsistency(t *testing.T) {
    // Given: Equivalent Go and Python codebases
    // When: Extract both
    // Then: Verify similar categorical structure
}
```

### 7.3 Test Data

**Location**: `testdata/python/`

```
testdata/python/
├── simple/
│   └── hello.py               # Basic module, no dependencies
├── inheritance/
│   ├── base.py                # Base class
│   └── derived.py             # Derived class
├── decorators/
│   ├── auth.py                # Decorator definitions
│   └── routes.py              # Decorated functions
└── realistic/
    ├── __init__.py            # Package
    ├── models/
    │   ├── __init__.py
    │   └── user.py            # User model
    ├── repositories/
    │   ├── __init__.py
    │   └── user_repository.py # Repository
    └── services/
        ├── __init__.py
        └── user_service.py    # Service with dependencies
```

### 7.4 Quality Gates

- ✅ **Unit test coverage**: ≥80%
- ✅ **Category axioms**: 100% pass rate
- ✅ **Morphism accuracy**: ≥90% (lower than Java due to dynamic typing)
- ✅ **Performance**: <50ms per file for typical modules

---

## 8. Comparison with GoExtractor

### 8.1 Structural Equivalence

| Aspect | GoExtractor | PythonExtractor |
|--------|-------------|-----------------|
| Package object | 1 per Go package | 1 per Python package |
| Module object | N/A (Go files ≠ modules) | 1 per .py file |
| Type object | Go struct/interface | Python class |
| Function object | Go func | Python function/method |
| Import morphism | Go import | Python import/from |
| Call morphism | Go func call | Python function/method call |
| Inheritance morphism | N/A (Go has no inheritance) | Python class inheritance |

### 8.2 Python-Specific Additions

**Additional morphism types not in GoExtractor**:
1. `"inheritance"` - class inheritance
2. `"decorator"` - decorator application
3. `"type_dependency"` - type hint dependencies

**Additional metadata**:
- Decorators list for classes/functions
- Type hints for parameters/returns
- Multiple inheritance base list

---

## 9. Error Handling

### 9.1 Python Subprocess Errors

```go
if err := cmd.Run(); err != nil {
    return fmt.Errorf("python ast parser failed for %s: %w", filePath, err)
}
```

### 9.2 Malformed Python Code

- **Strategy**: Skip file, log error, continue extraction
- **Rationale**: Syntax errors common during development, partial extraction better than failure

### 9.3 Missing Python Interpreter

```go
func (e *PythonExtractor) checkPython() error {
    cmd := exec.Command("python3", "--version")
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("python3 not found in PATH")
    }
    return nil
}
```

---

## 10. Configuration

### 10.1 Environment Variables

- `PYTHON_BIN`: Path to Python interpreter (default: `python3`)
- `PYTHON_PARSER_SCRIPT`: Path to parse_python.py (default: `./tools/python-parser/parse_python.py`)

### 10.2 Extraction Options

```go
type PythonExtractorOptions struct {
    SkipTests          bool   // Skip test files (test_*.py, *_test.py)
    PythonBin          string // Python interpreter path
    ParserScript       string // Path to parse_python.py
    IncludePrivate     bool   // Include private methods/functions (_name)
    ExtractTypeHints   bool   // Extract type hints (default: true)
}
```

---

## 11. Success Criteria

### 11.1 Functional Requirements

- [x] Implements Extractor interface
- [ ] Extracts all Python constructs to categorical objects
- [ ] Creates all dependency morphisms accurately
- [ ] Preserves category axioms (associativity, identity)
- [ ] Handles Python-specific features (decorators, duck typing, multiple inheritance)

### 11.2 Quality Requirements

- [ ] Unit test coverage ≥80%
- [ ] Integration tests pass on realistic Python projects
- [ ] Performance: <50ms per file
- [ ] Category axiom verification: 100% pass rate

### 11.3 Documentation Requirements

- [x] Specification complete (this document)
- [ ] Implementation guide for contributors
- [ ] API documentation (godoc)
- [ ] Example usage in README

---

## 12. Future Enhancements (v1.2+)

### 12.1 Advanced Features

- **Type stubs (.pyi)**: Extract from stub files for better typing
- **Async/await**: Track async function calls and coroutines
- **Context managers**: Extract __enter__/__exit__ relationships
- **Metaclasses**: Handle metaclass relationships
- **Descriptors**: Track property/descriptor usage

### 12.2 Performance Optimizations

- **Parallel parsing**: Process multiple files concurrently
- **Caching**: Cache parsed results for incremental analysis
- **Compiled AST**: Use Python's __pycache__ for faster parsing

---

**References**:
- GoExtractor implementation: `pkg/extractor/go_extractor.go`
- Category theory foundation: `pkg/category/types.go`
- Python ast module: https://docs.python.org/3/library/ast.html
- PEP 484 Type Hints: https://peps.python.org/pep-0484/
