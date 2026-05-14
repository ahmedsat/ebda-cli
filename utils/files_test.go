package utils_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ahmedsat/ebda-cli/utils"
)

func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(path, []byte("hi"), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	if !utils.FileExists(path) {
		t.Fatal("FileExists returned false for an existing file")
	}
	if utils.FileExists(filepath.Join(dir, "missing.txt")) {
		t.Fatal("FileExists returned true for a missing file")
	}
}

func TestSaveFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")
	if err := utils.SaveFile(path, []byte("hello")); err != nil {
		t.Fatalf("SaveFile: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "hello" {
		t.Fatalf("file content = %q, want hello", got)
	}
}
