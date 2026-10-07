package wiringmap

import (
	"os"
	"path/filepath"
)

// ModuleRoot returns the directory containing go.mod for path, or path itself if none is found.
func ModuleRoot(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	dir := abs
	if !info.IsDir() {
		dir = filepath.Dir(abs)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return abs, nil
		}
		dir = parent
	}
}

// RelToModule returns p as a slash-separated path relative to moduleRoot.
func RelToModule(moduleRoot, p string) (string, error) {
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(moduleRoot, abs)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(rel), nil
}
