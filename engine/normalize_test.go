package engine

import (
	"testing"
)

func TestNormalization(t *testing.T) {
	t.Run("missing edition and mode", func(t *testing.T) {
		t.Parallel()
		reviewpadFile, err := Load([]byte(testNormalizeNoEditionNoMode))
		if err != nil {
			t.Fatalf("err Load: %s", err.Error())
		}

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
		if reviewpadFile.Edition != "team" {
			t.Fatalf("expected edition: %s, got edition: %s", "team", reviewpadFile.Edition)
		}
	})

	//debug output
	/*y, err := yaml.Marshal(reviewpadFile)
	  if err != nil {
	  	t.Fatalf("err marshalling reviewpadFile: %s", err.Error())
	  }
	  fmt.Println(string(y))*/
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
