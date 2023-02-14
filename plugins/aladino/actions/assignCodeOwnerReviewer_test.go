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

var assignCodeOwnerReviewerCode = plugins_aladino.PluginBuiltIns().Actions["assignCodeOwnerReviewer"].Code

func TestAssignCodeOwnerReviewerCode(t *testing.T) {
	jackOpenReviewsQuery := `{
		"query": "query($query:String! $searchType:SearchType!){
			search(type: $searchType, query: $query) {
				issueCount
			}
		}",
		"variables":{
			"query": "repo:foobar/default-mock-repo is:open type:pr review-requested:jack",
			"searchType":"ISSUE"
		}
	}`
	janeOpenReviewsQuery := `{
		"query": "query($query:String! $searchType:SearchType!){
			search(type: $searchType, query: $query) {
				issueCount
			}
		}",
		"variables":{
			"query": "repo:foobar/default-mock-repo is:open type:pr review-requested:jane",
			"searchType":"ISSUE"
		}
	}`
	johnOpenReviewsQuery := `{
		"query": "query($query:String! $searchType:SearchType!){
			search(type: $searchType, query: $query) {
				issueCount
			}
		}",
		"variables":{
			"query": "repo:foobar/default-mock-repo is:open type:pr review-requested:john",
			"searchType":"ISSUE"
		}
	}`

	reviewRequestedFrom := ""
	tests := map[string]struct {
		maxReviews          *aladino.IntValue
		excludedReviewers   *aladino.ArrayValue
		wantErr             error
		mockBackendOptions  []mock.MockBackendOption
		graphQLHandler      http.HandlerFunc
		reviewRequestedFrom string
	}{
		"when pr is already assigned reviewers": {
			maxReviews:        aladino.BuildIntValue(3),
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
			maxReviews:        aladino.BuildIntValue(3),
			excludedReviewers: aladino.BuildArrayValue([]aladino.Value{}),
			mockBackendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
					&github.Reviewers{},
				),
			},
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			},
			wantErr: fmt.Errorf("error getting git blame information: %w", fmt.Errorf("error executing blame query: Message: 403 Forbidden; body: \"\", Locations: []")),
		},
		"when no blame information is found": {
			maxReviews:        aladino.BuildIntValue(3),
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
			},
			wantErr: fmt.Errorf("error getting git blame information: %w", fmt.Errorf("no blame information found")),
		},
		"when all files are owned by bot": {
			maxReviews:        aladino.BuildIntValue(3),
			excludedReviewers: aladino.BuildArrayValue([]aladino.Value{}),
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
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
						reviewRequestedFrom = reviewersRequest.Reviewers[0]
					}),
				),
			},
			reviewRequestedFrom: "test",
		},
		"when all files are owned pr creator": {
			maxReviews:        aladino.BuildIntValue(3),
			excludedReviewers: aladino.BuildArrayValue([]aladino.Value{}),
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
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
						reviewRequestedFrom = reviewersRequest.Reviewers[0]
					}),
				),
			},
			reviewRequestedFrom: "jane",
		},
		"when all code owners have too many open reviews": {
			maxReviews:        aladino.BuildIntValue(3),
			excludedReviewers: aladino.BuildArrayValue([]aladino.Value{}),
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				graphQLQuery := utils.MustRead(r.Body)
				switch utils.MinifyQuery(graphQLQuery) {
				case utils.MinifyQuery(jackOpenReviewsQuery),
					utils.MinifyQuery(janeOpenReviewsQuery),
					utils.MinifyQuery(johnOpenReviewsQuery):
					utils.MustWrite(w, `{
						"data": {
							"search": {
								"issueCount": 5
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
						reviewRequestedFrom = reviewersRequest.Reviewers[0]
					}),
				),
			},
			reviewRequestedFrom: "test",
		},
		"when first code owner is available": {
			maxReviews:        aladino.BuildIntValue(3),
			excludedReviewers: aladino.BuildArrayValue([]aladino.Value{}),
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				graphQLQuery := utils.MustRead(r.Body)
				switch utils.MinifyQuery(graphQLQuery) {
				case utils.MinifyQuery(janeOpenReviewsQuery),
					utils.MinifyQuery(johnOpenReviewsQuery):
					utils.MustWrite(w, `{
						"data": {
							"search": {
								"issueCount": 5
							}
						}
					}`)
					return
				case utils.MinifyQuery(jackOpenReviewsQuery):
					utils.MustWrite(w, `{
							"data": {
								"search": {
									"issueCount": 2
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
						reviewRequestedFrom = reviewersRequest.Reviewers[0]
					}),
				),
			},
			reviewRequestedFrom: "jack",
		},
		"when first code owner is excluded": {
			maxReviews:        aladino.BuildIntValue(3),
			excludedReviewers: aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("jack")}),
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				graphQLQuery := utils.MustRead(r.Body)
				switch utils.MinifyQuery(graphQLQuery) {
				case utils.MinifyQuery(johnOpenReviewsQuery):
					utils.MustWrite(w, `{
						"data": {
							"search": {
								"issueCount": 5
							}
						}
					}`)
					return
				case utils.MinifyQuery(janeOpenReviewsQuery), utils.MinifyQuery(jackOpenReviewsQuery):
					utils.MustWrite(w, `{
							"data": {
								"search": {
									"issueCount": 2
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
						reviewRequestedFrom = reviewersRequest.Reviewers[0]
					}),
				),
			},
			reviewRequestedFrom: "jane",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			reviewRequestedFrom = ""
			env := aladino.MockDefaultEnv(t, test.mockBackendOptions, test.graphQLHandler, nil, nil)

			err := assignCodeOwnerReviewerCode(env, []aladino.Value{test.maxReviews, test.excludedReviewers})

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.reviewRequestedFrom, reviewRequestedFrom)
		})
	}
}
