// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var hasUnaddressedThreads = plugins_aladino.PluginBuiltIns().Functions["hasUnaddressedThreads"].Code

func TestHasUnaddressedThreads_WhenRequestFails(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		},
		aladino.MockBuiltIns(),
		nil,
	)

	_, err := hasUnaddressedThreads(mockedEnv, []aladino.Value{})

	assert.NotNil(t, err)
}

func TestHasUnaddressedThreads(t *testing.T) {
	mockedGraphQLQuery := fmt.Sprintf(
		"{\"query\":\"query($pullRequestNumber:Int!$repositoryName:String!$repositoryOwner:String!$reviewThreadsCursor:String){repository(owner: $repositoryOwner, name: $repositoryName){pullRequest(number: $pullRequestNumber){reviewThreads(first: 10, after: $reviewThreadsCursor){nodes{isResolved,isOutdated},pageInfo{endCursor,hasNextPage}}}}}\",\"variables\":{\"pullRequestNumber\":%d,\"repositoryName\":\"%s\",\"repositoryOwner\":\"%s\",\"reviewThreadsCursor\":null}}\n",
		aladino.DefaultMockPrNum,
		aladino.DefaultMockPrRepoName,
		aladino.DefaultMockPrOwner,
	)

	tests := map[string]struct {
		reviewThreadsResponse string
		wantVal               aladino.Value
	}{
		"when none": {
			reviewThreadsResponse: `[]`,
			wantVal:               aladino.BuildFalseValue(),
		},
		"when unresolved": {
			reviewThreadsResponse: `[{"isResolved":false, "isOutdated":false}]`,
			wantVal:               aladino.BuildTrueValue(),
		},
		"when resolved": {
			reviewThreadsResponse: `[{"isResolved":true, "isOutdated":false}]`,
			wantVal:               aladino.BuildFalseValue(),
		},
		"when outdated": {
			reviewThreadsResponse: `[{"isResolved":false, "isOutdated":true}]`,
			wantVal:               aladino.BuildFalseValue(),
		},
		"when up to date": {
			reviewThreadsResponse: `[{"isResolved":false, "isOutdated":false}]`,
			wantVal:               aladino.BuildTrueValue(),
		},
		"when more than one thread": {
			reviewThreadsResponse: `[{"isResolved":true, "isOutdated":true}, {"isResolved":true, "isOutdated":true}]`,
			wantVal:               aladino.BuildFalseValue(),
		},
	}

	for name, test := range tests {
		mockedEnv := aladino.MockDefaultEnv(
			t,
			[]mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.Write(mock.MustMarshal(
							aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
								Number: github.Int(aladino.DefaultMockPrNum),
								Base: &github.PullRequestBranch{
									Repo: &github.Repository{
										Owner: &github.User{
											Login: github.String(aladino.DefaultMockPrOwner),
										},
										Name: github.String(aladino.DefaultMockPrRepoName),
									},
								},
							}),
						))
					}),
				),
			},
			func(w http.ResponseWriter, req *http.Request) {
				query := aladino.MustRead(req.Body)
				switch query {
				case mockedGraphQLQuery:
					aladino.MustWrite(
						w,
						fmt.Sprintf(
							`{"data": {
                                "repository": {
                                    "pullRequest": {
                                        "reviewThreads": {
                                            "nodes": %v
                                        }
                                    }
                                }
                            }}`,
							test.reviewThreadsResponse,
						),
					)
				}
			},
			aladino.MockBuiltIns(),
			nil,
		)

		t.Run(name, func(t *testing.T) {
			gotVal, err := hasUnaddressedThreads(mockedEnv, []aladino.Value{})

			assert.Nil(t, err)
			assert.Equal(t, test.wantVal, gotVal)
		})
	}
}
