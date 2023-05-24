// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var size = plugins_aladino.PluginBuiltIns().Functions["size"].Code

func TestSize_WhenRegexMatchFails(t *testing.T) {
	tests := map[string]struct {
		args    []lang.Value
		wantVal lang.Value
	}{
		"no exclusions": {
			args:    []lang.Value{lang.BuildArrayValue([]lang.Value{})},
			wantVal: lang.BuildIntValue(6),
		},
		"exclude yml": {
			args:    []lang.Value{lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("*.yml")})},
			wantVal: lang.BuildIntValue(3),
		},
		"invalid regex": {
			args:    []lang.Value{lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("*[.yml")})},
			wantVal: lang.BuildIntValue(6),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			changes := int64(3)
			mockedCodeReviewFileList := []*pbc.File{
				{
					Filename:     "reviewpad.yml",
					ChangesCount: changes,
				},
				{
					Filename:     "main.go",
					ChangesCount: changes,
				},
			}
			additions := 6
			deletions := 0
			mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
				AdditionsCount: int64(additions),
				DeletionsCount: int64(deletions),
			})

			mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
				t,
				nil,
				nil,
				mockedCodeReview,
				mockedCodeReviewFileList,
				aladino.MockBuiltIns(),
				nil,
			)

			gotSize, err := size(mockedEnv, test.args)

			assert.Nil(t, err)
			assert.Equal(t, test.wantVal, gotSize)
		})
	}
}
