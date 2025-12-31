#!/usr/bin/env python3
"""
Categorical Python Code Extractor
Extracts categorical models (Objects + Morphisms) from Python codebases.
Output format matches catreview-go's Category JSON structure.
"""

import ast
import json
import sys
from pathlib import Path
from typing import Dict, List, Set, Any
import time

class PythonCategoricalExtractor:
    def __init__(self):
        self.objects = {}
        self.morphisms = {}
        self.file_count = 0
        self.errors = []

    def extract_from_directory(self, root_path: Path, project_name: str) -> Dict:
        """Extract categorical model from a Python project."""
        print(f"�� Analyzing: {root_path}")
        print(f"📊 Project: {project_name}\n")

        start_time = time.time()

        # Find all Python files
        py_files = list(root_path.rglob("*.py"))
        self.file_count = len(py_files)

        print(f"📈 Python Files Found: {self.file_count}\n")
        print(f"🔍 Extracting categorical model...")

        # Process each file
        for py_file in py_files:
            try:
                self._extract_from_file(py_file, root_path)
            except Exception as e:
                self.errors.append(f"{py_file}: {str(e)}")

        # Add identity morphisms
        self._add_identity_morphisms()

        elapsed = time.time() - start_time
        print(f"✅ Extraction completed in {elapsed:.2f}s\n")

        # Print statistics
        self._print_statistics()

        return {
            "name": project_name,
            "objects": self.objects,
            "morphisms": self.morphisms
        }

    def _extract_from_file(self, file_path: Path, root: Path):
        """Extract categorical structures from a single Python file."""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                source = f.read()

            tree = ast.parse(source, filename=str(file_path))
            module_name = self._path_to_module_name(file_path, root)

            # Create module object
            module_id = module_name
            self.objects[module_id] = {
                "id": module_id,
                "type": "module",
                "name": file_path.stem,
                "metadata": {
                    "file": str(file_path.relative_to(root)),
                    "language": "python"
                }
            }

            # Extract top-level constructs
            visitor = CategoricalVisitor(module_name, self.objects, self.morphisms)
            visitor.visit(tree)

        except SyntaxError:
            pass  # Skip files with syntax errors
        except Exception as e:
            pass  # Skip files with other errors

    def _path_to_module_name(self, file_path: Path, root: Path) -> str:
        """Convert file path to Python module name."""
        try:
            rel_path = file_path.relative_to(root)
            parts = list(rel_path.parts)
            if parts[-1] == '__init__.py':
                parts = parts[:-1]
            else:
                parts[-1] = parts[-1].replace('.py', '')
            return '.'.join(parts).lstrip('.')
        except:
            return str(file_path.stem)

    def _add_identity_morphisms(self):
        """Add identity morphisms for all objects."""
        for obj_id in self.objects:
            morph_id = f"id_{obj_id}"
            self.morphisms[morph_id] = {
                "id": morph_id,
                "source": obj_id,
                "target": obj_id,
                "type": "identity"
            }

    def _print_statistics(self):
        """Print extraction statistics."""
        print("═══════════════════════════════════════════════════════════════")
        print("  Categorical Model Statistics")
        print("═══════════════════════════════════════════════════════════════\n")

        print(f"📦 Total Objects: {len(self.objects):,}")
        print(f"🔗 Total Morphisms: {len(self.morphisms):,}")
        if len(self.objects) > 0:
            print(f"📐 Ratio: {len(self.morphisms)/len(self.objects):.2f} morphisms per object\n")

        # Count by type
        obj_types = {}
        for obj in self.objects.values():
            obj_types[obj["type"]] = obj_types.get(obj["type"], 0) + 1

        print("📊 Objects by Type:")
        for obj_type, count in sorted(obj_types.items(), key=lambda x: -x[1]):
            pct = (count / len(self.objects)) * 100 if len(self.objects) > 0 else 0
            print(f"  • {obj_type:15s}: {count:5,} ({pct:5.1f}%)")

        morph_types = {}
        for morph in self.morphisms.values():
            morph_types[morph["type"]] = morph_types.get(morph["type"], 0) + 1

        print("\n🔗 Morphisms by Type:")
        for morph_type, count in sorted(morph_types.items(), key=lambda x: -x[1]):
            pct = (count / len(self.morphisms)) * 100 if len(self.morphisms) > 0 else 0
            print(f"  • {morph_type:15s}: {count:5,} ({pct:5.1f}%)")

        print(f"\n⚠️  Errors: {len(self.errors)}")
        if self.errors and len(self.errors) < 10:
            for error in self.errors[:5]:
                print(f"    {error}")

class CategoricalVisitor(ast.NodeVisitor):
    """AST visitor that extracts categorical structures."""

    def __init__(self, module_name: str, objects: Dict, morphisms: Dict):
        self.module_name = module_name
        self.objects = objects
        self.morphisms = morphisms
        self.current_class = None

    def visit_ClassDef(self, node):
        """Extract class as categorical object."""
        class_id = f"{self.module_name}.{node.name}"

        # Create class object
        self.objects[class_id] = {
            "id": class_id,
            "type": "class",
            "name": node.name,
            "metadata": {
                "module": self.module_name,
                "line": node.lineno,
                "language": "python",
                "bases": [self._get_name(base) for base in node.bases]
            }
        }

        # Create defines morphism
        self._add_morphism(self.module_name, class_id, "defines")

        # Create inheritance morphisms
        for base in node.bases:
            base_name = self._get_name(base)
            if base_name and '.' not in base_name:
                # Try to qualify with module
                base_id = f"{self.module_name}.{base_name}"
            else:
                base_id = base_name

            if base_id:
                self._add_morphism(class_id, base_id, "inheritance")

        # Visit methods
        old_class = self.current_class
        self.current_class = class_id
        self.generic_visit(node)
        self.current_class = old_class

    def visit_FunctionDef(self, node):
        """Extract function/method as categorical object."""
        if self.current_class:
            func_id = f"{self.current_class}.{node.name}"
            parent_id = self.current_class
        else:
            func_id = f"{self.module_name}.{node.name}"
            parent_id = self.module_name

        # Create function object
        self.objects[func_id] = {
            "id": func_id,
            "type": "function",
            "name": node.name,
            "metadata": {
                "module": self.module_name,
                "line": node.lineno,
                "language": "python"
            }
        }

        # Create defines morphism
        self._add_morphism(parent_id, func_id, "defines")

        self.generic_visit(node)

    def visit_Import(self, node):
        """Extract import as morphism."""
        for alias in node.names:
            imported_id = f"import:{alias.name}"

            # Create imported module object
            if imported_id not in self.objects:
                self.objects[imported_id] = {
                    "id": imported_id,
                    "type": "imported_module",
                    "name": alias.name,
                    "metadata": {"language": "python"}
                }

            # Create import morphism
            self._add_morphism(self.module_name, imported_id, "import")

    def visit_ImportFrom(self, node):
        """Extract from-import as morphism."""
        if node.module:
            imported_id = f"import:{node.module}"

            # Create imported module object
            if imported_id not in self.objects:
                self.objects[imported_id] = {
                    "id": imported_id,
                    "type": "imported_module",
                    "name": node.module,
                    "metadata": {"language": "python"}
                }

            # Create import morphism
            self._add_morphism(self.module_name, imported_id, "import")

    def _get_name(self, node):
        """Extract name from AST node."""
        if isinstance(node, ast.Name):
            return node.id
        elif isinstance(node, ast.Attribute):
            return f"{self._get_name(node.value)}.{node.attr}"
        return None

    def _add_morphism(self, source: str, target: str, morph_type: str):
        """Add a morphism between objects."""
        morph_id = f"{morph_type}:{source}->{target}"

        if morph_id not in self.morphisms:
            self.morphisms[morph_id] = {
                "id": morph_id,
                "source": source,
                "target": target,
                "type": morph_type
            }

def main():
    if len(sys.argv) < 3:
        print("Usage: python3 python_categorical_extractor.py <project_path> <project_name>")
        print("Example: python3 python_categorical_extractor.py /path/to/project my-project")
        sys.exit(1)

    project_path = Path(sys.argv[1])
    project_name = sys.argv[2]

    if not project_path.exists():
        print(f"❌ Error: Path does not exist: {project_path}")
        sys.exit(1)

    # Extract categorical model
    extractor = PythonCategoricalExtractor()
    category = extractor.extract_from_directory(project_path, project_name)

    # Save to JSON
    output_file = f"{project_name}-analysis.json"
    print(f"\n💾 Exporting to: {output_file}")

    with open(output_file, 'w') as f:
        json.dump(category, f, indent=2)

    print(f"✅ Export successful\n")

    # Save summary
    summary_file = f"{project_name}-summary.txt"
    with open(summary_file, 'w') as f:
        f.write("═══════════════════════════════════════════════════════════════\n")
        f.write("  Categorical Code Analysis Summary\n")
        f.write("═══════════════════════════════════════════════════════════════\n\n")
        f.write(f"Project: {project_path}\n")
        f.write(f"Python Files: {extractor.file_count}\n\n")
        f.write(f"Objects: {len(category['objects'])}\n")
        f.write(f"Morphisms: {len(category['morphisms'])}\n")
        if len(category['objects']) > 0:
            f.write(f"Ratio: {len(category['morphisms'])/len(category['objects']):.2f} morphisms per object\n")

    print(f"✅ Summary saved to: {summary_file}\n")

if __name__ == "__main__":
    main()
