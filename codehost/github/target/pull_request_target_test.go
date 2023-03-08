// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package target_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/reviewpad/reviewpad/v4/codehost"
	host "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetLatestReviewFromReviewer(t *testing.T) {
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

	reviewSubmissionTime, err := time.Parse(time.RFC3339, "2011-01-26T19:01:12Z")
	if err != nil {
		assert.Fail(t, err.Error())
	}

	tests := map[string]struct {
		ghGraphQLHandler func(http.ResponseWriter, *http.Request)
		wantReview       *codehost.Review
		wantErr          string
	}{
		"when the request for the last review from the reviewer fails": {
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				if query == utils.MinifyQuery(mockedLatestReviewFromReviewerGQLQuery) {
					http.Error(w, "GetLatestReviewRequestFail", http.StatusNotFound)
				}
			},
			wantErr: "non-200 OK status code: 404 Not Found body: \"GetLatestReviewRequestFail\\n\"",
		},
		"when the user has no reviews": {
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				switch query {
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
				}
			},
		},
		"when the user has reviews": {
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				switch query {
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
				}
			},
			wantReview: &codehost.Review{
				User: &codehost.User{
					Login: "test",
				},
				Body:        "test",
				State:       "COMMENTED",
				SubmittedAt: &reviewSubmissionTime,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(
				t,
				nil,
				test.ghGraphQLHandler,
				aladino.MockBuiltIns(),
				nil,
			)

			gotReview, gotErr := mockedEnv.GetTarget().(*target.PullRequestTarget).GetLatestReviewFromReviewer("test")

			if gotErr != nil {
				assert.EqualError(t, gotErr, test.wantErr)
				assert.Nil(t, gotReview)
			} else {
				assert.Equal(t, test.wantReview, gotReview)
			}
		})
	}
}

func TestGetApprovedReviewers(t *testing.T) {
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetails()

	mockedPullRequestNumber := host.GetPullRequestNumber(mockedPullRequest)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedPullRequest)
	mockRepo := host.GetPullRequestBaseRepoName(mockedPullRequest)

	mockedLatestReviewsGQLQuery := fmt.Sprintf(`{
		"query":"query($pullRequestNumber:Int!$repositoryName:String!$repositoryOwner:String!){
			repository(owner:$repositoryOwner,name:$repositoryName){
				pullRequest(number:$pullRequestNumber){
					latestReviews(last:100){
						nodes{
							author{login},
							state
						}
					}
				}
			}
		}",
		"variables":{
			"pullRequestNumber":%d,
			"repositoryName":"%s",
			"repositoryOwner":"%s"
		}
	}`, mockedPullRequestNumber, mockRepo, mockOwner)

	tests := map[string]struct {
		ghGraphQLHandler      func(http.ResponseWriter, *http.Request)
		wantApprovedReviewers []string
		wantErr               error
	}{
		"when pull request latest reviews request fails": {
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				if query == utils.MinifyQuery(mockedLatestReviewsGQLQuery) {
					http.Error(w, "GetLatestReviewsRequestFail", http.StatusNotFound)
				}
			},
			wantErr: fmt.Errorf("non-200 OK status code: 404 Not Found body: \"GetLatestReviewsRequestFail\\n\""),
		},
		"when pull request has no reviews": {
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				if query == utils.MinifyQuery(mockedLatestReviewsGQLQuery) {
					utils.MustWrite(w, `{
						"data": {
							"repository": {
								"pullRequest": {
									"latestReviews": {
										"nodes": []
									}
								}
							}
						}
					}`)
				}
			},
			wantApprovedReviewers: []string{},
		},
		"when pull request has no approvals": {
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				if query == utils.MinifyQuery(mockedLatestReviewsGQLQuery) {
					utils.MustWrite(w, `{
						"data": {
							"repository": {
								"pullRequest": {
									"latestReviews": {
										"nodes": [
											{
												"state": "CHANGES_REQUESTED",
												"author": {
													"login": "test"
												}
											}
										]
									}
								}
							}
						}
					}`)
				}
			},
			wantApprovedReviewers: []string{},
		},
		"when pull request has approvals": {
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				if query == utils.MinifyQuery(mockedLatestReviewsGQLQuery) {
					utils.MustWrite(w, `{
						"data": {
							"repository": {
								"pullRequest": {
									"latestReviews": {
										"nodes": [
											{
												"state": "APPROVED",
												"author": {
													"login": "test"
												}
											}
										]
									}
								}
							}
						}
					}`)
				}
			},
			wantApprovedReviewers: []string{"test"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(
				t,
				nil,
				test.ghGraphQLHandler,
				aladino.MockBuiltIns(),
				nil,
			)

			gotApprovedReviewers, gotErr := mockedEnv.GetTarget().(*target.PullRequestTarget).GetApprovedReviewers()

			assert.Equal(t, test.wantErr, gotErr)
			assert.Equal(t, test.wantApprovedReviewers, gotApprovedReviewers)
		})
	}
}
