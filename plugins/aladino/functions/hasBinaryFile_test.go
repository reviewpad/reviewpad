// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var hasBinaryFile = plugins_aladino.PluginBuiltIns().Functions["hasBinaryFile"].Code

func TestHasBinaryFile(t *testing.T) {
	tests := map[string]struct {
		wantResult     aladino.Value
		wantErr        error
		graphqlHandler func(http.ResponseWriter, *http.Request)
	}{
		"when graphql query errors": {
			wantResult: (aladino.Value)(nil),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: errors.New(`non-200 OK status code: 500 Internal Server Error body: ""`),
		},
		"when object is not binary file": {
			wantResult: aladino.BuildBoolValue(false),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				aladino.MustWrite(w, `{
					"data": {
						"repository": {
							"object": {
								"isBinary": false
							}
						}
					}
				}`)
			},
		},
		"when object is binary file": {
			wantResult: aladino.BuildBoolValue(true),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				aladino.MustWrite(w, `{
					"data": {
						"repository": {
							"object": {
								"isBinary": true
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

			res, err := hasBinaryFile(env, []aladino.Value{})

			assert.Equal(t, test.wantResult, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
