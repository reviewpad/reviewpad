// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var isUpdatedWithBaseBranch = plugins_aladino.PluginBuiltIns().Functions["isUpdatedWithBaseBranch"].Code

func TestIsUpdatedWithBaseBranch(t *testing.T) {
	tests := map[string]struct {
		graphQLHandler http.HandlerFunc
		wantErr        error
		wantIsUpdated  aladino.Value
	}{
		"when get pull request update to date query fails": {
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnprocessableEntity)
			},
			wantErr: errors.New("error getting pull request outdated information: non-200 OK status code: 422 Unprocessable Entity body: \"\""),
		},
		"when head is behind": {
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"baseRefOid": "oldBaseRefOid",
								"baseRef": {
									"target": {
										"history": {
											"nodes": [
												{
													"oid": "baseRefOid"
												}
											]
										}
									}
								}
							}
						}
					}
				}`)
			},
			wantIsUpdated: aladino.BuildBoolValue(false),
		},
		"when head is updated": {
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"baseRefOid": "baseRefOid",
								"baseRef": {
									"target": {
										"history": {
											"nodes": [
												{
													"oid": "baseRefOid"
												}
											]
										}
									}
								}
							}
						}
					}
				}`)
			},
			wantIsUpdated: aladino.BuildBoolValue(true),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, nil, test.graphQLHandler, nil, nil)
			isUpdated, err := isUpdatedWithBaseBranch(env, nil)

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantIsUpdated, isUpdated)
		})
	}
}
