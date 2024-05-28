package directus_test

import (
	"testing"

	h "github.com/Piitschy/drcts/test/testhelpers"
)

func TestCreateCollection(t *testing.T) {
	ctx, container, _ := h.NewDirectusContainerWithCollection(t, "latest",
		h.LoadTestCollection(t, "article.json"))
	defer container.Terminate(ctx)
}

func TestGetCollection(t *testing.T) {
	ctx, container, d := h.NewDirectusContainerWithCollection(t, "latest",
		h.LoadTestCollection(t, "article.json"))
	defer container.Terminate(ctx)

	c, err := d.GetCollection("articles")
	if err != nil {
		t.Fatalf("Failed to get collection: %s", err)
	}
	if c.Collection != "articles" {
		t.Fatalf("Collection name should be 'articles'")
	}
}
