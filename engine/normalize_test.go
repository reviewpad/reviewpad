package engine

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNormalize(t *testing.T) {
	t.Run("missing edition and mode", func(t *testing.T) {
		t.Parallel()
		reviewpadFile, err := Load([]byte(testNormalizeNoEditionNoMode))
		if err != nil {
			t.Fatalf("err Load: %s", err.Error())
		}

		_ = normalize(reviewpadFile)

		if reviewpadFile.Edition != defaultEdition {
			t.Fatalf("expected edition: %s, got edition: %s", defaultEdition, reviewpadFile.Edition)
		}
		if reviewpadFile.Mode != defaultMode {
			t.Fatalf("expected mode: %s, got mode: %s", defaultMode, reviewpadFile.Mode)
		}
	})

	t.Run("incorrect api-version number", func(t *testing.T) {
		t.Parallel()
		reviewpadFile, err := Load([]byte(testNormalizeIncorrectApiVersionNumber))
		if err != nil {
			t.Fatalf("err Load: %s", err.Error())
		}

		_ = normalize(reviewpadFile)

		if reviewpadFile.Version != defaultApiVersion {
			t.Fatalf("expected api-version: %s, got api-version: %s", defaultApiVersion, reviewpadFile.Version)
		}
	})

	t.Run("correct lower api-version number", func(t *testing.T) {
		t.Parallel()
		reviewpadFile, err := Load([]byte(testNormalizeCorrectLowerApiVersionNumber))
		if err != nil {
			t.Fatalf("err Load: %s", err.Error())
		}

		_ = normalize(reviewpadFile)

		if reviewpadFile.Version != testNormalizeCorrectExpectedVersion {
			t.Fatalf("expected api-version: %s, got api-version: %s", testNormalizeCorrectExpectedVersion, reviewpadFile.Version)
		}
	})

	t.Run("incorrect string case in edition", func(t *testing.T) {
		t.Parallel()
		reviewpadFile, err := Load([]byte(testNormalizeIncorrectCase))
		if err != nil {
			t.Fatalf("err Load: %s", err.Error())
		}

		_ = normalize(reviewpadFile)

		if reviewpadFile.Edition != "team" {
			t.Fatalf("expected edition: %s, got edition: %s", "team", reviewpadFile.Edition)
		}
	})
}

func TestNormalize_WithCustomNormalizers(t *testing.T) {
	var newNormalizeRuleCustom = func(defaultVal string, allowed ...string) *NormalizeRule {
		n := &NormalizeRule{Default: defaultVal, Allowed: allowed}
		return n
	}

	const (
		customApiVersion = "reviewpad.com/v4.0.0-alpha"
		customEdition    = "Ultimate" // only if "professional"
	)

	/*
	   api-version: reviewpad.com/v4.0.0-alpha
	   edition: Ultimate
	   mode: SILENT
	*/
	var customNormalizers = map[propertyKey]*NormalizeRule{
		// always use version, if it's empty or not: "reviewpad.com/v4.0.0-alpha"
		_version: (&NormalizeRule{}).WithModificators(func(_ string) string { return customApiVersion }),

		// replace edition: professional, to: ultimate
		// and Title, so if current edition == "professional", it should become: "Ultimate"
		_edition: newNormalizeRuleCustom(customEdition, "Ultimate", "Team").
			WithModificators(func(val string) string { return strings.Replace(val, "professional", "ultimate", 1) }).
			WithModificators(strings.ToLower, strings.Title),

		// our mode should be in Upper Case. So we:
		// - define default in UpperCase (if Mode empty)
		// - specify allowed values in UpperCase (this long strings split join construction just for testing purpose)
		// - specify modificator to UpperCase value
		_mode: newNormalizeRuleCustom(strings.ToUpper(defaultMode), strings.Split(strings.ToUpper(strings.Join(allowedModes, " ")), " ")...),
	}

	t.Run("missing edition and mode", func(t *testing.T) {
		t.Parallel()
		reviewpadFile, err := Load([]byte(testNormalizeNoEditionNoMode))
		if err != nil {
			t.Fatalf("err Load: %s", err.Error())
		}

		//debugOutput(reviewpadFile)
		_ = normalize(reviewpadFile, customNormalizers)
		//debugOutput(reviewpadFile)

		if reviewpadFile.Version != customApiVersion {
			t.Fatalf("expected api-version: %s, got api-version: %s", customApiVersion, reviewpadFile.Version)
		}
		if reviewpadFile.Edition != customEdition {
			t.Fatalf("expected edition: %s, got edition: %s", customEdition, reviewpadFile.Edition)
		}
		if reviewpadFile.Mode != strings.ToUpper(defaultMode) {
			t.Fatalf("expected mode: %s, got mode: %s", strings.ToUpper(defaultMode), reviewpadFile.Mode)
		}
	})

	t.Run("incorrect api-version number", func(t *testing.T) {
		t.Parallel()
		reviewpadFile, err := Load([]byte(testNormalizeIncorrectApiVersionNumber))
		if err != nil {
			t.Fatalf("err Load: %s", err.Error())
		}

		_ = normalize(reviewpadFile, customNormalizers)

		if reviewpadFile.Version != customApiVersion {
			t.Fatalf("expected api-version: %s, got api-version: %s", customApiVersion, reviewpadFile.Version)
		}
	})

	t.Run("incorrect string case in edition", func(t *testing.T) {
		t.Parallel()
		reviewpadFile, err := Load([]byte(testNormalizeIncorrectCase))
		if err != nil {
			t.Fatalf("err Load: %s", err.Error())
		}

		//debugOutput(reviewpadFile)
		_ = normalize(reviewpadFile, customNormalizers)
		//debugOutput(reviewpadFile)

		if reviewpadFile.Edition != "Team" {
			t.Fatalf("expected edition: %s, got edition: %s", "Team", reviewpadFile.Edition)
		}
	})

}

func debugOutput(f *ReviewpadFile) {
	y, err := yaml.Marshal(f)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(y))
}

const (
	testNormalizeCorrectExpectedVersion       = "reviewpad.com/v2.x"
	testNormalizeCorrectLowerApiVersionNumber = "api-version: " + testNormalizeCorrectExpectedVersion
	testNormalizeIncorrectApiVersionNumber    = "api-version: reviewpad.com/v6.x"
	testNormalizeIncorrectCase                = "edition: TeaM"
	testNormalizeNoEditionNoMode              = `api-version: reviewpad.com/v3.x

labels:
  small:
    description: Pull requests with less than 90 LOC
    color: "f15d22"

rules:
  - name: isSmallPatch
    kind: patch
    description: Patch has less than 90 LOC
    spec: $size() < 90

workflows:
  - name: labelSmall
    description: Label small pull requests
    if:
      - rule: isSmallPatch
    then:
      - $addLabel("small")`
)
