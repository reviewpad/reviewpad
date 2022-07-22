// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"fmt"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/stretchr/testify/assert"
)

// aladino.MockDefaultEnv calls aladino.NewFile upon creating the new eval env
func TestNewFile_WhenErrorInFilePatch(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	mockedFile := &github.CommitFile{
		Patch:    github.String("@@"),
		Filename: github.String(fileName),
	}

	gotFile, err := aladino.NewFile(mockedFile)

	assert.Nil(t, gotFile)
	assert.EqualError(t, err, fmt.Sprintf("error in file patch %s: error in chunk lines parsing (1): missing lines info: @@\npatch: @@", fileName))
}

func TestNewFile(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	mockedFile := &github.CommitFile{
		Patch:    github.String("@@ -2,9 +2,11 @@ package main\n- func previous() {\n+ func new() {\n+\nreturn"),
		Filename: github.String(fileName),
	}

	wantFile := &aladino.File{
		Repr: mockedFile,
	}
	wantFile.AppendToDiff(false, 2, 2, 2, 3, " func previous() {", " func new() {\n")

	gotFile, err := aladino.NewFile(mockedFile)

	assert.Nil(t, err)
	assert.Equal(t, wantFile, gotFile)
}

func TestQuery_WhenCompileFails(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	mockedFile := &aladino.File{
		Repr: &github.CommitFile{
			Patch:    github.String("@@ -2,9 +2,11 @@ package main\n- func previous() {\n+ func new() {\n+\nreturn"),
			Filename: github.String(fileName),
		},
	}
	mockedFile.AppendToDiff(false, 2, 2, 2, 3, " func previous() {", " func new() {\n")

	gotVal, err := mockedFile.Query("previous(")

	assert.False(t, gotVal)
	assert.EqualError(t, err, "query: compile error error parsing regexp: missing closing ): `previous(`")
}

func TestQuery_WhenFound(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	mockedFile := &aladino.File{
		Repr: &github.CommitFile{
			Patch:    github.String("@@ -2,9 +2,11 @@ package main\n- func previous() {\n+ func new() {\n+\nreturn"),
			Filename: github.String(fileName),
		},
	}
	mockedFile.AppendToDiff(false, 2, 2, 2, 3, " func previous() {", " func new() {\n")

	gotVal, err := mockedFile.Query("new()")

	assert.Nil(t, err)
	assert.True(t, gotVal)
}

func TestQuery_WhenNotFound(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	mockedFile := &aladino.File{
		Repr: &github.CommitFile{
			Patch:    github.String("@@ -2,9 +2,11 @@ package main\n- func previous() {\n+ func new() {\n+\nreturn"),
			Filename: github.String(fileName),
		},
	}
	mockedFile.AppendToDiff(false, 2, 2, 2, 3, " func previous() {", " func new() {\n")

	gotVal, err := mockedFile.Query("previous()")

	assert.Nil(t, err)
	assert.False(t, gotVal)
}
