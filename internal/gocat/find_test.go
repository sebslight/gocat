package gocat

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	// You might also need "reflect" or a helper library for comparing slices later.
)

func TestFindFiles(t *testing.T) {
	// 1. ARRANGE (Outer): The shared setup.
	tmpDir := t.TempDir()

	// Create a predictable, multi-level file structure.
	rootFileGo := filepath.Join(tmpDir, "root.go")
	os.WriteFile(rootFileGo, []byte("go"), 0644)

	rootFileTxt := filepath.Join(tmpDir, "root.txt")
	os.WriteFile(rootFileTxt, []byte("txt"), 0644)

	subDir := filepath.Join(tmpDir, "sub")
	os.Mkdir(subDir, 0755)

	subFileGo := filepath.Join(subDir, "sub.go")
	os.WriteFile(subFileGo, []byte("go"), 0644)

	// 2. DEFINE THE TABLE of test cases.
	testCases := []struct {
		name          string   // A name for the sub-test
		extensions    []string // Input: extensions to filter by
		maxDepth      int      // Input: the depth limit
		expectedFiles []string // Expected output
	}{
		{
			name:          "Basic .go filter, infinite depth",
			extensions:    []string{".go"},
			maxDepth:      -1, // -1 for infinite
			expectedFiles: []string{rootFileGo, subFileGo},
		},
		{
			name:          "Basic .go filter, 0 depth",
			extensions:    []string{".go"},
			maxDepth:      0,
			expectedFiles: []string{rootFileGo},
		},
		{
			name:          "Basic .go filter, 1 depth",
			extensions:    []string{".go"},
			maxDepth:      1,
			expectedFiles: []string{rootFileGo, subFileGo},
		},
		{
			name:     "No filter, infinite depth",
			maxDepth: -1,
			expectedFiles: []string{
				rootFileGo,
				rootFileTxt,
				subFileGo,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			foundFiles, err := FindFiles(tmpDir, tc.extensions, tc.maxDepth)
			if err != nil {
				t.Errorf("FindFiles() error = %v", err)
				return
			}
			if len(foundFiles) != len(tc.expectedFiles) {
				t.Errorf("FindFiles() expected %d files, got %d, %+v", len(tc.expectedFiles), len(foundFiles), foundFiles)
				return
			}
			sort.Strings(foundFiles)
			sort.Strings(tc.expectedFiles)
			if !reflect.DeepEqual(foundFiles, tc.expectedFiles) {
				t.Errorf("FindFiles() expected %v, got %v", tc.expectedFiles, foundFiles)
				return
			}
		})
	}
}

