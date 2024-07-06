package main

import (
	"fmt"
	"os"
	"testing"

	h "github.com/Piitschy/drcts/test/testhelpers"
)

var TestVersions = []string{
	"10.7",
	"10.6",
	"10.5",
	"10.4",
	"10.0",
	"9",
}

func testMigrateCommand(t *testing.T, baseVersion, targetVersion string) {
	baseCtx, baseContainer, baseDirectus := h.NewDirectusContainerWithCollection(t, baseVersion,
		h.LoadTestCollection(t, "article.json"))
	defer baseContainer.Terminate(baseCtx)

	targetCtx, targetContainer, targetDirectus := h.NewDirectusContainer(t, targetVersion)
	defer targetContainer.Terminate(targetCtx)

	_, err := targetDirectus.GetCollection("articles")
	if err == nil {
		t.Fatalf("Collection should not exist in target in version %s", targetVersion)
	}

	_, err = baseDirectus.GetCollection("articles")
	if err != nil {
		t.Fatalf("Collection should exist in base in version %s", baseVersion)
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
	args = append(args, "-f")

	fmt.Println(args)

	if err := app.Run(args); err != nil {
		t.Fatalf("Failed to run migration (v%s -> v%s): %s", baseVersion, targetVersion, err)
	}

	if err = targetDirectus.Login(h.AdminEmail, h.AdminPassword); err != nil {
		t.Fatalf("Failed to login: %s", err)
	}
	_, err = targetDirectus.GetCollection("articles")
	if err != nil {
		t.Fatalf("Collection should exist in target (v%s): %s", targetVersion, err)
	}
}

func TestMigrateLatest(t *testing.T) {
	testMigrateCommand(t, "latest", "latest")
}

func TestMigrateVersions(t *testing.T) {
	for _, version := range TestVersions {
		testMigrateCommand(t, version, version)
	}
}

func TestMigrateUpgrade(t *testing.T) {
	for i := 0; i < len(TestVersions)-1; i++ {
		testMigrateCommand(t, TestVersions[i], TestVersions[i+1])
	}
}

func TestMigrateDowngrade(t *testing.T) {
	for i := len(TestVersions) - 1; i > 0; i-- {
		testMigrateCommand(t, TestVersions[i], TestVersions[i-1])
	}
}
