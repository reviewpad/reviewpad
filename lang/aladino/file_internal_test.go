// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"fmt"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/stretchr/testify/assert"
)

const testPatch = `@@ -2,9 +2,11 @@ package main
- func previous() {
+ func new() {
+
return
}`

func TestAppendToDiff(t *testing.T) {
    fileName := "default-mock-repo/file1.ts"
	mockedFile := &github.CommitFile{
		Patch:    github.String(testPatch),
		Filename: github.String(fileName),
	}

    isContext := false
    oldStart := 2
    oldEnd := 2
    newStart := 2
    newEnd := 3
    oldLine := " func previous() {"
    newLine := " func new() {\n"

	file := &File{
		Repr: mockedFile,
	}
	file.AppendToDiff(
        isContext,
        oldStart,
        oldEnd,
        newStart,
        newEnd,
        oldLine,
        newLine,
    )

    gotDiff := file.Diff

    wantDiff := []*diffBlock{
        {
            isContext: isContext,
            Old: &diffSpan{
                int32(oldStart),
                int32(oldEnd),
            },
            New: &diffSpan{
                int32(newStart),
                int32(newEnd),
            },
            oldLine: oldLine,
            newLine: newLine,
        },
	}

	assert.Equal(t, wantDiff, gotDiff)
}

func TestNewFile_WhenErrorInFilePatch(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	mockedFile := &github.CommitFile{
		Patch:    github.String("@@"),
		Filename: github.String(fileName),
	}

	gotFile, err := NewFile(mockedFile)

	assert.Nil(t, gotFile)
	assert.EqualError(t, err, fmt.Sprintf("error in file patch %s: error in chunk lines parsing (1): missing lines info: @@\npatch: @@", fileName))
}

func TestNewFile(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	mockedFile := &github.CommitFile{
		Patch:    github.String(testPatch),
		Filename: github.String(fileName),
	}

	wantFile := &File{
		Repr: mockedFile,
	}
	wantFile.AppendToDiff(false, 2, 2, 2, 3, " func previous() {", " func new() {\n")

	gotFile, err := NewFile(mockedFile)

	assert.Nil(t, err)
	assert.Equal(t, wantFile, gotFile)
}

func TestQuery_WhenCompileFails(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	mockedFile := &File{
		Repr: &github.CommitFile{
			Patch:    github.String(testPatch),
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
	mockedFile := &File{
		Repr: &github.CommitFile{
			Patch:    github.String(testPatch),
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
	mockedFile := &File{
		Repr: &github.CommitFile{
			Patch:    github.String(testPatch),
			Filename: github.String(fileName),
		},
	}
	mockedFile.AppendToDiff(false, 2, 2, 2, 3, " func previous() {", " func new() {\n")

	gotVal, err := mockedFile.Query("previous()")

	assert.Nil(t, err)
	assert.False(t, gotVal)
}
