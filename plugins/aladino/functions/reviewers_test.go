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
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v2/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var reviewers = plugins_aladino.PluginBuiltIns().Functions["reviewers"].Code

func TestReviewers(t *testing.T) {
	ghUsersReviewers := []*github.User{
		{Login: github.String("mary")},
	}
	ghTeamReviewers := []*github.Team{
		{Slug: github.String("reviewpad")},
	}
	mockedPullRequest := mocks_aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User: &github.User{
			Login: github.String("john"),
		},
		RequestedReviewers: ghUsersReviewers,
		RequestedTeams:     ghTeamReviewers,
	})
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatchHandler(
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write(mock.MustMarshal(mockedPullRequest))
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	allReviewers := make([]aladino.Value, len(ghUsersReviewers)+len(ghTeamReviewers))
	for i, ghUser := range ghUsersReviewers {
		allReviewers[i] = aladino.BuildStringValue(ghUser.GetLogin())
	}

	for i, ghTeam := range ghTeamReviewers {
		allReviewers[i+len(ghUsersReviewers)] = aladino.BuildStringValue(ghTeam.GetSlug())
	}

	wantReviewers := aladino.BuildArrayValue(allReviewers)

	args := []aladino.Value{}
	gotReviewers, err := reviewers(mockedEnv, args)

	assert.Nil(t, err)
    assert.Equal(t, wantReviewers, gotReviewers)
}
