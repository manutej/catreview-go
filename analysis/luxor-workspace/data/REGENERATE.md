# Regenerating Large JSON Files

## Why Some Files Are Missing

Four large JSON files (>50 MB) are excluded from Git to avoid GitHub's file size limits:
- `ai-dialogue-analysis.json` (97 MB)
- `barque-analysis.json` (67 MB)
- `hekat-analysis.json` (64 MB)
- `hyperglyph-analysis.json` (64 MB)

These files can be regenerated locally using the extraction scripts.

---

## How to Regenerate

### Prerequisites
- Python 3.7+
- Access to LUXOR workspace at `/Users/manu/Documents/LUXOR/`

### Regenerate Individual Project

```bash
cd /Users/manu/Documents/LUXOR/catreview-go

# Regenerate ai-dialogue
python3 analysis/luxor-workspace/scripts/python_categorical_extractor.py \
    /Users/manu/Documents/LUXOR/PROJECTS/ai-dialogue \
    ai-dialogue

# Move output to data directory
mv ai-dialogue-analysis.json analysis/luxor-workspace/data/

# Repeat for other projects
python3 analysis/luxor-workspace/scripts/python_categorical_extractor.py \
    /Users/manu/Documents/LUXOR/PROJECTS/BARQUE \
    barque

python3 analysis/luxor-workspace/scripts/python_categorical_extractor.py \
    /Users/manu/Documents/LUXOR/PROJECTS/hekat \
    hekat

python3 analysis/luxor-workspace/scripts/python_categorical_extractor.py \
    /Users/manu/Documents/LUXOR/PROJECTS/hyperglyph \
    hyperglyph
```

### Regenerate All Missing Files

```bash
cd /Users/manu/Documents/LUXOR/catreview-go/analysis/luxor-workspace/scripts

# Run batch analysis
./batch_analyze.sh
```

---

## Verification

After regeneration, verify file sizes:

```bash
cd analysis/luxor-workspace/data
ls -lh *-analysis.json | awk '{print $5, $9}'
```

Expected output:
```
97M ai-dialogue-analysis.json
67M barque-analysis.json
64M hekat-analysis.json
64M hyperglyph-analysis.json
35M categorical-codebase-review-analysis.json
23M nanobanana-repo-analysis.json
22M HALCON-analysis.json
17M LUMOS-analysis.json
11M docrag-analysis.json
```

---

## Files Included in Git

These smaller files (<50 MB) ARE included in the repository:
- ✅ `categorical-codebase-review-analysis.json` (35 MB)
- ✅ `nanobanana-repo-analysis.json` (23 MB)
- ✅ `HALCON-analysis.json` (22 MB)
- ✅ `LUMOS-analysis.json` (17 MB)
- ✅ `docrag-analysis.json` (11 MB)

---

## Alternative: Git LFS

If you have Git LFS installed, you can track the large files:

```bash
# Install Git LFS (macOS)
brew install git-lfs
git lfs install

# Track large JSON files
git lfs track "analysis/luxor-workspace/data/*-analysis.json"

# Commit and push
git add .gitattributes
git commit -m "Track large JSON files with Git LFS"
git push
```

With Git LFS, all files will be available in the repository.

---

**Note**: The analysis documentation, visualizations, and summaries are all included in Git regardless of size. Only the raw JSON categorical models are excluded.
