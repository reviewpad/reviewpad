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
	mockMergedPR := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Merged:   github.Bool(true),
		ClosedAt: nil,
	})
	mockClosedPR := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Merged:   github.Bool(false),
		ClosedAt: &now,
	})
	mockMergedAndClosedPR := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Merged:   github.Bool(true),
		ClosedAt: &now,
	})
	mockNotMergedOrClosedPR := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Merged:   github.Bool(false),
		ClosedAt: nil,
	})
	isDeleteHeadBranchRequestPerformed := false

	testCases := []struct {
		name                    string
		mockPRHandler           mock.MockBackendOption
		mockDeleteRefHandler    mock.MockBackendOption
		deleteShouldBePerformed bool
		err                     error
	}{
		{
			name: "success: pull request is closed",
			mockPRHandler: mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write(mock.MustMarshal(mockClosedPR))
				}),
			),
			mockDeleteRefHandler: mock.WithRequestMatchHandler(
				mock.DeleteReposGitRefsByOwnerByRepoByRef,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					isDeleteHeadBranchRequestPerformed = true
					w.WriteHeader(http.StatusNoContent)
				}),
			),
			deleteShouldBePerformed: true,
		},
		{
			name: "success: pull request is merged",
			mockPRHandler: mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write(mock.MustMarshal(mockMergedPR))
				}),
			),
			mockDeleteRefHandler: mock.WithRequestMatchHandler(
				mock.DeleteReposGitRefsByOwnerByRepoByRef,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusNoContent)
					isDeleteHeadBranchRequestPerformed = true
				}),
			),
			deleteShouldBePerformed: true,
		},
		{
			name: "success: pull request is merged and closed",
			mockPRHandler: mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write(mock.MustMarshal(mockMergedAndClosedPR))
				}),
			),
			mockDeleteRefHandler: mock.WithRequestMatchHandler(
				mock.DeleteReposGitRefsByOwnerByRepoByRef,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusNoContent)
					isDeleteHeadBranchRequestPerformed = true
				}),
			),
			deleteShouldBePerformed: true,
		},
		{
			name: "success: pull request is closed but not merged",
			mockPRHandler: mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write(mock.MustMarshal(mockClosedPR))
				}),
			),
			mockDeleteRefHandler: mock.WithRequestMatchHandler(
				mock.DeleteReposGitRefsByOwnerByRepoByRef,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusNoContent)
					isDeleteHeadBranchRequestPerformed = true
				}),
			),
			deleteShouldBePerformed: true,
		},
		{
			name: "success: pull request is not merged or closed",
			mockPRHandler: mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write(mock.MustMarshal(mockNotMergedOrClosedPR))
				}),
			),
			mockDeleteRefHandler: mock.WithRequestMatchHandler(
				mock.DeleteReposGitRefsByOwnerByRepoByRef,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusNoContent)
				}),
			),
			err: nil,
		},
		{
			name: "error: request failed",
			mockPRHandler: mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write(mock.MustMarshal(mockMergedPR))
				}),
			),
			mockDeleteRefHandler: mock.WithRequestMatchHandler(
				mock.DeleteReposGitRefsByOwnerByRepoByRef,
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusUnprocessableEntity)
					json.NewEncoder(w).Encode(map[string]string{
						"message":           "Reference does not exist",
						"documentation_url": "https://docs.github.com/rest/reference/git#delete-a-reference",
					})
				}),
			),
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
				[]mock.MockBackendOption{
					testCase.mockPRHandler,
					testCase.mockDeleteRefHandler,
				},
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
