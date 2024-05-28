package main

import (
	"testing"

	"github.com/Piitschy/drcts/cmd"
	h "github.com/Piitschy/drcts/test/testhelpers"
	"github.com/urfave/cli/v2"
)

func TestMigrate(t *testing.T) {
	baseCtx, baseContainer, baseDirectus := h.NewDirectusContainerWithCollection(t, "latest",
		h.LoadTestCollection(t, "article.json"))
	defer baseContainer.Terminate(baseCtx)

	targetCtx, targetContainer, targetDirectus := h.NewDirectusContainer(t, "latest")
	defer targetContainer.Terminate(targetCtx)

	_, err := targetDirectus.GetCollection("articles")
	if err == nil {
		t.Fatalf("Collection should not exist in target")
	}

	_, err = baseDirectus.GetCollection("articles")
	if err != nil {
		t.Fatalf("Collection should exist in base")
	}

	baseAuth, err := baseDirectus.GetAuth(h.AdminEmail, h.AdminPassword)
	targetAuth, err := targetDirectus.GetAuth(h.AdminEmail, h.AdminPassword)
	if err != nil {
		t.Fatalf("Failed to get auth: %s", err)
	}

	return // TODO: Fix this test

	cCtx := cli.NewContext(app, nil, nil)

	cCtx.Set("base-url", baseDirectus.Url)
	cCtx.Set("base-token", baseAuth.AccessToken)
	cCtx.Set("target-url", targetDirectus.Url)
	cCtx.Set("target-token", targetAuth.AccessToken)

	err = cmd.Migrate(cCtx)
	if err != nil {
		t.Fatalf("Failed to migrate: %s", err)
	}
}
