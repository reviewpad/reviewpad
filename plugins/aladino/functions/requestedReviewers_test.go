// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var requestedReviewers = plugins_aladino.PluginBuiltIns().Functions["requestedReviewers"].Code

func TestRequestedReviewers(t *testing.T) {
	ghRequestedUsersReviewers := []*github.User{
		{Login: github.String("mary")},
	}
	ghRequestedTeamReviewers := []*github.Team{
		{Slug: github.String("reviewpad")},
	}
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User: &github.User{
			Login: github.String("john"),
		},
		RequestedReviewers: ghRequestedUsersReviewers,
		RequestedTeams:     ghRequestedTeamReviewers,
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequest))
				}),
			),
		},
		nil,
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
