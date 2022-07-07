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

var assignReviewer = plugins_aladino.PluginBuiltIns().Actions["assignReviewer"].Code

func TestAssignReviewer_WhenTotalRequiredReviewersIsZero(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv()
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue("jane"),
				aladino.BuildStringValue("john"),
			},
		),
		aladino.BuildIntValue(0),
	}
	err = assignReviewer(mockedEnv, args)

	assert.EqualError(t, err, "assignReviewer: total required reviewers can't be 0")
}

func TestAssignReviewer_WhenListOfReviewersIsEmpty(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv()
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{}), aladino.BuildIntValue(1)}
	err = assignReviewer(mockedEnv, args)

	assert.EqualError(t, err, "assignReviewer: list of reviewers can't be empty")
}

func TestAssignReviewer_WhenListReviewsRequestFails(t *testing.T) {
	failMessage := "ListReviewsRequestFail"
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatchHandler(
			mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
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

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue("jane"),
				aladino.BuildStringValue("john"),
			},
		),
		aladino.BuildIntValue(3),
	}
	err = assignReviewer(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

// The mocked pull request has an assigned reviewer ("jane") which hasn't made any review yet
// The provided reviewers list contains the author of the pull request ("john") and an already assigned reviewer ("jane")
// Since the elements mentioned previously will be discarded, there will be no users which can be requested a review
func TestAssignReviewer_WhenPullRequestAlreadyHasReviewers(t *testing.T) {
	var pullRequestHasReviewers bool
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatch(
			mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
			[]*github.PullRequestReview{},
		),
		mock.WithRequestMatchHandler(
			mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// If the request reviewers request was performed then the reviewers were assigned to the pull request
				pullRequestHasReviewers = true
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue("jane"),
				aladino.BuildStringValue("john"),
			},
		),
		aladino.BuildIntValue(2),
	}
	err = assignReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.False(t, pullRequestHasReviewers, "The pull request should not have new reviewers assigned")
}

func TestAssignReviewer(t *testing.T) {
	wantReviewers := []string{
		"mary",
		"steve",
	}
	var gotReviewers []string
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatch(
			mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
			[]*github.PullRequestReview{
				{
					User: &github.User{
						Login: github.String("steve"),
					},
				},
			},
		),
		mock.WithRequestMatchHandler(
			mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				rawBody, _ := ioutil.ReadAll(r.Body)
				body := ReviewersRequestPostBody{}

				json.Unmarshal(rawBody, &body)

				gotReviewers = body.Reviewers
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue("jane"),
				aladino.BuildStringValue("john"),
				aladino.BuildStringValue("mary"),
				aladino.BuildStringValue("steve"),
			},
		),
		aladino.BuildIntValue(4),
	}
	err = assignReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, len(wantReviewers), len(gotReviewers))
	assert.Contains(t, wantReviewers, gotReviewers[0])
	assert.Contains(t, wantReviewers, gotReviewers[1])
}
