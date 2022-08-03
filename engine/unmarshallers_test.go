package engine

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNormalization(t *testing.T) {
	reviewpadFile, err := Load([]byte(testYaml1))
	if err != nil {
		t.Fatalf("err Load: %s", err.Error())
	}
	y, err := yaml.Marshal(reviewpadFile)
	if err != nil {
		t.Fatalf("err marshalling reviewpadFile: %s", err.Error())
	}
	fmt.Println(string(y))
}

const testYaml1 = `api-version: reviewpad.com/v3.x

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
