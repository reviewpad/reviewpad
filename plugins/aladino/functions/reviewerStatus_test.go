// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var reviewerStatus = plugins_aladino.PluginBuiltIns(nil).Functions["reviewerStatus"].Code

func TestReviewerStatus_WhenRequestFails(t *testing.T) {
	failMessage := "ReviewerStatusRequestFail"
	mockedEnv := aladino.MockDefaultEnv(
		t,
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
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue("mary")}
	gotReviewState, err := reviewerStatus(mockedEnv, args)

	assert.Nil(t, gotReviewState)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestReviewerStatus_WhenUserIsNil(t *testing.T) {
	reviews := []*github.PullRequestReview{
		{
			State: github.String("COMMENTED"),
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				reviews,
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	wantReviewState := aladino.BuildStringValue("")

	args := []aladino.Value{aladino.BuildStringValue("mary")}
	gotReviewState, err := reviewerStatus(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewState, gotReviewState)
}

func TestReviewerStatus_WhenStateIsNil(t *testing.T) {
	reviews := []*github.PullRequestReview{
		{
			User: &github.User{
				Login: github.String("john"),
			},
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				reviews,
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	wantReviewState := aladino.BuildStringValue("")

	args := []aladino.Value{aladino.BuildStringValue("mary")}
	gotReviewState, err := reviewerStatus(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewState, gotReviewState)
}

func TestReviewerStatus_WhenStateUnknown(t *testing.T) {
	reviews := []*github.PullRequestReview{
		{
			State: github.String("COMMENTED"),
			User: &github.User{
				Login: github.String("john"),
			},
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				reviews,
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	wantReviewState := aladino.BuildStringValue("")

	args := []aladino.Value{aladino.BuildStringValue("mary")}
	gotReviewState, err := reviewerStatus(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewState, gotReviewState)
}

func TestReviewerStatus_WhenStateCommented(t *testing.T) {
	reviews := []*github.PullRequestReview{
		{
			State: github.String("COMMENTED"),
			User: &github.User{
				Login: github.String("mary"),
			},
		},
		{
			State: github.String("COMMENTED"),
			User: &github.User{
				Login: github.String("john"),
			},
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				reviews,
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	wantReviewState := aladino.BuildStringValue("COMMENTED")

	args := []aladino.Value{aladino.BuildStringValue("mary")}
	gotReviewState, err := reviewerStatus(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewState, gotReviewState)
}

func TestReviewerStatus_WhenStateApproved(t *testing.T) {
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
				Login: github.String("john"),
			},
		},
		{
			State: github.String("CHANGES_REQUESTED"),
			User: &github.User{
				Login: github.String("mary"),
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
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				reviews,
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	wantReviewState := aladino.BuildStringValue("APPROVED")

	args := []aladino.Value{aladino.BuildStringValue("mary")}
	gotReviewState, err := reviewerStatus(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewState, gotReviewState)
}

func TestReviewerStatus_WhenStateRequestedChanges(t *testing.T) {
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
		{
			State: github.String("COMMENTED"),
			User: &github.User{
				Login: github.String("mary"),
			},
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				reviews,
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	wantReviewState := aladino.BuildStringValue("CHANGES_REQUESTED")

	args := []aladino.Value{aladino.BuildStringValue("mary")}
	gotReviewState, err := reviewerStatus(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewState, gotReviewState)
}
