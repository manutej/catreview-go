# Categorical Repository Analysis - RMP Meta-Prompt with Ralph Loop

**Version**: 1.0
**Tier**: L4 (Parallel Consensus)
**Quality Threshold**: ≥0.90
**Execution Mode**: Parallel with Ralph Loop Orchestration

---

## Meta-Prompt Specification

### Categorical Structure

```
F: RepoList → [AnalysisTask]      (Functor - map repos to tasks)
||: [Task] → [Result]              (Parallel execution - coproduct)
M: Result →ⁿ Result                (Monad - iterative refinement)
Σ: [Result] → ComparativeReport    (Colimit - merge results)
```

### Execution Flow

```
1. F(repos) = [task₁, task₂, ..., taskₙ]     (Map repos to analysis tasks)
2. task₁ || task₂ || ... || taskₙ             (Execute in parallel)
3. M.bind(quality_check)                       (Verify ≥0.90 for each)
4. Σ(results) = comparative_analysis           (Synthesize findings)
```

---

## Parameterized Template

### Input Parameters

```yaml
parameters:
  repositories:
    - url: string           # GitHub URL
      name: string          # Short name for output
      language: string      # Primary language (go, typescript, python)
      expected_size: string # small | medium | large

  output_config:
    base_dir: string        # Base output directory (default: examples/)
    format: string          # json | markdown | both
    include_visualization: boolean

  quality_gates:
    min_quality: float      # Minimum quality score (default: 0.90)
    max_cycles: int         # Maximum dependency cycles allowed
    fail_on_violation: boolean

  execution:
    parallel: boolean       # Run repos in parallel (default: true)
    timeout_per_repo: int   # Timeout in seconds (default: 300)
    retry_on_failure: boolean
```

### Example Configuration

```yaml
repositories:
  - url: https://github.com/charmbracelet/crush
    name: crush
    language: go
    expected_size: small

  - url: https://github.com/charmbracelet/bubbletea
    name: bubbletea
    language: go
    expected_size: medium

  - url: https://github.com/charmbracelet/soft-serve
    name: soft-serve
    language: go
    expected_size: medium

  - url: https://github.com/charmbracelet/glow
    name: glow
    language: go
    expected_size: small

  - url: https://github.com/anthropics/claude-code
    name: claude-code
    language: typescript
    expected_size: large

  - url: https://github.com/manutej/lumina-ccn
    name: lumina-ccn
    language: go
    expected_size: small

output_config:
  base_dir: examples/
  format: both
  include_visualization: true

quality_gates:
  min_quality: 0.90
  max_cycles: 5
  fail_on_violation: false

execution:
  parallel: true
  timeout_per_repo: 300
  retry_on_failure: true
```

---

## Ralph Loop Integration

### Ralph Configuration

```yaml
ralph_config:
  mode: parallel_orchestration
  stop_hook: quality_gate_check
  max_iterations: 3
  quality_threshold: 0.90

  per_repo_workflow:
    - clone_repository
    - extract_categorical_model
    - analyze_complexity
    - verify_axioms
    - generate_report
    - quality_assessment

  convergence_criteria:
    - all_repos_analyzed: true
    - min_quality_met: true
    - comparative_analysis_complete: true
```

### Stop-Hook Script

```bash
#!/bin/bash
# ~/.claude/validators/categorical-analysis-gate.sh

QUALITY_THRESHOLD=0.90
REPORT_FILE=$1

# Extract quality score from report
quality=$(jq -r '.quality_metrics.aggregate_score' "$REPORT_FILE")

if (( $(echo "$quality < $QUALITY_THRESHOLD" | bc -l) )); then
  echo "❌ Quality gate failed: $quality < $QUALITY_THRESHOLD"
  exit 1
fi

echo "✅ Quality gate passed: $quality ≥ $QUALITY_THRESHOLD"
exit 0
```

---

## RMP Prompt Template

### Phase 1: Repository Cloning (Parallel)

```bash
For each repository in [REPO_LIST]:

  Agent: git-genius
  Task: Clone repository

  Execute:
    git clone [REPO_URL] [TEMP_DIR]/[REPO_NAME]
    cd [TEMP_DIR]/[REPO_NAME]
    git rev-parse --short HEAD > commit.txt
    git log -1 --format="%ai" > timestamp.txt

  Output:
    - Clone status: success | failure
    - Commit SHA: [hash]
    - Clone timestamp: [ISO-8601]

  Quality Gate:
    - Repository accessible: true
    - .git directory exists: true
```

### Phase 2: Categorical Extraction (Parallel)

```bash
For each cloned repository:

  Agent: practical-programmer
  Task: Extract categorical model

  Execute:
    cd [TEMP_DIR]/[REPO_NAME]

    # Detect primary source directory
    SRC_DIR=$(find . -type d -name 'pkg' -o -name 'src' -o -name 'lib' | head -1)
    if [ -z "$SRC_DIR" ]; then
      SRC_DIR="."
    fi

    # Run extraction
    catreview extract "$SRC_DIR" \
      --output [OUTPUT_DIR]/[REPO_NAME]-model.json \
      --pretty

  Output:
    - Model file: [REPO_NAME]-model.json
    - Object count: N
    - Morphism count: M
    - Extraction time: T seconds

  Quality Gate:
    - Model file exists: true
    - Object count > 0: true
    - Valid JSON: true
    - Extraction errors: 0
```

### Phase 3: Complexity Analysis (Parallel)

```bash
For each extracted model:

  Agent: practical-programmer
  Task: Analyze categorical complexity

  Execute:
    catreview analyze \
      [OUTPUT_DIR]/[REPO_NAME]-model.json \
      --output [OUTPUT_DIR]/[REPO_NAME]-analysis.json \
      --pretty

  Output:
    - Diagram complexity: C_D
    - Kolmogorov complexity: K bytes
    - Cycle count: N_cycles
    - Top unstable components: [list]
    - Top coupled components: [list]

  Quality Gate:
    - Analysis complete: true
    - Complexity score > 0: true
    - Cycles ≤ max_cycles: true
```

### Phase 4: Axiom Verification (Parallel)

```bash
For each categorical model:

  Agent: practical-programmer
  Task: Verify category axioms

  Execute:
    catreview verify \
      [OUTPUT_DIR]/[REPO_NAME]-model.json \
      --max-cycles [MAX_CYCLES] \
      --fail-on-violation=false \
      > [OUTPUT_DIR]/[REPO_NAME]-verification.txt

  Output:
    - Associativity verified: true | false
    - Identity verified: true | false
    - Cycle violations: N
    - Axiom errors: [list]

  Quality Gate:
    - Category axioms valid: true
    - Cycle count ≤ threshold: true
```

### Phase 5: Package Abstraction (Parallel)

```bash
For each file-level model:

  Agent: practical-programmer
  Task: Create package-level abstraction

  Execute:
    catreview abstract \
      [OUTPUT_DIR]/[REPO_NAME]-model.json \
      --output [OUTPUT_DIR]/[REPO_NAME]-packages.json \
      --pretty

  Output:
    - Package count: N_pkg
    - Package dependencies: M_pkg
    - Functor laws verified: true | false

  Quality Gate:
    - Abstraction successful: true
    - Functor laws hold: true
```

### Phase 6: Quality Assessment (Parallel)

```bash
For each repository analysis:

  Agent: MERCURIO
  Task: Multi-dimensional quality assessment

  Evaluate:
    Dimensions:
      1. Correctness (40%):
         - Category axioms verified
         - No extraction errors
         - Valid categorical structure

      2. Completeness (20%):
         - All source files analyzed
         - Dependencies captured
         - Metadata enriched

      3. Clarity (25%):
         - Report readability
         - Visualization quality
         - Documentation completeness

      4. Efficiency (15%):
         - Extraction performance
         - Analysis speed
         - Resource usage

  Formula:
    quality = 0.40×correctness + 0.25×clarity + 0.20×completeness + 0.15×efficiency

  Output:
    - Aggregate quality: Q ∈ [0, 1]
    - Per-dimension scores
    - Improvement recommendations

  Quality Gate:
    - Aggregate quality ≥ 0.90: true
```

### Phase 7: Comparative Analysis (Sequential - after all parallel complete)

```bash
Agent: MARS
Task: Multi-repository synthesis

Input:
  - All [REPO_NAME]-analysis.json files
  - All [REPO_NAME]-verification.txt results
  - All quality assessments

Execute:
  1. Load all analysis results
  2. Normalize metrics for comparison
  3. Identify patterns across repositories
  4. Generate comparative insights

  Comparative Metrics:
    - Complexity distribution: [min, median, max, stddev]
    - Cycle prevalence: [repos_with_cycles / total]
    - Architecture styles: [patterns detected]
    - Language impact: [go vs typescript vs python]
    - Size correlation: [complexity vs LOC]

Output:
  - Comparative report: comparative-analysis.md
  - Visualization: comparison-chart.json
  - Summary table: repo-comparison.csv
  - Insights: architectural-patterns.md

Quality Gate:
  - All repos included: true
  - Statistical validity: true
  - Actionable insights: ≥5
```

---

## RMP Loop - Iterative Refinement

### Iteration Logic

```python
def rmp_loop(repos, quality_threshold=0.90, max_iterations=3):
    iteration = 0
    results = {}

    while iteration < max_iterations:
        iteration += 1

        # Phase 1-6: Parallel execution
        parallel_results = execute_parallel([
            analyze_repo(repo) for repo in repos
        ])

        # Check quality for each repo
        quality_scores = {}
        for repo, result in parallel_results.items():
            quality = assess_quality(result)
            quality_scores[repo] = quality

            if quality < quality_threshold:
                print(f"⚠️  {repo}: Quality {quality:.2f} < {quality_threshold}")
                # Mark for re-analysis with refined parameters
                results[repo] = refine_analysis(repo, result, iteration)
            else:
                print(f"✅ {repo}: Quality {quality:.2f} ≥ {quality_threshold}")
                results[repo] = result

        # Check convergence
        if all(q >= quality_threshold for q in quality_scores.values()):
            print(f"🎯 Convergence achieved after {iteration} iterations")
            break

    # Phase 7: Comparative analysis
    comparative = synthesize_results(results)

    return {
        'individual_results': results,
        'comparative_analysis': comparative,
        'iterations': iteration,
        'final_quality': quality_scores
    }
```

### Refinement Strategies

If quality < threshold on iteration N, apply refinement:

```yaml
refinement_strategies:
  low_correctness:
    - Re-run extraction with verbose logging
    - Increase AST parsing depth
    - Manual verification of edge cases

  low_completeness:
    - Expand source directory search
    - Include vendor/third-party code
    - Add metadata enrichment

  low_clarity:
    - Enhance report formatting
    - Add more visualizations
    - Improve documentation

  low_efficiency:
    - Optimize extraction parameters
    - Cache intermediate results
    - Parallelize sub-operations
```

---

## Output Structure

```
examples/
├── comparative-analysis.md          # Cross-repo synthesis
├── comparison-chart.json            # Visualization data
├── repo-comparison.csv              # Tabular comparison
├── architectural-patterns.md        # Pattern analysis
│
├── crush/
│   ├── model.json                   # Categorical model
│   ├── analysis.json                # Complexity metrics
│   ├── verification.txt             # Axiom verification
│   ├── packages.json                # Package abstraction
│   ├── report.md                    # Human-readable report
│   └── quality-assessment.json      # Quality scores
│
├── bubbletea/
│   ├── model.json
│   ├── analysis.json
│   ├── verification.txt
│   ├── packages.json
│   ├── report.md
│   └── quality-assessment.json
│
├── soft-serve/
│   └── [same structure]
│
├── glow/
│   └── [same structure]
│
├── claude-code/
│   ├── LIMITATION.md               # Language not supported
│   └── alternative-analysis.md     # Manual analysis
│
└── lumina-ccn/
    └── [same structure]
```

---

## Example Report Template

### Individual Repository Report

```markdown
# Categorical Analysis: [REPO_NAME]

**Repository**: [URL]
**Commit**: [SHA]
**Analysis Date**: [ISO-8601]
**Language**: [Primary Language]
**Quality Score**: [Q]/1.0 ✅ | ⚠️

---

## Overview

| Metric | Value |
|--------|-------|
| Objects | [N] |
| Morphisms | [M] |
| Packages | [P] |
| Cycles | [C] |
| Diagram Complexity | [C_D] |
| Kolmogorov Complexity | [K] bytes |

---

## Complexity Analysis

### Basu-Isik Diagram Complexity: [C_D]

```
c(D) = Σc_obj + Σc_morph + c_comp
     = [X] + [Y] + [Z]
     = [C_D]
```

**Interpretation**: [Low | Medium | High] complexity for codebase of this size.

### Kolmogorov Estimate: [K] bytes

Compressed categorical model size: [K] bytes
**Interpretation**: [High | Medium | Low] architectural regularity (better compression = more regularity)

---

## Dependency Analysis

### Cycles Detected: [N]

[If N > 0, list cycles]

### Top 5 Most Unstable Components

| Component | Instability (I) | Ce | Ca |
|-----------|----------------|----|----|
| [name] | [I] | [Ce] | [Ca] |
...

**Instability**: I = Ce / (Ce + Ca) where 0 = stable, 1 = unstable

### Top 5 Most Coupled Components

| Component | Total Coupling | Ce | Ca |
|-----------|---------------|----|----|
| [name] | [Ce+Ca] | [Ce] | [Ca] |
...

---

## Category Axiom Verification

- ✅ Associativity: (h ∘ g) ∘ f = h ∘ (g ∘ f)
- ✅ Identity: f ∘ id_A = f and id_B ∘ f = f

[If violations, list them]

---

## Package-Level Abstraction

Via PackageAbstractionFunctor F: FileCat → PkgCat

| Metric | File-Level | Package-Level |
|--------|-----------|---------------|
| Objects | [N_files] | [N_packages] |
| Morphisms | [M_files] | [M_packages] |
| Cycles | [C_files] | [C_packages] |

**Functor Laws**:
- ✅ Composition: F(g ∘ f) = F(g) ∘ F(f)
- ✅ Identity: F(id_A) = id_{F(A)}

---

## Quality Assessment

| Dimension | Score | Weight | Contribution |
|-----------|-------|--------|--------------|
| Correctness | [0-1] | 40% | [X] |
| Clarity | [0-1] | 25% | [Y] |
| Completeness | [0-1] | 20% | [Z] |
| Efficiency | [0-1] | 15% | [W] |
| **Aggregate** | **[Q]** | **100%** | **[Q]** |

[If Q < 0.90, list improvement recommendations]

---

## Architectural Insights

[Pattern detection based on categorical structure]

- **Pattern 1**: [Description]
- **Pattern 2**: [Description]
- **Recommendations**: [Actionable items]

---

*Generated by catreview-go v1.0 with categorical meta-prompting*
```

### Comparative Report Template

```markdown
# Categorical Codebase Comparison

**Repositories Analyzed**: [N]
**Analysis Date**: [ISO-8601]
**Aggregate Quality**: [Q_avg]/1.0

---

## Overview Table

| Repository | Language | Objects | Morphisms | Cycles | Complexity | Quality |
|-----------|----------|---------|-----------|--------|------------|---------|
| [name] | [lang] | [N] | [M] | [C] | [C_D] | [Q] |
...

---

## Complexity Distribution

```
Min:    [C_min]
Q1:     [C_q1]
Median: [C_median]
Q3:     [C_q3]
Max:    [C_max]
StdDev: [σ]
```

**Insight**: [Interpretation of distribution]

---

## Cross-Repository Patterns

### Pattern 1: [Name]

**Observed in**: [repo1, repo2, ...]
**Categorical Structure**: [Description]
**Interpretation**: [What this means architecturally]

### Pattern 2: [Name]

...

---

## Language Impact Analysis

| Language | Avg Complexity | Avg Cycles | Avg Coupling |
|----------|---------------|-----------|--------------|
| Go | [C_go] | [Cy_go] | [Co_go] |
| TypeScript | [C_ts] | [Cy_ts] | [Co_ts] |
| Python | [C_py] | [Cy_py] | [Co_py] |

**Insight**: [Language-specific observations]

---

## Size vs Complexity Correlation

```
Correlation coefficient (Objects vs Complexity): r = [r]
```

**Insight**: [Linear | Non-linear | No correlation]

---

## Recommendations

### Architectural Best Practices (from high-quality repos)

1. [Practice from repo X]
2. [Practice from repo Y]

### Anti-Patterns (from problematic repos)

1. [Issue in repo Z]
2. [Issue in repo W]

---

*Synthesized by MARS agent using categorical colimit Σ: [Result] → ComparativeReport*
```

---

## Execution Command

### Using Ralph Loop

```bash
/ralph --rmp @quality:0.90 "
Execute categorical repository analysis using RMP meta-prompt:

Parameters:
  - Repositories: crush, bubbletea, soft-serve, glow, claude-code, lumina-ccn
  - Output: examples/
  - Quality threshold: 0.90
  - Parallel execution: true

Phases:
  1. Clone all 6 repos in parallel (git-genius agent × 6)
  2. Extract categorical models in parallel (practical-programmer × 6)
  3. Analyze complexity in parallel (practical-programmer × 6)
  4. Verify axioms in parallel (practical-programmer × 6)
  5. Create package abstractions in parallel (practical-programmer × 6)
  6. Quality assessment in parallel (MERCURIO × 6)
  7. Comparative synthesis (MARS × 1)

Quality Gates:
  - Per-repo quality ≥ 0.90
  - All category axioms verified
  - Cycles ≤ 5 per repo
  - Comparative analysis complete

Stop Condition:
  - All repos analyzed with Q ≥ 0.90
  - Comparative report generated
  - Architectural patterns identified
"
```

### Direct Execution (without Ralph)

```bash
# Set parameters
export REPOS="crush bubbletea soft-serve glow claude-code lumina-ccn"
export OUTPUT_DIR="examples"
export QUALITY_THRESHOLD=0.90

# Execute meta-prompt
/meta @mode:iterative @tier:L4 @quality:0.90 "
Analyze repositories: $REPOS
Output to: $OUTPUT_DIR
Generate comparative analysis
"
```

---

## Quality Gates Summary

| Phase | Gate | Threshold | Action on Failure |
|-------|------|-----------|-------------------|
| Clone | Repo accessible | 100% | Retry with auth |
| Extract | Model valid | Objects > 0 | Re-run with verbose |
| Analyze | Complexity > 0 | Always | Debug extractor |
| Verify | Axioms hold | Required | Investigate violations |
| Abstract | Functor laws | Required | Check mapping logic |
| Quality | Aggregate ≥ threshold | 0.90 | Refine and re-run |
| Comparative | All repos included | 100% | Wait for stragglers |

---

## Error Handling

### Language Not Supported

For repositories in unsupported languages (TypeScript, Python):

```yaml
fallback_strategy:
  - Create placeholder report
  - Document limitation
  - Suggest alternative analysis methods
  - Include in comparative analysis with caveats

example:
  - claude-code (TypeScript):
      status: "Language not supported by catreview-go v1.0"
      alternative: "Manual architectural review recommended"
      future: "TypeScript extractor planned for v1.1"
```

### Extraction Failures

```yaml
retry_strategy:
  max_retries: 3
  backoff: exponential
  fallback:
    - Try alternative source directory
    - Reduce extraction scope
    - Manual intervention required
```

---

## Success Criteria

All criteria must be met for Ralph loop to exit:

- ✅ All 6 repositories cloned successfully
- ✅ Categorical models extracted (or documented limitation)
- ✅ Complexity analysis complete for all supported languages
- ✅ Category axioms verified
- ✅ Quality scores ≥ 0.90 for each analysis
- ✅ Comparative analysis generated
- ✅ Architectural patterns identified
- ✅ Reports published to examples/ directory

---

## Reusability

### To Use with Different Repos

1. Edit the `repositories` section in the YAML config
2. Adjust `quality_gates` if needed
3. Run the same Ralph loop command

### To Extend

- Add new extractors for languages (Java, Python, Rust)
- Add new metrics (cohesion, modularity)
- Add new visualizations (D3.js, Graphviz)
- Integrate with CI/CD pipelines

---

*Meta-Prompt Version: 1.0*
*Compatible with: catreview-go v1.0+*
*Ralph Loop: Enabled*
*Quality Threshold: ≥0.90*
