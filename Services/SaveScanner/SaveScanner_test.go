package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func logTestFailure(errMsg string, result string, expected string) string {
	// Note to self: make sure that values from a test are converted to a string for logging
	return fmt.Sprintf("\n\n-- Error message: %s\n\n\tRESULT: %s\n\n\tEXPECTED: %s\n\n", errMsg, result, expected)
}

func slicesEqualIgnoreOrder(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aCopy := append([]string(nil), a...)
	bCopy := append([]string(nil), b...)
	sort.Strings(aCopy)
	sort.Strings(bCopy)
	return reflect.DeepEqual(aCopy, bCopy)
}

func TestSeek(t *testing.T) {
	//TODO: test for filtering (ignore specific directories)

	tmpDir := t.TempDir()

	dirs := []string{
		"emptyDir1",
		"emptyDir2",
		"dirWithFiles1",
		"dirWithFiles2/subdir",
	}

	files := map[string][]string{
		"dirWithFiles1":        {"file1.txt", "file2.go", "file3.md"},
		"dirWithFiles2/subdir": {"image.png", "script.sh", "data.json"},
	}

	// Generate temp directories given dirs array
	for _, d := range dirs {
		dirPath := filepath.Join(tmpDir, d)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dirPath, err)
		}
	}

	// given the key value pairs in files, generate respective files in directories using join
	for dir, filenames := range files {
		dirPath := filepath.Join(tmpDir, dir)
		for _, fname := range filenames {
			fpath := filepath.Join(dirPath, fname)
			content := []byte("Non important data")
			if err := os.WriteFile(fpath, content, 0644); err != nil {
				t.Fatalf("Failed to write file %s: %v", fpath, err)
			}
		}
	}

	tests := []struct {
		name      string
		directory string
		expected  []string
	}{
		{
			name:      "Test 1 - Finding files in dirWithFiles1",
			directory: "dirWithFiles1",
			expected: []string{
				filepath.Join("dirWithFiles1", "file1.txt"),
				filepath.Join("dirWithFiles1", "file2.go"),
				filepath.Join("dirWithFiles1", "file3.md"),
			},
		},
		{
			name:      "Test 2 - Finding files in dirWithFiles2",
			directory: "dirWithFiles2",
			expected: []string{
				filepath.Join("dirWithFiles2/subdir", "data.json"),
				filepath.Join("dirWithFiles2/subdir", "image.png"),
				filepath.Join("dirWithFiles2/subdir", "script.sh"),
			},
		},
		{
			name:      "Test 3 - Finding files in root directory \"tmpDir\"",
			directory: ".",
			expected: []string{
				filepath.Join("dirWithFiles1", "file1.txt"),
				filepath.Join("dirWithFiles1", "file2.go"),
				filepath.Join("dirWithFiles1", "file3.md"),
				filepath.Join("dirWithFiles2/subdir", "data.json"),
				filepath.Join("dirWithFiles2/subdir", "image.png"),
				filepath.Join("dirWithFiles2/subdir", "script.sh"),
			},
		},
		{
			name:      "Test 4 - Using an invalid directory",
			directory: "thisIsAnInvalidDirectory",
			expected:  []string{},
		},
	}

	for _, tc := range tests {
		path, err := os.Getwd()
		if err != nil {
			log.Println(err)
			return
		}

		t.Run(tc.name, func(t *testing.T) {

			os.Chdir(tmpDir)
			result := seek(tc.directory)

			if !(slicesEqualIgnoreOrder(result, tc.expected)) {
				t.Error(
					logTestFailure(
						fmt.Sprintf("When seeking for files in directory, %s, an error ocurred\n", tc.directory),
						strings.Join(result, ", "),
						strings.Join(tc.expected, ", "),
					),
				)
			}

		})
		os.Chdir(path) // this is here so golang can remove tmp data
	}

}
