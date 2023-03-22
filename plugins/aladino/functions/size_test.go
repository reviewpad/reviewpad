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

var size = plugins_aladino.PluginBuiltIns().Functions["size"].Code

func TestSize_WhenRegexMatchFails(t *testing.T) {
	tests := map[string]struct {
		args    []aladino.Value
		wantVal aladino.Value
	}{
		"no exclusions": {
			args:    []aladino.Value{aladino.BuildArrayValue([]aladino.Value{})},
			wantVal: aladino.BuildIntValue(6),
		},
		"exclude yml": {
			args:    []aladino.Value{aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("*.yml")})},
			wantVal: aladino.BuildIntValue(3),
		},
		"invalid regex": {
			args:    []aladino.Value{aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("*[.yml")})},
			wantVal: aladino.BuildIntValue(6),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			changes := int64(3)
			mockedCodeReviewFileList := []*pbe.CommitFile{
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
			mockedCodeReview := aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
				AdditionsCount: int64(additions),
				DeletionsCount: int64(deletions),
				Files:          mockedCodeReviewFileList,
			})

			mockedEnv := aladino.MockDefaultEnvWithCodeReview(
				t,
				nil,
				nil,
				mockedCodeReview,
				aladino.MockBuiltIns(),
				nil,
			)

			gotSize, err := size(mockedEnv, test.args)

			assert.Nil(t, err)
			assert.Equal(t, test.wantVal, gotSize)
		})
	}
}
