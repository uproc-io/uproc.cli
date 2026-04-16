package processes

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveUploadPathsWithSingleFileAndGlob(t *testing.T) {
	tempDir := t.TempDir()
	fileA := filepath.Join(tempDir, "a.pdf")
	fileB := filepath.Join(tempDir, "b.pdf")
	if err := os.WriteFile(fileA, []byte("a"), 0o644); err != nil {
		t.Fatalf("write fileA: %v", err)
	}
	if err := os.WriteFile(fileB, []byte("b"), 0o644); err != nil {
		t.Fatalf("write fileB: %v", err)
	}

	paths, err := resolveUploadPaths([]string{fileA, filepath.Join(tempDir, "*.pdf")})
	if err != nil {
		t.Fatalf("resolveUploadPaths error: %v", err)
	}

	if len(paths) != 2 {
		t.Fatalf("expected 2 unique paths, got %d: %v", len(paths), paths)
	}
}

func TestResolveUploadPathsNoMatches(t *testing.T) {
	tempDir := t.TempDir()
	_, err := resolveUploadPaths([]string{filepath.Join(tempDir, "*.pdf")})
	if err == nil {
		t.Fatal("expected error when no glob matches")
	}
}

func TestResolveUploadPathsRejectsDirectory(t *testing.T) {
	tempDir := t.TempDir()
	_, err := resolveUploadPaths([]string{tempDir})
	if err == nil {
		t.Fatal("expected error when input is directory")
	}
}
