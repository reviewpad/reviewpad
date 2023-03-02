// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	host "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var review = plugins_aladino.PluginBuiltIns().Actions["review"].Code

type ReviewEvent struct {
	Body  string
	Event string
}

func TestReview_WhenAuthenticatedUserLoginRequestFails(t *testing.T) {
	failMessage := "GetAuthenticatedUserLoginRequestFail"
	mockedAuthenticatedUserLoginGQLQuery := `{"query":"{viewer{login}}"}`
	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			if query == utils.MinifyQuery(mockedAuthenticatedUserLoginGQLQuery) {
				http.Error(w, failMessage, http.StatusNotFound)
			}
		},
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue("APPROVE"), aladino.BuildStringValue("test")}
	gotErr := review(mockedEnv, args)

	assert.EqualError(t, gotErr, fmt.Sprintf("non-200 OK status code: 404 Not Found body: \"%s\\n\"", failMessage))
}

func TestReview_WhenLatestReviewRequestFails(t *testing.T) {
	failMessage := "GetLatestReviewRequestFail"
	mockedAuthenticatedUserLoginGQLQuery := `{"query":"{viewer{login}}"}`
	mockedAuthenticatedUserLoginGQLQueryBody := `{
		"data": {
			"viewer": {
				"login": "test"
			}
		}
	}`

	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedAuthenticatedUserLoginGQLQuery):
				utils.MustWrite(w, mockedAuthenticatedUserLoginGQLQueryBody)
			default:
				http.Error(w, failMessage, http.StatusNotFound)
			}
		},
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue("APPROVE"), aladino.BuildStringValue("test")}
	gotErr := review(mockedEnv, args)

	assert.EqualError(t, gotErr, fmt.Sprintf("non-200 OK status code: 404 Not Found body: \"%s\\n\"", failMessage))
}

func TestReview_WhenLastPushDateRequestFails(t *testing.T) {
	failMessage := "GetPullRequestLastPushDate"
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetails()

	mockedPullRequestNumber := host.GetPullRequestNumber(mockedPullRequest)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedPullRequest)
	mockRepo := host.GetPullRequestBaseRepoName(mockedPullRequest)

	mockedAuthenticatedUserLoginGQLQuery := `{"query":"{viewer{login}}"}`
	mockedAuthenticatedUserLoginGQLQueryBody := `{
		"data": {
			"viewer": {
				"login": "test"
			}
		}
	}`

	mockedLatestReviewFromReviewerGQLQuery := fmt.Sprintf(`{
		"query":"query($author:String!$pullRequestNumber:Int!$repositoryName:String!$repositoryOwner:String!){
			repository(owner:$repositoryOwner,name:$repositoryName){
				pullRequest(number:$pullRequestNumber){
					reviews(last:1,author:$author){
						nodes{
							author{login},
							body,
							state,
							submittedAt
						}
					}
				}
			}
		}",
		"variables":{
			"author":"test",
			"pullRequestNumber":%d,
			"repositoryName":"%s",
			"repositoryOwner":"%s"
		}
	}`, mockedPullRequestNumber, mockRepo, mockOwner)

	mockedLatestReviewFromReviewerGQLQueryBody := `{
		"data": {
			"repository": {
				"pullRequest": {
					"reviews": {
						"nodes": [{
							"author": {
								"login": "test"
							},
							"body": "test",
							"state": "COMMENTED",
							"submittedAt": "2011-01-26T19:01:12Z"
						}]
					}
				}
			}
		}
	}`

	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedAuthenticatedUserLoginGQLQuery):
				utils.MustWrite(w, mockedAuthenticatedUserLoginGQLQueryBody)
			case utils.MinifyQuery(mockedLatestReviewFromReviewerGQLQuery):
				utils.MustWrite(w, mockedLatestReviewFromReviewerGQLQueryBody)
			default:
				http.Error(w, failMessage, http.StatusNotFound)
			}
		},
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue("APPROVE"), aladino.BuildStringValue("test")}
	gotErr := review(mockedEnv, args)

	assert.EqualError(t, gotErr, fmt.Sprintf("non-200 OK status code: 404 Not Found body: \"%s\\n\"", failMessage))
}

func TestReview_WhenPostReviewRequestFail(t *testing.T) {
	failMessage := "PostReviewRequestFail"
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetails()

	mockedPullRequestNumber := host.GetPullRequestNumber(mockedPullRequest)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedPullRequest)
	mockRepo := host.GetPullRequestBaseRepoName(mockedPullRequest)

	mockedLatestReviewFromReviewerGQLQuery := fmt.Sprintf(`{
		"query":"query($author:String!$pullRequestNumber:Int!$repositoryName:String!$repositoryOwner:String!){
			repository(owner:$repositoryOwner,name:$repositoryName){
				pullRequest(number:$pullRequestNumber){
					reviews(last:1,author:$author){
						nodes{
							author{login},
							body,
							state,
							submittedAt
						}
					}
				}
			}
		}",
		"variables":{
			"author":"test",
			"pullRequestNumber":%d,
			"repositoryName":"%s",
			"repositoryOwner":"%s"
		}
	}`, mockedPullRequestNumber, mockRepo, mockOwner)

	mockedLatestReviewFromReviewerGQLQueryBody := `{
		"data": {
			"repository": {
				"pullRequest": {
					"reviews": {
						"nodes": [{
							"author": {
								"login": "test"
							},
							"body": "test",
							"state": "COMMENTED",
							"submittedAt": "2011-01-26T19:01:12Z"
						}]
					}
				}
			}
		}
	}`

	mockedAuthenticatedUserLoginGQLQuery := `{"query":"{viewer{login}}"}`

	mockedAuthenticatedUserLoginGQLQueryBody := `{
		"data": {
			"viewer": {
				"login": "test"
			}
		}
	}`

	mockedLastPullRequestPushDateGQLQuery := fmt.Sprintf(`{
		"query":"query($pullRequestNumber:Int! $repositoryName:String! $repositoryOwner:String!){
			repository(owner: $repositoryOwner, name: $repositoryName){
				pullRequest(number: $pullRequestNumber){
					timelineItems(last: 1, itemTypes: [HEAD_REF_FORCE_PUSHED_EVENT, PULL_REQUEST_COMMIT]){
						nodes{
							__typename,
							...on HeadRefForcePushedEvent {
								createdAt
							},
							...on PullRequestCommit {
								commit {
									pushedDate,
									committedDate
								}
							}
						}
					}
				}
			}
		}",
		"variables":{
			"pullRequestNumber": %d,
			"repositoryName": "%s",
			"repositoryOwner": "%s"
		}
	}`, mockedPullRequestNumber, mockRepo, mockOwner)

	mockedLastPullRequestPushDateGQLQueryBody := `{
		"data": {
			"repository": {
				"pullRequest": {
					"timelineItems": {
						"nodes": [{
							"__typename": "PullRequestCommit",
							"commit": {
								"pushedDate": "2022-11-11T13:36:05Z",
								"committedDate": "2022-11-11T13:36:05Z"
							}
						}]
					}
				}
			}
		}
	}`

	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.PostReposPullsReviewsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						failMessage,
					)
				}),
			),
		},
		func(w http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedAuthenticatedUserLoginGQLQuery):
				utils.MustWrite(w, mockedAuthenticatedUserLoginGQLQueryBody)
			case utils.MinifyQuery(mockedLatestReviewFromReviewerGQLQuery):
				utils.MustWrite(w, mockedLatestReviewFromReviewerGQLQueryBody)
			case utils.MinifyQuery(mockedLastPullRequestPushDateGQLQuery):
				utils.MustWrite(w, mockedLastPullRequestPushDateGQLQueryBody)
			}
		},
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue("APPROVE"), aladino.BuildStringValue("test-2")}
	gotErr := review(mockedEnv, args)

	assert.Equal(t, gotErr.(*github.ErrorResponse).Message, failMessage)
}

func TestReview(t *testing.T) {
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetails()

	mockedPullRequestNumber := host.GetPullRequestNumber(mockedPullRequest)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedPullRequest)
	mockRepo := host.GetPullRequestBaseRepoName(mockedPullRequest)

	mockedLatestReviewFromReviewerGQLQuery := fmt.Sprintf(`{
		"query":"query($author:String!$pullRequestNumber:Int!$repositoryName:String!$repositoryOwner:String!){
			repository(owner:$repositoryOwner,name:$repositoryName){
				pullRequest(number:$pullRequestNumber){
					reviews(last:1,author:$author){
						nodes{
							author{login},
							body,
							state,
							submittedAt
						}
					}
				}
			}
		}",
		"variables":{
			"author":"test",
			"pullRequestNumber":%d,
			"repositoryName":"%s",
			"repositoryOwner":"%s"
		}
	}`, mockedPullRequestNumber, mockRepo, mockOwner)

	mockedAuthenticatedUserLoginGQLQuery := `{"query":"{viewer{login}}"}`

	mockedAuthenticatedUserLoginGQLQueryBody := `{
		"data": {
			"viewer": {
				"login": "test"
			}
		}
	}`

	mockedLastPullRequestPushDateGQLQuery := fmt.Sprintf(`{
		"query":"query($pullRequestNumber:Int! $repositoryName:String! $repositoryOwner:String!){
			repository(owner: $repositoryOwner, name: $repositoryName){
				pullRequest(number: $pullRequestNumber){
					timelineItems(last: 1, itemTypes: [HEAD_REF_FORCE_PUSHED_EVENT, PULL_REQUEST_COMMIT]){
						nodes{
							__typename,
							...on HeadRefForcePushedEvent {
								createdAt
							},
							...on PullRequestCommit {
								commit {
									pushedDate,
									committedDate
								}
							}
						}
					}
				}
			}
		}",
		"variables":{
			"pullRequestNumber": %d,
			"repositoryName": "%s",
			"repositoryOwner": "%s"
		}
	}`, mockedPullRequestNumber, mockRepo, mockOwner)

	tests := map[string]struct {
		clientOptions                              []mock.MockBackendOption
		mockedLatestReviewFromReviewerGQLQueryBody string
		mockedLastPullRequestPushDateGQLQueryBody  string
		inputReviewEvent                           string
		inputReviewBody                            string
		wantReview                                 *github.PullRequestReview
		wantErr                                    error
	}{
		"when pull request is closed": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						utils.MustWriteBytes(w, mock.MustMarshal(aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
							State: github.String("closed"),
						})))
					}),
				),
			},
			inputReviewEvent: "COMMENT",
			inputReviewBody:  "test",
			wantErr:          nil,
		},
		"when review event is not supported": {
			clientOptions:    []mock.MockBackendOption{},
			inputReviewEvent: "NOT_SUPPORTED_EVENT",
			inputReviewBody:  "test",
			wantErr:          fmt.Errorf("review: unsupported review state NOT_SUPPORTED_EVENT"),
		},
		"when review event is not an APPROVE and its body is empty": {
			clientOptions:    []mock.MockBackendOption{},
			inputReviewEvent: "COMMENT",
			wantErr:          fmt.Errorf("review: comment required in COMMENT state"),
		},
		"when authenticated user has not made any review": {
			clientOptions: []mock.MockBackendOption{},
			mockedLatestReviewFromReviewerGQLQueryBody: `{
				"data": {
					"repository": {
						"pullRequest": {
							"reviews": {
								"nodes": []
							}
						}
					}
				}
			}`,
			mockedLastPullRequestPushDateGQLQueryBody: `{
				"data": {
					"repository": {
						"pullRequest": {
							"timelineItems": {
								"nodes": [{
									"__typename": "PullRequestCommit",
									"commit": {
										"pushedDate": "2022-11-11T13:36:05Z",
										"committedDate": "2022-11-11T13:36:05Z"
									}
								}]
							}
						}
					}
				}
			}`,
			inputReviewEvent: "COMMENT",
			inputReviewBody:  "test",
			wantReview: &github.PullRequestReview{
				User: &github.User{
					Login: github.String("test"),
				},
				State: github.String("COMMENTED"),
				Body:  github.String("test"),
			},
		},
		"when there has been no pull request pushes since authenticated user's last review": {
			clientOptions: []mock.MockBackendOption{},
			mockedLatestReviewFromReviewerGQLQueryBody: `{
				"data": {
					"repository": {
						"pullRequest": {
							"reviews": {
								"nodes": [{
									"author": {
										"login": "test"
									},
									"body": "test",
									"state": "CHANGES_REQUESTED",
									"submittedAt": "2022-11-26T19:01:12Z"
								}]
							}
						}
					}
				}
			}`,
			mockedLastPullRequestPushDateGQLQueryBody: `{
				"data": {
					"repository": {
						"pullRequest": {
							"timelineItems": {
								"nodes": [{
									"__typename": "PullRequestCommit",
									"commit": {
										"pushedDate": "2022-11-11T13:36:05Z",
										"committedDate": "2022-11-11T13:36:05Z"
									}
								}]
							}
						}
					}
				}
			}`,
			inputReviewEvent: "REQUEST_CHANGES",
			inputReviewBody:  "test-2",
			wantReview: &github.PullRequestReview{
				User: &github.User{
					Login: github.String("test"),
				},
				State: github.String("CHANGES_REQUESTED"),
				Body:  github.String("test-2"),
			},
		},
		"when authenticated user has made a review whose state is not valid": {
			clientOptions: []mock.MockBackendOption{},
			mockedLatestReviewFromReviewerGQLQueryBody: `{
				"data": {
					"repository": {
						"pullRequest": {
							"reviews": {
								"nodes": [{
									"author": {
										"login": "test"
									},
									"body": "test",
									"state": "NOT_SUPPORTED",
									"submittedAt": "2022-10-26T19:01:12Z"
								}]
							}
						}
					}
				}
			}`,
			mockedLastPullRequestPushDateGQLQueryBody: `{
				"data": {
					"repository": {
						"pullRequest": {
							"timelineItems": {
								"nodes": [{
									"__typename": "PullRequestCommit",
									"commit": {
										"pushedDate": "2022-11-11T13:36:05Z",
										"committedDate": "2022-11-11T13:36:05Z"
									}
								}]
							}
						}
					}
				}
			}`,
			inputReviewEvent: "APPROVE",
			inputReviewBody:  "test",
			wantErr:          fmt.Errorf("review: unsupported review state NOT_SUPPORTED"),
		},
		"when latest review is the same as the one we want to create": {
			clientOptions: []mock.MockBackendOption{},
			mockedLatestReviewFromReviewerGQLQueryBody: `{
				"data": {
					"repository": {
						"pullRequest": {
							"reviews": {
								"nodes": [{
									"author": {
										"login": "test"
									},
									"body": "test",
									"state": "APPROVED",
									"submittedAt": "2022-11-26T19:01:12Z"
								}]
							}
						}
					}
				}
			}`,
			mockedLastPullRequestPushDateGQLQueryBody: `{
				"data": {
					"repository": {
						"pullRequest": {
							"timelineItems": {
								"nodes": [{
									"__typename": "PullRequestCommit",
									"commit": {
										"pushedDate": "2022-11-11T13:36:05Z",
										"committedDate": "2022-11-11T13:36:05Z"
									}
								}]
							}
						}
					}
				}
			}`,
			inputReviewEvent: "APPROVE",
			inputReviewBody:  "test",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(
				t,
				append(
					[]mock.MockBackendOption{
						mock.WithRequestMatchHandler(
							mock.PostReposPullsReviewsByOwnerByRepoByPullNumber,
							http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
								var gotReview ReviewEvent
								rawBody, _ := io.ReadAll(r.Body)
								utils.MustUnmarshal(rawBody, &gotReview)

								assert.Equal(t, *test.wantReview.Body, gotReview.Body)
								switch gotReview.Event {
								case "REQUEST_CHANGES":
									assert.Equal(t, *test.wantReview.State, "CHANGES_REQUESTED")
								case "COMMENT":
									assert.Equal(t, *test.wantReview.State, "COMMENTED")
								case "APPROVE":
									assert.Equal(t, *test.wantReview.State, "APPROVED")
								}

								w.WriteHeader(http.StatusOK)
							}),
						),
					},
					test.clientOptions...,
				),
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					switch query {
					case utils.MinifyQuery(mockedAuthenticatedUserLoginGQLQuery):
						utils.MustWrite(w, mockedAuthenticatedUserLoginGQLQueryBody)
					case utils.MinifyQuery(mockedLatestReviewFromReviewerGQLQuery):
						utils.MustWrite(w, test.mockedLatestReviewFromReviewerGQLQueryBody)
					case utils.MinifyQuery(mockedLastPullRequestPushDateGQLQuery):
						utils.MustWrite(w, test.mockedLastPullRequestPushDateGQLQueryBody)
					}
				},
				aladino.MockBuiltIns(),
				nil,
			)

			args := []aladino.Value{aladino.BuildStringValue(test.inputReviewEvent), aladino.BuildStringValue(test.inputReviewBody)}
			gotErr := review(mockedEnv, args)

			assert.Equal(t, test.wantErr, gotErr)
		})
	}
}
