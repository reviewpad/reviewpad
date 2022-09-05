package engine

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testNormalizeCorrectExpectedVersion = "reviewpad.com/v3.x"
)

func TestNormalize(t *testing.T) {
	t.Run("missing edition, mode, and version", func(t *testing.T) {
		t.Parallel()

		reviewpadFile, err := normalize(&ReviewpadFile{})
		assert.Nil(t, err)

		if reviewpadFile.Edition != defaultEdition {
			t.Fatalf("expected edition: %s, got edition: %s", defaultEdition, reviewpadFile.Edition)
		}
		if reviewpadFile.Mode != defaultMode {
			t.Fatalf("expected mode: %s, got mode: %s", defaultMode, reviewpadFile.Mode)
		}
		if reviewpadFile.Version != defaultApiVersion {
			t.Fatalf("expected version: %s, got version: %s", defaultApiVersion, reviewpadFile.Version)
		}
	})

	t.Run("correct lower api-version number", func(t *testing.T) {
		t.Parallel()

		reviewpadFile, err := normalize(&ReviewpadFile{
			Version: "reviewpad.com/v3.x",
		})
		assert.Nil(t, err)

		if reviewpadFile.Version != testNormalizeCorrectExpectedVersion {
			t.Fatalf("expected api-version: %s, got api-version: %s", testNormalizeCorrectExpectedVersion, reviewpadFile.Version)
		}
	})
}

func TestNormalize_WithCustomNormalizers(t *testing.T) {
	const (
		customApiVersion = "reviewpad.com/v4.0.0-alpha"
	)

	/*
	   api-version: reviewpad.com/v4.0.0-alpha
	   edition: Ultimate
	   mode: SILENT
	*/
	var customNormalizers = []*NormalizeRule{
		// always use version, if it's empty or not: "reviewpad.com/v4.0.0-alpha"
		(&NormalizeRule{}).WithModificators(func(file *ReviewpadFile) (*ReviewpadFile, error) {
			file.Version = customApiVersion
			return file, nil
		}),
		(&NormalizeRule{}).WithModificators(func(file *ReviewpadFile) (*ReviewpadFile, error) {
			file.Mode = strings.ToUpper(defaultMode)
			return file, nil
		}),
	}

	t.Run("missing mode", func(t *testing.T) {
		t.Parallel()

		reviewpadFile, err := normalize(&ReviewpadFile{
			Version: "reviewpad.com/v3.x",
		}, customNormalizers...)
		assert.Nil(t, err)

		if reviewpadFile.Version != customApiVersion {
			t.Fatalf("expected api-version: %s, got api-version: %s", customApiVersion, reviewpadFile.Version)
		}
		if reviewpadFile.Mode != strings.ToUpper(defaultMode) {
			t.Fatalf("expected mode: %s, got mode: %s", strings.ToUpper(defaultMode), reviewpadFile.Mode)
		}
	})
}
