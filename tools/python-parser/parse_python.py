#!/usr/bin/env python3
"""
Python AST Parser for catreview

Extracts categorical structures from Python source code:
- Modules/Packages → Module objects
- Classes → Class objects
- Functions/Methods → Function objects
- Imports → Import morphisms
- Calls → Call morphisms
- Inheritance → Inheritance morphisms
- Type hints → Type dependency morphisms

Outputs JSON to stdout for consumption by Go PythonExtractor.
"""

import ast
import json
import sys
from typing import Any, Dict, List, Optional


class PythonCategoricalExtractor(ast.NodeVisitor):
    """
    Extracts categorical model from Python AST.

    Uses the Visitor pattern to traverse the AST and collect:
    - Class declarations (with inheritance and decorators)
    - Function/method declarations (with type hints and decorators)
    - Import statements (import and from...import)
    - Function calls (for morphism creation)
    """

    def __init__(self, filename: str):
        self.filename = filename
        self.module_name = self._compute_module_name(filename)

        # Collected structures
        self.imports: List[Dict[str, Any]] = []
        self.classes: List[Dict[str, Any]] = []
        self.functions: List[Dict[str, Any]] = []
        self.calls: List[Dict[str, Any]] = []

        # Context tracking
        self.current_class: Optional[str] = None
        self.current_function: Optional[str] = None

    def _compute_module_name(self, filename: str) -> str:
        """Convert file path to module name: path/to/file.py → path.to.file"""
        return filename.replace('.py', '').replace('/', '.').replace('\\', '.')

    def _get_name(self, node: ast.AST) -> str:
        """Extract name from AST node (handles Name, Attribute, etc.)"""
        if isinstance(node, ast.Name):
            return node.id
        elif isinstance(node, ast.Attribute):
            value = self._get_name(node.value)
            return f"{value}.{node.attr}" if value else node.attr
        elif isinstance(node, ast.Call):
            return self._get_name(node.func)
        elif isinstance(node, ast.Constant):
            return str(node.value)
        else:
            return ""

    def _extract_type_hint(self, annotation: Optional[ast.expr]) -> Optional[str]:
        """Extract type hint from annotation node"""
        if annotation is None:
            return None
        return self._get_name(annotation)

    def visit_Import(self, node: ast.Import) -> None:
        """Extract: import module [as alias]"""
        for alias in node.names:
            self.imports.append({
                "type": "import",
                "name": alias.name,
                "asname": alias.asname,
                "line": node.lineno,
            })
        self.generic_visit(node)

    def visit_ImportFrom(self, node: ast.ImportFrom) -> None:
        """Extract: from module import name [as alias]"""
        module = node.module if node.module else ""
        for alias in node.names:
            self.imports.append({
                "type": "from_import",
                "module": module,
                "name": alias.name,
                "asname": alias.asname,
                "line": node.lineno,
            })
        self.generic_visit(node)

    def visit_ClassDef(self, node: ast.ClassDef) -> None:
        """Extract class declarations with bases, decorators, and methods"""
        # Save previous context
        prev_class = self.current_class
        self.current_class = node.name

        class_info = {
            "name": node.name,
            "bases": [self._get_name(base) for base in node.bases],
            "decorators": [self._get_name(d) for d in node.decorator_list],
            "methods": [],
            "line": node.lineno,
            "docstring": ast.get_docstring(node),
        }

        # Extract methods
        for item in node.body:
            if isinstance(item, ast.FunctionDef) or isinstance(item, ast.AsyncFunctionDef):
                method_info = self._extract_function_info(item, is_method=True)
                class_info["methods"].append(method_info)

        self.classes.append(class_info)

        # Visit children (for nested classes, etc.)
        self.generic_visit(node)

        # Restore context
        self.current_class = prev_class

    def visit_FunctionDef(self, node: ast.FunctionDef) -> None:
        """Extract top-level function declarations"""
        # Only process if not already handled as a method
        if self.current_class is None:
            func_info = self._extract_function_info(node, is_method=False)
            self.functions.append(func_info)

        # Save context and visit children
        prev_function = self.current_function
        self.current_function = node.name
        self.generic_visit(node)
        self.current_function = prev_function

    def visit_AsyncFunctionDef(self, node: ast.AsyncFunctionDef) -> None:
        """Extract async function declarations"""
        if self.current_class is None:
            func_info = self._extract_function_info(node, is_method=False)
            func_info["async"] = True
            self.functions.append(func_info)

        prev_function = self.current_function
        self.current_function = node.name
        self.generic_visit(node)
        self.current_function = prev_function

    def _extract_function_info(self, node: ast.FunctionDef, is_method: bool) -> Dict[str, Any]:
        """Extract function/method information including parameters and type hints"""
        # Extract parameters with type hints
        parameters = []
        for arg in node.args.args:
            param = {
                "name": arg.arg,
                "type": self._extract_type_hint(arg.annotation),
            }
            parameters.append(param)

        return {
            "name": node.name,
            "parameters": parameters,
            "return_type": self._extract_type_hint(node.returns),
            "decorators": [self._get_name(d) for d in node.decorator_list],
            "line": node.lineno,
            "is_method": is_method,
            "async": isinstance(node, ast.AsyncFunctionDef),
            "docstring": ast.get_docstring(node),
        }

    def visit_Call(self, node: ast.Call) -> None:
        """Extract function calls for morphism creation"""
        call_info = {
            "func": self._get_name(node.func),
            "line": node.lineno,
            "context_class": self.current_class,
            "context_function": self.current_function,
        }
        self.calls.append(call_info)
        self.generic_visit(node)

    def extract(self, source_code: str) -> Dict[str, Any]:
        """
        Parse Python source code and extract categorical structures.

        Returns:
            Dictionary with module, imports, classes, functions, and calls
        """
        try:
            tree = ast.parse(source_code, filename=self.filename)
            self.visit(tree)

            return {
                "module": self.module_name,
                "file": self.filename,
                "imports": self.imports,
                "classes": self.classes,
                "functions": self.functions,
                "calls": self.calls,
            }
        except SyntaxError as e:
            return {
                "error": f"Syntax error in {self.filename}: {e}",
                "module": self.module_name,
                "file": self.filename,
            }


def main():
    """
    CLI entry point for Python AST parser.

    Usage: python3 parse_python.py <file.py>
    Outputs: JSON to stdout
    """
    if len(sys.argv) < 2:
        print(json.dumps({"error": "Usage: parse_python.py <file.py>"}), file=sys.stderr)
        sys.exit(1)

    filepath = sys.argv[1]

    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            source_code = f.read()
    except FileNotFoundError:
        print(json.dumps({"error": f"File not found: {filepath}"}), file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(json.dumps({"error": f"Failed to read {filepath}: {e}"}), file=sys.stderr)
        sys.exit(1)

    extractor = PythonCategoricalExtractor(filepath)
    result = extractor.extract(source_code)

    # Output JSON to stdout (consumed by Go PythonExtractor)
    print(json.dumps(result, indent=2))


if __name__ == "__main__":
    main()
