// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var hasFileName = plugins_aladino.PluginBuiltIns().Functions["hasFileName"].Code

func TestHasFileName_WhenTrue(t *testing.T) {
	defaultMockPrFileName := "default-mock-repo/file1.ts"
	mockedPullRequestFileList := []*github.CommitFile{
		{
			Filename: github.String(defaultMockPrFileName),
		},
	}
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		nil,
		nil,
		aladino.GetDefaultPullRequestDetails(),
		mockedPullRequestFileList,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{lang.BuildStringValue(defaultMockPrFileName)}
	gotVal, err := hasFileName(mockedEnv, args)

	wantVal := lang.BuildBoolValue(true)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestHasFileName_WhenFalse(t *testing.T) {
	defaultMockPrFileName := "default-mock-repo/file1.ts"
	mockedPullRequestFileList := []*github.CommitFile{
		{
			Filename: github.String("default-mock-repo/file2.ts"),
		},
	}
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		nil,
		nil,
		aladino.GetDefaultPullRequestDetails(),
		mockedPullRequestFileList,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{lang.BuildStringValue(defaultMockPrFileName)}
	gotVal, err := hasFileName(mockedEnv, args)

	wantVal := lang.BuildBoolValue(false)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
