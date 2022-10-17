// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var deleteHeadBranch = plugins_aladino.PluginBuiltIns().Actions["deleteHeadBranch"].Code

func TestDeleteHeadBranch_WhenRequestFails(t *testing.T) {
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Merged: github.Bool(true),
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatchHandler(
				mock.DeleteReposGitRefsByOwnerByRepoByRef,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					http.Error(w, "reference not found", http.StatusUnprocessableEntity)
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	err := deleteHeadBranch(mockedEnv, []aladino.Value{})

	assert.NotNil(t, err)
}

func TestDeleteHeadBranch_WhenPRIsNotMergedOrClosed(t *testing.T) {
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Merged:   github.Bool(false),
		ClosedAt: nil,
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	err := deleteHeadBranch(mockedEnv, []aladino.Value{})

	assert.NotNil(t, err)

	assert.Equal(t, err, errors.New("pull request should be merged or closed before deleting head branch"))
}

func TestDeleteHeadBranch_SuccessMerged(t *testing.T) {
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Merged:   github.Bool(true),
		ClosedAt: nil,
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatchHandler(
				mock.DeleteReposGitRefsByOwnerByRepoByRef,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusNoContent)
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	err := deleteHeadBranch(mockedEnv, []aladino.Value{})

	assert.Nil(t, err)
}

func TestDeleteHeadBranch_SuccessClosed(t *testing.T) {
	now := time.Now()
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Merged:   github.Bool(false),
		ClosedAt: &now,
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatchHandler(
				mock.DeleteReposGitRefsByOwnerByRepoByRef,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusNoContent)
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	err := deleteHeadBranch(mockedEnv, []aladino.Value{})

	assert.Nil(t, err)
}
