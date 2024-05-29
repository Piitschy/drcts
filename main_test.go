package main

import (
	"fmt"
	"os"
	"testing"

	h "github.com/Piitschy/drcts/test/testhelpers"
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

	args := os.Args[:1]
	args = append(args, "--base-url", baseDirectus.Url)
	args = append(args, "--base-email", h.AdminEmail)
	args = append(args, "--base-password", h.AdminPassword)
	args = append(args, "--target-url", targetDirectus.Url)
	args = append(args, "--target-email", h.AdminEmail)
	args = append(args, "--target-password", h.AdminPassword)
	args = append(args, "migrate")
	args = append(args, "-y")

	fmt.Println(args)

	if err := app.Run(args); err != nil {
		t.Fatalf("Failed to run migration: %s", err)
	}

	if err = targetDirectus.Login(h.AdminEmail, h.AdminPassword); err != nil {
		t.Fatalf("Failed to login: %s", err)
	}
	_, err = targetDirectus.GetCollection("articles")
	if err != nil {
		t.Fatalf("Collection should exist in target: %s", err)
	}
}
