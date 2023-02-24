// Copyright 2023 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var removeFromProjectCode = plugins_aladino.PluginBuiltIns().Actions["removeFromProject"].Code

func TestRemoveFromProject(t *testing.T) {
	getProjectByNameGraphQLQuery := `{
		"query":"query($name:String! $repositoryName:String! $repositoryOwner:String!) {
			repository(owner: $repositoryOwner, name: $repositoryName) {
				projectsV2(query: $name, first: 50, orderBy: {field: TITLE, direction: ASC}) {
					nodes{
						id,
						number,
						title
					}
				}
			}
		}",
		"variables":{
		   "name":"test",
		   "repositoryName":"default-mock-repo",
		   "repositoryOwner":"foobar"
		}
	 }`

	getProjectItemIDGraphQLQuery := `{
		"query":"query($after:String! $name:String! $number:Int! $owner:String!){
			repository(owner: $owner, name: $name){
				pullRequest(number: $number){
					projectItems(first: 100, after: $after){
						pageInfo{
							hasNextPage,
							endCursor
						},
						nodes{
							project{
								id
							},
							id
						}
					}
				}
			}
		}",
		"variables":{
		   "after":"",
		   "name":"default-mock-repo",
		   "number":6,
		   "owner":"foobar"
		}
	 }`

	deleteProjectItemGraphQLQuery := `{
		"query":"mutation($input:DeleteProjectV2ItemInput!){
			deleteProjectV2Item(input: $input){
				clientMutationId
			}
		}",
		"variables":{
		   "input":{
			  "projectId":"test",
			  "itemId":"test"
		   }
		}
	 }
	 `

	tests := map[string]struct {
		graphQLHandler func(w http.ResponseWriter, r *http.Request)
		wantErr        error
		projectName    string
	}{
		"when get project by name fails": {
			projectName: "test",
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				graphQLQuery := utils.MustRead(r.Body)
				if utils.MinifyQuery(graphQLQuery) == utils.MinifyQuery(getProjectByNameGraphQLQuery) {
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
			wantErr: errors.New("failed to get project by name: non-200 OK status code: 500 Internal Server Error body: \"\""),
		},
		"when the project is not found": {
			projectName: "test",
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				graphQLQuery := utils.MustRead(r.Body)
				if utils.MinifyQuery(graphQLQuery) == utils.MinifyQuery(getProjectByNameGraphQLQuery) {
					w.WriteHeader(http.StatusOK)
					utils.MustWrite(w, `{
						"data":{
							"repository": {
								"projectsV2":{
									"nodes":[]
								}
							}
						}
					}`)
				}
			},
			wantErr: errors.New("failed to get project by name: project not found"),
		},
		"when get project item id fails": {
			projectName: "test",
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				graphQLQuery := utils.MustRead(r.Body)
				if utils.MinifyQuery(graphQLQuery) == utils.MinifyQuery(getProjectByNameGraphQLQuery) {
					w.WriteHeader(http.StatusOK)
					utils.MustWrite(w, `{
						"data":{
							"repository": {
								"projectsV2":{
									"nodes":[
										{
											"id":"test",
											"number":1,
											"title":"test"
										}
									]
								}
							}
						}
					}`)
				}

				if utils.MinifyQuery(graphQLQuery) == utils.MinifyQuery(getProjectItemIDGraphQLQuery) {
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
			wantErr: errors.New("failed to get project item id: non-200 OK status code: 500 Internal Server Error body: \"\""),
		},
		"when the project item is not found": {
			projectName: "test",
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				graphQLQuery := utils.MustRead(r.Body)
				if utils.MinifyQuery(graphQLQuery) == utils.MinifyQuery(getProjectByNameGraphQLQuery) {
					w.WriteHeader(http.StatusOK)
					utils.MustWrite(w, `{
						"data":{
							"repository": {
								"projectsV2":{
									"nodes":[
										{
											"id":"test",
											"number":1,
											"title":"test"
										}
									]
								}
							}
						}
					}`)
				}

				if utils.MinifyQuery(graphQLQuery) == utils.MinifyQuery(getProjectItemIDGraphQLQuery) {
					w.WriteHeader(http.StatusOK)
					utils.MustWrite(w, `{
						"data":{
							"repository": {
								"pullRequest":{
									"projectItems":{
										"pageInfo":{
											"hasNextPage":false,
											"endCursor":null
										},
										"nodes":[]
									}
								}
							}
						}
					}`)
				}
			},
		},
		"when remove from project fails": {
			projectName: "test",
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				graphQLQuery := utils.MustRead(r.Body)
				if utils.MinifyQuery(graphQLQuery) == utils.MinifyQuery(getProjectByNameGraphQLQuery) {
					w.WriteHeader(http.StatusOK)
					utils.MustWrite(w, `{
						"data":{
							"repository": {
								"projectsV2":{
									"nodes":[
										{
											"id":"test",
											"number":1,
											"title":"test"
										}
									]
								}
							}
						}
					}`)
				}

				if utils.MinifyQuery(graphQLQuery) == utils.MinifyQuery(getProjectItemIDGraphQLQuery) {
					w.WriteHeader(http.StatusOK)
					utils.MustWrite(w, `{
						"data":{
							"repository": {
								"pullRequest":{
									"projectItems":{
										"pageInfo":{
											"hasNextPage":false,
											"endCursor":null
										},
										"nodes":[
											{
												"project":{
													"id":"test"
												},
												"id":"test"
											}
										]
									}
								}
							}
						}
					}`)
				}

				if utils.MinifyQuery(graphQLQuery) == utils.MinifyQuery(deleteProjectItemGraphQLQuery) {
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
			wantErr: errors.New("failed to remove from project: non-200 OK status code: 500 Internal Server Error body: \"\""),
		},
		"when remove from project succeeds": {
			projectName: "test",
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				graphQLQuery := utils.MustRead(r.Body)
				if utils.MinifyQuery(graphQLQuery) == utils.MinifyQuery(getProjectByNameGraphQLQuery) {
					w.WriteHeader(http.StatusOK)
					utils.MustWrite(w, `{
						"data":{
							"repository": {
								"projectsV2":{
									"nodes":[
										{
											"id":"test",
											"number":1,
											"title":"test"
										}
									]
								}
							}
						}
					}`)
				}

				if utils.MinifyQuery(graphQLQuery) == utils.MinifyQuery(getProjectItemIDGraphQLQuery) {
					w.WriteHeader(http.StatusOK)
					utils.MustWrite(w, `{
						"data":{
							"repository": {
								"pullRequest":{
									"projectItems":{
										"pageInfo":{
											"hasNextPage":false,
											"endCursor":null
										},
										"nodes":[
											{
												"project":{
													"id":"test"
												},
												"id":"test"
											}
										]
									}
								}
							}
						}
					}`)
				}

				if utils.MinifyQuery(graphQLQuery) == utils.MinifyQuery(deleteProjectItemGraphQLQuery) {
					w.WriteHeader(http.StatusOK)
					utils.MustWrite(w, `{
						"data":{
							"deleteProjectV2Item":{
								"clientMutationId":"test"
							}
						}
					}`)
				}
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, nil, test.graphQLHandler, nil, nil)

			gotErr := removeFromProjectCode(env, []aladino.Value{aladino.BuildStringValue(test.projectName)})

			assert.Equal(t, test.wantErr, gotErr)
		})
	}
}
