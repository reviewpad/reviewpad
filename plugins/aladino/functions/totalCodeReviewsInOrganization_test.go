// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var totalCodeReviewsInOrganization = plugins_aladino.PluginBuiltIns().Functions["totalCodeReviewsInOrganization"].Code

func TestTotalCodeReviewsInOrganization_WhenQueryFails(t *testing.T) {
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

	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedReviewsCountQuery):
				http.Error(w, "GetReviewsCountFail", http.StatusNotFound)
			}
		},
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue("test")}
	gotResult, gotErr := totalCodeReviewsInOrganization(mockedEnv, args)

	assert.Nil(t, gotResult)
	assert.EqualError(t, gotErr, "non-200 OK status code: 404 Not Found body: \"GetReviewsCountFail\\n\"")
}

func TestTotalCodeReviewsInOrganization(t *testing.T) {
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
		inputUsername               string
		mockedReviewsCountQueryData string
		clientOptions               []mock.MockBackendOption
		wantResult                  aladino.Value
		wantErr                     error
	}{
		"when the user has created reviews in several repositories open PRs of an organization/user": {
			inputUsername: "test",
			mockedReviewsCountQueryData: `{
				"data": {
					"search": {
						"pageInfo": {
							"endCursor": "Y3Vyc29yOjI=",
							"hasNextPage": false
						},
						"nodes": [
							{
								"reviews": {
									"totalCount": 1
								}
							},
							{
								"reviews": {
									"totalCount": 2
								}
							}
						]
					}
				}
			}`,
			wantResult: aladino.BuildIntValue(3),
		},
		"when the user has created reviews in only one repository open PRs of an organization/user": {
			inputUsername: "test",
			mockedReviewsCountQueryData: `{
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
			}`,
			wantResult: aladino.BuildIntValue(5),
		},
		"when the user has not created reviews in repositories open PRs of an organization/user": {
			inputUsername: "test",
			mockedReviewsCountQueryData: `{
				"data": {
					"search": {
						"pageInfo": {
							"endCursor": "Y3Vyc29yOjI=",
							"hasNextPage": false
						},
						"nodes": []
					}
				}
			}`,
			wantResult: aladino.BuildIntValue(0),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(
				t,
				test.clientOptions,
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					switch query {
					case utils.MinifyQuery(mockedReviewsCountQuery):
						utils.MustWrite(w, test.mockedReviewsCountQueryData)
					}
				},
				aladino.MockBuiltIns(),
				nil,
			)

			args := []aladino.Value{aladino.BuildStringValue(test.inputUsername)}
			gotResult, gotErr := totalCodeReviewsInOrganization(mockedEnv, args)

			assert.Equal(t, test.wantErr, gotErr)
			assert.Equal(t, test.wantResult, gotResult)
		})
	}
}
