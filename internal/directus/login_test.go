package directus_test

import (
	"testing"

	h "github.com/Piitschy/drcts/test/testhelpers"
)

func TestGetToken(t *testing.T) {
	ctx, container, d := h.NewDirectusContainer(t, "latest")
	defer container.Terminate(ctx)

	auth, err := d.GetAuth(h.AdminEmail, h.AdminPassword)
	if err != nil {
		t.Fatalf("Failed to get token: %s", err)
	}
	if auth.AccessToken == "" {
		t.Fatalf("Token is empty")
	}
}

func TestGetTokenInvalid(t *testing.T) {
	ctx, container, d := h.NewDirectusContainer(t, "latest")
	defer container.Terminate(ctx)

	_, err := d.GetAuth("invalid", "invalid")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
