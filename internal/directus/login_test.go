package directus_test

import (
	"testing"

	"github.com/Piitschy/drcts/test/testhelpers"
)

func TestGetToken(t *testing.T) {
	ctx, container, d := testhelpers.NewDirectusContainer(t, "latest")
	defer container.Terminate(ctx)

	auth, err := d.GetAuth(testhelpers.AdminEmail, testhelpers.AdminPassword)
	if err != nil {
		t.Fatalf("Failed to get token: %s", err)
	}
	if auth.AccessToken == "" {
		t.Fatalf("Token is empty")
	}
}

func TestGetTokenInvalid(t *testing.T) {
	ctx, container, d := testhelpers.NewDirectusContainer(t, "latest")
	defer container.Terminate(ctx)

	_, err := d.GetAuth("invalid", "invalid")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
