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

var approvalsCount = plugins_aladino.PluginBuiltIns().Functions["approvalsCount"].Code

func TestApprovalsCount(t *testing.T) {
	tests := map[string]struct {
		graphQLHandler http.HandlerFunc
		wantRes        aladino.Value
		wantErr        error
	}{
		"when graphql query errors": {
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				utils.MustWrite(w, `{"message": "internal error"}`)
			},
			wantErr: errors.New("non-200 OK status code: 500 Internal Server Error body: \"{\\\"message\\\": \\\"internal error\\\"}\""),
		},
		"when successful with total count 0": {
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `
				{
					"data": {
						"repository": {
							"pullRequest": {
								"reviews": {
									"totalCount": 0
								}
							}
						}
					}
				}
				`)
			},
			wantRes: aladino.BuildIntValue(0),
		},
		"when successful": {
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `
				{
					"data": {
						"repository": {
							"pullRequest": {
								"reviews": {
									"totalCount": 3
								}
							}
						}
					}
				}
				`)
			},
			wantRes: aladino.BuildIntValue(3),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, nil, test.graphQLHandler, nil, nil)

			res, err := approvalsCount(env, []aladino.Value{})

			assert.Equal(t, test.wantRes, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
