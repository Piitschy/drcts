package directus_test

import (
	"testing"

	h "github.com/Piitschy/drcts/test/testhelpers"
)

func TestGetSnapshot(t *testing.T) {
	ctx, container, d := h.NewDirectusContainer(t, "latest")
	defer container.Terminate(ctx)

	err := d.Login(h.AdminEmail, h.AdminPassword)
	if err != nil {
		t.Fatalf("Failed to login: %s", err)
	}

	_, err = d.GetSnapshot()
	if err != nil {
		t.Fatalf("Failed to get snapshot: %s", err)
	}
}

func TestGetRawSnapshot(t *testing.T) {
	ctx, container, d := h.NewDirectusContainer(t, "latest")
	defer container.Terminate(ctx)

	err := d.Login(h.AdminEmail, h.AdminPassword)
	if err != nil {
		t.Fatalf("Failed to login: %s", err)
	}

	bytes, err := d.GetRawSnapshot("json")
	if err != nil {
		t.Fatalf("Failed to get raw snapshot: %s", err)
	}

	if len(bytes) == 0 {
		t.Fatalf("Snapshot is empty")
	}
}
