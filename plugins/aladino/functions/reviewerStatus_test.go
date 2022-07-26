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
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var reviewerStatus = plugins_aladino.PluginBuiltIns().Functions["reviewerStatus"].Code

func TestReviewerStatusRequestFails(t *testing.T) {
	failMessage := "ReviewerStatusRequestFail"
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

	args := []aladino.Value{aladino.BuildStringValue("mary")}
	gotReviewState, err := reviewerStatus(mockedEnv, args)

	assert.Nil(t, gotReviewState)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestReviewerStatusNeutral(t *testing.T) {
	reviews := []*github.PullRequestReview{
		{
			State: github.String("COMMENTED"),
			User: &github.User{
				Login: github.String("john"),
			},
		},
	}
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				reviews,
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantReviewState := aladino.BuildStringValue("neutral")

	args := []aladino.Value{aladino.BuildStringValue("mary")}
	gotReviewState, err := reviewerStatus(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewState, gotReviewState)
}

func TestReviewerStatusApproved(t *testing.T) {
	reviews := []*github.PullRequestReview{
		{
			State: github.String("COMMENTED"),
			User: &github.User{
				Login: github.String("john"),
			},
		},
		{
			State: github.String("APPROVED"),
			User: &github.User{
				Login: github.String("mary"),
			},
		},
		{
			State: github.String("COMMENTED"),
			User: &github.User{
				Login: github.String("mary"),
			},
		},
	}
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				reviews,
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantReviewState := aladino.BuildStringValue("approved")

	args := []aladino.Value{aladino.BuildStringValue("mary")}
	gotReviewState, err := reviewerStatus(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewState, gotReviewState)
}

func TestReviewerStatusRequestedChanges(t *testing.T) {
	reviews := []*github.PullRequestReview{
		{
			State: github.String("COMMENTED"),
			User: &github.User{
				Login: github.String("john"),
			},
		},
		{
			State: github.String("APPROVED"),
			User: &github.User{
				Login: github.String("mary"),
			},
		},
		{
			State: github.String("CHANGES_REQUESTED"),
			User: &github.User{
				Login: github.String("mary"),
			},
		},
	}
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				reviews,
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantReviewState := aladino.BuildStringValue("requested_changes")

	args := []aladino.Value{aladino.BuildStringValue("mary")}
	gotReviewState, err := reviewerStatus(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewState, gotReviewState)
}
