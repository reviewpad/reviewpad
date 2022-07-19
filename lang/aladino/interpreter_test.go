// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v2/engine"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	"github.com/stretchr/testify/assert"
)

func TestReport_WhenFindReportCommentFails(t *testing.T) {
	mockedPullRequest := mocks_aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User: &github.User{Login: github.String("john")},
		Base: &github.PullRequestBranch{
			Repo: &github.Repository{
				Owner: &github.User{
					Login: github.String("john"),
				},
				Name: github.String("default-mock-repo"),
			},
			Ref: github.String("master"),
		},
	})
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedInterpreter := &aladino.Interpreter{
		Env: mockedEnv,
	}

	err = mockedInterpreter.Report(engine.SILENT_MODE)

	assert.EqualError(t, err, "[report] error getting issues mock response not found for /repos/john/default-mock-repo/issues/6/comments")
}

func TestReport_OnSilentMode_WhenThereIsAlreadyAReviewpadComment(t *testing.T) {
	var isDeletedCommentRequested bool
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{
					{
						ID:   github.Int64(1234),
						Body: github.String("<!--@annotation-reviewpad-report-->\n**Reviewpad Report**\n\n:scroll: **Explanation**\nNo workflows activated"),
					},
				},
			),
			mock.WithRequestMatchHandler(
				mock.DeleteReposIssuesCommentsByOwnerByRepoByCommentId,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// When a pull request has a reviewpad comment then this comment is deleted
					isDeletedCommentRequested = true
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedInterpreter := &aladino.Interpreter{
		Env: mockedEnv,
	}

	err = mockedInterpreter.Report(engine.SILENT_MODE)

	assert.Nil(t, err)
	assert.True(t, isDeletedCommentRequested)
}

func TestReport_OnSilentMode_WhenNoReviewpadCommentIsFound(t *testing.T) {
	var isDeletedCommentRequested bool
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{},
			),
			mock.WithRequestMatchHandler(
				mock.DeleteReposIssuesCommentsByOwnerByRepoByCommentId,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// When a pull request has a reviewpad comment then this comment is deleted
					isDeletedCommentRequested = true
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedInterpreter := &aladino.Interpreter{
		Env: mockedEnv,
	}

	err = mockedInterpreter.Report(engine.SILENT_MODE)

	assert.Nil(t, err)
	assert.False(t, isDeletedCommentRequested)
}

func TestReport_OnVerboseMode_WhenNoReviewpadCommentIsFound(t *testing.T) {
	var addedComment string
	commentToBeAdded := "<!--@annotation-reviewpad-report-->\n**Reviewpad Report**\n\n:scroll: **Explanation**\nNo workflows activated"
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := github.IssueComment{}

					json.Unmarshal(rawBody, &body)

					addedComment = *body.Body
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedInterpreter := &aladino.Interpreter{
		Env: mockedEnv,
	}

	err = mockedInterpreter.Report(engine.VERBOSE_MODE)

	assert.Nil(t, err)
	assert.Equal(t, commentToBeAdded, addedComment)
}

func TestReport_OnVerboseMode_WhenThereIsAlreadyAReviewpadComment(t *testing.T) {
	var updatedComment string
	commentUpdated := "<!--@annotation-reviewpad-report-->\n**Reviewpad Report**\n\n:scroll: **Explanation**\nNo workflows activated"
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{
					{
						ID:   github.Int64(1234),
						Body: github.String("<!--@annotation-reviewpad-report-->\n**Reviewpad Report**\n**:information_source: Messages**\n* Changes the README.md"),
					},
				},
			),
			mock.WithRequestMatchHandler(
				mock.PatchReposIssuesCommentsByOwnerByRepoByCommentId,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := github.IssueComment{}

					json.Unmarshal(rawBody, &body)

					updatedComment = *body.Body
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedInterpreter := &aladino.Interpreter{
		Env: mockedEnv,
	}

	err = mockedInterpreter.Report(engine.VERBOSE_MODE)

	assert.Nil(t, err)
	assert.Equal(t, commentUpdated, updatedComment)
}
