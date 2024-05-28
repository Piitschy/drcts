package directus_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Piitschy/drcts/internal/directus"
	"github.com/Piitschy/drcts/test/testhelpers"
)

func LoadTestCollection(t *testing.T, name string) *directus.Collection {
	f, err := os.Open(filepath.Join("..", "..", "test", "testdata", name))
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
	f, err := os.Open(filepath.Join("..", "..", "test", "testdata", name))
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

func TestNewDirectus(t *testing.T) {
	ctx, container, _ := testhelpers.NewDirectusContainer(t, "latest")
	defer container.Terminate(ctx)
}
