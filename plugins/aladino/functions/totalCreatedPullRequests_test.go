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
	"github.com/stretchr/testify/assert"
)

var totalCreatedPullRequests = plugins_aladino.PluginBuiltIns().Functions["totalCreatedPullRequests"].Code

func TestTotalCreatedPullRequests_WhenListIssuesByRepoRequestFails(t *testing.T) {
	devName := "steve"
	failMessage := "ListListIssuesByRepoRequestFail"
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
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
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue(devName)}
	gotTotal, err := totalCreatedPullRequests(mockedEnv, args)

	assert.Nil(t, gotTotal)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestTotalCreatedPullRequests_WhenThereIsPullRequestIssues(t *testing.T) {
	devName := "steve"
	ghIssues := []*github.Issue{
		{
			Title:            github.String("First Issue"),
			PullRequestLinks: nil,
		},
		{
			Title: github.String("Second Issue"),
			PullRequestLinks: &github.PullRequestLinks{
				URL: github.String("pull-request-link"),
			},
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesByOwnerByRepo,
				ghIssues,
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	wantTotal := aladino.BuildIntValue(1)

	args := []aladino.Value{aladino.BuildStringValue(devName)}
	gotTotal, err := totalCreatedPullRequests(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantTotal, gotTotal)
}
