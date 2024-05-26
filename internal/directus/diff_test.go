package directus_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Piitschy/drctsdm/internal/directus"
)

func TestGetDiffNoDifference(t *testing.T) {
	ctx, container, d := NewDirectusContainer(t, "latest")
	defer container.Terminate(ctx)

	err := d.Login(adminEmail, adminPassword)
	if err != nil {
		t.Fatalf("Failed to login: %s", err)
	}

	s, err := d.GetSnapshot()
	if err != nil {
		t.Fatalf("Failed to get snapshot: %s", err)
	}

	diff, err := d.GetDiff(s, true)
	if err != nil {
		t.Fatalf("Failed to get diff: %s", err)
	}

	if diff != nil {
		t.Fatalf("Diff should be nil")
	}
}

func TestGetDiffWithDifference(t *testing.T) {
	ctx, container, base := NewDirectusContainer(t, "latest")
	defer container.Terminate(ctx)

	err := base.Login(adminEmail, adminPassword)
	if err != nil {
		t.Fatalf("Failed to login: %s", err)
	}

	f, err := os.Open(filepath.Join("..", "..", "test", "testdata", "article.json"))
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

	err = base.CreateCollection(c)
	if err != nil {
		t.Fatalf("Failed to create collection: %s", err)
	}

	s, err := base.GetSnapshot()
	if err != nil {
		t.Fatalf("Failed to get snapshot: %s", err)
	}

	_, err = base.GetDiff(s, true)
	if err != nil {
		t.Fatalf("Failed to get diff: %s", err)
	}
}
