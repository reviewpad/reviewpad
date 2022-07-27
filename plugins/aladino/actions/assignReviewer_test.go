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
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var assignReviewer = plugins_aladino.PluginBuiltIns().Actions["assignReviewer"].Code

func TestAssignReviewer_WhenTotalRequiredReviewersIsZero(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue("jane"),
			},
		),
		aladino.BuildIntValue(0),
	}
	err = assignReviewer(mockedEnv, args)

	assert.EqualError(t, err, "assignReviewer: total required reviewers can't be 0")
}

func TestAssignReviewer_WhenListOfReviewersIsEmpty(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{}), aladino.BuildIntValue(1)}
	err = assignReviewer(mockedEnv, args)

	assert.EqualError(t, err, "assignReviewer: list of reviewers can't be empty")
}

func TestAssignReviewer_WhenAuthorIsInListOfReviewers(t *testing.T) {
	var gotReviewers []string
	authorLogin := "john"
	reviewerLogin := "mary"
	wantReviewers := []string{
		reviewerLogin,
	}
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User:               &github.User{Login: github.String(authorLogin)},
		RequestedReviewers: []*github.User{},
	})
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				[]*github.PullRequestReview{},
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
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue(authorLogin),
				aladino.BuildStringValue(reviewerLogin),
			},
		),
		aladino.BuildIntValue(1),
	}
	err = assignReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantReviewers, gotReviewers, "pr author shouldn't be assigned as a reviewer")
}

func TestAssignReviewer_WhenTotalRequiredReviewersIsMoreThanTotalAvailableReviewers(t *testing.T) {
	var gotReviewers []string
	reviewerLogin := "mary"
	authorLogin := "john"
	totalRequiredReviewers := 2
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User:               &github.User{Login: github.String(authorLogin)},
		RequestedReviewers: []*github.User{},
	})
	wantReviewers := []string{
		reviewerLogin,
	}
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				[]*github.PullRequestReview{},
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
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue(reviewerLogin),
			},
		),
		aladino.BuildIntValue(totalRequiredReviewers),
	}
	err = assignReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantReviewers, gotReviewers, "the list of assign reviewers should be all provided reviewers")
}

func TestAssignReviewer_WhenListReviewsRequestFails(t *testing.T) {
	failMessage := "ListReviewsRequestFail"
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
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
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue("jane"),
			},
		),
		aladino.BuildIntValue(3),
	}
	err = assignReviewer(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestAssignReviewer_WhenPullRequestAlreadyHasReviews(t *testing.T) {
	var gotReviewers []string
	authorLogin := "john"
	reviewerLogin := "mary"
	wantReviewers := []string{
		reviewerLogin,
	}
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User:               &github.User{Login: github.String(authorLogin)},
		RequestedReviewers: []*github.User{},
	})
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				[]*github.PullRequestReview{
					{
						User: &github.User{
							Login: github.String(reviewerLogin),
						},
						State: github.String("COMMENTED"),
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
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue(reviewerLogin),
			},
		),
		aladino.BuildIntValue(1),
	}
	err = assignReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantReviewers, gotReviewers)
}

func TestAssignReviewer_WhenPullRequestAlreadyHasRequestedReviewers(t *testing.T) {
	var gotReviewers []string
	authorLogin := "john"
	reviewerA := "mary"
	reviewerB := "steve"
	wantReviewers := []string{
		reviewerB,
	}
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User: &github.User{Login: github.String(authorLogin)},
		RequestedReviewers: []*github.User{
			{Login: github.String(reviewerA)},
		},
	})
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				[]*github.PullRequestReview{},
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
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue(reviewerA),
				aladino.BuildStringValue(reviewerB),
			},
		),
		aladino.BuildIntValue(2),
	}
	err = assignReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantReviewers, gotReviewers, "when a reviewer already has a requested review, then it shouldn't be re-requested")
}

// Test scenario description:
// The mocked pull request has an assigned reviewer ("mary") which hasn't made any review yet
// The provided reviewers list contains an already assigned reviewer ("mary")
// Since a review has already been requested to the already assigned reviewer, so there's no available reviewers left
func TestAssignReviewer_HasNoAvailableReviewers(t *testing.T) {
	var isRequestReviewersRequestPerformed bool
	authorLogin := "john"
	reviewerLogin := "mary"
	totalRequiredReviewers := 1
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User: &github.User{Login: github.String(authorLogin)},
		RequestedReviewers: []*github.User{
			{Login: github.String(reviewerLogin)},
		},
	})
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				[]*github.PullRequestReview{},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// If the request reviewers request was performed then the reviewers were assigned to the pull request
					isRequestReviewersRequestPerformed = true
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue(reviewerLogin),
			},
		),
		aladino.BuildIntValue(totalRequiredReviewers),
	}
	err = assignReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.False(t, isRequestReviewersRequestPerformed, "the action shouldn't request for reviewers")
}

func TestAssignReviewer_WhenPullRequestAlreadyApproved(t *testing.T) {
	var isRequestReviewersRequestPerformed bool
	authorLogin := "john"
	reviewerA := "mary"
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User: &github.User{Login: github.String(authorLogin)},
		RequestedReviewers: []*github.User{
			{Login: github.String(reviewerA)},
		},
	})
	reviewID := int64(1)
	reviewState := "APPROVED"
	mockedReviews := []*github.PullRequestReview{
		{
			ID:    &reviewID,
			State: &reviewState,
			User:  &github.User{Login: github.String(reviewerA)},
		},
	}
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				mockedReviews,
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// If the request reviewers request was performed then the reviewers were assigned to the pull request
					isRequestReviewersRequestPerformed = true
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue(reviewerA),
			},
		),
		aladino.BuildIntValue(1),
	}
	err = assignReviewer(mockedEnv, args)
	assert.Nil(t, err)
	assert.False(t, isRequestReviewersRequestPerformed, "the action shouldn't request for reviewers")
}
