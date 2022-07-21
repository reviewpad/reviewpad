// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"log"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v3/mocks/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var milestone = plugins_aladino.PluginBuiltIns().Functions["milestone"].Code

func TestMilestone(t *testing.T) {
	milestoneTitle := "v1.0"
	mockedPullRequest := mocks_aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Milestone: &github.Milestone{
			Title: github.String(milestoneTitle),
		},
	})
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{}
	gotMilestoneTitle, err := milestone(mockedEnv, args)

	wantMilestoneTitle := aladino.BuildStringValue(milestoneTitle)

	assert.Nil(t, err)
	assert.Equal(t, wantMilestoneTitle, gotMilestoneTitle)
}
