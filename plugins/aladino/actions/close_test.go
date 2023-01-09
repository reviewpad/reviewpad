// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var close = plugins_aladino.PluginBuiltIns().Actions["close"].Code

func TestClose_WhenCloseRequestFails(t *testing.T) {
	failMessage := "ClosePullRequestRequestFail"
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetails()
	mockedClosePullRequestMutation := fmt.Sprintf(`{
        "query": "mutation($input:ClosePullRequestInput!) {
            closePullRequest(input: $input) {
                clientMutationId
            }
        }",
        "variables":{
            "input":{
                "pullRequestId": "%v"
            }
        }
    }`, mockedPullRequest.GetNodeID())

	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			if query == utils.MinifyQuery(mockedClosePullRequestMutation) {
				http.Error(w, failMessage, http.StatusNotFound)
			}
		},
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue(""), aladino.BuildStringValue("")}
	gotErr := close(mockedEnv, args)

	assert.EqualError(t, gotErr, fmt.Sprintf("non-200 OK status code: 404 Not Found body: \"%s\\n\"", failMessage))
}

func TestClose(t *testing.T) {
	var gotState, closeComment, gotStateReason string
	var commentCreated bool

	entityNodeID := aladino.DefaultMockEntityNodeID

	mockedClosePullRequestMutation := fmt.Sprintf(`{
        "query": "mutation($input:ClosePullRequestInput!) {
            closePullRequest(input: $input) {
                clientMutationId
            }
        }",
        "variables":{
            "input":{
                "pullRequestId": "%v"
            }
        }
    }`, entityNodeID)

	mockedClosePullRequestBody := `{"data": {"closePullRequest": {"clientMutationId": "client_mutation_id"}}}}`

	// mockedCloseIssuetMutation := fmt.Sprintf(`{
	//     "query": "mutation($input:CloseIssueInput!) {
	//         closeIssue(input: $input) {
	//             clientMutationId
	//         }
	//     }",
	//     "variables":{
	//         "input":{
	//             "issueId": "%v"
	//         }
	//     }
	// }`, entityNodeID)

	tests := map[string]struct {
		mockedEnv        aladino.Env
		inputComment     string
		inputStateReason string
		wantState        string
		wantStateReason  string
		wantComment      bool
		wantErr          string
	}{
		"when pull request is closed with comment": {
			mockedEnv: aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							rawBody, _ := io.ReadAll(r.Body)
							body := github.IssueComment{}

							utils.MustUnmarshal(rawBody, &body)

							commentCreated = true
							closeComment = *body.Body
						}),
					),
				},
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					if query == utils.MinifyQuery(mockedClosePullRequestMutation) {
						gotState = "closed"
						utils.MustWrite(w, mockedClosePullRequestBody)
					}
				},
				aladino.MockBuiltIns(),
				nil,
			),
			inputComment: "Lorem Ipsum",
			wantComment:  true,
			wantState:    "closed",
		},
		"when pull request is close without comment": {
			mockedEnv: aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							// If the create comment request was performed then the comment was created
							commentCreated = true
						}),
					),
				},
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					if query == utils.MinifyQuery(mockedClosePullRequestMutation) {
						gotState = "closed"
						utils.MustWrite(w, mockedClosePullRequestBody)
					}
				},
				aladino.MockBuiltIns(),
				nil,
			),
			wantState: "closed",
		},
		// "when issue is closed with comment and with reason completed": {
		// 	mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
		// 		t,
		// 		[]mock.MockBackendOption{
		// 			mock.WithRequestMatchHandler(
		// 				mock.PatchReposIssuesByOwnerByRepoByIssueNumber,
		// 				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 					rawBody, _ := io.ReadAll(r.Body)
		// 					body := struct {
		// 						State       string `json:"state"`
		// 						StateReason string `json:"state_reason"`
		// 					}{}

		// 					utils.MustUnmarshal(rawBody, &body)

		// 					gotState = body.State
		// 					gotStateReason = body.StateReason
		// 				}),
		// 			),
		// 			mock.WithRequestMatchHandler(
		// 				mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
		// 				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 					rawBody, _ := io.ReadAll(r.Body)
		// 					body := github.IssueComment{}

		// 					utils.MustUnmarshal(rawBody, &body)

		// 					commentCreated = true
		// 					closeComment = *body.Body
		// 				}),
		// 			),
		// 		},
		// 		nil,
		// 		aladino.MockBuiltIns(),
		// 		nil,
		// 		&handler.TargetEntity{
		// 			Owner:  aladino.DefaultMockPrOwner,
		// 			Repo:   aladino.DefaultMockPrRepoName,
		// 			Number: aladino.DefaultMockPrNum,
		// 			Kind:   handler.Issue,
		// 		},
		// 	),
		// 	inputComment:     "done",
		// 	inputStateReason: "completed",
		// 	wantState:        "closed",
		// 	wantStateReason:  "completed",
		// 	wantComment:      true,
		// },
		// "when issue is closed with comment and with reason not_planned": {
		// 	mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
		// 		t,
		// 		[]mock.MockBackendOption{
		// 			mock.WithRequestMatchHandler(
		// 				mock.PatchReposIssuesByOwnerByRepoByIssueNumber,
		// 				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 					rawBody, _ := io.ReadAll(r.Body)
		// 					body := struct {
		// 						State       string `json:"state"`
		// 						StateReason string `json:"state_reason"`
		// 					}{}

		// 					utils.MustUnmarshal(rawBody, &body)

		// 					gotState = body.State
		// 					gotStateReason = body.StateReason
		// 				}),
		// 			),
		// 			mock.WithRequestMatchHandler(
		// 				mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
		// 				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 					rawBody, _ := io.ReadAll(r.Body)
		// 					body := github.IssueComment{}

		// 					utils.MustUnmarshal(rawBody, &body)

		// 					commentCreated = true
		// 					closeComment = *body.Body
		// 				}),
		// 			),
		// 		},
		// 		nil,
		// 		aladino.MockBuiltIns(),
		// 		nil,
		// 		&handler.TargetEntity{
		// 			Owner:  aladino.DefaultMockPrOwner,
		// 			Repo:   aladino.DefaultMockPrRepoName,
		// 			Number: aladino.DefaultMockPrNum,
		// 			Kind:   handler.Issue,
		// 		},
		// 	),
		// 	inputComment:     "wont do",
		// 	inputStateReason: "not_planned",
		// 	wantState:        "closed",
		// 	wantStateReason:  "not_planned",
		// 	wantComment:      true,
		// },
		// "when issue is closed with no comment and with reason completed": {
		// 	mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
		// 		t,
		// 		[]mock.MockBackendOption{
		// 			mock.WithRequestMatchHandler(
		// 				mock.PatchReposIssuesByOwnerByRepoByIssueNumber,
		// 				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 					rawBody, _ := io.ReadAll(r.Body)
		// 					body := struct {
		// 						State       string `json:"state"`
		// 						StateReason string `json:"state_reason"`
		// 					}{}

		// 					utils.MustUnmarshal(rawBody, &body)

		// 					gotState = body.State
		// 					gotStateReason = body.StateReason
		// 				}),
		// 			),
		// 		},
		// 		nil,
		// 		aladino.MockBuiltIns(),
		// 		nil,
		// 		&handler.TargetEntity{
		// 			Owner:  aladino.DefaultMockPrOwner,
		// 			Repo:   aladino.DefaultMockPrRepoName,
		// 			Number: aladino.DefaultMockPrNum,
		// 			Kind:   handler.Issue,
		// 		},
		// 	),
		// 	inputStateReason: "completed",
		// 	wantState:        "closed",
		// 	wantStateReason:  "completed",
		// 	wantComment:      false,
		// },
		// "when issue is closed with no comment and with reason not_planned": {
		// 	mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
		// 		t,
		// 		[]mock.MockBackendOption{
		// 			mock.WithRequestMatchHandler(
		// 				mock.PatchReposIssuesByOwnerByRepoByIssueNumber,
		// 				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 					rawBody, _ := io.ReadAll(r.Body)
		// 					body := struct {
		// 						State       string `json:"state"`
		// 						StateReason string `json:"state_reason"`
		// 					}{}

		// 					utils.MustUnmarshal(rawBody, &body)

		// 					gotState = body.State
		// 					gotStateReason = body.StateReason
		// 				}),
		// 			),
		// 		},
		// 		nil,
		// 		aladino.MockBuiltIns(),
		// 		nil,
		// 		&handler.TargetEntity{
		// 			Owner:  aladino.DefaultMockPrOwner,
		// 			Repo:   aladino.DefaultMockPrRepoName,
		// 			Number: aladino.DefaultMockPrNum,
		// 			Kind:   handler.Issue,
		// 		},
		// 	),
		// 	inputStateReason: "not_planned",
		// 	wantState:        "closed",
		// 	wantStateReason:  "not_planned",
		// 	wantComment:      false,
		// },
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			args := []aladino.Value{aladino.BuildStringValue(test.inputComment), aladino.BuildStringValue(test.inputStateReason)}
			gotErr := close(test.mockedEnv, args)

			if gotErr != nil && gotErr.(*github.ErrorResponse).Message != test.wantErr {
				assert.FailNow(t, "Close() error = %v, wantErr %v", gotErr, test.wantErr)
			}

			assert.Equal(t, test.wantState, gotState)
			assert.Equal(t, test.inputComment, closeComment)
			assert.Equal(t, test.wantComment, commentCreated)
			assert.Equal(t, test.wantStateReason, gotStateReason)

			// Since these are variables common to all tests we need to reset their values at the end of each test
			gotState = ""
			closeComment = ""
			commentCreated = false
			gotStateReason = ""
		})
	}
}
