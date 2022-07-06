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
	failMessage := "ListReviewerslRequestFail"
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

func TestAssignRandomReviewer_WhenPullRequestHasAssignedReviewers(t *testing.T) {
	var noAssignedReviewers bool
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatch(
			mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
			github.Reviewers{
				Users: []*github.User{
					{Login: github.String("jane")},
				},
			},
		),
		mock.WithRequestMatchHandler(
			mock.GetReposAssigneesByOwnerByRepo,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// After checking if the pull request has already assigned reviewers, if not, we perform a request to list available assignees
				noAssignedReviewers = true
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{}
	err = assignRandomReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.False(t, noAssignedReviewers, "The pull request should not have already assigned reviewers")
}

func TestAssignRandomReviewer_WhenListAssigneesRequestFails(t *testing.T) {
	failMessage := "ListAssigneeslRequestFail"
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

func TestAssignRandomReviewer_WhenThereIsNoUsers(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatch(
			mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
			github.Reviewers{},
		),
		mock.WithRequestMatch(
			mock.GetReposAssigneesByOwnerByRepo,
			[]*github.User{
				// Author of the mocked pull request
				{Login: github.String("john")},
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

func TestAssignRandomReviewer(t *testing.T) {
	selectedReviewer := []string{}
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatch(
			mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
			github.Reviewers{},
		),
		mock.WithRequestMatch(
			mock.GetReposAssigneesByOwnerByRepo,
			[]*github.User{
				// Author of the mocked pull request
				{Login: github.String("john")},
				{Login: github.String("mary")},
			},
		),
		mock.WithRequestMatchHandler(
			mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				rawBody, _ := ioutil.ReadAll(r.Body)
				body := ReviewersRequestPostBody{}

				json.Unmarshal(rawBody, &body)

				selectedReviewer = body.Reviewers
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{}
	err = assignRandomReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(selectedReviewer))
	assert.Equal(t, "mary", selectedReviewer[0])
}
