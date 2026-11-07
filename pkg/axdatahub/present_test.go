package axdatahub

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestPresentFindsTheLibraryWithoutLoadingIt(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("LD_LIBRARY_PATH", "/nonexistent:"+dir)
	if err := Present(); !errors.Is(err, ErrUnavailable) {
		if _, statErr := os.Stat("/usr/lib/" + libraryName); statErr != nil {
			t.Fatalf("Present without the file: %v, want ErrUnavailable", err)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, libraryName), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Present(); err != nil {
		t.Fatalf("Present with the file in LD_LIBRARY_PATH: %v", err)
	}
}
