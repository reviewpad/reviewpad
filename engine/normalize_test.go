package engine

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	})
}

func TestNormalize_WithCustomNormalizers(t *testing.T) {
	defaultNormalizedRuleMode := NewNormalizeRule()
	defaultNormalizedRuleMode.WithModificators(func(file *ReviewpadFile) (*ReviewpadFile, error) {
		file.Mode = strings.ToUpper(defaultMode)
		return file, nil
	})

	var customNormalizers = []*NormalizeRule{
		defaultNormalizedRuleMode,
	}

	t.Run("missing mode", func(t *testing.T) {
		t.Parallel()

		reviewpadFile, err := normalize(&ReviewpadFile{}, customNormalizers...)
		assert.Nil(t, err)

		if reviewpadFile.Mode != strings.ToUpper(defaultMode) {
			t.Fatalf("expected mode: %s, got mode: %s", strings.ToUpper(defaultMode), reviewpadFile.Mode)
		}
	})
}
