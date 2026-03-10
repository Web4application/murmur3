package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	src := "source.txt"
	dst := "clone.txt"

	// Create a test source file if it doesn't exist
	if _, err := os.Stat(src); os.IsNotExist(err) {
		if err := os.WriteFile(src, []byte("Hello, clonefile/reflink!\n"), 0644); err != nil {
			log.Fatalf("failed to create source file: %v", err)
		}
	}

	// Perform the clone (platform-specific implementation)
	if err := cloneFile(src, dst); err != nil {
		log.Fatalf("clone failed: %v", err)
	}

	fmt.Printf("Successfully cloned %s to %s\n", src, dst)
}
