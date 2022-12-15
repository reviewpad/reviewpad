// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/jarcoal/httpmock"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/stretchr/testify/assert"
)

func buildPayload(payload []byte) *json.RawMessage {
	rawPayload := json.RawMessage(payload)
	return &rawPayload
}

func TestParseEvent_Failure(t *testing.T) {
	event := `{"type": "ping",}`
	gotEvent, err := handler.ParseEvent(event)

	assert.NotNil(t, err)
	assert.Nil(t, gotEvent)
}

func TestParseEvent(t *testing.T) {
	event := `{"action": "ping"}`
	wantEvent := &handler.ActionEvent{
		ActionName: github.String("ping"),
	}

	gotEvent, err := handler.ParseEvent(event)

	assert.Nil(t, err)
	assert.Equal(t, wantEvent, gotEvent)
}

func TestProcessEvent_Failure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

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
			gotTargets, gotEvents, gotErr := handler.ProcessEvent(test.event)

			assert.Nil(t, gotTargets)
			assert.Nil(t, gotEvents)
			assert.NotNil(t, gotErr)
		})
	}
}

func TestProcessEvent(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

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
		event       *handler.ActionEvent
		wantTargets []*handler.TargetEntity
		wantEvents  []*handler.EventData
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
			wantEvents: []*handler.EventData{
				{
					EventName:   "pull_request",
					EventAction: "opened",
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
			wantEvents: []*handler.EventData{
				{
					EventName:   "pull_request_target",
					EventAction: "opened",
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
			wantEvents: []*handler.EventData{
				{
					EventName:   "pull_request_review",
					EventAction: "opened",
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
			wantEvents: []*handler.EventData{
				{
					EventName:   "pull_request_review_comment",
					EventAction: "created",
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
			wantEvents: []*handler.EventData{
				{
					EventName: "schedule",
				},
				{
					EventName: "schedule",
				},
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
			wantEvents: []*handler.EventData{
				{
					EventName:   "workflow_run",
					EventAction: "completed",
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
			wantEvents:  []*handler.EventData{},
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
			wantEvents: []*handler.EventData{
				{
					EventName:   "issues",
					EventAction: "opened",
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
			wantEvents: []*handler.EventData{
				{
					EventName:   "issue_comment",
					EventAction: "opened",
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
			wantEvents: []*handler.EventData{
				{
					EventName: "status",
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
			wantEvents:  []*handler.EventData{},
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
			wantEvents: []*handler.EventData{
				{
					EventName: "push",
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
			wantEvents:  []*handler.EventData{},
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
					Repo:   "testrepo",
					Owner:  "testowner",
					Kind:   "",
					Number: 0,
				},
				{
					Repo:   "testrepo2",
					Owner:  "testowner2",
					Kind:   "",
					Number: 0,
				},
			},
			wantEvents: []*handler.EventData{
				{
					EventName:   "installation",
					EventAction: "created",
				},
				{
					EventName:   "installation",
					EventAction: "created",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotTargets, gotEvents, err := handler.ProcessEvent(test.event)

			assert.Nil(t, err)
			assert.ElementsMatch(t, test.wantEvents, gotEvents)
			assert.ElementsMatch(t, test.wantTargets, gotTargets)
		})
	}
}
