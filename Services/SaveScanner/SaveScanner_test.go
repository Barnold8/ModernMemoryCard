package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func logTestFailure(errMsg string, result string, expected string) string {
	// Note to self: make sure that values from a test are converted to a string for logging
	return fmt.Sprintf("\n\n-- Error message: %s\n\n\tRESULT: %s\n\n\tEXPECTED: %s\n\n", errMsg, result, expected)
}

func TestSeek(t *testing.T) {

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
				filepath.Join("dirWithFiles2", "image.png"),
				filepath.Join("dirWithFiles2", "script.sh"),
				filepath.Join("dirWithFiles2", "data.json"),
			},
		},
		{
			name:      "Test 3 - Finding files in root directory \"tmpDir\"",
			directory: "dirWithFiles2",
			expected: []string{
				filepath.Join("tmpDir", "dirWithFiles1", "file1.txt"),
				filepath.Join("tmpDir", "dirWithFiles1", "file2.go"),
				filepath.Join("tmpDir", "dirWithFiles1", "file3.md"),
				filepath.Join("tmpDir", "dirWithFiles2", "image.png"),
				filepath.Join("tmpDir", "dirWithFiles2", "script.sh"),
				filepath.Join("tmpDir", "dirWithFiles2", "data.json"),
			},
		},
		{
			name:      "Test 4 - Using an invalid directory",
			directory: "thisIsAnInvalidDirectory",
			expected:  []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			result := seek(tc.directory)

			if !(reflect.DeepEqual(result, tc.expected)) {
				t.Error(
					logTestFailure(
						fmt.Sprintf("When seeking for files in directory, %s, an error ocurred\n", tc.directory),
						strings.Join(result, ", "),
						strings.Join(tc.expected, ", "),
					),
				)
			}

		})
	}

}
