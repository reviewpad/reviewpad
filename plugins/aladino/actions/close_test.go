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
	"github.com/reviewpad/reviewpad/v3/handler"
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

func TestClose_WhenCommenteRequestFails(t *testing.T) {
	entityNodeID := aladino.GetDefaultMockPullRequestDetails().GetNodeID()
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

	mockedCloseIssueMutation := fmt.Sprintf(`{
        "query": "mutation($input:CloseIssueInput!) {
            closeIssue(input: $input) {
                clientMutationId
            }
        }",
        "variables":{
            "input":{
                "issueId": "%v",
				"stateReason": "COMPLETED"
            }
        }
    }`, entityNodeID)

	mockedCloseIssueBody := `{"data": {"closeIssue": {"clientMutationId": "client_mutation_id"}}}}`

	tests := map[string]struct {
		mockedEnv    aladino.Env
		inputComment string
		wantErr      string
	}{
		"when pull request is closed": {
			mockedEnv: aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							mock.WriteError(
								w,
								http.StatusInternalServerError,
								"failed to comment on pull request",
							)
						}),
					),
				},
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					if query == utils.MinifyQuery(mockedClosePullRequestMutation) {
						utils.MustWrite(w, mockedClosePullRequestBody)
					}
				},
				aladino.MockBuiltIns(),
				nil,
			),
			inputComment: "test",
			wantErr:      "failed to comment on pull request",
		},
		"when issue is closed": {
			mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.GetReposIssuesByOwnerByRepoByIssueNumber,
						&github.Issue{
							Number: github.Int(aladino.DefaultMockPrNum),
							NodeID: github.String(entityNodeID),
							Repository: &github.Repository{
								Name: github.String(aladino.DefaultMockPrRepoName),
							},
							User: &github.User{
								Login: github.String(aladino.DefaultMockPrOwner),
							},
						},
					),
					mock.WithRequestMatchHandler(
						mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							mock.WriteError(
								w,
								http.StatusInternalServerError,
								"failed to comment on issue",
							)
						}),
					),
				},
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					if query == utils.MinifyQuery(mockedCloseIssueMutation) {
						utils.MustWrite(w, mockedCloseIssueBody)
					}
				},
				aladino.MockBuiltIns(),
				nil,
				&handler.TargetEntity{
					Owner:  aladino.DefaultMockPrOwner,
					Repo:   aladino.DefaultMockPrRepoName,
					Number: aladino.DefaultMockPrNum,
					Kind:   handler.Issue,
				},
			),
			inputComment: "test",
			wantErr:      "failed to comment on issue",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			args := []aladino.Value{aladino.BuildStringValue(test.inputComment), aladino.BuildStringValue("completed")}
			gotErr := close(test.mockedEnv, args)

			assert.Equal(t, test.wantErr, gotErr.(*github.ErrorResponse).Message)
		})
	}
}

func TestClose(t *testing.T) {
	var closeComment string
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

	mockedCloseIssueAsCompletedMutation := fmt.Sprintf(`{
	    "query": "mutation($input:CloseIssueInput!) {
	        closeIssue(input: $input) {
	            clientMutationId
	        }
	    }",
	    "variables":{
	        "input":{
	            "issueId": "%v",
				"stateReason": "COMPLETED"
	        }
	    }
	}`, entityNodeID)

	mockedCloseIssueAsNotPlannedMutation := fmt.Sprintf(`{
	    "query": "mutation($input:CloseIssueInput!) {
	        closeIssue(input: $input) {
	            clientMutationId
	        }
	    }",
	    "variables":{
	        "input":{
	            "issueId": "%v",
				"stateReason": "NOT_PLANNED"
	        }
	    }
	}`, entityNodeID)

	mockedCloseIssueBody := `{"data": {"closeIssue": {"clientMutationId": "client_mutation_id"}}}}`

	tests := map[string]struct {
		mockedEnv        aladino.Env
		inputComment     string
		inputStateReason string
		wantComment      bool
		wantErr          error
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
						utils.MustWrite(w, mockedClosePullRequestBody)
					}
				},
				aladino.MockBuiltIns(),
				nil,
			),
			inputComment: "Lorem Ipsum",
			wantComment:  true,
			wantErr:      nil,
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
						utils.MustWrite(w, mockedClosePullRequestBody)
					}
				},
				aladino.MockBuiltIns(),
				nil,
			),
			wantErr: nil,
		},
		"when issue is closed with comment and with reason completed": {
			mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.GetReposIssuesByOwnerByRepoByIssueNumber,
						&github.Issue{
							Number: github.Int(aladino.DefaultMockPrNum),
							NodeID: github.String(entityNodeID),
							Repository: &github.Repository{
								Name: github.String(aladino.DefaultMockPrRepoName),
							},
							User: &github.User{
								Login: github.String(aladino.DefaultMockPrOwner),
							},
						},
					),
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
					if query == utils.MinifyQuery(mockedCloseIssueAsCompletedMutation) {
						utils.MustWrite(w, mockedCloseIssueBody)
					}
				},
				aladino.MockBuiltIns(),
				nil,
				&handler.TargetEntity{
					Owner:  aladino.DefaultMockPrOwner,
					Repo:   aladino.DefaultMockPrRepoName,
					Number: aladino.DefaultMockPrNum,
					Kind:   handler.Issue,
				},
			),
			inputComment:     "done",
			inputStateReason: "completed",
			wantComment:      true,
			wantErr:          nil,
		},
		"when issue is closed with comment and with reason not_planned": {
			mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.GetReposIssuesByOwnerByRepoByIssueNumber,
						&github.Issue{
							Number: github.Int(aladino.DefaultMockPrNum),
							NodeID: github.String(entityNodeID),
							Repository: &github.Repository{
								Name: github.String(aladino.DefaultMockPrRepoName),
							},
							User: &github.User{
								Login: github.String(aladino.DefaultMockPrOwner),
							},
						},
					),
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
					if query == utils.MinifyQuery(mockedCloseIssueAsNotPlannedMutation) {
						utils.MustWrite(w, mockedCloseIssueBody)
					}
				},
				aladino.MockBuiltIns(),
				nil,
				&handler.TargetEntity{
					Owner:  aladino.DefaultMockPrOwner,
					Repo:   aladino.DefaultMockPrRepoName,
					Number: aladino.DefaultMockPrNum,
					Kind:   handler.Issue,
				},
			),
			inputComment:     "wont do",
			inputStateReason: "not_planned",
			wantComment:      true,
			wantErr:          nil,
		},
		"when issue is closed with no comment and with reason completed": {
			mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.GetReposIssuesByOwnerByRepoByIssueNumber,
						&github.Issue{
							Number: github.Int(aladino.DefaultMockPrNum),
							NodeID: github.String(entityNodeID),
							Repository: &github.Repository{
								Name: github.String(aladino.DefaultMockPrRepoName),
							},
							User: &github.User{
								Login: github.String(aladino.DefaultMockPrOwner),
							},
						},
					),
				},
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					if query == utils.MinifyQuery(mockedCloseIssueAsCompletedMutation) {
						utils.MustWrite(w, mockedCloseIssueBody)
					}
				},
				aladino.MockBuiltIns(),
				nil,
				&handler.TargetEntity{
					Owner:  aladino.DefaultMockPrOwner,
					Repo:   aladino.DefaultMockPrRepoName,
					Number: aladino.DefaultMockPrNum,
					Kind:   handler.Issue,
				},
			),
			inputStateReason: "completed",
			wantComment:      false,
			wantErr:          nil,
		},
		"when issue is closed with no comment and with reason not_planned": {
			mockedEnv: aladino.MockDefaultEnvWithTargetEntity(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.GetReposIssuesByOwnerByRepoByIssueNumber,
						&github.Issue{
							Number: github.Int(aladino.DefaultMockPrNum),
							NodeID: github.String(entityNodeID),
							Repository: &github.Repository{
								Name: github.String(aladino.DefaultMockPrRepoName),
							},
							User: &github.User{
								Login: github.String(aladino.DefaultMockPrOwner),
							},
						},
					),
				},
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					if query == utils.MinifyQuery(mockedCloseIssueAsNotPlannedMutation) {
						utils.MustWrite(w, mockedCloseIssueBody)
					}
				},
				aladino.MockBuiltIns(),
				nil,
				&handler.TargetEntity{
					Owner:  aladino.DefaultMockPrOwner,
					Repo:   aladino.DefaultMockPrRepoName,
					Number: aladino.DefaultMockPrNum,
					Kind:   handler.Issue,
				},
			),
			inputStateReason: "not_planned",
			wantComment:      false,
			wantErr:          nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			args := []aladino.Value{aladino.BuildStringValue(test.inputComment), aladino.BuildStringValue(test.inputStateReason)}
			gotErr := close(test.mockedEnv, args)

			assert.Equal(t, test.wantErr, gotErr)
			assert.Equal(t, test.inputComment, closeComment)
			assert.Equal(t, test.wantComment, commentCreated)

			// Since these are variables common to all tests we need to reset their values at the end of each test
			closeComment = ""
			commentCreated = false
		})
	}
}
