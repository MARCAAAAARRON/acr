package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// skipDirs are folder names we never want to walk into.
var skipDirs = map[string]bool{
	".git":         true,
	"vendor":       true,
	"node_modules": true,
	"dist":         true,
	"build":        true,
	".idea":        true,
	".vscode":      true,
}
var allowedExt = map[string]bool{
	"go": true, "php": true, "js": true, "ts": true,
	"py": true, "java": true, "rb": true, "md": true,
}

// FileInfo describes one source file we found worth indexing.
type FileInfo struct {
	Path string // relative to repo root
	Ext  string
}

// Scan walks root and returns every file we consider relevant.
func Scan(root string) ([]FileInfo, error) {
	var results []FileInfo

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.TrimPrefix(filepath.Ext(path), ".")
		if !allowedExt[ext]{
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		results = append(results, FileInfo{Path: rel, Ext: ext})
		return nil
	})

	return results, err
}
