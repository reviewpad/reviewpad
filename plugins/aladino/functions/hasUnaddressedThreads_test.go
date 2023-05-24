// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"net/http"
	"testing"

	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
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

	_, err := hasUnaddressedThreads(mockedEnv, []lang.Value{})

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
		wantVal               lang.Value
	}{
		"when none": {
			reviewThreadsResponse: `[]`,
			wantVal:               lang.BuildFalseValue(),
		},
		"when unresolved": {
			reviewThreadsResponse: `[{"isResolved":false, "isOutdated":false}]`,
			wantVal:               lang.BuildTrueValue(),
		},
		"when resolved": {
			reviewThreadsResponse: `[{"isResolved":true, "isOutdated":false}]`,
			wantVal:               lang.BuildFalseValue(),
		},
		"when outdated": {
			reviewThreadsResponse: `[{"isResolved":false, "isOutdated":true}]`,
			wantVal:               lang.BuildFalseValue(),
		},
		"when up to date": {
			reviewThreadsResponse: `[{"isResolved":false, "isOutdated":false}]`,
			wantVal:               lang.BuildTrueValue(),
		},
		"when more than one thread": {
			reviewThreadsResponse: `[{"isResolved":true, "isOutdated":true}, {"isResolved":true, "isOutdated":true}]`,
			wantVal:               lang.BuildFalseValue(),
		},
	}

	for name, test := range tests {
		mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
			t,
			nil,
			func(w http.ResponseWriter, req *http.Request) {
				query := utils.MustRead(req.Body)
				switch query {
				case mockedGraphQLQuery:
					utils.MustWrite(
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
			aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
				Number: aladino.DefaultMockPrNum,
				Base: &pbc.Branch{
					Repo: &pbc.Repository{
						Owner: aladino.DefaultMockPrOwner,
						Name:  aladino.DefaultMockPrRepoName,
					},
				},
			}),
			aladino.GetDefaultPullRequestFileList(),
			aladino.MockBuiltIns(),
			nil,
		)

		t.Run(name, func(t *testing.T) {
			gotVal, err := hasUnaddressedThreads(mockedEnv, []lang.Value{})

			assert.Nil(t, err)
			assert.Equal(t, test.wantVal, gotVal)
		})
	}
}
