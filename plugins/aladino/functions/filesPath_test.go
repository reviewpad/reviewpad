// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	pbe "github.com/reviewpad/api/go/entities"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var filesPath = plugins_aladino.PluginBuiltIns().Functions["filesPath"].Code

func TestFilesPath(t *testing.T) {
	tests := map[string]struct {
		files      []*pbe.CommitFile
		wantResult aladino.Value
		wantErr    error
	}{
		"when successful": {
			files: []*pbe.CommitFile{
				{
					Filename: "go.mod",
				},
			},
			wantResult: aladino.BuildArrayValue([]aladino.Value{
				aladino.BuildStringValue("go.mod"),
			}),
		},
		"when successful with nil file": {
			files: []*pbe.CommitFile{
				{
					Filename: "go.mod",
				},
				nil,
			},
			wantResult: aladino.BuildArrayValue([]aladino.Value{
				aladino.BuildStringValue("go.mod"),
			}),
		},
		"when successful with empty file name": {
			files: []*pbe.CommitFile{
				{
					Filename: "go.sum",
				},
				{
					Filename: "",
				},
			},
			wantResult: aladino.BuildArrayValue([]aladino.Value{
				aladino.BuildStringValue("go.sum"),
			}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnvWithCodeReview(
				t,
				nil,
				nil,
				aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
					Files: test.files,
				}),
				aladino.MockBuiltIns(),
				nil,
			)

			res, err := filesPath(mockedEnv, []aladino.Value{})

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantResult, res)
		})
	}
}
