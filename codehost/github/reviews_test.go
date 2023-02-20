// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetReviewsCountByUserFromOpenPullRequests(t *testing.T) {
	mockedReviewsCountQuery := fmt.Sprintf(`{
		"query":"query($afterCursor:String$author:String!$perPage:Int!$query:String!){
			search(query: $query, type: ISSUE, first: $perPage, after: $afterCursor){
				pageInfo{
					endCursor,
					hasNextPage
				},
				nodes{
					... on PullRequest{
						reviews(author: $author){
							totalCount
						}
					}
				}
			}
		}",
		"variables":{
			"afterCursor": null,
			"author": "test",
			"perPage": 100,
			"query": "user:%s is:pr is:open reviewed-by:test"
		}
	}`, aladino.DefaultMockPrOwner)

	tests := map[string]struct {
		ghGraphQLHandler func(http.ResponseWriter, *http.Request)
		wantCount        int
		wantErr          string
	}{
		"when the request for the reviews count fails": {
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				if query == utils.MinifyQuery(mockedReviewsCountQuery) {
					http.Error(w, "GetReviewsCountFail", http.StatusNotFound)
				}
			},
			wantErr: "non-200 OK status code: 404 Not Found body: \"GetReviewsCountFail\\n\"",
		},
		"when the user has created reviews in open pull requests": {
			ghGraphQLHandler: func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(utils.MustRead(req.Body))
				switch query {
				case utils.MinifyQuery(mockedReviewsCountQuery):
					utils.MustWrite(w, `{
						"data": {
							"search": {
								"pageInfo": {
									"endCursor": "Y3Vyc29yOjI=",
									"hasNextPage": false
								},
								"nodes": [
									{
										"reviews": {
											"totalCount": 5
										}
									}
								]
							}
						}
					}`)
				}
			},
			wantCount: 5,
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

			gotCount, gotErr := mockedEnv.GetGithubClient().GetReviewsCountByUserFromOpenPullRequests(mockedEnv.GetCtx(), aladino.DefaultMockPrOwner, "test")

			if test.wantErr != "" {
				assert.EqualError(t, gotErr, test.wantErr)
				assert.Equal(t, 0, gotCount)
			} else {
				assert.Equal(t, test.wantCount, gotCount)
				assert.Nil(t, gotErr)
			}
		})
	}
}
