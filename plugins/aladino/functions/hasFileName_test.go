// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	pbe "github.com/reviewpad/api/go/entities"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var hasFileName = plugins_aladino.PluginBuiltIns().Functions["hasFileName"].Code

func TestHasFileName_WhenTrue(t *testing.T) {
	defaultMockPrFileName := "default-mock-repo/file1.ts"
	mockedPullRequestFileList := []*pbe.CommitFile{
		{
			Filename: defaultMockPrFileName,
		},
	}
	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		nil,
		nil,
		aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
			Files: mockedPullRequestFileList,
		}),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue(defaultMockPrFileName)}
	gotVal, err := hasFileName(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(true)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestHasFileName_WhenFalse(t *testing.T) {
	defaultMockPrFileName := "default-mock-repo/file1.ts"
	mockedPullRequestFileList := []*pbe.CommitFile{
		{
			Filename: "default-mock-repo/file2.ts",
		},
	}
	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		nil,
		nil,
		aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
			Files: mockedPullRequestFileList,
		}),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue(defaultMockPrFileName)}
	gotVal, err := hasFileName(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(false)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
