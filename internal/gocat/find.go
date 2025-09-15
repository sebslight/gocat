package gocat

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func FindFiles(root string, extensions []string, maxDepth int) ([]string, error) {
	paths := []string{}

	callbackFunc := func(path string, f fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			if maxDepth == -1 {
				return nil
			}

			relativePath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			if relativePath == "." {
				return nil
			}

			depth := strings.Count(relativePath, string(os.PathSeparator))
			if depth >= maxDepth {
				return filepath.SkipDir
			}
			return nil
		}
		shouldAdd := false
		if len(extensions) == 0 {
			shouldAdd = true
		} else {
			for _, ext := range extensions {
				if filepath.Ext(f.Name()) == ext {
					shouldAdd = true
					break
				}
			}
		}
		if shouldAdd {
			paths = append(paths, path)
		}
		return nil
	}

	err := filepath.WalkDir(root, callbackFunc)
	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	return paths, nil
}
