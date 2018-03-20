package oci // import "github.com/docker/docker/oci"

import (
	"fmt"
	"path/filepath"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// SetRawAccess removes the given raw access paths from the specs MaskedPaths and ReadonlyPaths.
func SetRawAccess(s *specs.Spec, rawaccess []string) error {
	// Iterate over the rawaccess paths.
	for _, rawpath := range rawaccess {
		// Iterate over the masked paths.
		for i, p := range s.Linux.MaskedPaths {
			matched, err := filepath.Match(rawpath, p)
			if err != nil {
				return fmt.Errorf("masked paths: checking if %s matches %s failed: %v", rawpath, p, err)
			}
			if matched {
				if len(s.Linux.MaskedPaths) > i+1 {
					s.Linux.MaskedPaths = append(s.Linux.MaskedPaths[:i], s.Linux.MaskedPaths[i+1:]...)
				} else {
					s.Linux.MaskedPaths = s.Linux.MaskedPaths[:i]
				}
			}
		}
		// Iterate over the readonly paths.
		for i, p := range s.Linux.ReadonlyPaths {
			matched, err := filepath.Match(rawpath, p)
			if err != nil {
				return fmt.Errorf("readonly paths: checking if %s matches %s failed: %v", rawpath, p, err)
			}
			if matched {
				if len(s.Linux.ReadonlyPaths) > i+1 {
					s.Linux.ReadonlyPaths = append(s.Linux.ReadonlyPaths[:i], s.Linux.ReadonlyPaths[i+1:]...)
				} else {
					s.Linux.ReadonlyPaths = s.Linux.ReadonlyPaths[:i]
				}
			}
		}
	}
	return nil
}
