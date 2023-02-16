// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var isUpdatedWithBaseBranch = plugins_aladino.PluginBuiltIns().Functions["isUpdatedWithBaseBranch"].Code

func TestIsUpdatedWithBaseBranch(t *testing.T) {
	tests := map[string]struct {
		graphQLHandler http.HandlerFunc
		wantErr        error
		wantIsUpdated  aladino.Value
	}{
		"when get head behind by query fails": {
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnprocessableEntity)
			},
			wantErr: errors.New("error getting head behind by information: non-200 OK status code: 422 Unprocessable Entity body: \"\""),
		},
		"when head is behind": {
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"pullRequest": {
								"baseRef": {
									"compare": {
										"behindBy": 1
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
								"baseRef": {
									"compare": {
										"behindBy": 0
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
