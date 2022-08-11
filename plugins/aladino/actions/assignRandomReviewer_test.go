// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var assignRandomReviewer = plugins_aladino.PluginBuiltIns(plugins_aladino.DefaultPluginConfig()).Actions["assignRandomReviewer"].Code

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
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User: &github.User{Login: github.String(authorLogin)},
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
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := ReviewersRequestPostBody{}

					json.Unmarshal(rawBody, &body)

					selectedReviewers = body.Reviewers
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
	assert.ElementsMatch(t, []string{assigneeLogin}, selectedReviewers)
}

func TestAssignRandomReviewer_WhenThereIsNoUsers(t *testing.T) {
	authorLogin := "maria"
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User: &github.User{Login: github.String(authorLogin)},
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
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{}
	err := assignRandomReviewer(mockedEnv, args)

	assert.EqualError(t, err, "can't assign a random user because there is no users")
}
