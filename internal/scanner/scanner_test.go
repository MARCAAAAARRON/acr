package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanSkipsIgnoredDirsAndExtensions(t *testing.T) {
	dir := t.TempDir()

	// allowed files — should show up in results
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(dir, "readme.md"), []byte("# hello"), 0644)

	// disallowed extension — should NOT show up
	os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("random notes"), 0644)

	// a skipped directory with a file inside it — neither the folder
	// nor anything inside it should show up
	vendorDir := filepath.Join(dir, "vendor")
	os.Mkdir(vendorDir, 0755)
	os.WriteFile(filepath.Join(vendorDir, "lib.go"), []byte("package lib"), 0644)

	files, err := Scan(dir)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// build a simple lookup so we can check presence/absence by path
	found := make(map[string]bool)
	for _, f := range files {
		found[f.Path] = true
	}

	if !found["main.go"] {
		t.Errorf("expected main.go to be found, but it wasn't")
	}
	if !found["readme.md"] {
		t.Errorf("expected readme.md to be found, but it wasn't")
	}
	if found["notes.txt"] {
		t.Errorf("expected notes.txt to be excluded (disallowed extension), but it was found")
	}
	if found["vendor\\lib.go"] || found["vendor/lib.go"] {
		t.Errorf("expected vendor/lib.go to be excluded (skipped directory), but it was found")
	}
}