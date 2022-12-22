// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var filesPath = plugins_aladino.PluginBuiltIns().Functions["filesPath"].Code

func TestFilesPath(t *testing.T) {
	tests := map[string]struct {
		files      *[]*github.CommitFile
		wantResult aladino.Value
		wantErr    error
	}{
		"when successful": {
			files: &[]*github.CommitFile{
				{
					Filename: github.String("go.mod"),
					Patch:    nil,
				},
				{
					Filename: github.String("go.sum"),
					Patch:    nil,
				},
				{
					Filename: github.String("cmd/main.go"),
					Patch:    nil,
				},
			},
			wantResult: aladino.BuildArrayValue([]aladino.Value{
				aladino.BuildStringValue("go.mod"),
				aladino.BuildStringValue("go.sum"),
				aladino.BuildStringValue("cmd/main.go"),
			}),
		},
		"when successful with nil file": {
			files: &[]*github.CommitFile{
				{
					Filename: github.String("go.mod"),
					Patch:    nil,
				},
				nil,
			},
			wantResult: aladino.BuildArrayValue([]aladino.Value{
				aladino.BuildStringValue("go.mod"),
			}),
		},
		"when successful with nil file name": {
			files: &[]*github.CommitFile{
				{
					Filename: github.String("go.sum"),
					Patch:    nil,
				},
				{
					Filename: nil,
				},
			},
			wantResult: aladino.BuildArrayValue([]aladino.Value{
				aladino.BuildStringValue("go.sum"),
			}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							utils.MustWriteBytes(w, mock.MustMarshal(test.files))
						}),
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
			)

			res, err := filesPath(mockedEnv, []aladino.Value{})

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantResult, res)
		})
	}
}
