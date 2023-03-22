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

var commitCount = plugins_aladino.PluginBuiltIns().Functions["commitCount"].Code

func TestCommitCount(t *testing.T) {
	totalCommits := 1
	mockedCodeReview := aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
		CommitsCount: int64(totalCommits),
	})
	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		nil,
		nil,
		mockedCodeReview,
		aladino.MockBuiltIns(),
		nil,
	)

	wantCommitCount := aladino.BuildIntValue(totalCommits)

	args := []aladino.Value{}
	gotCommitCount, err := commitCount(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantCommitCount, gotCommitCount, "it should get the pull request commit count")
}
