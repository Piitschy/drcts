package directus_test

import (
	"testing"
)

func TestGetToken(t *testing.T) {
	ctx, container, d := NewDirectusContainer(t, "latest")
	defer container.Terminate(ctx)

	auth, err := d.GetAuth(adminEmail, adminPassword)
	if err != nil {
		t.Fatalf("Failed to get token: %s", err)
	}
	if auth.AccessToken == "" {
		t.Fatalf("Token is empty")
	}
}

func TestGetTokenInvalid(t *testing.T) {
	ctx, container, d := NewDirectusContainer(t, "latest")
	defer container.Terminate(ctx)

	_, err := d.GetAuth("invalid", "invalid")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
