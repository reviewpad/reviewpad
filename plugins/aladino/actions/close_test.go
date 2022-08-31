// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var close = plugins_aladino.PluginBuiltIns().Actions["close"].Code

func TestClose(t *testing.T) {
	failMessage := "EditRequestFail"
	var gotState, closeComment string
	var commentCreated bool

	tests := map[string]struct {
		mockedEnv    aladino.Env
		inputComment string
		wantState    string
		wantComment  bool
		wantErr      string
	}{
		"when edit request fails": {
			mockedEnv: aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.PatchReposPullsByOwnerByRepoByPullNumber,
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
				handler.PullRequest,
			),
			wantErr: failMessage,
		},
		"when entity is closed with comment": {
			mockedEnv: aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.PatchReposPullsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							rawBody, _ := ioutil.ReadAll(r.Body)
							body := github.PullRequest{}

							json.Unmarshal(rawBody, &body)

							gotState = body.GetState()
						}),
					),
					mock.WithRequestMatchHandler(
						mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							rawBody, _ := ioutil.ReadAll(r.Body)
							body := github.IssueComment{}

							json.Unmarshal(rawBody, &body)

							commentCreated = true
							closeComment = *body.Body
						}),
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
				handler.PullRequest,
			),
			inputComment: "Lorem Ipsum",
			wantComment:  true,
			wantState:    "closed",
		},
		"when entity is close without comment": {
			mockedEnv: aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.PatchReposPullsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							rawBody, _ := ioutil.ReadAll(r.Body)
							body := github.PullRequest{}

							json.Unmarshal(rawBody, &body)

							gotState = body.GetState()
						}),
					),
					mock.WithRequestMatchHandler(
						mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							// If the create comment request was performed then the comment was created
							commentCreated = true
						}),
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
				handler.PullRequest,
			),
			wantState: "closed",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			args := []aladino.Value{aladino.BuildStringValue(test.inputComment)}
			gotErr := close(test.mockedEnv, args)

			if gotErr != nil && gotErr.(*github.ErrorResponse).Message != test.wantErr {
				assert.FailNow(t, "Close() error = %v, wantErr %v", gotErr, test.wantErr)
			}

			assert.Equal(t, test.wantState, gotState)
			assert.Equal(t, test.inputComment, closeComment)
			assert.Equal(t, test.wantComment, commentCreated)

			// Since these are variables common to all tests we need to reset their values at the end of each test
			gotState = ""
			closeComment = ""
			commentCreated = false
		})
	}
}
