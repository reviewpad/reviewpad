// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var fileCount = plugins_aladino.PluginBuiltIns().Functions["fileCount"].Code

func TestFileCount(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedPullRequestFileList := &[]*github.CommitFile{
		{
			Filename: github.String("default-mock-repo/file1.ts"),
			Patch:    nil,
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequestFileList))
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

	wantFileCount := aladino.BuildIntValue(len(*mockedPullRequestFileList))

	args := []aladino.Value{}
	gotFileCount, err := fileCount(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantFileCount, gotFileCount, "action should count the total pull request files")
}
