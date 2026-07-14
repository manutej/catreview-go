package extractor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/manu/catreview/pkg/category"
)

// writeJava writes src to <dir>/<rel>, creating parent directories.
func writeJava(t *testing.T, dir, rel, src string) string {
	t.Helper()
	full := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(full, []byte(src), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	return full
}

// objByName returns the first Object with the given simple name and type.
func objByName(cat *category.Category, name, typ string) *category.Object {
	for _, o := range cat.Objects() {
		if o.Name == name && o.Type == typ {
			return o
		}
	}
	return nil
}

func hasMorphism(cat *category.Category, source, target, typ string) bool {
	for _, m := range cat.Morphisms() {
		if m.Source == source && m.Target == target && m.Type == typ {
			return true
		}
	}
	return false
}

const shapeJava = `
package com.example.geometry;

import java.util.List;
import java.util.ArrayList;

/**
 * A geometric shape.
 */
public interface Shape {
    double area();
}
`

const circleJava = `
package com.example.geometry;

import java.lang.Math;

/** A circle. */
public class Circle implements Shape {
    private final double radius;

    public Circle(double radius) {
        this.radius = radius;
    }

    @Override
    public double area() {
        return Math.PI * squared(radius);
    }

    private double squared(double x) {
        return x * x;
    }
}
`

func TestJavaExtractor_BasicShape(t *testing.T) {
	dir := t.TempDir()
	writeJava(t, dir, "com/example/geometry/Shape.java", shapeJava)
	writeJava(t, dir, "com/example/geometry/Circle.java", circleJava)

	ex := NewJavaExtractor()
	cat, err := ex.ExtractFromPath(dir)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}

	// Interface + class objects present with correct kinds.
	shape := objByName(cat, "Shape", "interface")
	if shape == nil {
		t.Fatal("Shape interface object not found")
	}
	if shape.ID != "com.example.geometry.Shape" {
		t.Errorf("Shape id = %q, want com.example.geometry.Shape", shape.ID)
	}
	circle := objByName(cat, "Circle", "class")
	if circle == nil {
		t.Fatal("Circle class object not found")
	}

	// Package metadata propagated.
	if got := shape.Metadata["package"]; got != "com.example.geometry" {
		t.Errorf("Shape package = %v, want com.example.geometry", got)
	}
	// Javadoc summary captured.
	if got := circle.Metadata["doc"]; got != "A circle." {
		t.Errorf("Circle doc = %q, want %q", got, "A circle.")
	}

	// Inheritance: Circle implements Shape, resolved to the declared object.
	if !hasMorphism(cat, "com.example.geometry.Circle", "com.example.geometry.Shape", "inheritance") {
		t.Error("expected inheritance morphism Circle -> Shape")
	}

	// Methods: area() and squared() and the constructor exist.
	area := objByName(cat, "area", "function")
	if area == nil {
		t.Fatal("area method object not found")
	}
	if got := area.Metadata["class"]; got != "com.example.geometry.Circle" && got != "com.example.geometry.Shape" {
		t.Errorf("area class metadata unexpected: %v", got)
	}
	if objByName(cat, "squared", "function") == nil {
		t.Error("squared method object not found")
	}
	if ctor := objByName(cat, "Circle", "function"); ctor == nil {
		t.Error("Circle constructor object not found")
	} else if ctor.Metadata["kind"] != "constructor" {
		t.Errorf("Circle ctor kind = %v, want constructor", ctor.Metadata["kind"])
	}

	// Same-class call: area() calls squared().
	areaID := "com.example.geometry.Circle.area"
	squaredID := "com.example.geometry.Circle.squared"
	if !hasMorphism(cat, areaID, squaredID, "function_call") {
		t.Errorf("expected function_call %s -> %s", areaID, squaredID)
	}

	// Imports became import morphisms.
	if !hasMorphism(cat, filepath.Join(dir, "com/example/geometry/Circle.java"), "import:java.lang.Math", "import") {
		t.Error("expected import morphism for java.lang.Math")
	}

	// The extracted category satisfies category axioms.
	if err := cat.VerifyAxioms(); err != nil {
		t.Errorf("category axioms failed: %v", err)
	}
}

const nestedJava = `
package app;

public class Outer extends Base {
    static int TABLE[] = {1, 2, 3};

    public enum Color { RED, GREEN, BLUE }

    private static class Inner implements Runnable {
        public void run() {
            if (true) { System.out.println("x"); }
        }
    }

    record Point(int x, int y) {}
}

@interface Marker {}
`

func TestJavaExtractor_NestedAndKinds(t *testing.T) {
	dir := t.TempDir()
	writeJava(t, dir, "app/Outer.java", nestedJava)

	ex := NewJavaExtractor()
	cat, err := ex.ExtractFromPath(dir)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}

	cases := []struct {
		id   string
		kind string
	}{
		{"app.Outer", "class"},
		{"app.Outer.Color", "enum"},
		{"app.Outer.Inner", "class"},
		{"app.Outer.Point", "record"},
		{"app.Marker", "annotation"},
	}
	for _, c := range cases {
		obj, ok := cat.GetObject(c.id)
		if !ok {
			t.Errorf("missing object %s", c.id)
			continue
		}
		if obj.Type != c.kind {
			t.Errorf("%s kind = %q, want %q", c.id, obj.Type, c.kind)
		}
	}

	// Array initializer braces must not be misread as a type/method.
	if _, ok := cat.GetObject("app.Outer.TABLE"); ok {
		t.Error("array field TABLE should not be an object")
	}

	// extends Base resolves to an external_type placeholder (Base is undefined).
	if !hasMorphism(cat, "app.Outer", "external_type:Base", "inheritance") {
		t.Error("expected inheritance Outer -> external_type:Base")
	}

	// run() is a real method; the if-block inside it is not.
	if objByName(cat, "run", "function") == nil {
		t.Error("run method not found")
	}
	for _, o := range cat.Objects() {
		if o.Name == "if" || o.Name == "true" {
			t.Errorf("control-flow leaked as object: %s (%s)", o.ID, o.Type)
		}
	}

	if err := cat.VerifyAxioms(); err != nil {
		t.Errorf("category axioms failed: %v", err)
	}
}

func TestJavaExtractor_CommentsAndStringsIgnored(t *testing.T) {
	const tricky = `
package t;
// class Ghost {}  -- should be ignored
/* class AlsoGhost { void hidden(){} } */
public class Real {
    String code = "public class NotAType { }";
    String block = """
        interface StillNotAType {}
        """;
    void method() {
        String s = "brace } inside string";
    }
}
`
	dir := t.TempDir()
	writeJava(t, dir, "t/Real.java", tricky)

	ex := NewJavaExtractor()
	cat, err := ex.ExtractFromPath(dir)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}

	if objByName(cat, "Real", "class") == nil {
		t.Fatal("Real class not found")
	}
	for _, bad := range []string{"Ghost", "AlsoGhost", "NotAType", "StillNotAType"} {
		if o := objByName(cat, bad, "class"); o != nil {
			t.Errorf("commented/string type %q leaked as object", bad)
		}
	}
	if objByName(cat, "method", "function") == nil {
		t.Error("method() not found")
	}
	if err := cat.VerifyAxioms(); err != nil {
		t.Errorf("category axioms failed: %v", err)
	}
}

func TestJavaExtractor_InterfaceSatisfaction(t *testing.T) {
	// Compile-time check that JavaExtractor implements Extractor.
	var _ Extractor = (*JavaExtractor)(nil)

	ex := NewJavaExtractor()
	if ex.Language() != "java" {
		t.Errorf("Language() = %q, want java", ex.Language())
	}
	exts := ex.FileExtensions()
	if len(exts) != 1 || exts[0] != ".java" {
		t.Errorf("FileExtensions() = %v, want [.java]", exts)
	}
}
