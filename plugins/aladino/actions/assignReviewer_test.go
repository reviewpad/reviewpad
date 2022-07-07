// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

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

func TestAssignReviewer_WhenAuthorIsRemovedFromProvidedReviewersList(t *testing.T) {
	var gotReviewers []string
	authorLogin := "john"
	reviewerLogin := "mary"
	wantReviewers := []string{
		reviewerLogin,
	}
	defaultPullRequestDetails := mocks_aladino.GetDefaultMockPullRequestDetails()
	defaultPullRequestDetails.User = &github.User{
		Login: github.String(authorLogin),
	}
	defaultPullRequestDetails.RequestedReviewers = []*github.User{}
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatchHandler(
			// Overwrite default mock to pull request request details
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write(mock.MustMarshal(defaultPullRequestDetails))
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
	assert.ElementsMatch(t, wantReviewers, gotReviewers)
}

func TestAssignReviewer_WhenTotalRequiredReviewersIsMoreThanTotalAvailableReviewers(t *testing.T) {
	var buf bytes.Buffer
	var gotReviewers []string
	reviewerLogin := "mary"
	authorLogin := "john"
	totalRequiredReviewers := 2
	defaultPullRequestDetails := mocks_aladino.GetDefaultMockPullRequestDetails()
	defaultPullRequestDetails.User = &github.User{
		Login: github.String(authorLogin),
	}
	defaultPullRequestDetails.RequestedReviewers = []*github.User{}
	wantReviewers := []string{
		reviewerLogin,
	}
	wantLog := fmt.Sprintf("%v assignReviewer: total required reviewers %v exceeds the total available reviewers %v\n", time.Now().Format("2006/01/02 15:04:05"), totalRequiredReviewers, len(wantReviewers))
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatchHandler(
			// Overwrite default mock to pull request request details
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write(mock.MustMarshal(defaultPullRequestDetails))
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
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	log.SetOutput(&buf)

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue(reviewerLogin),
			},
		),
		aladino.BuildIntValue(totalRequiredReviewers),
	}
	err = assignReviewer(mockedEnv, args)

	gotLog := buf.String()

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantReviewers, gotReviewers)
	assert.Equal(t, wantLog, gotLog)
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
	defaultPullRequestDetails := mocks_aladino.GetDefaultMockPullRequestDetails()
	defaultPullRequestDetails.User = &github.User{
		Login: github.String(authorLogin),
	}
	defaultPullRequestDetails.RequestedReviewers = []*github.User{}
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatchHandler(
			// Overwrite default mock to pull request request details
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write(mock.MustMarshal(defaultPullRequestDetails))
			}),
		),
		mock.WithRequestMatch(
			mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
			[]*github.PullRequestReview{
				{
					User: &github.User{
						Login: github.String(reviewerLogin),
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
	defaultPullRequestDetails := mocks_aladino.GetDefaultMockPullRequestDetails()
	defaultPullRequestDetails.User = &github.User{
		Login: github.String(authorLogin),
	}
	defaultPullRequestDetails.RequestedReviewers = []*github.User{
		{Login: github.String(reviewerA)},
	}
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatchHandler(
			// Overwrite default mock to pull request request details
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write(mock.MustMarshal(defaultPullRequestDetails))
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
	assert.ElementsMatch(t, wantReviewers, gotReviewers)
}

// The mocked pull request has an assigned reviewer ("mary") which hasn't made any review yet
// The provided reviewers list contains an already assigned reviewer ("mary")
// Since a review has already been requested to the already assigned reviewer, so there's no available reviewers left
func TestAssignReviewer_HasNoAvailableReviewers(t *testing.T) {
	var pullRequestHasReviewers bool
	authorLogin := "john"
	reviewerLogin := "mary"
	totalRequiredReviewers := 1
	defaultPullRequestDetails := mocks_aladino.GetDefaultMockPullRequestDetails()
	defaultPullRequestDetails.User = &github.User{
		Login: github.String(authorLogin),
	}
	defaultPullRequestDetails.RequestedReviewers = []*github.User{
		{Login: github.String(reviewerLogin)},
	}
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatchHandler(
			// Overwrite default mock to pull request request details
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write(mock.MustMarshal(defaultPullRequestDetails))
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
				aladino.BuildStringValue(reviewerLogin),
			},
		),
		aladino.BuildIntValue(totalRequiredReviewers),
	}
	err = assignReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.False(t, pullRequestHasReviewers, "The pull request should not have new reviewers assigned")
}
