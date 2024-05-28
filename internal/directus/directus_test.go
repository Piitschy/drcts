package directus_test

import (
	"testing"

	h "github.com/Piitschy/drcts/test/testhelpers"
)

func TestNewDirectus(t *testing.T) {
	ctx, container, _ := h.NewDirectusContainer(t, "latest")
	defer container.Terminate(ctx)
}
