package directus_test

import (
	"testing"

	"github.com/Piitschy/drcts/internal/directus"
	"github.com/Piitschy/drcts/test/testhelpers"
)

func TestCreateCollection(t *testing.T) {
	ctx, container, _ := testhelpers.NewDirectusContainerWithCollection(t, "latest", LoadTestCollection(t, "article.json"), []*directus.Field{LoadTestField(t, "id_field.json")})
	defer container.Terminate(ctx)
}

func TestGetCollection(t *testing.T) {
	ctx, container, d := testhelpers.NewDirectusContainerWithCollection(t, "latest", LoadTestCollection(t, "article.json"), []*directus.Field{LoadTestField(t, "id_field.json")})
	defer container.Terminate(ctx)

	c, err := d.GetCollection("articles")
	if err != nil {
		t.Fatalf("Failed to get collection: %s", err)
	}
	if c.Collection != "articles" {
		t.Fatalf("Collection name should be 'articles'")
	}
}
