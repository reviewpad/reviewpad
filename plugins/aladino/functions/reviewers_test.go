// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var reviewers = plugins_aladino.PluginBuiltIns(plugins_aladino.DefaultPluginConfig()).Functions["reviewers"].Code

func TestReviewers(t *testing.T) {
	ghUsersReviewers := []*github.User{
		{Login: github.String("mary")},
	}
	ghTeamReviewers := []*github.Team{
		{Slug: github.String("reviewpad")},
	}
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User: &github.User{
			Login: github.String("john"),
		},
		RequestedReviewers: ghUsersReviewers,
		RequestedTeams:     ghTeamReviewers,
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	wantReviewers := aladino.BuildArrayValue([]aladino.Value{
		aladino.BuildStringValue("mary"),
		aladino.BuildStringValue("reviewpad"),
	})

	args := []aladino.Value{}
	gotReviewers, err := reviewers(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewers, gotReviewers)
}
