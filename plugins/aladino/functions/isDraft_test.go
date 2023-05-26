// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var isDraft = plugins_aladino.PluginBuiltIns().Functions["isDraft"].Code

func TestIsDraft_WhenTrue(t *testing.T) {
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		IsDraft: true,
	})
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		nil,
		nil,
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{}
	gotVal, err := isDraft(mockedEnv, args)

	wantVal := lang.BuildBoolValue(true)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal, "pull request should be in draft")
}

func TestIsDraft_WhenFalse(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{}
	gotVal, err := isDraft(mockedEnv, args)

	wantVal := lang.BuildBoolValue(false)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal, "pull request should not be in draft")
}
