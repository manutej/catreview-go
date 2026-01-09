#!/usr/bin/env python3
"""Scan LUXOR/PROJECTS/ for Python projects and count files."""

import os
from pathlib import Path

projects_dir = Path("/Users/manu/Documents/LUXOR/PROJECTS")

projects = []
for item in sorted(projects_dir.iterdir()):
    if item.is_dir() and not item.name.startswith('.'):
        py_files = list(item.rglob("*.py"))
        if py_files:
            projects.append({
                'name': item.name,
                'path': str(item),
                'py_files': len(py_files)
            })

# Sort by file count descending
projects.sort(key=lambda x: x['py_files'], reverse=True)

print("=" * 80)
print(f"Python Projects in LUXOR/PROJECTS/")
print("=" * 80)
print()

total_files = 0
for i, proj in enumerate(projects, 1):
    print(f"{i:2d}. {proj['name']:40s} {proj['py_files']:6,} Python files")
    total_files += proj['py_files']

print()
print("=" * 80)
print(f"Total: {len(projects)} projects with {total_files:,} Python files")
print("=" * 80)
print()

# Already analyzed
analyzed = {'BARQUE', 'ai-dialogue'}
print("Already Analyzed:")
for proj in projects:
    if proj['name'] in analyzed:
        print(f"  ✅ {proj['name']:40s} {proj['py_files']:6,} files")
print()

print("Remaining to Analyze:")
remaining = [p for p in projects if p['name'] not in analyzed]
for proj in remaining:
    print(f"  ⏳ {proj['name']:40s} {proj['py_files']:6,} files")
print()

print(f"Progress: {len(analyzed)}/{len(projects)} projects analyzed ({len(analyzed)/len(projects)*100:.1f}%)")
