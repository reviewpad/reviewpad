// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	host "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var approve = plugins_aladino.PluginBuiltIns().Actions["approve"].Code

func TestApprove(t *testing.T) {
	mockedAuthenticatedUserLoginGQLQuery := `{"query":"{viewer{login}}"}`

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
		clientOptions    []mock.MockBackendOption
		ghGraphQLHandler func(http.ResponseWriter, *http.Request)
		inputReviewBody  string
		wantReview       *github.PullRequestReview
		wantErr          error
	}{
		"when the pull request is closed": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						utils.MustWriteBytes(w, mock.MustMarshal(aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
							State: github.String("closed"),
						})))
					}),
				),
			},
		},
		"when the pull request is in draft": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						utils.MustWriteBytes(w, mock.MustMarshal(aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
							State: github.String("open"),
							Draft: github.Bool(true),
						})))
					}),
				),
			},
		},
		"when the request for the authenticated user login fails": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						utils.MustWriteBytes(w, mock.MustMarshal(aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
							State: github.String("open"),
							Draft: github.Bool(false),
						})))
					}),
				),
			},
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				if query == utils.MinifyQuery(mockedAuthenticatedUserLoginGQLQuery) {
					http.Error(w, "GetAuthenticatedUserLoginRequestFail", http.StatusNotFound)
				}
			},
			wantErr: errors.New("non-200 OK status code: 404 Not Found body: \"GetAuthenticatedUserLoginRequestFail\\n\""),
		},
		"when the request for the last review made by the user fails": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						utils.MustWriteBytes(w, mock.MustMarshal(aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
							State: github.String("open"),
							Draft: github.Bool(false),
						})))
					}),
				),
			},
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				switch query {
				case utils.MinifyQuery(mockedAuthenticatedUserLoginGQLQuery):
					utils.MustWrite(w, `{
						"data": {
							"viewer": {
								"login": "test"
							}
						}
					}`)
				case utils.MinifyQuery(mockedLatestReviewFromReviewerGQLQuery):
					http.Error(w, "GetLatestReviewRequestFail", http.StatusNotFound)
				}
			},
			wantErr: errors.New("non-200 OK status code: 404 Not Found body: \"GetLatestReviewRequestFail\\n\""),
		},
		"when the request for the pull request's last push date fails": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						utils.MustWriteBytes(w, mock.MustMarshal(aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
							State: github.String("open"),
							Draft: github.Bool(false),
						})))
					}),
				),
			},
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				switch query {
				case utils.MinifyQuery(mockedAuthenticatedUserLoginGQLQuery):
					utils.MustWrite(w, `{
						"data": {
							"viewer": {
								"login": "test"
							}
						}
					}`)
				case utils.MinifyQuery(mockedLatestReviewFromReviewerGQLQuery):
					utils.MustWrite(w, `{
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
					}`)
				case utils.MinifyQuery(mockedLastPullRequestPushDateGQLQuery):
					http.Error(w, "GetLastPullRequestPushDate", http.StatusNotFound)
				}
			},
			wantErr: errors.New("non-200 OK status code: 404 Not Found body: \"GetLastPullRequestPushDate\\n\""),
		},
		"when the user hasn't made a review yet": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						utils.MustWriteBytes(w, mock.MustMarshal(aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
							State: github.String("open"),
							Draft: github.Bool(false),
						})))
					}),
				),
			},
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				switch query {
				case utils.MinifyQuery(mockedAuthenticatedUserLoginGQLQuery):
					utils.MustWrite(w, `{
						"data": {
							"viewer": {
								"login": "test"
							}
						}
					}`)
				case utils.MinifyQuery(mockedLatestReviewFromReviewerGQLQuery):
					utils.MustWrite(w, `{
						"data": {
							"repository": {
								"pullRequest": {
									"reviews": {
										"nodes": []
									}
								}
							}
						}
					}`)
				case utils.MinifyQuery(mockedLastPullRequestPushDateGQLQuery):
					utils.MustWrite(w, `{
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
					}`)
				}
			},
			inputReviewBody: "test",
			wantReview: &github.PullRequestReview{
				User: &github.User{
					Login: github.String("test"),
				},
				State: github.String("APPROVED"),
				Body:  github.String("test"),
			},
		},
		"when the pull request hasn't has any updates since the last review made by the user": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						utils.MustWriteBytes(w, mock.MustMarshal(aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
							State: github.String("open"),
							Draft: github.Bool(false),
						})))
					}),
				),
			},
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				switch query {
				case utils.MinifyQuery(mockedAuthenticatedUserLoginGQLQuery):
					utils.MustWrite(w, `{
						"data": {
							"viewer": {
								"login": "test"
							}
						}
					}`)
				case utils.MinifyQuery(mockedLatestReviewFromReviewerGQLQuery):
					utils.MustWrite(w, `{
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
					}`)
				case utils.MinifyQuery(mockedLastPullRequestPushDateGQLQuery):
					utils.MustWrite(w, `{
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
					}`)
				}
			},
			inputReviewBody: "test-2",
			wantReview: &github.PullRequestReview{
				User: &github.User{
					Login: github.String("test"),
				},
				State: github.String("APPROVED"),
				Body:  github.String("test-2"),
			},
		},
		"when the approval is the exact same as the last approval made by the user": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						utils.MustWriteBytes(w, mock.MustMarshal(aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
							State: github.String("open"),
							Draft: github.Bool(false),
						})))
					}),
				),
			},
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				switch query {
				case utils.MinifyQuery(mockedAuthenticatedUserLoginGQLQuery):
					utils.MustWrite(w, `{
						"data": {
							"viewer": {
								"login": "test"
							}
						}
					}`)
				case utils.MinifyQuery(mockedLatestReviewFromReviewerGQLQuery):
					utils.MustWrite(w, `{
						"data": {
							"repository": {
								"pullRequest": {
									"reviews": {
										"nodes": [{
											"author": {
												"login": "test"
											},
											"body": "test-3",
											"state": "APPROVED",
											"submittedAt": "2022-10-26T19:01:12Z"
										}]
									}
								}
							}
						}
					}`)
				case utils.MinifyQuery(mockedLastPullRequestPushDateGQLQuery):
					utils.MustWrite(w, `{
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
					}`)
				}
			},
			inputReviewBody: "test-3",
			wantReview: &github.PullRequestReview{
				User: &github.User{
					Login: github.String("test"),
				},
				State: github.String("APPROVED"),
				Body:  github.String("test-3"),
			},
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

								w.WriteHeader(http.StatusOK)
							}),
						),
					},
					test.clientOptions...,
				),
				test.ghGraphQLHandler,
				aladino.MockBuiltIns(),
				nil,
			)

			args := []aladino.Value{aladino.BuildStringValue(test.inputReviewBody)}
			gotErr := approve(mockedEnv, args)

			assert.Equal(t, test.wantErr, gotErr)
		})
	}
}
