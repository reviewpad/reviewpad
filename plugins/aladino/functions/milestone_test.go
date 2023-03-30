// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var milestone = plugins_aladino.PluginBuiltIns().Functions["milestone"].Code

func TestMilestone(t *testing.T) {
	milestoneTitle := "v1.0"
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Milestone: &pbc.Milestone{
			Title: milestoneTitle,
		},
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

	args := []aladino.Value{}
	gotMilestoneTitle, err := milestone(mockedEnv, args)

	wantMilestoneTitle := aladino.BuildStringValue(milestoneTitle)

	assert.Nil(t, err)
	assert.Equal(t, wantMilestoneTitle, gotMilestoneTitle)
}
