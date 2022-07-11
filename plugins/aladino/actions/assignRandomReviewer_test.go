// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"encoding/json"
	"io/ioutil"
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

var assignRandomReviewer = plugins_aladino.PluginBuiltIns().Actions["assignRandomReviewer"].Code

type ReviewersRequestPostBody struct {
	Reviewers []string `json:"reviewers"`
}

func TestAssignRandomReviewer_WhenListReviewersRequestFails(t *testing.T) {
	failMessage := "ListReviewersRequestFail"
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
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
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{}
	err = assignRandomReviewer(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestAssignRandomReviewer_WhenPullRequestAlreadyHasReviewers(t *testing.T) {
	var isListAssigneesFetched bool
	requestedReviewer := "jane"
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
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
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{}
	err = assignRandomReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.False(t, isListAssigneesFetched, "shouldn't fetch the list of available assignees since the pull request already has reviewers")
}

func TestAssignRandomReviewer_WhenListAssigneesRequestFails(t *testing.T) {
	failMessage := "ListAssigneesRequestFail"
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
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
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{}
	err = assignRandomReviewer(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestAssignRandomReviewer_ShouldFilterPullRequestAuthor(t *testing.T) {
	selectedReviewers := []string{}
	authorLogin := "maria"
	assigneeLogin := "peter"
    mockedPullRequest := mocks_aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User:               &github.User{Login: github.String(authorLogin)},
	})
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatchHandler(
			// Overwrite default mock to pull request request details
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
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{}
	err = assignRandomReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{assigneeLogin}, selectedReviewers)
}

func TestAssignRandomReviewer_WhenThereIsNoUsers(t *testing.T) {
	authorLogin := "maria"
    mockedPullRequest := mocks_aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User:               &github.User{Login: github.String(authorLogin)},
	})
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
        mock.WithRequestMatchHandler(
			// Overwrite default mock to pull request request details
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
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{}
	err = assignRandomReviewer(mockedEnv, args)

	assert.EqualError(t, err, "can't assign a random user because there is no users")
}
