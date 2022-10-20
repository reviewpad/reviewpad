// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

// import (
// 	"encoding/json"
// 	"io/ioutil"
// 	"net/http"
// 	"testing"
// 	"time"

// 	"github.com/google/go-github/v45/github"
// 	"github.com/migueleliasweb/go-github-mock/src/mock"

// 	"github.com/reviewpad/reviewpad/v3/lang/aladino"
// 	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
// 	"github.com/stretchr/testify/assert"
// )

// var review = plugins_aladino.PluginBuiltIns().Actions["review"].Code

// type ReviewRequestPostBody struct {
// 	Event string `json:"event"`
// 	Body  string `json:"body"`
// }

// func TestReview_WhenReviewMethodIsUnsupported(t *testing.T) {
// 	mockedEnv := aladino.MockDefaultEnv(
// 		t,
// 		nil,
// 		nil,
// 		aladino.MockBuiltIns(),
// 		&github.WorkflowRunEvent{
// 			WorkflowRun: &github.WorkflowRun{
// 				Actor: &github.User{
// 					Login: github.String("reviewpad-bot"),
// 				},
// 			},
// 		},
// 	)

// 	args := []aladino.Value{aladino.BuildStringValue("INVALID"), aladino.BuildStringValue("some comment")}
// 	err := review(mockedEnv, args)

// 	assert.EqualError(t, err, "review: unsupported review event INVALID")
// }

// func TestReview_WhenListReviewsRequestFails(t *testing.T) {
// 	failMessage := "ListReviewsRequestFail"
// 	mockedEnv := aladino.MockDefaultEnv(
// 		t,
// 		[]mock.MockBackendOption{
// 			mock.WithRequestMatchHandler(
// 				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
// 				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 					mock.WriteError(
// 						w,
// 						http.StatusInternalServerError,
// 						failMessage,
// 					)
// 				}),
// 			),
// 		},
// 		nil,
// 		aladino.MockBuiltIns(),
// 		&github.WorkflowRunEvent{
// 			WorkflowRun: &github.WorkflowRun{
// 				Actor: &github.User{
// 					Login: github.String("reviewpad-bot"),
// 				},
// 			},
// 		},
// 	)

// 	args := []aladino.Value{aladino.BuildStringValue("APPROVE"), aladino.BuildStringValue("some comment")}
// 	err := review(mockedEnv, args)

// 	assert.Equal(t, failMessage, err.(*github.ErrorResponse).Message)
// }

// func TestReview_WhenBodyIsNotEmpty(t *testing.T) {
// 	reviewBody := "test comment"

// 	tests := map[string]struct {
// 		inputReviewEvent string
// 	}{
// 		"when review event is APPROVE": {
// 			inputReviewEvent: "APPROVE",
// 		},
// 		"when review event is REQUEST_CHANGES": {
// 			inputReviewEvent: "REQUEST_CHANGES",
// 		},
// 		"when review event is COMMENT": {
// 			inputReviewEvent: "COMMENT",
// 		},
// 	}

// 	for _, test := range tests {
// 		var gotReviewEvent, gotReviewBody string

// 		mockedEnv := aladino.MockDefaultEnv(
// 			t,
// 			[]mock.MockBackendOption{
// 				mock.WithRequestMatch(
// 					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
// 					[]*github.PullRequestReview{},
// 				),
// 				mock.WithRequestMatchHandler(
// 					mock.PostReposPullsReviewsByOwnerByRepoByPullNumber,
// 					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 						rawBody, _ := ioutil.ReadAll(r.Body)

// 						body := ReviewRequestPostBody{}

// 						json.Unmarshal(rawBody, &body)

// 						gotReviewEvent = body.Event
// 						gotReviewBody = body.Body
// 					}),
// 				),
// 			},
// 			nil,
// 			aladino.MockBuiltIns(),
// 			&github.WorkflowRunEvent{
// 				WorkflowRun: &github.WorkflowRun{
// 					Actor: &github.User{
// 						Login: github.String("reviewpad-bot"),
// 					},
// 				},
// 			},
// 		)

// 		args := []aladino.Value{aladino.BuildStringValue(test.inputReviewEvent), aladino.BuildStringValue(reviewBody)}

// 		err := review(mockedEnv, args)

// 		assert.Nil(t, err)
// 		assert.Equal(t, test.inputReviewEvent, gotReviewEvent)
// 		assert.Equal(t, reviewBody, gotReviewBody)
// 	}
// }

// func TestReview_WhenBodyIsEmpty(t *testing.T) {
// 	tests := map[string]struct {
// 		inputReviewEvent string
// 		wantErr          string
// 	}{
// 		"when review event is APPROVE": {
// 			inputReviewEvent: "APPROVE",
// 			wantErr:          "",
// 		},
// 		"when review event is REQUEST_CHANGES": {
// 			inputReviewEvent: "REQUEST_CHANGES",
// 			wantErr:          "review: comment required in REQUEST_CHANGES event",
// 		},
// 		"when review event is COMMENT": {
// 			inputReviewEvent: "COMMENT",
// 			wantErr:          "review: comment required in COMMENT event",
// 		},
// 	}

// 	for name, test := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			mockedEnv := aladino.MockDefaultEnv(
// 				t,
// 				[]mock.MockBackendOption{
// 					mock.WithRequestMatch(
// 						mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
// 						[]*github.PullRequestReview{},
// 					),
// 					mock.WithRequestMatch(
// 						mock.PostReposPullsReviewsByOwnerByRepoByPullNumber,
// 						&ReviewRequestPostBody{},
// 					),
// 				},
// 				nil,
// 				aladino.MockBuiltIns(),
// 				&github.WorkflowRunEvent{
// 					WorkflowRun: &github.WorkflowRun{
// 						Actor: &github.User{
// 							Login: github.String("reviewpad-bot"),
// 						},
// 					},
// 				},
// 			)

// 			args := []aladino.Value{aladino.BuildStringValue(test.inputReviewEvent), aladino.BuildStringValue("")}

// 			gotErr := review(mockedEnv, args)

// 			if gotErr != nil && gotErr.Error() != test.wantErr {
// 				assert.FailNow(t, "review() error = %v, wantErr %v", gotErr, test.wantErr)
// 			}
// 		})
// 	}
// }

// func TestReview_WhenExistsPreviousAutomaticReview(t *testing.T) {
// 	reviewBody := "test comment"
// 	reviewerLogin := "reviewpad-bot"

// 	pullRequestUpdate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
// 	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
// 		UpdatedAt: &pullRequestUpdate,
// 	})

// 	timeBeforePullRequestUpdate := time.Date(2019, 12, 29, 0, 0, 0, 0, time.UTC)
// 	timeAfterPullRequestUpdate := time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC)

// 	tests := map[string]struct {
// 		inputReviewEvent   string
// 		clientOptions      []mock.MockBackendOption
// 		shouldReviewBeDone bool
// 	}{
// 		"when last review event was APPROVE": {
// 			inputReviewEvent: "APPROVE",
// 			clientOptions: []mock.MockBackendOption{
// 				mock.WithRequestMatch(
// 					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
// 					[]*github.PullRequestReview{
// 						{
// 							ID:   github.Int64(1),
// 							Body: github.String(""),
// 							User: &github.User{
// 								Login: github.String(reviewerLogin),
// 							},
// 							State:       github.String("APPROVED"),
// 							SubmittedAt: &timeAfterPullRequestUpdate,
// 						},
// 					},
// 				),
// 			},
// 			shouldReviewBeDone: false,
// 		},
// 		"when last review event was REQUEST_CHANGES and the pull request has updates since then": {
// 			inputReviewEvent: "APPROVE",
// 			clientOptions: []mock.MockBackendOption{
// 				mock.WithRequestMatch(
// 					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
// 					[]*github.PullRequestReview{
// 						{
// 							ID:   github.Int64(1),
// 							Body: github.String(""),
// 							User: &github.User{
// 								Login: github.String(reviewerLogin),
// 							},
// 							State:       github.String("REQUEST_CHANGES"),
// 							SubmittedAt: &timeBeforePullRequestUpdate,
// 						},
// 					},
// 				),
// 			},
// 			shouldReviewBeDone: true,
// 		},
// 		"when last review event was REQUEST_CHANGES and the pull request has no updates since then": {
// 			inputReviewEvent: "APPROVE",
// 			clientOptions: []mock.MockBackendOption{
// 				mock.WithRequestMatch(
// 					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
// 					[]*github.PullRequestReview{
// 						{
// 							ID:   github.Int64(1),
// 							Body: github.String(""),
// 							User: &github.User{
// 								Login: github.String(reviewerLogin),
// 							},
// 							State:       github.String("REQUEST_CHANGES"),
// 							SubmittedAt: &timeAfterPullRequestUpdate,
// 						},
// 					},
// 				),
// 			},
// 			shouldReviewBeDone: false,
// 		},
// 		"when last review event was COMMENT and the pull request has updates since then": {
// 			inputReviewEvent: "APPROVE",
// 			clientOptions: []mock.MockBackendOption{
// 				mock.WithRequestMatch(
// 					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
// 					[]*github.PullRequestReview{
// 						{
// 							ID:   github.Int64(1),
// 							Body: github.String(""),
// 							User: &github.User{
// 								Login: github.String(reviewerLogin),
// 							},
// 							State:       github.String("COMMENTED"),
// 							SubmittedAt: &timeBeforePullRequestUpdate,
// 						},
// 					},
// 				),
// 			},
// 			shouldReviewBeDone: true,
// 		},
// 		"when last review event was COMMENT and the pull request has no updates since then": {
// 			inputReviewEvent: "APPROVE",
// 			clientOptions: []mock.MockBackendOption{
// 				mock.WithRequestMatch(
// 					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
// 					[]*github.PullRequestReview{
// 						{
// 							ID:   github.Int64(1),
// 							Body: github.String(""),
// 							User: &github.User{
// 								Login: github.String(reviewerLogin),
// 							},
// 							State:       github.String("COMMENTED"),
// 							SubmittedAt: &timeAfterPullRequestUpdate,
// 						},
// 					},
// 				),
// 			},
// 			shouldReviewBeDone: false,
// 		},
// 	}

// 	for _, test := range tests {
// 		var isPostReviewRequestPerformed bool
// 		mockedEnv := aladino.MockDefaultEnv(
// 			t,
// 			append(
// 				[]mock.MockBackendOption{
// 					mock.WithRequestMatchHandler(
// 						mock.GetReposPullsByOwnerByRepoByPullNumber,
// 						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
// 							w.Write(mock.MustMarshal(mockedPullRequest))
// 						}),
// 					),
// 					mock.WithRequestMatchHandler(
// 						mock.PostReposPullsReviewsByOwnerByRepoByPullNumber,
// 						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 							isPostReviewRequestPerformed = true
// 						}),
// 					),
// 				},
// 				test.clientOptions...,
// 			),
// 			nil,
// 			aladino.MockBuiltIns(),
// 			&github.WorkflowRunEvent{
// 				WorkflowRun: &github.WorkflowRun{
// 					Actor: &github.User{
// 						Login: github.String("reviewpad-bot"),
// 					},
// 				},
// 			},
// 		)

// 		args := []aladino.Value{aladino.BuildStringValue(test.inputReviewEvent), aladino.BuildStringValue(reviewBody)}

// 		err := review(mockedEnv, args)

// 		assert.Nil(t, err)
// 		assert.Equal(t, test.shouldReviewBeDone, isPostReviewRequestPerformed)
// 	}
// }
