package directus_test

import (
	"fmt"
	"testing"

	"github.com/Piitschy/drcts/internal/directus"
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
	baseCtx, baseContainer, baseDirectus := NewDirectusContainerWithCollection(t, "latest", LoadTestCollection(t, "article.json"), []*directus.Field{LoadTestField(t, "id_field.json")})
	defer baseContainer.Terminate(baseCtx)

	targetCtx, targetContainer, targetDirectus := NewDirectusContainer(t, "latest")
	defer targetContainer.Terminate(targetCtx)

	err := baseDirectus.Login(adminEmail, adminPassword)
	err = targetDirectus.Login(adminEmail, adminPassword)
	if err != nil {
		t.Fatalf("Failed to login: %s", err)
	}

	s, err := baseDirectus.GetSnapshot()
	if err != nil {
		t.Fatalf("Failed to get snapshot: %s", err)
	}

	diff, err := targetDirectus.GetDiff(s, true)
	if err != nil {
		t.Fatalf("Failed to get diff: %s", err)
	}

	if diff == nil {
		t.Fatalf("Diff should not be nil")
	}
}

func TestApplyDiff(t *testing.T) {
	baseCtx, baseContainer, baseDirectus := NewDirectusContainerWithCollection(t, "latest", LoadTestCollection(t, "article.json"), []*directus.Field{LoadTestField(t, "id_field.json")})
	defer baseContainer.Terminate(baseCtx)

	targetCtx, targetContainer, targetDirectus := NewDirectusContainer(t, "latest")
	defer targetContainer.Terminate(targetCtx)

	err := baseDirectus.Login(adminEmail, adminPassword)
	err = targetDirectus.Login(adminEmail, adminPassword)
	if err != nil {
		t.Fatalf("Failed to login: %s", err)
	}

	s, err := baseDirectus.GetSnapshot()
	if err != nil {
		t.Fatalf("Failed to get snapshot: %s", err)
	}

	diff, err := targetDirectus.GetDiff(s, true)
	if err != nil {
		t.Fatalf("Failed to get diff: %s", err)
	}

	if diff == nil {
		t.Fatalf("Diff should not be nil")
	}

	tSchemaOriginal, err := targetDirectus.GetSnapshot()
	if err != nil {
		t.Fatalf("Failed to get snapshot: %s", err)
	}

	fmt.Println("len Collections:", len(tSchemaOriginal.Collections))

	// Apply the diff
	err = targetDirectus.ApplyDiff(diff)
	if err != nil {
		t.Fatalf("Failed to apply diff: %s", err)
	}

	tSchema, err := targetDirectus.GetSnapshot()
	if err != nil {
		t.Fatalf("Failed to get snapshot: %s", err)
	}

	fmt.Println("len Collections:", len(tSchema.Collections))

	if len(tSchemaOriginal.Collections) == len(tSchema.Collections) {
		t.Fatalf("Collections should not be equal")
	}

	baseArticles, err := baseDirectus.GetCollection("articles")
	targetArticles, err := targetDirectus.GetCollection("articles")
	if err != nil {
		t.Fatalf("Failed to get collection: %s", err)
	}

	if baseArticles.Collection != targetArticles.Collection {
		t.Fatalf("Collections should be equal")
	}

	if len(baseArticles.Fields) != len(targetArticles.Fields) {
		t.Fatalf("Fields should be equal")
	}
}
