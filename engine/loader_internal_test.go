// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	data := []byte("hello")

	wantVal := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

	gotVal := hash(data)

	assert.Equal(t, wantVal, gotVal)
}

func TestLoad_WhenDataParseFails(t *testing.T) {
	data := []byte("hello")

	gotReviewpadFile, err := Load(data)

	assert.Nil(t, gotReviewpadFile)
	assert.EqualError(t, err, "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `hello` into engine.ReviewpadFile")
}

func TestLoad(t *testing.T) {
	data := []byte(`
api-version:  reviewpad.com/v1alpha
mode:         silent
edition:      professional
ignore-errors: false
groups:
    - name: seniors
      description: Senior developers
      kind: developers
      spec: '["john"]'
rules:
    - name: test-rule
      kind: patch
      description: testing rule
      spec: 1 == 1
labels:
    bug:
      color: f29513
      description: Something isn't working
workflows:
    - name: test
      description: Test process
      alwaysRun:   true
      if:
        - rule: tautology
      then:
        - $action()
`)

	wantReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports:      []PadImport{},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   false,
				Rules: []PadWorkflowRule{
					{Rule:         "tautology"},
				},
				Actions: []string{
					"$action()",
				},
			},
		},
	}

	gotReviewpadFile, err := Load(data)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewpadFile, gotReviewpadFile)
}

func TestParse_WhenProvidedANonYamlFormat(t *testing.T) {
	data := []byte("hello")

	gotReviewpadFile, err := parse(data)

	assert.Nil(t, gotReviewpadFile)
	assert.EqualError(t, err, "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `hello` into engine.ReviewpadFile")
}

func TestParse(t *testing.T) {
	data := []byte(`
api-version:  reviewpad.com/v1alpha
mode:         silent
edition:      professional
ignore-errors: false
groups:
    - name: seniors
      description: Senior developers
      kind: developers
      spec: '["john"]'
rules:
    - name: test-rule
      kind: patch
      description: testing rule
      spec: 1 == 1
labels:
    bug:
      color: f29513
      description: Something isn't working
workflows:
    - name: test
      description: Test process
      alwaysRun:   true
      if:
        - rule: tautology
      then:
        - $action()
`)

	wantReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   false,
				Rules: []PadWorkflowRule{
					{Rule:         "tautology"},
				},
				Actions: []string{
					"$action()",
				},
			},
		},
	}

	gotReviewpadFile, err := parse(data)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewpadFile, gotReviewpadFile)
}
