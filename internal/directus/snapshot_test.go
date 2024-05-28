package directus_test

import (
	"testing"

	"github.com/Piitschy/drcts/test/testhelpers"
)

func TestGetSnapshot(t *testing.T) {
	ctx, container, d := testhelpers.NewDirectusContainer(t, "latest")
	defer container.Terminate(ctx)

	err := d.Login(testhelpers.AdminEmail, testhelpers.AdminPassword)
	if err != nil {
		t.Fatalf("Failed to login: %s", err)
	}

	_, err = d.GetSnapshot()
	if err != nil {
		t.Fatalf("Failed to get snapshot: %s", err)
	}
}

func TestGetRawSnapshot(t *testing.T) {
	ctx, container, d := testhelpers.NewDirectusContainer(t, "latest")
	defer container.Terminate(ctx)

	err := d.Login(testhelpers.AdminEmail, testhelpers.AdminPassword)
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
