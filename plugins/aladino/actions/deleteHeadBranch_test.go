// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"encoding/json"
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

func TestDeleteHeadBranch(t *testing.T) {
	now := time.Now()
	isDeleteHeadBranchRequestPerformed := false
	testCases := []struct {
		name                    string
		clientOptions           []mock.MockBackendOption
		deleteShouldBePerformed bool
		err                     error
	}{
		{
			name: "success: pull request is closed",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.Write(mock.MustMarshal(aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
							Merged:   github.Bool(false),
							ClosedAt: &now,
						})))
					}),
				),
				mock.WithRequestMatchHandler(
					mock.DeleteReposGitRefsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						isDeleteHeadBranchRequestPerformed = true
						w.WriteHeader(http.StatusNoContent)
					}),
				),
			},
			deleteShouldBePerformed: true,
		},
		{
			name: "success: pull request is merged",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.Write(mock.MustMarshal(aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
							Merged:   github.Bool(false),
							ClosedAt: &now,
						})))
					}),
				),
				mock.WithRequestMatchHandler(
					mock.DeleteReposGitRefsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.WriteHeader(http.StatusNoContent)
						isDeleteHeadBranchRequestPerformed = true
					}),
				),
			},
			deleteShouldBePerformed: true,
		},
		{
			name: "success: pull request is merged and closed",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.Write(mock.MustMarshal(aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
							Merged:   github.Bool(true),
							ClosedAt: &now,
						})))
					}),
				),
				mock.WithRequestMatchHandler(
					mock.DeleteReposGitRefsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.WriteHeader(http.StatusNoContent)
						isDeleteHeadBranchRequestPerformed = true
					}),
				),
			},
			deleteShouldBePerformed: true,
		},
		{
			name: "success: pull request is closed but not merged",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.Write(mock.MustMarshal(aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
							Merged:   github.Bool(false),
							ClosedAt: &now,
						})))
					}),
				),
				mock.WithRequestMatchHandler(
					mock.DeleteReposGitRefsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.WriteHeader(http.StatusNoContent)
						isDeleteHeadBranchRequestPerformed = true
					}),
				),
			},
			deleteShouldBePerformed: true,
		},
		{
			name: "success: pull request is not merged or closed",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.Write(mock.MustMarshal(aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
							Merged:   github.Bool(false),
							ClosedAt: nil,
						})))
					}),
				),
				mock.WithRequestMatchHandler(
					mock.DeleteReposGitRefsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.WriteHeader(http.StatusNoContent)
					}),
				),
			},
			err: nil,
		},
		{
			name: "error: request failed",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.Write(mock.MustMarshal(aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
							Merged: github.Bool(true),
						})))
					}),
				),
				mock.WithRequestMatchHandler(
					mock.DeleteReposGitRefsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.WriteHeader(http.StatusUnprocessableEntity)
						json.NewEncoder(w).Encode(map[string]string{
							"message":           "Reference does not exist",
							"documentation_url": "https://docs.github.com/rest/reference/git#delete-a-reference",
						})
					}),
				),
			},
			err: &github.ErrorResponse{
				Message:          "Reference does not exist",
				DocumentationURL: "https://docs.github.com/rest/reference/git#delete-a-reference",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(
				t,
				testCase.clientOptions,
				nil,
				aladino.MockBuiltIns(),
				nil,
			)

			err := deleteHeadBranch(mockedEnv, []aladino.Value{})

			// this allows simplified checking of github error response equality
			if e, ok := err.(*github.ErrorResponse); ok {
				e.Response = nil
			}

			assert.Equal(t, testCase.err, err)

			assert.Equal(t, testCase.deleteShouldBePerformed, isDeleteHeadBranchRequestPerformed)

			isDeleteHeadBranchRequestPerformed = false
		})
	}
}
