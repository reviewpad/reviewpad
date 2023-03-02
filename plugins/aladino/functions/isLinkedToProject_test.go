// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var isLinkedToProject = plugins_aladino.PluginBuiltIns().Functions["isLinkedToProject"].Code

func TestIsLinkedToProject(t *testing.T) {
	tests := map[string]struct {
		args           []aladino.Value
		wantResult     aladino.Value
		wantErr        error
		graphqlHandler func(http.ResponseWriter, *http.Request)
	}{
		"when graphql query errors": {
			args:       []aladino.Value{aladino.BuildStringValue("project title")},
			wantResult: (aladino.Value)(nil),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: errors.New(`non-200 OK status code: 500 Internal Server Error body: ""`),
		},
		"when graphql query errors with project items not found": {
			args:       []aladino.Value{aladino.BuildStringValue("project title")},
			wantResult: (aladino.Value)(nil),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository":{
							"pullRequest":{
							}
						}
					}
				}`)
			},
			wantErr: errors.New(`project items not found`),
		},
		"when linked project is false": {
			args:       []aladino.Value{aladino.BuildStringValue("project title")},
			wantResult: aladino.BuildBoolValue(false),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository":{
							"pullRequest":{
								"projectItems":{
									"nodes":[
										{
											"project":{
												"id": "1",
												"number": 1,
												"title": "Project 1"
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
		},
		"when linked project is true": {
			args:       []aladino.Value{aladino.BuildStringValue("project title")},
			wantResult: aladino.BuildBoolValue(true),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository":{
							"pullRequest":{
								"projectItems":{
									"nodes":[
										{
											"project":{
												"id": "1",
												"number": 1,
												"title": "project"
											},
											"project":{
												"id": "2",
												"number": 2,
												"title": "project title"
											}
										}
									]
								}
							}
						}
					}
				}`)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, []mock.MockBackendOption{}, test.graphqlHandler, aladino.MockBuiltIns(), nil)

			res, err := isLinkedToProject(env, test.args)

			assert.Equal(t, test.wantResult, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
