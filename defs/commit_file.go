// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package defs

import "github.com/google/go-github/v49/github"

type CommitFile struct {
	SHA              string
	Filename         string
	Additions        int
	Deletions        int
	Changes          int
	Status           string
	Patch            string
	PreviousFilename string
}

type CommitFiles []CommitFile

func FromGitHubCommitFiles(files []*github.CommitFile) CommitFiles {
	commitFiles := make(CommitFiles, len(files))
	for i, file := range files {
		commitFiles[i] = FromGitHubCommitFile(file)
	}
	return commitFiles
}

func FromGitHubCommitFile(file *github.CommitFile) CommitFile {
	return CommitFile{
		SHA:              file.GetSHA(),
		Filename:         file.GetFilename(),
		Additions:        file.GetAdditions(),
		Deletions:        file.GetDeletions(),
		Changes:          file.GetChanges(),
		Status:           file.GetStatus(),
		Patch:            file.GetPatch(),
		PreviousFilename: file.GetPreviousFilename(),
	}
}
