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

	"github.com/google/go-github/v52/github"
	"github.com/jarcoal/httpmock"
	"github.com/reviewpad/go-lib/entities"
	log "github.com/reviewpad/go-lib/logrus"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func buildPayload(payload []byte) *json.RawMessage {
	rawPayload := json.RawMessage(payload)
	return &rawPayload
}

func TestParseEvent_Failure(t *testing.T) {
	log := log.NewLogger(logrus.DebugLevel)

	event := `{"type": "ping",}`
	gotEvent, err := handler.ParseEvent(log, event)

	assert.NotNil(t, err)
	assert.Nil(t, gotEvent)
}

func TestParseEvent(t *testing.T) {
	log := log.NewLogger(logrus.DebugLevel)

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
	log := log.NewLogger(logrus.DebugLevel)

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
	log := log.NewLogger(logrus.DebugLevel)

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
								Type:  github.String("Organization"),
							},
							Visibility: github.String("public"),
						},
					},
					Head: &github.PullRequestBranch{
						Ref: github.String("feat"),
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
								Type:  github.String("Organization"),
							},
							Visibility: github.String("public"),
						},
					},
					Head: &github.PullRequestBranch{
						Ref: github.String("bug"),
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
					Repository: &github.Repository{
						Owner: &github.User{
							Type: github.String("User"),
						},
						Visibility: github.String("public"),
					},
				},
				{
					Number: github.Int(130),
					PullRequestLinks: &github.PullRequestLinks{
						HTMLURL: github.String(fmt.Sprintf("https://api.github.com/repos/%v/%v/pull/%v", owner, repo, 130)),
					},
					Repository: &github.Repository{
						Owner: &github.User{
							Type: github.String("User"),
						},
						Visibility: github.String("public"),
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

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.github.com/repos/%v/%v", owner, repo),
		func(req *http.Request) (*http.Response, error) {
			b, err := json.Marshal(github.Repository{
				Owner: &github.User{
					Type: github.String("Organization"),
				},
				Visibility: github.String("public"),
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
		wantTargets      []*entities.TargetEntity
		wantEventDetails *entities.EventDetails
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
							"login": "reviewpad",
							"type": "Organization"
						},
						"visibility": "public"
					},
					"pull_request": {
						"body": "## Description",
						"number": 130
					}
				}`)),
			},
			wantTargets: []*entities.TargetEntity{
				{
					Kind:        entities.PullRequest,
					Number:      130,
					Owner:       owner,
					Repo:        repo,
					AccountType: "Organization",
					Visibility:  "public",
				},
			},
			wantEventDetails: &entities.EventDetails{
				EventName:   "pull_request",
				EventAction: "opened",
				Payload: &github.PullRequestEvent{
					Action: github.String("opened"),
					Number: github.Int(130),
					Repo: &github.Repository{
						Name: github.String("reviewpad"),
						Owner: &github.User{
							Login: github.String("reviewpad"),
							Type:  github.String("Organization"),
						},
						Visibility: github.String("public"),
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
							"login": "reviewpad",
							"type": "Organization"
						},
						"visibility": "public"
					},
					"pull_request": {
						"body": "## Description",
						"number": 130
					}
				}`)),
			},
			wantTargets: []*entities.TargetEntity{
				{
					Kind:        entities.PullRequest,
					Number:      130,
					Owner:       owner,
					Repo:        repo,
					AccountType: "Organization",
					Visibility:  "public",
				},
			},
			wantEventDetails: &entities.EventDetails{
				EventName:   "pull_request_target",
				EventAction: "opened",
				Payload: &github.PullRequestTargetEvent{
					Action: github.String("opened"),
					Number: github.Int(130),
					Repo: &github.Repository{
						Name: github.String("reviewpad"),
						Owner: &github.User{
							Login: github.String("reviewpad"),
							Type:  github.String("Organization"),
						},
						Visibility: github.String("public"),
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
							"login": "reviewpad",
							"type": "Organization"
						},
						"visibility": "public"
					},
					"pull_request": {
						"body": "## Description",
						"number": 130
					}
				}`)),
			},
			wantTargets: []*entities.TargetEntity{
				{
					Kind:        entities.PullRequest,
					Number:      130,
					Owner:       owner,
					Repo:        repo,
					AccountType: "Organization",
					Visibility:  "public",
				},
			},
			wantEventDetails: &entities.EventDetails{
				EventName:   "pull_request_review",
				EventAction: "opened",
				Payload: &github.PullRequestReviewEvent{
					Action: github.String("opened"),
					Repo: &github.Repository{
						Name: github.String("reviewpad"),
						Owner: &github.User{
							Login: github.String("reviewpad"),
							Type:  github.String("Organization"),
						},
						Visibility: github.String("public"),
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
							"login": "reviewpad",
							"type": "Organization"
						},
						"visibility": "public"
					},
					"pull_request": {
						"body": "## Description",
						"number": 130
					}
				}`)),
			},
			wantTargets: []*entities.TargetEntity{
				{
					Kind:        entities.PullRequest,
					Number:      130,
					Owner:       owner,
					Repo:        repo,
					AccountType: "Organization",
					Visibility:  "public",
				},
			},
			wantEventDetails: &entities.EventDetails{
				EventName:   "pull_request_review_comment",
				EventAction: "created",
				Payload: &github.PullRequestReviewCommentEvent{
					Action: github.String("created"),
					Repo: &github.Repository{
						Name: github.String("reviewpad"),
						Owner: &github.User{
							Login: github.String("reviewpad"),
							Type:  github.String("Organization"),
						},
						Visibility: github.String("public"),
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
			wantTargets: []*entities.TargetEntity{
				{
					Kind:        entities.PullRequest,
					Number:      130,
					Owner:       owner,
					Repo:        repo,
					AccountType: "Organization",
					Visibility:  "public",
				},
				{
					Kind:        entities.PullRequest,
					Number:      aladino.DefaultMockPrNum,
					Owner:       owner,
					Repo:        repo,
					AccountType: "Organization",
					Visibility:  "public",
				},
			},
			wantEventDetails: &entities.EventDetails{
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
			wantTargets: []*entities.TargetEntity{
				{
					Kind:        entities.PullRequest,
					Number:      aladino.DefaultMockPrNum,
					Owner:       owner,
					Repo:        repo,
					AccountType: "Organization",
					Visibility:  "public",
				},
			},
			wantEventDetails: &entities.EventDetails{
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
			wantTargets: []*entities.TargetEntity{},
			wantEventDetails: &entities.EventDetails{
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
			wantTargets: []*entities.TargetEntity{
				{
					Kind:   entities.Issue,
					Number: 130,
					Owner:  owner,
					Repo:   owner,
				},
			},
			wantEventDetails: &entities.EventDetails{
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
			wantTargets: []*entities.TargetEntity{
				{
					Kind:   entities.Issue,
					Number: 130,
					Owner:  owner,
					Repo:   repo,
				},
			},
			wantEventDetails: &entities.EventDetails{
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
			wantTargets: []*entities.TargetEntity{
				{
					Kind:        entities.PullRequest,
					Number:      aladino.DefaultMockPrNum,
					Owner:       owner,
					Repo:        repo,
					AccountType: "Organization",
					Visibility:  "public",
				},
			},
			wantEventDetails: &entities.EventDetails{
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
			wantTargets: []*entities.TargetEntity{},
			wantEventDetails: &entities.EventDetails{
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
						"full_name": "reviewpad/reviewpad",
						"owner": {
							"login": "reviewpad",
							"type": "Organization"
						},
						"name": "reviewpad"
					},
					"ref": "refs/heads/feat"
				}`)),
			},
			wantTargets: []*entities.TargetEntity{
				{
					Kind:        entities.PullRequest,
					Number:      aladino.DefaultMockPrNum,
					Owner:       owner,
					Repo:        repo,
					AccountType: "Organization",
					Visibility:  "public",
				},
				{
					Kind:        entities.Repository,
					Owner:       owner,
					Repo:        repo,
					AccountType: "Organization",
					Visibility:  "public",
				},
			},
			wantEventDetails: &entities.EventDetails{
				EventName: "push",
				Payload: &github.PushEvent{
					Repo: &github.PushEventRepository{
						FullName: github.String("reviewpad/reviewpad"),
						Name:     github.String("reviewpad"),
						Owner: &github.User{
							Login: github.String("reviewpad"),
							Type:  github.String("Organization"),
						},
					},
					Ref: github.String("refs/heads/feat"),
				},
			},
		},
		"push_no_match": {
			event: &handler.ActionEvent{
				EventName: github.String("push"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"repository": {
						"full_name": "reviewpad/reviewpad",
						"private": true,
						"owner": {
							"login": "reviewpad",
							"type": "Organization"
						},
						"name": "reviewpad"
					},
					"ref": "refs/heads/docs"
				}`)),
			},
			wantTargets: []*entities.TargetEntity{
				{
					Kind:        entities.Repository,
					Owner:       owner,
					Repo:        repo,
					AccountType: "Organization",
					Visibility:  "private",
				},
			},
			wantEventDetails: &entities.EventDetails{
				EventName: "push",
				Payload: &github.PushEvent{
					Repo: &github.PushEventRepository{
						FullName: github.String("reviewpad/reviewpad"),
						Name:     github.String("reviewpad"),
						Private:  github.Bool(true),
						Owner: &github.User{
							Login: github.String("reviewpad"),
							Type:  github.String("Organization"),
						},
					},
					Ref: github.String("refs/heads/docs"),
				},
			},
		},
		"installation": {
			event: &handler.ActionEvent{
				EventName: github.String("installation"),
				EventPayload: buildPayload([]byte(`{
					"action": "created",
					"installation": {
						"account": {
							"type": "User"
						}
					},
					"repositories": [
						{
							"full_name": "testowner/testrepo",
							"private": true
						},
						{
							"full_name": "testowner2/testrepo2",
							"private": false
						}
					]
				}`)),
			},
			wantTargets: []*entities.TargetEntity{
				{
					Repo:        "testrepo",
					Owner:       "testowner",
					AccountType: "User",
					Visibility:  "private",
				},
				{
					Repo:        "testrepo2",
					Owner:       "testowner2",
					AccountType: "User",
					Visibility:  "public",
				},
			},
			wantEventDetails: &entities.EventDetails{
				EventName:   "installation",
				EventAction: "created",
				Payload: &github.InstallationEvent{
					Action: github.String("created"),
					Installation: &github.Installation{
						Account: &github.User{
							Type: github.String("User"),
						},
					},
					Repositories: []*github.Repository{
						{
							FullName: github.String("testowner/testrepo"),
							Private:  github.Bool(true),
						},
						{
							FullName: github.String("testowner2/testrepo2"),
							Private:  github.Bool(false),
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
					"installation": {
						"account": {
							"type": "Organization"
						}
					},
					"repositories_added": [
						{
							"full_name": "testowner/testrepo",
							"private": true
						},
						{
							"full_name": "testowner2/testrepo2",
							"private": false
						}
					]
				}`)),
			},
			wantTargets: []*entities.TargetEntity{
				{
					Repo:        "testrepo",
					Owner:       "testowner",
					AccountType: "Organization",
					Visibility:  "private",
				},
				{
					Repo:        "testrepo2",
					Owner:       "testowner2",
					AccountType: "Organization",
					Visibility:  "public",
				},
			},
			wantEventDetails: &entities.EventDetails{
				EventName:   "installation_repositories",
				EventAction: "added",
				Payload: &github.InstallationRepositoriesEvent{
					Action: github.String("added"),
					Installation: &github.Installation{
						Account: &github.User{
							Type: github.String("Organization"),
						},
					},
					RepositoriesAdded: []*github.Repository{
						{
							FullName: github.String("testowner/testrepo"),
							Private:  github.Bool(true),
						},
						{
							FullName: github.String("testowner2/testrepo2"),
							Private:  github.Bool(false),
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
					"installation": {
						"account": {
							"type": "Organization"
						}
					},
					"repositories_removed": [
						{
							"full_name": "testowner/testrepo",
							"private": true
						},
						{
							"full_name": "testowner2/testrepo2",
							"private": true
						}
					]
				}`)),
			},
			wantTargets: []*entities.TargetEntity{
				{
					Repo:        "testrepo",
					Owner:       "testowner",
					AccountType: "Organization",
					Visibility:  "private",
				},
				{
					Repo:        "testrepo2",
					Owner:       "testowner2",
					AccountType: "Organization",
					Visibility:  "private",
				},
			},
			wantEventDetails: &entities.EventDetails{
				EventName:   "installation_repositories",
				EventAction: "removed",
				Payload: &github.InstallationRepositoriesEvent{
					Action: github.String("removed"),
					Installation: &github.Installation{
						Account: &github.User{
							Type: github.String("Organization"),
						},
					},
					RepositoriesRemoved: []*github.Repository{
						{
							FullName: github.String("testowner/testrepo"),
							Private:  github.Bool(true),
						},
						{
							FullName: github.String("testowner2/testrepo2"),
							Private:  github.Bool(true),
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
							"login": "reviewpad",
							"type": "User"
						},
						"visibility": "public"
					}
				}`)),
			},
			wantTargets: []*entities.TargetEntity{
				{
					Kind:        entities.PullRequest,
					Owner:       "reviewpad",
					Repo:        "reviewpad",
					Number:      1,
					AccountType: "User",
					Visibility:  "public",
				},
			},
			wantEventDetails: &entities.EventDetails{
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
							Type:  github.String("User"),
						},
						Visibility: github.String("public"),
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
			wantTargets: []*entities.TargetEntity{
				{
					Kind:        entities.PullRequest,
					Owner:       "reviewpad",
					Repo:        "reviewpad",
					Number:      6,
					AccountType: "Organization",
					Visibility:  "public",
				},
			},
			wantEventDetails: &entities.EventDetails{
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
		"check_suite when pr is from same repo": {
			event: &handler.ActionEvent{
				EventName: github.String("check_suite"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"action": "rerequested",
					"check_suite": {
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
			wantTargets: []*entities.TargetEntity{
				{
					Kind:   entities.PullRequest,
					Owner:  "reviewpad",
					Repo:   "reviewpad",
					Number: 1,
				},
			},
			wantEventDetails: &entities.EventDetails{
				EventName:   "check_suite",
				EventAction: "rerequested",
				Payload: &github.CheckSuiteEvent{
					Action: github.String("rerequested"),
					CheckSuite: &github.CheckSuite{
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
		"check_suite when pr is from a forked repo": {
			event: &handler.ActionEvent{
				EventName: github.String("check_suite"),
				Token:     github.String("test-token"),
				EventPayload: buildPayload([]byte(`{
					"action": "created",
					"check_suite": {
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
			wantTargets: []*entities.TargetEntity{
				{
					Kind:        entities.PullRequest,
					Owner:       "reviewpad",
					Repo:        "reviewpad",
					Number:      6,
					AccountType: "Organization",
					Visibility:  "public",
				},
			},
			wantEventDetails: &entities.EventDetails{
				EventName:   "check_suite",
				EventAction: "created",
				Payload: &github.CheckSuiteEvent{
					Action: github.String("created"),
					CheckSuite: &github.CheckSuite{
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
