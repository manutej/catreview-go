# Analysis Directory

This directory contains categorical analysis results for various codebases analyzed using catreview-go.

## Structure

```
analysis/
├── README.md                    # This file
└── luxor-workspace/             # LUXOR workspace analysis (2025-12-30)
    ├── README.md                # Entry point for LUXOR analysis
    ├── docs/                    # 10 comprehensive analysis documents
    ├── data/                    # 9 JSON categorical models
    ├── visualizations/          # 24 SVG + 24 DOT graph files
    ├── summaries/               # 9 project summary text files
    └── scripts/                 # Python analysis scripts
```

## Analyses

### LUXOR Workspace (2025-12-30)
- **Coverage**: 96% (12,750 files across 9 Python projects)
- **Documentation**: 10 comprehensive guides
- **Visualizations**: 24 SVG graphs (inheritance, modules, composition)
- **Key Finding**: "LUXOR Architectural Fingerprint" - functional-first, graph-centric, type-safe
- **Read**: [luxor-workspace/README.md](luxor-workspace/README.md)

---

## Adding New Analyses

To add analysis of a new codebase:

1. Create directory: `analysis/{codebase-name}/`
2. Run extraction: `python3 scripts/python_categorical_extractor.py /path/to/codebase {codebase-name}`
3. Generate visualizations: `go run examples/python/visualize_project.go --input {codebase-name}-analysis.json`
4. Organize results in the appropriate subdirectories
5. Create README.md documenting findings

---

## Tools

Analysis tools are located in the repository root:
- `python_categorical_extractor.py` - Python AST analyzer (copied to luxor-workspace/scripts/)
- `examples/python/visualize_project.go` - Visualization generator
- `cmd/catreview/` - Go-based categorical extractor

---

**Note**: This directory is tracked in Git to preserve analysis results alongside the analysis tools.
