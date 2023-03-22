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

var fileCount = plugins_aladino.PluginBuiltIns().Functions["fileCount"].Code

func TestFileCount(t *testing.T) {
	mockedCodeReview := aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
		Files: []*pbe.CommitFile{
			{
				Filename: "default-mock-repo/file1.ts",
				Patch:    "",
			},
		},
	})

	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		nil,
		nil,
		mockedCodeReview,
		aladino.MockBuiltIns(),
		nil,
	)

	wantFileCount := aladino.BuildIntValue(len(mockedCodeReview.Files))

	args := []aladino.Value{}
	gotFileCount, err := fileCount(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantFileCount, gotFileCount, "action should count the total pull request files")
}
