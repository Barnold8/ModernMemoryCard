package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// dir := "../."

	// f := seek(".")

	// fmt.Println("ABSOLUTEFILES AFTER IS %s\n", strings.Join(f, " | "))

}

func seek(directory string) []string {

	absoluteFiles := []string{}

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() { // is a file
			fmt.Printf("Found %s\n\n", path)
			absoluteFiles = append(absoluteFiles, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking the path:", err)
	}

	// fmt.Println("ABSOLUTEFILES BEFORE IS %s\n", strings.Join(absoluteFiles, " | "))

	return absoluteFiles
}
