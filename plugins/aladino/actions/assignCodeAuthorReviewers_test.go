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
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var assignCodeAuthorReviewer = plugins_aladino.PluginBuiltIns().Actions["assignCodeAuthorReviewers"].Code

func TestAssignCodeAuthorReviewerCode(t *testing.T) {
	getOpenPullRequestsAsReviewerQuery := `{
		"query":"query($name:String!$owner:String!){
			repository(owner: $owner, name: $name){
				pullRequests(states: OPEN, last: 50){
					nodes{reviewRequests(first: 50){
						nodes{
							requestedReviewer{
								__typename,
								... on User{
									login
								}
							}
						}
					},
					reviews(first: 50){
						nodes{
							author{
								login
							}
						}
					}
				}
			}
		}
	}",
	"variables":{
		"name":"default-mock-repo",
		"owner":"foobar"
	}
	}`

	reviewRequestedFrom := []string{}
	tests := map[string]struct {
		totalReviewers      int
		maxReviews          int
		excludedReviewers   *aladino.ArrayValue
		wantErr             error
		mockBackendOptions  []mock.MockBackendOption
		graphQLHandler      http.HandlerFunc
		reviewRequestedFrom []string
	}{
		"when the pull request is already assigned reviewers": {
			totalReviewers:    1,
			maxReviews:        3,
			excludedReviewers: aladino.BuildArrayValue([]aladino.Value{}),
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
					&github.Reviewers{
						Users: []*github.User{
							{
								Login: github.String("test"),
							},
						},
					},
				),
			},
		},
		"when error getting git blame information": {
			totalReviewers:    1,
			maxReviews:        3,
			excludedReviewers: aladino.BuildArrayValue([]aladino.Value{}),
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
					&github.Reviewers{},
				),
				mock.WithRequestMatch(
					mock.GetReposAssigneesByOwnerByRepo,
					[]*github.User{
						{
							Login: github.String("test"),
						},
					},
				),
			},
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			},
			wantErr: fmt.Errorf("error getting authors from git blame: error executing blame query: Message: 403 Forbidden; body: \"\", Locations: []"),
		},
		"when no blame information is found": {
			totalReviewers:    1,
			maxReviews:        3,
			excludedReviewers: aladino.BuildArrayValue([]aladino.Value{}),
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"object": {
							}
						}
					}
				}`)
			},
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
					&github.Reviewers{},
				),
				mock.WithRequestMatch(
					mock.GetReposAssigneesByOwnerByRepo,
					[]*github.User{
						{
							Login: github.String("test"),
						},
					},
				),
			},
			wantErr: fmt.Errorf("error getting authors from git blame: no blame information found"),
		},
		"when get open pull requests as reviewer query fails": {
			totalReviewers:    1,
			maxReviews:        3,
			excludedReviewers: aladino.BuildArrayValue([]aladino.Value{}),
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				graphQLQuery := utils.MustRead(r.Body)
				if utils.MinifyQuery(getOpenPullRequestsAsReviewerQuery) == utils.MinifyQuery(graphQLQuery) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"object": {
								"blame0": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 10,
											"commit": {
												"author": {
													"user": {
														"login": "reviewpad[bot]"
													}
												}
											}
										}
									]
								},
								"blame1": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 5,
											"commit": {
												"author": {
													"user": {
														"login": "reviewpad[bot]"
													}
												}
											}
										}
									]
								},
								"blame3": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 10,
											"commit": {
												"author": {
													"user": {
														"login": "reviewpad[bot]"
													}
												}
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
					&github.Reviewers{},
				),
				mock.WithRequestMatch(
					mock.GetReposAssigneesByOwnerByRepo,
					[]*github.User{
						{
							Login: github.String("test"),
						},
					},
				),
			},
			wantErr: fmt.Errorf("error getting open pull requests as reviewer: non-200 OK status code: 400 Bad Request body: \"\""),
		},
		"when all files are owned by bot": {
			totalReviewers:    1,
			maxReviews:        3,
			excludedReviewers: aladino.BuildArrayValue([]aladino.Value{}),
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				graphQLQuery := utils.MustRead(r.Body)
				if utils.MinifyQuery(getOpenPullRequestsAsReviewerQuery) == utils.MinifyQuery(graphQLQuery) {
					utils.MustWrite(w, `{
						"data": {
							"repository": {
								"pullRequests": {
									"nodes": []
								}
							}
						}
					}`)
					return
				}

				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"object": {
								"blame0": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 10,
											"commit": {
												"author": {
													"user": {
														"login": "reviewpad[bot]"
													}
												}
											}
										}
									]
								},
								"blame1": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 5,
											"commit": {
												"author": {
													"user": {
														"login": "reviewpad[bot]"
													}
												}
											}
										}
									]
								},
								"blame3": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 10,
											"commit": {
												"author": {
													"user": {
														"login": "reviewpad[bot]"
													}
												}
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
					&github.Reviewers{},
					&github.Reviewers{},
				),
				mock.WithRequestMatch(
					mock.GetReposAssigneesByOwnerByRepo,
					[]*github.User{
						{
							Login: github.String("test"),
						},
					},
					[]*github.User{
						{
							Login: github.String("test"),
						},
					},
				),
				mock.WithRequestMatchHandler(
					mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						body, _ := io.ReadAll(r.Body)
						reviewersRequest := &github.ReviewersRequest{}
						utils.MustUnmarshal(body, reviewersRequest)
						reviewRequestedFrom = reviewersRequest.Reviewers
					}),
				),
			},
			reviewRequestedFrom: []string{"test"},
		},
		"when all files are owned pull request author": {
			totalReviewers:    1,
			maxReviews:        3,
			excludedReviewers: aladino.BuildArrayValue([]aladino.Value{}),
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				graphQLQuery := utils.MustRead(r.Body)
				if utils.MinifyQuery(getOpenPullRequestsAsReviewerQuery) == utils.MinifyQuery(graphQLQuery) {
					utils.MustWrite(w, `{
						"data": {
							"repository": {
								"pullRequests": {
									"nodes": []
								}
							}
						}
					}`)
					return
				}

				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"object": {
								"blame0": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 10,
											"commit": {
												"author": {
													"user": {
														"login": "john"
													}
												}
											}
										}
									]
								},
								"blame1": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 5,
											"commit": {
												"author": {
													"user": {
														"login": "john"
													}
												}
											}
										}
									]
								},
								"blame3": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 10,
											"commit": {
												"author": {
													"user": {
														"login": "john"
													}
												}
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
					&github.Reviewers{},
					&github.Reviewers{},
				),
				mock.WithRequestMatch(
					mock.GetReposAssigneesByOwnerByRepo,
					[]*github.User{
						{
							Login: github.String("john"),
						},
						{
							Login: github.String("jane"),
						},
					},
					[]*github.User{
						{
							Login: github.String("john"),
						},
						{
							Login: github.String("jane"),
						},
					},
				),
				mock.WithRequestMatchHandler(
					mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						body, _ := io.ReadAll(r.Body)
						reviewersRequest := &github.ReviewersRequest{}
						utils.MustUnmarshal(body, reviewersRequest)
						reviewRequestedFrom = reviewersRequest.Reviewers
					}),
				),
			},
			reviewRequestedFrom: []string{"jane"},
		},
		"when all code owners are handling too many open pull requests": {
			totalReviewers:    1,
			maxReviews:        1,
			excludedReviewers: aladino.BuildArrayValue([]aladino.Value{}),
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				graphQLQuery := utils.MustRead(r.Body)
				if utils.MinifyQuery(getOpenPullRequestsAsReviewerQuery) == utils.MinifyQuery(graphQLQuery) {
					utils.MustWrite(w, `{
						"data": {
							"repository": {
								"pullRequests": {
									"nodes": [
										{
											"reviewRequests": {
												"nodes": [
													{
														"requestedReviewer": {
															"login": "john"
														}
													},
													{
														"requestedReviewer": {
															"login": "jane"
														}
													},
													{
														"requestedReviewer": {
															"login": "jack"
														}
													}
												]
											}
										}
									]
								}
							}
						}
					}`)
					return
				}

				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"object": {
								"blame0": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 10,
											"commit": {
												"author": {
													"user": {
														"login": "john"
													}
												}
											}
										},
										{
											"startingLine": 11,
											"endingLine": 30,
											"commit": {
												"author": {
													"user": {
														"login": "jane"
													}
												}
											}
										},
										{
											"startingLine": 31,
											"endingLine": 100,
											"commit": {
												"author": {
													"user": {
														"login": "jack"
													}
												}
											}
										}
									]
								},
								"blame1": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 5,
											"commit": {
												"author": {
													"user": {
														"login": "john"
													}
												}
											}
										},
										{
											"startingLine": 6,
											"endingLine": 35,
											"commit": {
												"author": {
													"user": {
														"login": "jane"
													}
												}
											}
										},
										{
											"startingLine": 35,
											"endingLine": 111,
											"commit": {
												"author": {
													"user": {
														"login": "jack"
													}
												}
											}
										}
									]
								},
								"blame3": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 10,
											"commit": {
												"author": {
													"user": {
														"login": "john"
													}
												}
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
					&github.Reviewers{},
					&github.Reviewers{},
				),
				mock.WithRequestMatch(
					mock.GetReposAssigneesByOwnerByRepo,
					[]*github.User{
						{
							Login: github.String("test"),
						},
					},
					[]*github.User{
						{
							Login: github.String("test"),
						},
					},
				),
				mock.WithRequestMatchHandler(
					mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						body, _ := io.ReadAll(r.Body)
						reviewersRequest := &github.ReviewersRequest{}
						utils.MustUnmarshal(body, reviewersRequest)
						reviewRequestedFrom = reviewersRequest.Reviewers
					}),
				),
			},
			reviewRequestedFrom: []string{"test"},
		},
		"when first code owner is available": {
			totalReviewers:    1,
			maxReviews:        1,
			excludedReviewers: aladino.BuildArrayValue([]aladino.Value{}),
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				graphQLQuery := utils.MustRead(r.Body)
				if utils.MinifyQuery(getOpenPullRequestsAsReviewerQuery) == utils.MinifyQuery(graphQLQuery) {
					utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequests": {
								"nodes": [
									{
										"reviewRequests": {
											"nodes": [
												{
													"requestedReviewer": {
														"login": "john"
													}
												},
												{
													"requestedReviewer": {
														"login": "jane"
													}
												}
											]
										}
									}
								]
							}
						}
					}
				}`)
					return
				}

				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"object": {
								"blame0": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 10,
											"commit": {
												"author": {
													"user": {
														"login": "john"
													}
												}
											}
										},
										{
											"startingLine": 11,
											"endingLine": 30,
											"commit": {
												"author": {
													"user": {
														"login": "jane"
													}
												}
											}
										},
										{
											"startingLine": 31,
											"endingLine": 100,
											"commit": {
												"author": {
													"user": {
														"login": "jack"
													}
												}
											}
										}
									]
								},
								"blame1": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 5,
											"commit": {
												"author": {
													"user": {
														"login": "john"
													}
												}
											}
										},
										{
											"startingLine": 6,
											"endingLine": 35,
											"commit": {
												"author": {
													"user": {
														"login": "jane"
													}
												}
											}
										},
										{
											"startingLine": 35,
											"endingLine": 111,
											"commit": {
												"author": {
													"user": {
														"login": "jack"
													}
												}
											}
										}
									]
								},
								"blame3": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 10,
											"commit": {
												"author": {
													"user": {
														"login": "john"
													}
												}
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
					&github.Reviewers{},
				),
				mock.WithRequestMatch(
					mock.GetReposAssigneesByOwnerByRepo,
					[]*github.User{
						{
							Login: github.String("jane"),
						},
						{
							Login: github.String("jack"),
						},
						{
							Login: github.String("john"),
						},
					},
				),
				mock.WithRequestMatchHandler(
					mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						body, _ := io.ReadAll(r.Body)
						reviewersRequest := &github.ReviewersRequest{}
						utils.MustUnmarshal(body, reviewersRequest)
						reviewRequestedFrom = reviewersRequest.Reviewers
					}),
				),
			},
			reviewRequestedFrom: []string{"jack"},
		},
		"when first code owner is excluded": {
			totalReviewers:    3,
			maxReviews:        3,
			excludedReviewers: aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("jack")}),
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				graphQLQuery := utils.MustRead(r.Body)
				if utils.MinifyQuery(getOpenPullRequestsAsReviewerQuery) == utils.MinifyQuery(graphQLQuery) {
					utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequests": {
								"nodes": []
							}
						}
					}
				}`)
					return
				}

				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"object": {
								"blame0": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 10,
											"commit": {
												"author": {
													"user": {
														"login": "james"
													}
												}
											}
										},
										{
											"startingLine": 11,
											"endingLine": 30,
											"commit": {
												"author": {
													"user": {
														"login": "jane"
													}
												}
											}
										},
										{
											"startingLine": 31,
											"endingLine": 100,
											"commit": {
												"author": {
													"user": {
														"login": "jack"
													}
												}
											}
										}
									]
								},
								"blame1": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 5,
											"commit": {
												"author": {
													"user": {
														"login": "james"
													}
												}
											}
										},
										{
											"startingLine": 6,
											"endingLine": 35,
											"commit": {
												"author": {
													"user": {
														"login": "jane"
													}
												}
											}
										},
										{
											"startingLine": 35,
											"endingLine": 111,
											"commit": {
												"author": {
													"user": {
														"login": "jack"
													}
												}
											}
										}
									]
								},
								"blame3": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 10,
											"commit": {
												"author": {
													"user": {
														"login": "james"
													}
												}
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
					&github.Reviewers{},
				),
				mock.WithRequestMatch(
					mock.GetReposAssigneesByOwnerByRepo,
					[]*github.User{
						{
							Login: github.String("jane"),
						},
						{
							Login: github.String("jack"),
						},
						{
							Login: github.String("john"),
						},
						{
							Login: github.String("james"),
						},
					},
				),
				mock.WithRequestMatchHandler(
					mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						body, _ := io.ReadAll(r.Body)
						reviewersRequest := &github.ReviewersRequest{}
						utils.MustUnmarshal(body, reviewersRequest)
						reviewRequestedFrom = reviewersRequest.Reviewers
					}),
				),
			},
			reviewRequestedFrom: []string{"jane", "james"},
		},
		"when code owner isn't an available assignee": {
			totalReviewers:    1,
			maxReviews:        1,
			excludedReviewers: aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("jack")}),
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				graphQLQuery := utils.MustRead(r.Body)
				if utils.MinifyQuery(getOpenPullRequestsAsReviewerQuery) == utils.MinifyQuery(graphQLQuery) {
					utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequests": {
								"nodes": [
									{
										"reviewRequests": {
											"nodes": [
												{
													"requestedReviewer": {
														"login": "jack"
													}
												},
												{
													"requestedReviewer": {
														"login": "james"
													}
												},
												{
													"requestedReviewer": {
														"login": "john"
													}
												}
											]
										}
									}
								]
							}
						}
					}
				}`)
					return
				}

				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"object": {
								"blame0": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 10,
											"commit": {
												"author": {
													"user": {
														"login": "james"
													}
												}
											}
										},
										{
											"startingLine": 11,
											"endingLine": 30,
											"commit": {
												"author": {
													"user": {
														"login": "jane"
													}
												}
											}
										},
										{
											"startingLine": 31,
											"endingLine": 100,
											"commit": {
												"author": {
													"user": {
														"login": "jack"
													}
												}
											}
										}
									]
								},
								"blame1": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 5,
											"commit": {
												"author": {
													"user": {
														"login": "james"
													}
												}
											}
										},
										{
											"startingLine": 6,
											"endingLine": 35,
											"commit": {
												"author": {
													"user": {
														"login": "jane"
													}
												}
											}
										},
										{
											"startingLine": 35,
											"endingLine": 111,
											"commit": {
												"author": {
													"user": {
														"login": "jack"
													}
												}
											}
										}
									]
								},
								"blame3": {
									"ranges": [
										{
											"startingLine": 1,
											"endingLine": 10,
											"commit": {
												"author": {
													"user": {
														"login": "james"
													}
												}
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
					&github.Reviewers{},
					&github.Reviewers{},
				),
				mock.WithRequestMatch(
					mock.GetReposAssigneesByOwnerByRepo,
					[]*github.User{
						{
							Login: github.String("test"),
						},
					},
					[]*github.User{
						{
							Login: github.String("test"),
						},
					},
				),
				mock.WithRequestMatchHandler(
					mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						body, _ := io.ReadAll(r.Body)
						reviewersRequest := &github.ReviewersRequest{}
						utils.MustUnmarshal(body, reviewersRequest)
						reviewRequestedFrom = reviewersRequest.Reviewers
					}),
				),
			},
			reviewRequestedFrom: []string{"test"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			reviewRequestedFrom = nil
			env := aladino.MockDefaultEnv(t, test.mockBackendOptions, test.graphQLHandler, nil, nil)

			err := assignCodeAuthorReviewer(env, []aladino.Value{aladino.BuildIntValue(test.totalReviewers), test.excludedReviewers, aladino.BuildIntValue(test.maxReviews)})

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.reviewRequestedFrom, reviewRequestedFrom)
		})
	}
}
