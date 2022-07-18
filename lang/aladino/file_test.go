// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	"github.com/stretchr/testify/assert"
)

// mocks_aladino.MockDefaultEnv calls aladino.NewFile upon creating the new eval env
func TestNewFile_WhenErrorInFilePatch(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	mockedFile := &github.CommitFile{
		Patch:    github.String("@@"),
		Filename: github.String(fileName),
	}
	mockedPullRequestFileList := &[]*github.CommitFile{mockedFile}
	_, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequestFileList))
				}),
			),
		},
		nil,
	)

	assert.EqualError(t, err, fmt.Sprintf("error in file patch %s: error in chunk lines parsing (1): missing lines info: @@\npatch: @@", fileName))
}

func TestNewFile(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	mockedFile := &github.CommitFile{
		Patch:    github.String("@@ -2,9 +2,11 @@ package main\n- func previous() {\n+ func new() {\n+\nreturn"),
		Filename: github.String(fileName),
	}
	mockedPullRequestFileList := &[]*github.CommitFile{mockedFile}
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequestFileList))
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantFile := &aladino.File{
		Repr: mockedFile,
	}
	wantFile.AppendToDiff(false, 2, 2, 2, 3, " func previous() {", " func new() {\n")

	gotFile := mockedEnv.GetPatch()[fileName]

	assert.Equal(t, wantFile, gotFile)
}

func TestQuery_WhenCompileFails(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	mockedCommitFile := &github.CommitFile{
		Patch:    github.String("@@ -2,9 +2,11 @@ package main\n- func previous() {\n+ func new() {\n+\nreturn"),
		Filename: github.String(fileName),
	}
	mockedFile := &aladino.File{
		Repr: mockedCommitFile,
	}
	mockedFile.AppendToDiff(false, 2, 2, 2, 3, " func previous() {", " func new() {\n")

	gotVal, err := mockedFile.Query("previous(")

	assert.False(t, gotVal)
	assert.EqualError(t, err, "query: compile error error parsing regexp: missing closing ): `previous(`")
}

func TestQuery_WhenFound(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	mockedCommitFile := &github.CommitFile{
		Patch:    github.String("@@ -2,9 +2,11 @@ package main\n- func previous() {\n+ func new() {\n+\nreturn"),
		Filename: github.String(fileName),
	}
	mockedFile := &aladino.File{
		Repr: mockedCommitFile,
	}
	mockedFile.AppendToDiff(false, 2, 2, 2, 3, " func previous() {", " func new() {\n")

	gotVal, err := mockedFile.Query("new()")

	assert.Nil(t, err)
	assert.True(t, gotVal)
}

func TestQuery_WhenNotFound(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	mockedCommitFile := &github.CommitFile{
		Patch:    github.String("@@ -2,9 +2,11 @@ package main\n- func previous() {\n+ func new() {\n+\nreturn"),
		Filename: github.String(fileName),
	}
	mockedFile := &aladino.File{
		Repr: mockedCommitFile,
	}
	mockedFile.AppendToDiff(false, 2, 2, 2, 3, " func previous() {", " func new() {\n")

	gotVal, err := mockedFile.Query("previous()")

	assert.Nil(t, err)
	assert.False(t, gotVal)
}
