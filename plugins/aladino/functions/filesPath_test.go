// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var filesPath = plugins_aladino.PluginBuiltIns().Functions["filesPath"].Code

func TestFilesPath(t *testing.T) {
	tests := map[string]struct {
		files      []*github.CommitFile
		wantResult lang.Value
		wantErr    error
	}{
		"when successful": {
			files: []*github.CommitFile{
				{
					Filename: github.String("go.mod"),
				},
			},
			wantResult: lang.BuildArrayValue([]lang.Value{
				lang.BuildStringValue("go.mod"),
			}),
		},
		"when successful with nil file": {
			files: []*github.CommitFile{
				{
					Filename: github.String("go.mod"),
				},
				nil,
			},
			wantResult: lang.BuildArrayValue([]lang.Value{
				lang.BuildStringValue("go.mod"),
			}),
		},
		"when successful with empty file name": {
			files: []*github.CommitFile{
				{
					Filename: github.String("go.sum"),
				},
				{
					Filename: github.String(""),
				},
			},
			wantResult: lang.BuildArrayValue([]lang.Value{
				lang.BuildStringValue("go.sum"),
			}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
				t,
				nil,
				nil,
				aladino.GetDefaultPullRequestDetails(),
				test.files,
				aladino.MockBuiltIns(),
				nil,
			)

			res, err := filesPath(mockedEnv, []lang.Value{})

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantResult, res)
		})
	}
}
