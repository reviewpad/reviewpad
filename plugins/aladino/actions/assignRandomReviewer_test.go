// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var assignRandomReviewer = plugins_aladino.PluginBuiltIns().Actions["assignRandomReviewer"].Code

type ReviewersRequestPostBody struct {
	Reviewers []string `json:"reviewers"`
}

func TestAssignRandomReviewer_WhenListReviewersRequestFails(t *testing.T) {
	failMessage := "ListReviewersRequestFail"
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
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

	args := []aladino.Value{}
	err := assignRandomReviewer(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestAssignRandomReviewer_WhenPullRequestAlreadyHasReviewers(t *testing.T) {
	var isListAssigneesFetched bool
	requestedReviewer := "jane"
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				github.Reviewers{
					Users: []*github.User{
						{Login: github.String(requestedReviewer)},
					},
				},
			),
			mock.WithRequestMatchHandler(
				mock.GetReposAssigneesByOwnerByRepo,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// When a pull request has no reviewers then the list of available assignees is fetched
					isListAssigneesFetched = true
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{}
	err := assignRandomReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.False(t, isListAssigneesFetched, "shouldn't fetch the list of available assignees since the pull request already has reviewers")
}

func TestAssignRandomReviewer_WhenListAssigneesRequestFails(t *testing.T) {
	failMessage := "ListAssigneesRequestFail"
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				github.Reviewers{},
			),
			mock.WithRequestMatchHandler(
				mock.GetReposAssigneesByOwnerByRepo,
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

	args := []aladino.Value{}
	err := assignRandomReviewer(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestAssignRandomReviewer_ShouldFilterPullRequestAuthor(t *testing.T) {
	selectedReviewers := []string{}
	authorLogin := "maria"
	assigneeLogin := "peter"
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Author: &pbc.User{Login: authorLogin},
	})
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				github.Reviewers{},
			),
			mock.WithRequestMatch(
				mock.GetReposAssigneesByOwnerByRepo,
				[]*github.User{
					{Login: github.String(authorLogin)},
					{Login: github.String(assigneeLogin)},
				},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := io.ReadAll(r.Body)
					body := ReviewersRequestPostBody{}

					utils.MustUnmarshal(rawBody, &body)

					selectedReviewers = body.Reviewers
				}),
			),
		},
		nil,
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{}
	err := assignRandomReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{assigneeLogin}, selectedReviewers)
}

func TestAssignRandomReviewer_WhenThereIsNoUsers(t *testing.T) {
	selectedReviewers := []string{}
	authorLogin := "maria"
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Author: &pbc.User{Login: authorLogin},
	})
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				github.Reviewers{},
			),
			mock.WithRequestMatch(
				mock.GetReposAssigneesByOwnerByRepo,
				[]*github.User{
					{Login: github.String(authorLogin)},
				},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := io.ReadAll(r.Body)
					body := ReviewersRequestPostBody{}

					utils.MustUnmarshal(rawBody, &body)

					selectedReviewers = body.Reviewers
				}),
			),
		},
		nil,
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{}
	err := assignRandomReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{}, selectedReviewers)
}
