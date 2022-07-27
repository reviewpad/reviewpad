// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var mockedReviewpadFile = &ReviewpadFile{
	Version:      "reviewpad.com/v1alpha",
	Edition:      "professional",
	Mode:         "silent",
	IgnoreErrors: false,
	Imports: []PadImport{
		{Url: "https://foo.bar/mockedImportedReviewpadFile.yml"},
	},
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
			Name:        "tautology",
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
			Name:        "test-workflow-A",
			Description: "Test process",
			AlwaysRun:   false,
			Rules: []PadWorkflowRule{
				{Rule: "tautology"},
				{ExtraActions: []string{"$merge()"}},
			},
			Actions: []string{
				"$merge()",
			},
		},
	},
}

var mockedReviewpadFileData = []byte(`
api-version:  reviewpad.com/v1alpha
mode:         silent
edition:      professional
ignore-errors: false
imports:
    - url: https://foo.bar/mockedImportedReviewpadFile.yml
groups:
    - name: seniors
      description: Senior developers
      kind: developers
      spec: '["john"]'
rules:
    - name: tautology
      kind: patch
      description: testing rule
      spec: 1 == 1
labels:
    bug:
      color: f29513
      description: Something isn't working
workflows:
    - name: test-workflow-A
      description: Test process
      alwaysRun:   true
      if:
        - rule: tautology
          extra-actions:
            - $merge()
      then:
        - $merge()
`)

var mockedImportedReviewpadFile = &ReviewpadFile{
	Version:      "reviewpad.com/v1alpha",
	Edition:      "professional",
	Mode:         "silent",
	IgnoreErrors: false,
	Imports:      []PadImport{},
	Rules: []PadRule{
		{
			Name:        "tautology",
			Kind:        "patch",
			Description: "testing rule",
			Spec:        "1 == 1",
		},
	},
	Labels: map[string]PadLabel{
		"enhancement": {
			Color:       "a2eeef",
			Description: "New feature or request",
		},
	},
	Workflows: []PadWorkflow{
		{
			Name:        "test-workflow-B",
			Description: "Test process",
			AlwaysRun:   false,
			Rules: []PadWorkflowRule{
				{Rule: "tautology"},
			},
			Actions: []string{
				"$action()",
			},
		},
	},
}

var mockedImportedReviewpadFileData = []byte(`
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
    - name: tautology
      kind: patch
      description: testing rule
      spec: 1 == 1
labels:
    bug:
      color: f29513
      description: Something isn't working
workflows:
    - name: test-workflow-A
      description: Test process
      alwaysRun:   true
      if:
        - rule: tautology
      then:
        - $action()
`)

var simpleReviewpadFileData = []byte(`
api-version:  reviewpad.com/v1alpha
mode:         silent
edition:      professional
rules:
    - name: tautology
      kind: patch
      description: testing rule
      spec: true
workflows:
    - name: simple-workflow
      description: Test process
      alwaysRun:   true
      if:
        - rule: tautology
      then:
        - $merge()
`)

var simpleReviewpadFileWithImportsData = []byte(`
api-version:  reviewpad.com/v1alpha
mode:         silent
edition:      professional
imports:
    - url: https://foo.bar/anotherSimpleReviewpadFile.yml
rules:
    - name: tautology
      kind: patch
      description: testing rule
      spec: true
workflows:
    - name: simple-workflow
      description: Test process
      alwaysRun:   true
      if:
        - rule: tautology
      then:
        - $merge()
`)

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
    - name: tautology
      kind: patch
      description: testing rule
      spec: 1 == 1
labels:
    bug:
      color: f29513
      description: Something isn't working
workflows:
    - name: test-workflow-A
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
				Name:        "tautology",
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
				Name:        "test-workflow-A",
				Description: "Test process",
				AlwaysRun:   false,
				Rules: []PadWorkflowRule{
					{Rule: "tautology"},
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
	wantReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/mockedImportedReviewpadFile.yml"},
		},
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
				Name:        "tautology",
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
				Name:        "test-workflow-A",
				Description: "Test process",
				AlwaysRun:   false,
				Rules: []PadWorkflowRule{
					{
						Rule: "tautology",
						ExtraActions: []string{
							"$merge()",
						},
					},
				},
				Actions: []string{
					"$merge()",
				},
			},
		},
	}

	gotReviewpadFile, err := parse(mockedReviewpadFileData)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewpadFile, gotReviewpadFile)
}

func TestTransform(t *testing.T) {
	reviewpadFile := &ReviewpadFile{}
	*reviewpadFile = *mockedReviewpadFile
	wantReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/mockedImportedReviewpadFile.yml"},
		},
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
				Name:        "tautology",
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
				Name:        "test-workflow-A",
				Description: "Test process",
				AlwaysRun:   false,
				Rules: []PadWorkflowRule{
					{Rule: "tautology"},
					{ExtraActions: []string{"$merge(\"merge\")"}},
				},
				Actions: []string{
					"$merge(\"merge\")",
				},
			},
		},
	}

	gotReviewpadFile := transform(reviewpadFile)

	assert.Equal(t, wantReviewpadFile, gotReviewpadFile)
}

func TestLoadImport_WhenInvalidUrl(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	reviewpadImport := PadImport{Url: "https://foo.bar/invalidUrl"}

	httpmock.RegisterResponder("GET", reviewpadImport.Url,
		httpmock.NewErrorResponder(fmt.Errorf("invalid import url")),
	)

	gotReviewpadFile, content, err := loadImport(reviewpadImport)

	assert.Nil(t, gotReviewpadFile)
	assert.Equal(t, "", content)
	assert.EqualError(t, err, "Get \"https://foo.bar/invalidUrl\": invalid import url")
}

func TestLoadImport_WhenContentIsInInvalidYamlFormat(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	reviewpadImport := PadImport{Url: "https://foo.bar/file.yml"}

	httpmock.RegisterResponder("GET", reviewpadImport.Url,
		httpmock.NewBytesResponder(200, []byte("invalid")),
	)

	gotReviewpadFile, content, err := loadImport(reviewpadImport)

	assert.Nil(t, gotReviewpadFile)
	assert.Equal(t, "", content)
	assert.EqualError(t, err, "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `invalid` into engine.ReviewpadFile")
}

func TestLoadImport(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	simpleReviewpadFileImport := PadImport{
		Url: "https://foo.bar/simpleReviewpadFile.yml",
	}

	httpmock.RegisterResponder("GET", simpleReviewpadFileImport.Url,
		httpmock.NewBytesResponder(200, simpleReviewpadFileData),
	)

	wantReviewpadFile := &ReviewpadFile{
		Version: "reviewpad.com/v1alpha",
		Mode:    "silent",
		Edition: "professional",
		Rules: []PadRule{
			{
				Name:        "tautology",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "true",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "simple-workflow",
				Description: "Test process",
				AlwaysRun:   false,
				Rules: []PadWorkflowRule{
					{Rule: "tautology"},
				},
				Actions: []string{"$merge(\"merge\")"},
			},
		},
	}

	wantHashContent := fmt.Sprintf("%x", sha256.Sum256(simpleReviewpadFileData))

	gotReviewpadFile, gotHashContent, err := loadImport(simpleReviewpadFileImport)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewpadFile, gotReviewpadFile)
	assert.Equal(t, wantHashContent, gotHashContent)
}

func TestInlineImports_WhenPadImportUrlIsNotProvided(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	reviewpadFile := &ReviewpadFile{}
	*reviewpadFile = *mockedReviewpadFile
	invalidUrl := "https://foo.bar/invalidUrl"
	reviewpadFile.Imports = []PadImport{
		{Url: invalidUrl},
	}

	httpmock.RegisterResponder("GET", invalidUrl,
		httpmock.NewErrorResponder(fmt.Errorf("invalid import url")),
	)

	loadEnv := &LoadEnv{
		Visited: make(map[string]bool),
		Stack:   make(map[string]bool),
	}

	gotReviewpadFile, err := inlineImports(reviewpadFile, loadEnv)

	assert.Nil(t, gotReviewpadFile)
	assert.EqualError(t, err, "Get \"https://foo.bar/invalidUrl\": invalid import url")
}

func TestInlineImports_WhenThereIsCycleDependency(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	importUrl := "https://foo.bar/mockedImportedReviewpadFile.yml"
	reviewpadFile := &ReviewpadFile{}
	*reviewpadFile = *mockedImportedReviewpadFile
	reviewpadFile.Imports = []PadImport{
		{Url: importUrl},
	}

	httpmock.RegisterResponder("GET", importUrl,
		httpmock.NewBytesResponder(200, mockedImportedReviewpadFileData),
	)

	expectedFileHash := fmt.Sprintf("%x", sha256.Sum256(mockedImportedReviewpadFileData))

	mockedVisited := make(map[string]bool)
	mockedStack := make(map[string]bool)

	mockedVisited[expectedFileHash] = true
	mockedStack[expectedFileHash] = true

	loadEnv := &LoadEnv{
		Visited: mockedVisited,
		Stack:   mockedStack,
	}

	gotReviewpadFile, err := inlineImports(mockedReviewpadFile, loadEnv)

	assert.Nil(t, gotReviewpadFile)
	assert.EqualError(t, err, "loader: cyclic dependency")
}

func TestInlineImports_WhenVisitsAreOptimized(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockedImportedReviewpadFileUrl := "https://foo.bar/mockedImportedReviewpadFile.yml"
	reviewpadFile := &ReviewpadFile{}
	*reviewpadFile = *mockedReviewpadFile
	reviewpadFile.Imports = []PadImport{
		{Url: mockedImportedReviewpadFileUrl},
	}

	httpmock.RegisterResponder("GET", mockedImportedReviewpadFileUrl,
		httpmock.NewBytesResponder(200, mockedImportedReviewpadFileData),
	)

	expectedFileHash := fmt.Sprintf("%x", sha256.Sum256(mockedImportedReviewpadFileData))

	mockedVisited := make(map[string]bool)
	mockedStack := make(map[string]bool)

	mockedVisited[expectedFileHash] = true

	loadEnv := &LoadEnv{
		Visited: mockedVisited,
		Stack:   mockedStack,
	}

	wantReviewpadFile := &ReviewpadFile{}
	*wantReviewpadFile = *mockedReviewpadFile
	wantReviewpadFile.Imports = []PadImport{}

	gotReviewpadFile, err := inlineImports(reviewpadFile, loadEnv)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewpadFile, gotReviewpadFile)
}

func TestInlineImports_WhenSubTreeFileInlineImportsFails(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockedSimpleReviewpadFileUrl := "https://foo.bar/simpleReviewpadFile.yml"
	reviewpadFile := &ReviewpadFile{}
	*reviewpadFile = *mockedReviewpadFile
	reviewpadFile.Imports = []PadImport{
		{Url: mockedSimpleReviewpadFileUrl},
	}

	httpmock.RegisterResponder("GET", mockedSimpleReviewpadFileUrl,
		httpmock.NewBytesResponder(200, simpleReviewpadFileWithImportsData),
	)

	httpmock.RegisterResponder("GET", "https://foo.bar/anotherSimpleReviewpadFile.yml",
		httpmock.NewBytesResponder(200, simpleReviewpadFileWithImportsData),
	)

	mockedVisited := make(map[string]bool)
	mockedStack := make(map[string]bool)

	loadEnv := &LoadEnv{
		Visited: mockedVisited,
		Stack:   mockedStack,
	}

	gotReviewpadFile, err := inlineImports(reviewpadFile, loadEnv)

	assert.Nil(t, gotReviewpadFile)
	assert.EqualError(t, err, "loader: cyclic dependency")
}

func TestInlineImports(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockedSimpleReviewpadFileUrl := "https://foo.bar/simpleReviewpadFile.yml"
	reviewpadFile := &ReviewpadFile{}
	*reviewpadFile = *mockedReviewpadFile
	reviewpadFile.Imports = []PadImport{
		{Url: mockedSimpleReviewpadFileUrl},
	}

	httpmock.RegisterResponder("GET", mockedSimpleReviewpadFileUrl,
		httpmock.NewBytesResponder(200, simpleReviewpadFileData),
	)

	mockedVisited := make(map[string]bool)
	mockedStack := make(map[string]bool)

	loadEnv := &LoadEnv{
		Visited: mockedVisited,
		Stack:   mockedStack,
	}

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
				Name:        "tautology",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
			{
				Name:        "tautology",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "true",
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
				Name:        "test-workflow-A",
				Description: "Test process",
				AlwaysRun:   false,
				Rules: []PadWorkflowRule{
					{Rule: "tautology"},
					{ExtraActions: []string{"$merge()"}}},
				Actions: []string{"$merge()"}},
			{
				Name:        "simple-workflow",
				Description: "Test process",
				AlwaysRun:   false,
				Rules: []PadWorkflowRule{
					{Rule: "tautology"},
				},
				Actions: []string{"$merge(\"merge\")"},
			},
		},
	}

	gotReviewpadFile, err := inlineImports(reviewpadFile, loadEnv)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewpadFile, gotReviewpadFile)
}
