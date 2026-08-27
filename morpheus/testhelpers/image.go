// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// WritePNG writes a minimal valid 1x1 PNG image to dir/name and returns its
// absolute path. A real PNG is used (rather than a placeholder) so that servers
// which validate uploaded image content accept it.
func WritePNG(t *testing.T, dir, name string) string {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 0x11, G: 0x22, B: 0x33, A: 0xff})

	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("creating test image %s: %v", path, err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatalf("encoding test image %s: %v", path, err)
	}

	return path
}
