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

var requestedReviewers = plugins_aladino.PluginBuiltIns().Functions["requestedReviewers"].Code

func TestRequestedReviewers(t *testing.T) {
	mockedRequestedReviewers := &pbc.RequestedReviewers{
		Users: []*pbc.User{
			{
				Login: "mary",
			},
		},
		Teams: []*pbc.Team{
			{
				Slug: "reviewpad",
			},
		},
	}
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Author: &pbc.User{
			Login: "john",
		},
		RequestedReviewers: mockedRequestedReviewers,
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

	wantRequestedReviewers := aladino.BuildArrayValue([]aladino.Value{
		aladino.BuildStringValue("mary"),
		aladino.BuildStringValue("reviewpad"),
	})

	args := []aladino.Value{}
	gotRequestedReviewers, err := requestedReviewers(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantRequestedReviewers, gotRequestedReviewers)
}
