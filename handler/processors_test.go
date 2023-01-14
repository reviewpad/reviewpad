// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/jarcoal/httpmock"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func buildPayload(payload []byte) *json.RawMessage {
	rawPayload := json.RawMessage(payload)
	return &rawPayload
}

func TestParseEvent_Failure(t *testing.T) {
	log := utils.NewLogger(logrus.DebugLevel)

	event := `{"type": "ping",}`
	gotEvent, err := handler.ParseEvent(log, event)

	assert.NotNil(t, err)
	assert.Nil(t, gotEvent)
}

func TestParseEvent(t *testing.T) {
	log := utils.NewLogger(logrus.DebugLevel)

	event := `{"action": "ping"}`
	wantEvent := &handler.ActionEvent{
		ActionName: github.String("ping"),
	}

	gotEvent, err := handler.ParseEvent(log, event)

	assert.Nil(t, err)
	assert.Equal(t, wantEvent, gotEvent)
}

func TestProcessEvent_Failure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	log := utils.NewLogger(logrus.DebugLevel)

	owner := "reviewpad"
	repo := "reviewpad"
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.github.com/repos/%v/%v/pulls", owner, repo),
		func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("error")
		},
	)

	tests := map[string]struct {
		event *handler.ActionEvent
	}{
		"pull_request": {
			event: &handler.ActionEvent{
				EventName:    github.String("pull_request"),
				EventPayload: buildPayload([]byte(`{,}`)),
			},
		},
		"unsupported_event": {
			event: &handler.ActionEvent{
				EventName: github.String("branch_protection_rule"),
				EventPayload: buildPayload([]byte(`{
					"action": "branch_protection_rule"
				}`)),
			},
		},
		"cron": {
			event: &handler.ActionEvent{
				EventName:  github.String("schedule"),
				Token:      github.String("test-token"),
				Repository: github.String("reviewpad/reviewpad"),
			},
		},
		"workflow_run_match": {
			event: &handler.ActionEvent{
				EventName: github.String("workflow_run"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"action": "completed",
					"repository": {
						"name": "reviewpad",
						"owner": {
							"login": "reviewpad"
						}
					},
                    "pull_request": {
                        "number": 1
                    },
					"workflow_run": {
						"head_sha": "4bf24cc72f3a62423927a0ac8d70febad7c78e0g"
					}
				}`)),
			},
		},
		"status": {
			event: &handler.ActionEvent{
				EventName: github.String("status"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"repository": {
						"name": "reviewpad",
						"owner": {
							"login": "reviewpad"
						}
					},
                    "pull_request": {
                        "number": 1
                    }
					"status": {
						"sha": "4bf24cc72f3a62423927a0ac8d70febad7c78e0g"
					}
				}`)),
			},
		},
		"push": {
			event: &handler.ActionEvent{
				EventName: github.String("push"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"repository": {
						"full_name": "reviewpad/reviewpad"
					},
					"ref": "refs/heads/main"
				}`)),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotTargets, gotEventDetails, gotErr := handler.ProcessEvent(log, test.event)

			assert.Nil(t, gotTargets)
			assert.Nil(t, gotEventDetails)
			assert.NotNil(t, gotErr)
		})
	}
}

func TestProcessEvent(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	log := utils.NewLogger(logrus.DebugLevel)

	owner := "reviewpad"
	repo := "reviewpad"
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.github.com/repos/%v/%v/pulls", owner, repo),
		func(req *http.Request) (*http.Response, error) {
			b, err := json.Marshal([]*github.PullRequest{
				{
					Number: github.Int(aladino.DefaultMockPrNum),
					Base: &github.PullRequestBranch{
						Ref: github.String("refs/heads/main"),
						Repo: &github.Repository{
							Name: github.String(repo),
							Owner: &github.User{
								Login: github.String(owner),
							},
						},
					},
					Head: &github.PullRequestBranch{
						SHA: github.String("4bf24cc72f3a62423927a0ac8d70febad7c78e0g"),
					},
				},
				{
					Number: github.Int(130),
					Base: &github.PullRequestBranch{
						Ref: github.String("refs/heads/feat"),
						Repo: &github.Repository{
							Name: github.String(repo),
							Owner: &github.User{
								Login: github.String(owner),
							},
						},
					},
					Head: &github.PullRequestBranch{
						SHA: github.String("4bf24cc72f3a62423927a0ac8d70febad7c78e0k"),
					},
				},
			})
			if err != nil {
				return nil, err
			}

			resp := httpmock.NewBytesResponse(200, b)

			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.github.com/repos/%v/%v/issues", owner, repo),
		func(req *http.Request) (*http.Response, error) {
			b, err := json.Marshal([]*github.Issue{
				{
					Number: github.Int(aladino.DefaultMockPrNum),
					PullRequestLinks: &github.PullRequestLinks{
						HTMLURL: github.String(fmt.Sprintf("https://api.github.com/repos/%v/%v/pull/%v", owner, repo, aladino.DefaultMockPrNum)),
					},
				},
				{
					Number: github.Int(130),
					PullRequestLinks: &github.PullRequestLinks{
						HTMLURL: github.String(fmt.Sprintf("https://api.github.com/repos/%v/%v/pull/%v", owner, repo, 130)),
					},
				},
			})
			if err != nil {
				return nil, err
			}

			resp := httpmock.NewBytesResponse(200, b)

			return resp, nil
		},
	)

	tests := map[string]struct {
		event            *handler.ActionEvent
		wantTargets      []*handler.TargetEntity
		wantEventDetails *handler.EventDetails
		wantPayload      interface{}
		wantErr          error
	}{
		"pull_request": {
			event: &handler.ActionEvent{
				EventName: github.String("pull_request"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"action": "opened",
					"number": 130,
					"repository": {
						"name": "reviewpad",
						"owner": {
							"login": "reviewpad"
						}
					},
					"pull_request": {
						"body": "## Description",
						"number": 130
					}
				}`)),
			},
			wantTargets: []*handler.TargetEntity{
				{
					Kind:   handler.PullRequest,
					Number: 130,
					Owner:  owner,
					Repo:   repo,
				},
			},
			wantEventDetails: &handler.EventDetails{
				EventName:   "pull_request",
				EventAction: "opened",
				Payload: &github.PullRequestEvent{
					Action: github.String("opened"),
					Number: github.Int(130),
					Repo: &github.Repository{
						Name: github.String("reviewpad"),
						Owner: &github.User{
							Login: github.String("reviewpad"),
						},
					},
					PullRequest: &github.PullRequest{
						Body:   github.String("## Description"),
						Number: github.Int(130),
					},
				},
			},
		},
		"pull_request_target": {
			event: &handler.ActionEvent{
				EventName: github.String("pull_request_target"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"action": "opened",
					"number": 130,
					"repository": {
						"name": "reviewpad",
						"owner": {
							"login": "reviewpad"
						}
					},
					"pull_request": {
						"body": "## Description",
						"number": 130
					}
				}`)),
			},
			wantTargets: []*handler.TargetEntity{
				{
					Kind:   handler.PullRequest,
					Number: 130,
					Owner:  owner,
					Repo:   repo,
				},
			},
			wantEventDetails: &handler.EventDetails{
				EventName:   "pull_request_target",
				EventAction: "opened",
				Payload: &github.PullRequestTargetEvent{
					Action: github.String("opened"),
					Number: github.Int(130),
					Repo: &github.Repository{
						Name: github.String("reviewpad"),
						Owner: &github.User{
							Login: github.String("reviewpad"),
						},
					},
					PullRequest: &github.PullRequest{
						Body:   github.String("## Description"),
						Number: github.Int(130),
					},
				},
			},
		},
		"pull_request_review": {
			event: &handler.ActionEvent{
				EventName: github.String("pull_request_review"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"action": "opened",
					"number": 130,
					"repository": {
						"name": "reviewpad",
						"owner": {
							"login": "reviewpad"
						}
					},
					"pull_request": {
						"body": "## Description",
						"number": 130
					}
				}`)),
			},
			wantTargets: []*handler.TargetEntity{
				{
					Kind:   handler.PullRequest,
					Number: 130,
					Owner:  owner,
					Repo:   repo,
				},
			},
			wantEventDetails: &handler.EventDetails{
				EventName:   "pull_request_review",
				EventAction: "opened",
				Payload: &github.PullRequestReviewEvent{
					Action: github.String("opened"),
					Repo: &github.Repository{
						Name: github.String("reviewpad"),
						Owner: &github.User{
							Login: github.String("reviewpad"),
						},
					},
					PullRequest: &github.PullRequest{
						Body:   github.String("## Description"),
						Number: github.Int(130),
					},
				},
			},
		},
		"pull_request_review_comment": {
			event: &handler.ActionEvent{
				EventName: github.String("pull_request_review_comment"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"action": "created",
					"repository": {
						"name": "reviewpad",
						"owner": {
							"login": "reviewpad"
						}
					},
					"pull_request": {
						"body": "## Description",
						"number": 130
					}
				}`)),
			},
			wantTargets: []*handler.TargetEntity{
				{
					Kind:   handler.PullRequest,
					Number: 130,
					Owner:  owner,
					Repo:   repo,
				},
			},
			wantEventDetails: &handler.EventDetails{
				EventName:   "pull_request_review_comment",
				EventAction: "created",
				Payload: &github.PullRequestReviewCommentEvent{
					Action: github.String("created"),
					Repo: &github.Repository{
						Name: github.String("reviewpad"),
						Owner: &github.User{
							Login: github.String("reviewpad"),
						},
					},
					PullRequest: &github.PullRequest{
						Body:   github.String("## Description"),
						Number: github.Int(130),
					},
				},
			},
		},
		"cron": {
			event: &handler.ActionEvent{
				EventName:  github.String("schedule"),
				Token:      github.String("test-token"),
				Repository: github.String("reviewpad/reviewpad"),
			},
			wantTargets: []*handler.TargetEntity{
				{
					Kind:   handler.PullRequest,
					Number: 130,
					Owner:  owner,
					Repo:   repo,
				},
				{
					Kind:   handler.PullRequest,
					Number: aladino.DefaultMockPrNum,
					Owner:  owner,
					Repo:   repo,
				},
			},
			wantEventDetails: &handler.EventDetails{
				EventName: "schedule",
			},
		},
		"workflow_run_match": {
			event: &handler.ActionEvent{
				EventName: github.String("workflow_run"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"action": "completed",
					"repository": {
						"name": "reviewpad",
						"owner": {
							"login": "reviewpad"
						}
					},
					"workflow_run": {
						"head_sha": "4bf24cc72f3a62423927a0ac8d70febad7c78e0g"
					}
				}`)),
			},
			wantTargets: []*handler.TargetEntity{
				{
					Kind:   handler.PullRequest,
					Number: aladino.DefaultMockPrNum,
					Owner:  owner,
					Repo:   repo,
				},
			},
			wantEventDetails: &handler.EventDetails{
				EventName:   "workflow_run",
				EventAction: "completed",
				Payload: &github.WorkflowRunEvent{
					Action: github.String("completed"),
					Repo: &github.Repository{
						Name: github.String("reviewpad"),
						Owner: &github.User{
							Login: github.String("reviewpad"),
						},
					},
					WorkflowRun: &github.WorkflowRun{
						HeadSHA: github.String("4bf24cc72f3a62423927a0ac8d70febad7c78e0g"),
					},
				},
			},
		},
		"workflow_run_no_match": {
			event: &handler.ActionEvent{
				EventName: github.String("workflow_run"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"action": "completed",
					"repository": {
						"name": "reviewpad",
						"owner": {
							"login": "reviewpad"
						}
					},
					"workflow_run": {
						"head_sha": "4bf24cc72f3a62423927a0ac8d70febad7c78e0a"
					}
				}`)),
			},
			wantTargets: []*handler.TargetEntity{},
			wantEventDetails: &handler.EventDetails{
				EventName:   "workflow_run",
				EventAction: "completed",
				Payload: &github.WorkflowRunEvent{
					Action: github.String("completed"),
					Repo: &github.Repository{
						Name: github.String("reviewpad"),
						Owner: &github.User{
							Login: github.String("reviewpad"),
						},
					},
					WorkflowRun: &github.WorkflowRun{
						HeadSHA: github.String("4bf24cc72f3a62423927a0ac8d70febad7c78e0a"),
					},
				},
			},
		},
		"issues": {
			event: &handler.ActionEvent{
				EventName: github.String("issues"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"action": "opened",
					"number": 130,
					"repository": {
						"name": "reviewpad",
						"owner": {
							"login": "reviewpad"
						}
					},
					"issue": {
						"body": "## Description",
						"number": 130
					}
				}`)),
			},
			wantTargets: []*handler.TargetEntity{
				{
					Kind:   handler.Issue,
					Number: 130,
					Owner:  owner,
					Repo:   owner,
				},
			},
			wantEventDetails: &handler.EventDetails{
				EventName:   "issues",
				EventAction: "opened",
				Payload: &github.IssuesEvent{
					Action: github.String("opened"),
					Repo: &github.Repository{
						Name: github.String("reviewpad"),
						Owner: &github.User{
							Login: github.String("reviewpad"),
						},
					},
					Issue: &github.Issue{
						Body:   github.String("## Description"),
						Number: github.Int(130),
					},
				},
			},
		},
		"issue_comment": {
			event: &handler.ActionEvent{
				EventName: github.String("issue_comment"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"action": "opened",
					"number": 130,
					"repository": {
						"name": "reviewpad",
						"owner": {
							"login": "reviewpad"
						}
					},
					"issue": {
						"body": "## Description",
						"number": 130
					},
					"comment": {
						"body": "comment"
					}
				}`)),
			},
			wantTargets: []*handler.TargetEntity{
				{
					Kind:   handler.Issue,
					Number: 130,
					Owner:  owner,
					Repo:   repo,
				},
			},
			wantEventDetails: &handler.EventDetails{
				EventName:   "issue_comment",
				EventAction: "opened",
				Payload: &github.IssueCommentEvent{
					Action: github.String("opened"),
					Repo: &github.Repository{
						Name: github.String("reviewpad"),
						Owner: &github.User{
							Login: github.String("reviewpad"),
						},
					},
					Issue: &github.Issue{
						Body:   github.String("## Description"),
						Number: github.Int(130),
					},
					Comment: &github.IssueComment{
						Body: github.String("comment"),
					},
				},
			},
		},
		"status_match": {
			event: &handler.ActionEvent{
				EventName: github.String("status"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"repository": {
						"name": "reviewpad",
						"owner": {
							"login": "reviewpad"
						}
					},
                    "pull_request": {
                        "number": 6
                    },
					"sha": "4bf24cc72f3a62423927a0ac8d70febad7c78e0g"
				}`)),
			},
			wantTargets: []*handler.TargetEntity{
				{
					Kind:   handler.PullRequest,
					Number: aladino.DefaultMockPrNum,
					Owner:  owner,
					Repo:   repo,
				},
			},
			wantEventDetails: &handler.EventDetails{
				EventName: "status",
				Payload: &github.StatusEvent{
					Repo: &github.Repository{
						Name: github.String("reviewpad"),
						Owner: &github.User{
							Login: github.String("reviewpad"),
						},
					},
					SHA: github.String("4bf24cc72f3a62423927a0ac8d70febad7c78e0g"),
				},
			},
		},
		"status_no_match": {
			event: &handler.ActionEvent{
				EventName: github.String("status"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"repository": {
						"name": "reviewpad",
						"owner": {
							"login": "reviewpad"
						}
					},
					"sha": "4bf24cc72f3a62423927a0ac8d70febad7c78e0a"
				}`)),
			},
			wantTargets: []*handler.TargetEntity{},
			wantEventDetails: &handler.EventDetails{
				EventName: "status",
				Payload: &github.StatusEvent{
					Repo: &github.Repository{
						Name: github.String("reviewpad"),
						Owner: &github.User{
							Login: github.String("reviewpad"),
						},
					},
					SHA: github.String("4bf24cc72f3a62423927a0ac8d70febad7c78e0a"),
				},
			},
		},
		"push": {
			event: &handler.ActionEvent{
				EventName: github.String("push"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"repository": {
						"full_name": "reviewpad/reviewpad"
					},
					"ref": "refs/heads/main"
				}`)),
			},
			wantTargets: []*handler.TargetEntity{
				{
					Kind:   handler.PullRequest,
					Number: aladino.DefaultMockPrNum,
					Owner:  owner,
					Repo:   repo,
				},
			},
			wantEventDetails: &handler.EventDetails{
				EventName: "push",
				Payload: &github.PushEvent{
					Repo: &github.PushEventRepository{
						FullName: github.String("reviewpad/reviewpad"),
					},
					Ref: github.String("refs/heads/main"),
				},
			},
		},
		"push_no_match": {
			event: &handler.ActionEvent{
				EventName: github.String("push"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"repository": {
						"full_name": "reviewpad/reviewpad"
					},
					"ref": "refs/heads/master"
				}`)),
			},
			wantTargets: []*handler.TargetEntity{},
			wantEventDetails: &handler.EventDetails{
				EventName: "push",
				Payload: &github.PushEvent{
					Repo: &github.PushEventRepository{
						FullName: github.String("reviewpad/reviewpad"),
					},
					Ref: github.String("refs/heads/master"),
				},
			},
		},
		"installation": {
			event: &handler.ActionEvent{
				EventName: github.String("installation"),
				EventPayload: buildPayload([]byte(`{
					"action": "created",
					"repositories": [
						{
							"full_name": "testowner/testrepo"
						},
						{
							"full_name": "testowner2/testrepo2"
						}
					]
				}`)),
			},
			wantTargets: []*handler.TargetEntity{
				{
					Repo:  "testrepo",
					Owner: "testowner",
				},
				{
					Repo:  "testrepo2",
					Owner: "testowner2",
				},
			},
			wantEventDetails: &handler.EventDetails{
				EventName:   "installation",
				EventAction: "created",
				Payload: &github.InstallationEvent{
					Action: github.String("created"),
					Repositories: []*github.Repository{
						{
							FullName: github.String("testowner/testrepo"),
						},
						{
							FullName: github.String("testowner2/testrepo2"),
						},
					},
				},
			},
		},
		"installation with invalid repo full name": {
			event: &handler.ActionEvent{
				EventName: github.String("installation"),
				EventPayload: buildPayload([]byte(`{
					"action": "created",
					"repositories": [
						{
							"full_name": "testowner/testrepo"
						},
						{
							"full_name": "testowner2"
						}
					]
				}`)),
			},
			wantTargets:      nil,
			wantEventDetails: nil,
			wantErr:          errors.New("invalid full repository name: testowner2"),
		},
		"installation_repositories added": {
			event: &handler.ActionEvent{
				EventName: github.String("installation_repositories"),
				EventPayload: buildPayload([]byte(`{
					"action": "added",
					"repositories_added": [
						{
							"full_name": "testowner/testrepo"
						},
						{
							"full_name": "testowner2/testrepo2"
						}
					]
				}`)),
			},
			wantTargets: []*handler.TargetEntity{
				{
					Repo:  "testrepo",
					Owner: "testowner",
				},
				{
					Repo:  "testrepo2",
					Owner: "testowner2",
				},
			},
			wantEventDetails: &handler.EventDetails{
				EventName:   "installation_repositories",
				EventAction: "added",
				Payload: &github.InstallationRepositoriesEvent{
					Action: github.String("added"),
					RepositoriesAdded: []*github.Repository{
						{
							FullName: github.String("testowner/testrepo"),
						},
						{
							FullName: github.String("testowner2/testrepo2"),
						},
					},
				},
			},
		},
		"installation_repositories removed": {
			event: &handler.ActionEvent{
				EventName: github.String("installation_repositories"),
				EventPayload: buildPayload([]byte(`{
					"action": "removed",
					"repositories_removed": [
						{
							"full_name": "testowner/testrepo"
						},
						{
							"full_name": "testowner2/testrepo2"
						}
					]
				}`)),
			},
			wantTargets: []*handler.TargetEntity{
				{
					Repo:  "testrepo",
					Owner: "testowner",
				},
				{
					Repo:  "testrepo2",
					Owner: "testowner2",
				},
			},
			wantEventDetails: &handler.EventDetails{
				EventName:   "installation_repositories",
				EventAction: "removed",
				Payload: &github.InstallationRepositoriesEvent{
					Action: github.String("removed"),
					RepositoriesRemoved: []*github.Repository{
						{
							FullName: github.String("testowner/testrepo"),
						},
						{
							FullName: github.String("testowner2/testrepo2"),
						},
					},
				},
			},
		},
		"check_run when pr is from same repo": {
			event: &handler.ActionEvent{
				EventName: github.String("check_run"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"action": "rerequested",
					"check_run": {
						"id": 1,
						"head_sha": "4bf24cc72f3a62423927a0ac8d70febad7c78e0g",
						"pull_requests": [
							{
								"number": 1
							}
						]
					},
					"repository": {
						"name": "reviewpad",
						"owner": {
							"login": "reviewpad"
						}
					}
				}`)),
			},
			wantTargets: []*handler.TargetEntity{
				{
					Kind:   handler.PullRequest,
					Owner:  "reviewpad",
					Repo:   "reviewpad",
					Number: 1,
				},
			},
			wantEventDetails: &handler.EventDetails{
				EventName:   "check_run",
				EventAction: "rerequested",
				Payload: &github.CheckRunEvent{
					Action: github.String("rerequested"),
					CheckRun: &github.CheckRun{
						ID:      github.Int64(1),
						HeadSHA: github.String("4bf24cc72f3a62423927a0ac8d70febad7c78e0g"),
						PullRequests: []*github.PullRequest{
							{
								Number: github.Int(1),
							},
						},
					},
					Repo: &github.Repository{
						Name: github.String("reviewpad"),
						Owner: &github.User{
							Login: github.String("reviewpad"),
						},
					},
				},
			},
		},
		"check_run when pr is from a forked repo": {
			event: &handler.ActionEvent{
				EventName: github.String("check_run"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"action": "created",
					"check_run": {
						"id": 1,
						"head_sha": "4bf24cc72f3a62423927a0ac8d70febad7c78e0g",
						"pull_requests": []
					},
					"repository": {
						"name": "reviewpad",
						"full_name": "reviewpad/reviewpad",
						"owner": {
							"login": "reviewpad"
						}
					}
				}`)),
			},
			wantTargets: []*handler.TargetEntity{
				{
					Kind:   handler.PullRequest,
					Owner:  "reviewpad",
					Repo:   "reviewpad",
					Number: 6,
				},
			},
			wantEventDetails: &handler.EventDetails{
				EventName:   "check_run",
				EventAction: "created",
				Payload: &github.CheckRunEvent{
					Action: github.String("created"),
					CheckRun: &github.CheckRun{
						ID:           github.Int64(1),
						HeadSHA:      github.String("4bf24cc72f3a62423927a0ac8d70febad7c78e0g"),
						PullRequests: []*github.PullRequest{},
					},
					Repo: &github.Repository{
						Name:     github.String("reviewpad"),
						FullName: github.String("reviewpad/reviewpad"),
						Owner: &github.User{
							Login: github.String("reviewpad"),
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotTargets, gotEventDetails, err := handler.ProcessEvent(log, test.event)

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantEventDetails, gotEventDetails)
			assert.ElementsMatch(t, test.wantTargets, gotTargets)
		})
	}
}
