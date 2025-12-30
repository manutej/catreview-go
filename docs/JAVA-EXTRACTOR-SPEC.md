# Java Extractor Specification

**Version**: 1.1.0
**Status**: Implementation Ready
**Based on**: GoExtractor v1.0 production patterns

---

## 1. Purpose

The Java Extractor implements the `Extractor` interface to extract categorical models from Java source code, producing identical categorical structures as GoExtractor to ensure uniform analysis across languages.

## 2. Design Principles

### 2.1 Uniformity with GoExtractor

The Java Extractor must produce **identical categorical structure** to enable:
- Same complexity analysis algorithms (Basu-Isik, Kolmogorov, Coupling)
- Same functor transformations (file → package abstraction)
- Same axiom verification (associativity, identity)
- Same cycle detection algorithms

### 2.2 Language-Native AST Parsing

- **Primary Approach**: Use JavaParser library via subprocess
- **Alternative**: Eclipse JDT Core (if JavaParser insufficient)
- **Rationale**: Leverage Java's mature AST parsing ecosystem while keeping all analysis in Go

### 2.3 Category Axiom Preservation

All extracted categories MUST satisfy:
1. **Associativity**: (h ∘ g) ∘ f = h ∘ (g ∘ f)
2. **Identity**: f ∘ id_A = f and id_B ∘ f = f

---

## 3. Java-to-Category Mapping

### 3.1 Objects (Ob(C))

| Java Construct | Object Type | Object ID Format | Example |
|----------------|-------------|------------------|---------|
| Package | `"package"` | Package path | `"com.example.service"` |
| Class | `"class"` | Fully qualified name | `"com.example.service.UserService"` |
| Interface | `"interface"` | Fully qualified name | `"com.example.repository.Repository"` |
| Enum | `"enum"` | Fully qualified name | `"com.example.model.Status"` |
| Method | `"function"` | Class.method signature | `"com.example.service.UserService.save(User)"` |
| Constructor | `"function"` | Class.<init> | `"com.example.model.User.<init>(String,int)"` |

### 3.2 Morphisms (Hom(A,B))

| Java Dependency | Morphism Type | Source | Target | Example |
|-----------------|---------------|--------|--------|---------|
| import statement | `"import"` | File/Class | Imported class | `UserService → User` |
| Method call | `"function_call"` | Caller method | Called method | `save() → validate()` |
| Class inheritance (extends) | `"inheritance"` | Subclass | Superclass | `AdminUser → User` |
| Interface implementation | `"implementation"` | Class | Interface | `UserRepository → Repository` |
| Field type reference | `"field_type"` | Class | Field type | `UserService → UserRepository` |
| Method parameter type | `"param_type"` | Method | Parameter type | `save(User) → User` |
| Method return type | `"return_type"` | Method | Return type | `findUser() → User` |

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
│  JavaExtractor (Go)                     │
│  - Walks .java files                    │
│  - Calls subprocess per file            │
│  - Parses JSON output                   │
│  - Builds category.Category             │
└─────────────────────────────────────────┘
              │
              ▼ exec.Command
┌─────────────────────────────────────────┐
│  JavaParser CLI (Java)                  │
│  - Parses .java file to AST             │
│  - Extracts declarations & references   │
│  - Outputs JSON                         │
└─────────────────────────────────────────┘
              │
              ▼ JSON output
┌─────────────────────────────────────────┐
│  Structured Extraction Data             │
│  {                                      │
│    "package": "com.example",            │
│    "classes": [...],                    │
│    "methods": [...],                    │
│    "imports": [...],                    │
│    "calls": [...]                       │
│  }                                      │
└─────────────────────────────────────────┘
```

### 4.2 JavaParser CLI Tool

**Location**: `tools/java-parser/`

**Build**:
```bash
cd tools/java-parser
mvn clean package
# Produces: target/java-parser-cli.jar
```

**Usage**:
```bash
java -jar java-parser-cli.jar --file UserService.java --output json
```

**Output Format**:
```json
{
  "package": "com.example.service",
  "imports": [
    {"name": "com.example.model.User", "static": false},
    {"name": "com.example.repository.UserRepository", "static": false}
  ],
  "classes": [
    {
      "name": "UserService",
      "type": "class",
      "extends": "BaseService",
      "implements": ["Service"],
      "fields": [
        {"name": "repository", "type": "UserRepository"}
      ],
      "methods": [
        {
          "name": "save",
          "parameters": [{"name": "user", "type": "User"}],
          "returnType": "void",
          "calls": [
            {"target": "repository.save", "line": 15},
            {"target": "user.validate", "line": 14}
          ]
        }
      ]
    }
  ]
}
```

---

## 5. Extraction Algorithm

### 5.1 High-Level Flow

```
1. ExtractFromPath(root string)
   ├─> Walk directory tree
   ├─> For each .java file:
   │   ├─> Parse file via JavaParser subprocess
   │   ├─> Extract package → Create package Object
   │   ├─> Extract classes/interfaces/enums → Create type Objects
   │   ├─> Extract methods/constructors → Create function Objects
   │   ├─> Extract imports → Create import Morphisms
   │   ├─> Extract method calls → Create call Morphisms
   │   ├─> Extract inheritance/implements → Create inheritance Morphisms
   │   └─> Extract field/param/return types → Create type dependency Morphisms
   ├─> Create identity Morphisms for all Objects
   └─> Return category.Category
```

### 5.2 Detailed Extraction Steps

#### Step 1: Package Extraction
```go
packageObj := category.NewObject(
    packageName,           // ID: "com.example.service"
    "package",            // Type
    packageName,          // Name
    map[string]interface{}{
        "language": "java",
        "files":    []string{filePath},
    },
)
category.AddObject(packageObj)
```

#### Step 2: Class/Interface/Enum Extraction
```go
for _, classDecl := range parsed.Classes {
    fqn := packageName + "." + classDecl.Name
    classObj := category.NewObject(
        fqn,                    // ID: "com.example.service.UserService"
        classDecl.Type,         // "class" | "interface" | "enum"
        classDecl.Name,         // "UserService"
        map[string]interface{}{
            "language":   "java",
            "package":    packageName,
            "file":       filePath,
            "extends":    classDecl.Extends,
            "implements": classDecl.Implements,
            "fields":     classDecl.Fields,
        },
    )
    category.AddObject(classObj)
}
```

#### Step 3: Method Extraction
```go
for _, methodDecl := range classDecl.Methods {
    methodID := fqn + "." + methodDecl.Name + "(" + formatParams(methodDecl.Parameters) + ")"
    methodObj := category.NewObject(
        methodID,               // ID: "UserService.save(User)"
        "function",
        methodDecl.Name,        // "save"
        map[string]interface{}{
            "language":    "java",
            "class":       fqn,
            "parameters":  methodDecl.Parameters,
            "returnType":  methodDecl.ReturnType,
            "file":        filePath,
            "line":        methodDecl.Line,
        },
    )
    category.AddObject(methodObj)
}
```

#### Step 4: Import Morphism Creation
```go
for _, imp := range parsed.Imports {
    importMorphism := category.NewMorphism(
        generateMorphismID(),
        fqn,                    // Source: importing class
        imp.Name,               // Target: imported class
        "import",
        map[string]interface{}{
            "static": imp.Static,
            "file":   filePath,
        },
    )
    category.AddMorphism(importMorphism)
}
```

#### Step 5: Method Call Morphism Creation
```go
for _, call := range methodDecl.Calls {
    callMorphism := category.NewMorphism(
        generateMorphismID(),
        methodID,               // Source: calling method
        resolveCallTarget(call.Target), // Target: called method
        "function_call",
        map[string]interface{}{
            "line": call.Line,
            "file": filePath,
        },
    )
    category.AddMorphism(callMorphism)
}
```

#### Step 6: Inheritance/Implementation Morphism Creation
```go
// Extends
if classDecl.Extends != "" {
    extendsMorphism := category.NewMorphism(
        generateMorphismID(),
        fqn,                    // Source: subclass
        resolveType(classDecl.Extends), // Target: superclass
        "inheritance",
        map[string]interface{}{
            "kind": "extends",
        },
    )
    category.AddMorphism(extendsMorphism)
}

// Implements
for _, iface := range classDecl.Implements {
    implMorphism := category.NewMorphism(
        generateMorphismID(),
        fqn,                    // Source: implementing class
        resolveType(iface),     // Target: interface
        "implementation",
        map[string]interface{}{},
    )
    category.AddMorphism(implMorphism)
}
```

#### Step 7: Type Dependency Morphisms
```go
// Field types
for _, field := range classDecl.Fields {
    fieldMorphism := category.NewMorphism(
        generateMorphismID(),
        fqn,                    // Source: class with field
        resolveType(field.Type), // Target: field type
        "field_type",
        map[string]interface{}{
            "field": field.Name,
        },
    )
    category.AddMorphism(fieldMorphism)
}
```

---

## 6. Testing Requirements

### 6.1 Unit Tests

**File**: `pkg/extractor/java_extractor_test.go`

```go
func TestJavaExtractorBasic(t *testing.T) {
    // Given: Simple Java class
    // When: Extract categorical model
    // Then: Verify objects and morphisms extracted correctly
}

func TestJavaExtractorInheritance(t *testing.T) {
    // Given: Java class with extends and implements
    // When: Extract categorical model
    // Then: Verify inheritance morphisms created
}

func TestJavaExtractorMethodCalls(t *testing.T) {
    // Given: Java class with method calls
    // When: Extract categorical model
    // Then: Verify call morphisms created with correct targets
}

func TestJavaExtractorCategoryAxioms(t *testing.T) {
    // Given: Extracted Java category
    // When: Verify axioms
    // Then: Associativity and identity laws hold
}
```

### 6.2 Integration Tests

**File**: `pkg/extractor/java_integration_test.go`

```go
func TestJavaExtractorRealProject(t *testing.T) {
    // Given: Small real Java project (Spring Boot sample)
    // When: Extract categorical model
    // Then: Verify realistic object/morphism counts
}

func TestJavaVsGoConsistency(t *testing.T) {
    // Given: Equivalent Go and Java codebases
    // When: Extract both
    // Then: Verify similar categorical structure
}
```

### 6.3 Test Data

**Location**: `testdata/java/`

```
testdata/java/
├── simple/
│   └── HelloWorld.java       # Basic class, no dependencies
├── inheritance/
│   ├── Base.java              # Superclass
│   └── Derived.java           # Subclass with extends
├── interfaces/
│   ├── Repository.java        # Interface
│   └── UserRepository.java    # Implementation
└── realistic/
    ├── User.java              # Model
    ├── UserRepository.java    # Repository interface
    ├── UserService.java       # Service with dependencies
    └── UserController.java    # Controller
```

### 6.4 Quality Gates

- ✅ **Unit test coverage**: ≥80%
- ✅ **Category axioms**: 100% pass rate
- ✅ **Morphism accuracy**: ≥95% (compared to manual inspection)
- ✅ **Performance**: <100ms per file for typical classes

---

## 7. Comparison with GoExtractor

### 7.1 Structural Equivalence

| Aspect | GoExtractor | JavaExtractor |
|--------|-------------|---------------|
| Package object | 1 per Go package | 1 per Java package |
| Type object | Go struct/interface | Java class/interface/enum |
| Function object | Go func | Java method/constructor |
| Import morphism | Go import | Java import |
| Call morphism | Go func call | Java method call |
| Inheritance morphism | N/A (Go has no inheritance) | Java extends/implements |

### 7.2 Java-Specific Additions

**Additional morphism types not in GoExtractor**:
1. `"inheritance"` - extends relationship
2. `"implementation"` - implements relationship
3. `"field_type"` - field type dependency
4. `"param_type"` - parameter type dependency
5. `"return_type"` - return type dependency

**Rationale**: Java's OOP features require richer dependency tracking while maintaining categorical structure compatibility.

---

## 8. Error Handling

### 8.1 JavaParser Subprocess Errors

```go
if err := cmd.Run(); err != nil {
    return fmt.Errorf("javaparser failed for %s: %w", filePath, err)
}
```

### 8.2 Malformed Java Code

- **Strategy**: Skip file, log error, continue extraction
- **Rationale**: Partial extraction better than complete failure

### 8.3 Missing JavaParser CLI

```go
func (e *JavaExtractor) checkJavaParser() error {
    cmd := exec.Command("java", "-jar", javaParserPath, "--version")
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("JavaParser CLI not found at %s", javaParserPath)
    }
    return nil
}
```

---

## 9. Configuration

### 9.1 Environment Variables

- `JAVAPARSER_JAR`: Path to JavaParser CLI jar (default: `./tools/java-parser/target/java-parser-cli.jar`)
- `JAVA_HOME`: Java installation directory

### 9.2 Extraction Options

```go
type JavaExtractorOptions struct {
    SkipTests       bool   // Skip test files (*Test.java)
    JavaParserPath  string // Path to JavaParser CLI
    IncludePrivate  bool   // Include private methods/fields
}
```

---

## 10. Success Criteria

### 10.1 Functional Requirements

- [x] Implements Extractor interface
- [ ] Extracts all Java constructs to categorical objects
- [ ] Creates all dependency morphisms accurately
- [ ] Preserves category axioms (associativity, identity)
- [ ] Produces identical structure to GoExtractor for equivalent code

### 10.2 Quality Requirements

- [ ] Unit test coverage ≥80%
- [ ] Integration tests pass on realistic Java projects
- [ ] Performance: <100ms per file
- [ ] Category axiom verification: 100% pass rate

### 10.3 Documentation Requirements

- [x] Specification complete (this document)
- [ ] Implementation guide for contributors
- [ ] API documentation (godoc)
- [ ] Example usage in README

---

## 11. Future Enhancements (v1.2+)

### 11.1 Advanced Features

- **Generics**: Extract type parameters and constraints
- **Annotations**: Capture @annotations as metadata
- **Lambda expressions**: Track functional interfaces and closures
- **Inner classes**: Handle nested class declarations

### 11.2 Performance Optimizations

- **Parallel parsing**: Process multiple files concurrently
- **Caching**: Cache parsed results for incremental analysis
- **Lazy loading**: Stream parse results instead of loading all into memory

---

**References**:
- GoExtractor implementation: `pkg/extractor/go_extractor.go`
- Category theory foundation: `pkg/category/types.go`
- JavaParser documentation: https://javaparser.org/
