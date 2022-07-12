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

var totalCreatedPullRequests = plugins_aladino.PluginBuiltIns().Functions["totalCreatedPullRequests"].Code

func TestTotalCreatedPullRequests_WhenListIssuesByRepoRequestFails(t *testing.T) {
	devName := "steve"

	failMessage := "ListListIssuesByRepoRequestFail"
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatchHandler(
			mock.GetReposIssuesByOwnerByRepo,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mock.WriteError(
					w,
					http.StatusInternalServerError,
					failMessage,
				)
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue(devName)}
	gotTotal, err := totalCreatedPullRequests(mockedEnv, args)

	assert.Nil(t, gotTotal)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestTotalCreatedPullRequests_WhenThereIsNonPullRequestIssues(t *testing.T) {
	devName := "steve"
    ghIssues := []*github.Issue{
        {
            Title: github.String("First Issue"),
            PullRequestLinks: nil,
            
        },
        {
            Title: github.String("Second Issue"),
            PullRequestLinks: &github.PullRequestLinks{
                URL: github.String("pull-request-link"),
            },
        },
    }
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatch(
			mock.GetReposIssuesByOwnerByRepo,
			ghIssues,
		),
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

    wantTotal := aladino.BuildIntValue(1)

	args := []aladino.Value{aladino.BuildStringValue(devName)}
	gotTotal, err := totalCreatedPullRequests(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantTotal, gotTotal)
}
