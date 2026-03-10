code
//go:build darwin
// +build darwin

package main

import (
	"fmt"

	"golang.org/x/sys/unix"
)

// cloneFile uses macOS's clonefile syscall for fast copy-on-write cloning.
func cloneFile(from, to string) error {
	if err := unix.Clonefile(from, to, 0); err != nil {
		return fmt.Errorf("clonefile failed: %w", err)
	}
	return nil
}
