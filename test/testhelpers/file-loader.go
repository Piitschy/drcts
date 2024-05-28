package testhelpers

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Piitschy/drcts/internal/directus"
)

func genDots(n int) string {
	var path = []string{}
	for i := 0; i < n; i++ {
		path = append(path, "..")
	}
	return filepath.Join(path...)
}

func openTestDataFile(name string) (*os.File, error) {
	var f *os.File
	var err error
	for i := 0; i < 5; i++ {
		f, err = os.Open(filepath.Join(genDots(i), "test", "testdata", name))
		if err == nil {
			break
		}
	}
	return f, err
}

func LoadTestCollection(t *testing.T, name string) *directus.Collection {
	f, err := openTestDataFile(name)
	if err != nil {
		t.Fatalf("Failed to open file: %s", err)
	}
	defer f.Close()
	bytes, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("Failed to read file: %s", err)
	}

	c, err := directus.UnmarshalCollection(bytes)
	if err != nil {
		t.Fatalf("Failed to unmarshal collection: %s", err)
	}
	return c
}
func LoadTestField(t *testing.T, name string) *directus.Field {
	f, err := openTestDataFile(name)
	if err != nil {
		t.Fatalf("Failed to open file: %s", err)
	}
	defer f.Close()
	bytes, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("Failed to read file: %s", err)
	}

	field, err := directus.UnmarshalField(bytes)
	if err != nil {
		t.Fatalf("Failed to unmarshal field: %s", err)
	}
	return field
}
