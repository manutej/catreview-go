#!/bin/bash
# Batch analyze all Python projects in LUXOR/PROJECTS/

set -e

PROJECTS_DIR="/Users/manu/Documents/LUXOR/PROJECTS"
OUTPUT_DIR="/Users/manu/Documents/LUXOR/catreview-go"

cd "$OUTPUT_DIR"

echo "════════════════════════════════════════════════════════════════"
echo "  LUXOR Workspace Categorical Analysis - Batch Mode"
echo "════════════════════════════════════════════════════════════════"
echo ""

# Projects to analyze (ordered by size, descending)
declare -a PROJECTS=(
    "hekat:3102"
    "hyperglyph:2183"
    "nanobanana-repo:1784"
    "HALCON:940"
    "LUMOS:680"
    "docrag:418"
    "discopy:116"
    "paper2agent:22"
)

TOTAL_PROJECTS=${#PROJECTS[@]}
COMPLETED=0
FAILED=0

START_TIME=$(date +%s)

for project_info in "${PROJECTS[@]}"; do
    IFS=':' read -r project expected_files <<< "$project_info"

    COMPLETED=$((COMPLETED + 1))

    echo "────────────────────────────────────────────────────────────────"
    echo "[$COMPLETED/$TOTAL_PROJECTS] Analyzing: $project ($expected_files files)"
    echo "────────────────────────────────────────────────────────────────"
    echo ""

    PROJECT_PATH="$PROJECTS_DIR/$project"

    if [ ! -d "$PROJECT_PATH" ]; then
        echo "⚠️  Directory not found: $PROJECT_PATH"
        echo "   Skipping..."
        echo ""
        FAILED=$((FAILED + 1))
        continue
    fi

    # Check if already analyzed
    if [ -f "${project}-analysis.json" ]; then
        echo "✅ Already analyzed: ${project}-analysis.json exists"
        echo "   Skipping..."
        echo ""
        continue
    fi

    # Run analysis with timeout (5 minutes max per project)
    echo "🔍 Running categorical extraction..."
    if timeout 300 go run examples/python/analyze_project.go \
        "$PROJECT_PATH" \
        "$project" 2>&1 | tee "${project}-analysis.log"; then

        echo ""
        echo "✅ Analysis complete: ${project}-analysis.json"

        # Generate visualizations (with shorter timeout)
        echo "📊 Generating visualizations..."
        if timeout 120 go run examples/python/visualize_project.go \
            --input "${project}-analysis.json" \
            --output "$project" \
            --max-nodes 80 2>&1 | tee -a "${project}-analysis.log"; then

            # Convert to SVG
            echo "🎨 Converting to SVG..."
            for dot_file in ${project}-*.dot; do
                if [ -f "$dot_file" ]; then
                    svg_file="${dot_file%.dot}.svg"
                    dot -Tsvg "$dot_file" -o "$svg_file" 2>/dev/null && \
                        echo "   ✅ $svg_file"
                fi
            done
            echo ""
        else
            echo "⚠️  Visualization generation timed out or failed"
            echo ""
        fi
    else
        echo "❌ Analysis failed or timed out for: $project"
        echo "   Check ${project}-analysis.log for details"
        echo ""
        FAILED=$((FAILED + 1))
    fi

    # Progress update
    ELAPSED=$(($(date +%s) - START_TIME))
    echo "Progress: $COMPLETED/$TOTAL_PROJECTS projects processed"
    echo "Time elapsed: ${ELAPSED}s"
    echo ""
done

END_TIME=$(date +%s)
TOTAL_TIME=$((END_TIME - START_TIME))

echo "════════════════════════════════════════════════════════════════"
echo "  Batch Analysis Complete"
echo "════════════════════════════════════════════════════════════════"
echo ""
echo "Total projects: $TOTAL_PROJECTS"
echo "Completed: $((TOTAL_PROJECTS - FAILED))"
echo "Failed: $FAILED"
echo "Total time: ${TOTAL_TIME}s ($(($TOTAL_TIME / 60))m $(($TOTAL_TIME % 60))s)"
echo ""
echo "Next steps:"
echo "  1. Review: ls -lh *-analysis.json"
echo "  2. Open visualizations: open *-calls.svg"
echo "  3. Generate comparative report: python3 generate_comparative_report.py"
echo ""
