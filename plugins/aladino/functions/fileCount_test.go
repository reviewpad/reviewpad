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

var fileCount = plugins_aladino.PluginBuiltIns().Functions["fileCount"].Code

func TestFileCount(t *testing.T) {
	mockedFiles := []*github.CommitFile{
		{
			Filename: github.String("default-mock-repo/file1.ts"),
			Patch:    github.String(""),
		},
	}

	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		nil,
		nil,
		aladino.GetDefaultPullRequestDetails(),
		mockedFiles,
		aladino.MockBuiltIns(),
		nil,
	)

	wantFileCount := lang.BuildIntValue(len(mockedFiles))

	args := []lang.Value{}
	gotFileCount, err := fileCount(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantFileCount, gotFileCount, "action should count the total pull request files")
}
