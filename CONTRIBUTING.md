# Contributing to Catreview-Go

Thank you for considering contributing to catreview-go! This document provides guidelines for contributing to the project.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Adding Language Support](#adding-language-support)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)

---

## Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

---

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title** and description
- **Steps to reproduce** the issue
- **Expected behavior** vs. **actual behavior**
- **Go version** and **OS** information
- **Sample code** or repository link if applicable

**Template**:
```markdown
**Bug Description**: [Clear, concise description]

**Steps to Reproduce**:
1. Run `catreview extract ./repo`
2. See error

**Expected**: Should extract categorical model
**Actual**: Error: "..."

**Environment**:
- Go version: 1.21.5
- OS: macOS 14.1
- catreview version: 1.0.0
```

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. Include:

- **Clear title** and detailed description
- **Use case**: Why is this enhancement useful?
- **Alternatives considered**
- **Implementation ideas** (optional)

### Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`go test ./pkg/...`)
6. Commit with descriptive messages
7. Push to your fork
8. Open a Pull Request

---

## Development Setup

### Prerequisites

- **Go 1.21+** (required)
- **Git** (for version control)
- **Make** (optional, for build automation)

### Setup Steps

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/catreview-go
cd catreview-go

# Build the project
go build -o catreview ./cmd/catreview

# Run tests
go test ./pkg/...

# Run on sample project
./catreview extract ./pkg -o test-model.json
./catreview analyze test-model.json
```

### Project Structure

```
catreview-go/
├── cmd/catreview/       # CLI application
├── pkg/
│   ├── category/        # Core category theory
│   ├── functor/         # Functor system
│   ├── analysis/        # Complexity metrics
│   └── extractor/       # Language extractors
├── examples/            # Analysis examples
├── docs/                # Documentation
└── tests/               # Integration tests
```

---

## Pull Request Process

### Before Submitting

1. **Ensure tests pass**: `go test ./pkg/...`
2. **Run linter**: `golangci-lint run` (if available)
3. **Update documentation**: Modify README.md if adding features
4. **Add examples**: For new functionality, add to examples/
5. **Self-review**: Check diff for unintended changes

### PR Title Format

```
<type>(<scope>): <description>

Examples:
feat(extractor): Add Java AST extractor
fix(analysis): Correct Kolmogorov complexity calculation
docs(readme): Add installation instructions
test(category): Add identity axiom tests
```

**Types**: `feat`, `fix`, `docs`, `test`, `refactor`, `perf`, `chore`

### PR Description Template

```markdown
## Description
[Clear description of changes]

## Motivation
[Why is this change needed?]

## Changes Made
- [ ] Added/Modified X
- [ ] Fixed Y
- [ ] Updated documentation

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests pass
- [ ] Manual testing performed

## Checklist
- [ ] Code follows Go best practices
- [ ] Tests pass locally
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
```

### Review Process

1. **Automated checks** must pass (tests, linting)
2. **Maintainer review** within 2-3 business days
3. **Feedback addressed** via additional commits
4. **Approval** → Merge by maintainer

---

## Coding Standards

### Go Style Guide

Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines:

1. **Formatting**: Use `gofmt` or `goimports`
2. **Naming**:
   - Exported: `CamelCase`
   - Unexported: `camelCase`
   - Interfaces: `Reader`, `Writer` (not `IReader`)
3. **Error Handling**: Always check errors, wrap with context
4. **Comments**: Document all exported functions/types

**Example**:
```go
// ComplexityAnalyzer performs categorical complexity analysis.
// It implements the Basu-Isik diagram complexity metric.
type ComplexityAnalyzer struct {
    category *category.Category
}

// DiagramComplexity computes the total diagram complexity.
// Returns c(D) = Σc_obj + Σc_morph + c_comp.
func (a *ComplexityAnalyzer) DiagramComplexity() float64 {
    // Implementation
}
```

### Code Organization

- **Small functions**: < 30 lines preferred
- **Single responsibility**: Each function does one thing
- **DRY**: Extract common logic into helper functions
- **Interfaces**: Define for testability

---

## Adding Language Support

### Process

1. **Create extractor**: `pkg/extractor/{lang}_extractor.go`
2. **Implement interface**:
   ```go
   type Extractor interface {
       ExtractFromPath(root string) (*category.Category, error)
   }
   ```
3. **Map constructs** to categorical objects/morphisms
4. **Add tests**: `pkg/extractor/{lang}_extractor_test.go`
5. **Update CLI**: Add language detection logic
6. **Document**: Add to README.md

### Example: Java Extractor

```go
package extractor

import (
    "catreview/pkg/category"
    // Java parser library
)

// JavaExtractor extracts categorical models from Java code.
type JavaExtractor struct {
    category *category.Category
}

// ExtractFromPath extracts from Java source directory.
func (e *JavaExtractor) ExtractFromPath(root string) (*category.Category, error) {
    // Parse .java files
    // Extract classes → Objects
    // Extract method calls → Morphisms
    return e.category, nil
}
```

### Language Support Checklist

- [ ] AST parser integrated
- [ ] Objects extraction (classes, methods, functions)
- [ ] Morphisms extraction (dependencies, calls, imports)
- [ ] Identity morphisms created
- [ ] Unit tests with sample code
- [ ] Integration test with real repository
- [ ] Documentation updated
- [ ] Example analysis added to examples/

---

## Testing Guidelines

### Unit Tests

- **Coverage**: Aim for 80%+ coverage
- **Table-driven tests**: For multiple scenarios
- **Clear names**: `TestFunctionName_Scenario_Expected`

**Example**:
```go
func TestCategory_AddObject_Success(t *testing.T) {
    cat := category.NewCategory("test")
    obj := &category.Object{ID: "obj1", Type: "function"}

    err := cat.AddObject(obj)

    assert.NoError(t, err)
    assert.Equal(t, 1, len(cat.Objects()))
}
```

### Integration Tests

Located in `tests/`, these verify end-to-end workflows:

```bash
# Run integration tests
go test ./tests/... -v
```

### Test Data

- Use `testdata/` for sample repositories
- Keep test repos small (< 100 files)
- Include edge cases (cycles, high coupling, etc.)

---

## Documentation

### Code Documentation

- **Godoc**: All exported symbols must have comments
- **Examples**: Add `Example` functions for complex APIs
- **Inline comments**: Explain non-obvious logic

### User Documentation

When adding features, update:

1. **README.md**: Add to relevant section
2. **QUICK-START.md**: If user-facing change
3. **Examples**: Add practical usage examples
4. **CHANGELOG.md**: Document changes (see below)

### Changelog Format

Follow [Keep a Changelog](https://keepachangelog.com/):

```markdown
## [Unreleased]
### Added
- Java extractor for analyzing Java codebases (#123)

### Fixed
- Kolmogorov complexity calculation for large files (#456)

### Changed
- Improved cycle detection algorithm (#789)
```

---

## Questions?

- **GitHub Issues**: For bugs/features
- **Discussions**: For questions/ideas
- **Email**: [your-email] for security issues

---

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing to catreview-go!** 🧮
