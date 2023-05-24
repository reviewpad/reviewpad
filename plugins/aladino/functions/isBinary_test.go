// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var isBinary = plugins_aladino.PluginBuiltIns().Functions["isBinary"].Code

func TestIsBinary(t *testing.T) {
	tests := map[string]struct {
		filename       lang.Value
		wantResult     lang.Value
		wantErr        error
		graphqlHandler func(http.ResponseWriter, *http.Request)
	}{
		"when graphql query errors": {
			wantResult: (lang.Value)(nil),
			filename:   lang.BuildStringValue("file.txt"),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: errors.New(`non-200 OK status code: 500 Internal Server Error body: ""`),
		},
		"when file is not a binary file": {
			filename:   lang.BuildStringValue("file.txt"),
			wantResult: lang.BuildBoolValue(false),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
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
		"when file is a binary file": {
			filename:   lang.BuildStringValue("binary.exe"),
			wantResult: lang.BuildBoolValue(true),
			graphqlHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
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
			env := aladino.MockDefaultEnv(t, nil, test.graphqlHandler, nil, nil)

			res, err := isBinary(env, []lang.Value{test.filename})

			assert.Equal(t, test.wantResult, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
