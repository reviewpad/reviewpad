// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"crypto/sha256"
	"fmt"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/jinzhu/copier"
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
	data, err := os.ReadFile("../testdata/engine/simpleReviewpadFile.yml")
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Error reading reviewpad file: %v", err))
	}

	wantReviewpadFile := &ReviewpadFile{}
	copier.Copy(wantReviewpadFile, simpleReviewpadFile)

	wantReviewpadFile.Imports = []PadImport{}

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
	wantReviewpadFile := mockedReviewpadFile

	mockedReviewpadFileData, err := os.ReadFile("../testdata/engine/mockedReviewpadFile.yml")
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Error getting reviewpad file: %v", err))
	}

	gotReviewpadFile, err := parse(mockedReviewpadFileData)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewpadFile, gotReviewpadFile)
}

func TestTransform(t *testing.T) {
	reviewpadFile := &ReviewpadFile{}
	copier.Copy(reviewpadFile, mockedReviewpadFile)
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
						Rule:         "tautology",
						ExtraActions: []string{"$addLabel(\"test-workflow-a\")"},
					},
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
		httpmock.NewBytesResponder(200, []byte("invalid-content")),
	)

	gotReviewpadFile, content, err := loadImport(reviewpadImport)

	assert.Nil(t, gotReviewpadFile)
	assert.Equal(t, "", content)
	assert.EqualError(t, err, "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `invalid...` into engine.ReviewpadFile")
}

func TestLoadImport(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	simpleReviewpadFileData, err := os.ReadFile("../testdata/engine/simpleReviewpadFile.yml")
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Error reading reviewpad file: %v", err))
	}

	simpleReviewpadFileImport := PadImport{
		Url: "https://foo.bar/simpleReviewpadFile.yml",
	}

	httpmock.RegisterResponder("GET", simpleReviewpadFileImport.Url,
		httpmock.NewBytesResponder(200, simpleReviewpadFileData),
	)

	wantReviewpadFile := simpleReviewpadFile

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
	copier.Copy(reviewpadFile, mockedReviewpadFile)

	invalidUrl := "https://foo.bar/invalid-url"
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
	assert.EqualError(t, err, "Get \"https://foo.bar/invalid-url\": invalid import url")
}

func TestInlineImports_WhenThereIsCycleDependency(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockedImportedReviewpadFileData, err := os.ReadFile("../testdata/engine/mockedImportedReviewpadFile.yml")
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Error getting reviewpad file: %v", err))
	}

	importUrl := "https://foo.bar/mockedImportedReviewpadFile.yml"

	reviewpadFile := &ReviewpadFile{}
	copier.Copy(reviewpadFile, mockedImportedReviewpadFile)

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

	mockedImportedReviewpadFileData, err := os.ReadFile("../testdata/engine/mockedImportedReviewpadFile.yml")
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Error getting reviewpad file: %v", err))
	}

	mockedImportedReviewpadFileUrl := "https://foo.bar/mockedImportedReviewpadFile.yml"

	reviewpadFile := &ReviewpadFile{}
	copier.Copy(reviewpadFile, mockedReviewpadFile)

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
	copier.Copy(wantReviewpadFile, mockedReviewpadFile)

	wantReviewpadFile.Imports = []PadImport{}

	gotReviewpadFile, err := inlineImports(reviewpadFile, loadEnv)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewpadFile, gotReviewpadFile)
}

func TestInlineImports_WhenSubTreeFileInlineImportsFails(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	simpleReviewpadFileWithImportsData, err := os.ReadFile("../testdata/engine/simpleReviewpadFileWithImports.yml")
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Error getting reviewpad file: %v", err))
	}

	mockedSimpleReviewpadFileUrl := "https://foo.bar/simpleReviewpadFile.yml"

	reviewpadFile := &ReviewpadFile{}
	copier.Copy(reviewpadFile, mockedReviewpadFile)

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

	simpleReviewpadFileData, err := os.ReadFile("../testdata/engine/simpleReviewpadFile.yml")
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Error getting reviewpad file: %v", err))
	}

	mockedSimpleReviewpadFileUrl := "https://foo.bar/simpleReviewpadFile.yml"

	reviewpadFile := &ReviewpadFile{}
	copier.Copy(reviewpadFile, mockedReviewpadFile)

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
			mockedReviewpadFilePadWorkflow,
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
