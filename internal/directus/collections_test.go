package directus_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Piitschy/drcts/internal/directus"
	"github.com/testcontainers/testcontainers-go"
)

func NewDirectusContainerWithCollection(t *testing.T, version string, collection *directus.Collection, fields []*directus.Field) (context.Context, testcontainers.Container, *directus.Directus) {
	ctx, container, d := NewDirectusContainer(t, version)

	err := d.Login(adminEmail, adminPassword)
	if err != nil {
		container.Terminate(ctx)
		t.Fatalf("Failed to login: %s", err)
	}

	err = d.CreateCollection(collection)
	if err != nil {
		container.Terminate(ctx)
		t.Fatalf("Failed to create collection: %s", err)
	}

	if err != nil {
		container.Terminate(ctx)
		t.Fatalf("Failed to create collection: %s", err)
	}
	return ctx, container, d
}

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

func TestCreateCollection(t *testing.T) {
	ctx, container, _ := NewDirectusContainerWithCollection(t, "latest", LoadTestCollection(t, "article.json"), []*directus.Field{LoadTestField(t, "id_field.json")})
	defer container.Terminate(ctx)
}

func TestGetCollection(t *testing.T) {
	ctx, container, d := NewDirectusContainerWithCollection(t, "latest", LoadTestCollection(t, "article.json"), []*directus.Field{LoadTestField(t, "id_field.json")})
	defer container.Terminate(ctx)

	c, err := d.GetCollection("articles")
	if err != nil {
		t.Fatalf("Failed to get collection: %s", err)
	}
	if c.Collection != "articles" {
		t.Fatalf("Collection name should be 'articles'")
	}
}
