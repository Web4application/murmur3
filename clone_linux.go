//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// FICLONE ioctl constant (from linux/fs.h)
const ficlone = 0x40049409

// cloneFile tries to use FICLONE ioctl for instant reflink cloning.
// Falls back to normal copy if reflink is not supported.
func cloneFile(from, to string) error {
	src, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer src.Close()

	dst, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open destination: %w", err)
	}
	defer dst.Close()

	// Try reflink clone
	if err := unix.IoctlFileClone(int(dst.Fd()), int(src.Fd())); err == nil {
		return nil // success
	}

	// If reflink fails, fall back to normal copy
	if _, err := src.Seek(0, 0); err != nil {
		return fmt.Errorf("seek source: %w", err)
	}
	if _, err := dst.Seek(0, 0); err != nil {
		return fmt.Errorf("seek destination: %w", err)
	}

	if _, err := dst.ReadFrom(src); err != nil {
		return fmt.Errorf("fallback copy failed: %w", err)
	}

	return nil
}

// IoctlFileClone is a helper for FICLONE ioctl
func IoctlFileClone(dstFD, srcFD int) error {
	return unix.IoctlSetInt(dstFD, ficlone, srcFD)
}
clonefile_fallback.go (Other OSes)

Go

Copy code
//go:build !darwin && !linux
// +build !darwin,!linux

package main

import (
	"fmt"
	"io"
	"os"
)

// cloneFile falls back to a standard file copy for unsupported OSes.
func cloneFile(from, to string) error {
	src, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("create destination: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("copy data: %w", err)
	}

	return nil
}
