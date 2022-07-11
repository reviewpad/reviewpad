// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"log"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v2/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var hasCodePattern = plugins_aladino.PluginBuiltIns().Functions["hasCodePattern"].Code

func TestHasCodePattern_WhenPullRequestPatchHasNilFile(t *testing.T) {
	mockedPullRequestFileList := &[]*github.CommitFile{}
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatchHandler(
			// Overwrite default mock request to get pull request changed files
			mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write(mock.MustMarshal(mockedPullRequestFileList))
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedEnv.GetPatch()["file1.ts"] = nil

	args := []aladino.Value{aladino.BuildStringValue("placeBet\\(.*\\)")}
	gotVal, err := hasCodePattern(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(false)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestHasCodePattern_WhenFileHasQueryError(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv()
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue("a(")}
	gotVal, err := hasCodePattern(mockedEnv, args)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, "query: compile error error parsing regexp: missing closing ): `a(`")
}

func TestHasCodePattern(t *testing.T) {
    mockedPullRequestFileList := &[]*github.CommitFile{}
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
        mock.WithRequestMatchHandler(
			// Overwrite default mock request to get pull request changed files
			mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write(mock.MustMarshal(mockedPullRequestFileList))
			}),
		),
    )
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

    fileName := "default-mock-repo/file1.ts"
	ghFile := &github.CommitFile{
		SHA:      github.String("1234"),
		Patch:    github.String("@@ -2,9 +2,11 @@ package main\n- func previous() {\n+ func new() {\n+\nreturn"),
		Filename: github.String(fileName),
	}
	file, _ := aladino.NewFile(ghFile)

	mockedEnv.GetPatch()[fileName] = file

	args := []aladino.Value{aladino.BuildStringValue("new\\(.*\\)")}
	gotVal, err := hasCodePattern(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(true)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
