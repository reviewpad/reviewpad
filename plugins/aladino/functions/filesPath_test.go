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
	mockedPullRequestFileList := &[]*github.CommitFile{
		{
			Filename: github.String("go.mod"),
		},
		{
			Filename: github.String("go.sum"),
		},
		{
			Filename: github.String("cmd/main.go"),
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequestFileList))
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	wantFilesPath := aladino.BuildArrayValue([]aladino.Value{
		aladino.BuildStringValue("go.mod"),
		aladino.BuildStringValue("go.sum"),
		aladino.BuildStringValue("cmd/main.go"),
	})

	gotFilesPath, err := filesPath(mockedEnv, []aladino.Value{})

	assert.Nil(t, err)
	assert.Equal(t, wantFilesPath, gotFilesPath)
}
