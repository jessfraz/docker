package daemon // import "github.com/docker/docker/daemon"

import (
	"testing"

	"github.com/docker/docker/oci"
	"github.com/gotestyourself/gotestyourself/assert"
	is "github.com/gotestyourself/gotestyourself/assert/cmp"
)

// TestRawAccessProc checks that a user-specified RawAccess of /proc/*
// will not have those paths masked or set as readonly.
func TestRawAccessProc(t *testing.T) {
	// We can't call createSpec() so mimick the minimal part
	// of its code flow, just enough to reproduce the issue.
	rawAccess := []string{"/proc/*"}
	spec := oci.DefaultSpec()
	err := oci.SetRawAccess(&spec, rawAccess)
	assert.Check(t, err)

	// Check that the spec does not have certain ReadonlyPaths and MaskedPaths
	assert.Check(t, is.Equal(false, inSlice(spec.Linux.MaskedPaths, "/proc/kcore")))
	assert.Check(t, is.Equal(false, inSlice(spec.Linux.ReadonlyPaths, "/proc/sys")))
}
