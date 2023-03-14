// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package defs_test

import (
	"testing"

	"github.com/google/go-github/v49/github"
	"github.com/reviewpad/reviewpad/v4/defs"
	"github.com/stretchr/testify/assert"
)

func TestFromGitHubCommitFiles(t *testing.T) {
	tests := map[string]struct {
		githubCommitFiles []*github.CommitFile
		wantCommitFiles   defs.CommitFiles
	}{
		"when single file": {
			githubCommitFiles: []*github.CommitFile{
				{
					SHA:              github.String("sha"),
					Filename:         github.String("filename"),
					Additions:        github.Int(1),
					Deletions:        github.Int(1),
					Changes:          github.Int(1),
					Status:           github.String("status"),
					Patch:            github.String("patch"),
					PreviousFilename: github.String("previous_filename"),
				},
			},
			wantCommitFiles: defs.CommitFiles{
				{
					SHA:              "sha",
					Filename:         "filename",
					Additions:        1,
					Deletions:        1,
					Changes:          1,
					Status:           "status",
					Patch:            "patch",
					PreviousFilename: "previous_filename",
				},
			},
		},
		"when multiple files": {
			githubCommitFiles: []*github.CommitFile{
				{
					SHA:              github.String("sha"),
					Filename:         github.String("filename"),
					Additions:        github.Int(1),
					Deletions:        github.Int(1),
					Changes:          github.Int(1),
					Status:           github.String("status"),
					Patch:            github.String("patch"),
					PreviousFilename: github.String("previous_filename"),
				},
				{
					SHA:              github.String("sha"),
					Filename:         github.String("filename"),
					Additions:        github.Int(1),
					Deletions:        github.Int(1),
					Changes:          github.Int(1),
					Status:           github.String("status"),
					Patch:            github.String("patch"),
					PreviousFilename: github.String("previous_filename"),
				},
			},
			wantCommitFiles: defs.CommitFiles{
				{
					SHA:              "sha",
					Filename:         "filename",
					Additions:        1,
					Deletions:        1,
					Changes:          1,
					Status:           "status",
					Patch:            "patch",
					PreviousFilename: "previous_filename",
				},
				{
					SHA:              "sha",
					Filename:         "filename",
					Additions:        1,
					Deletions:        1,
					Changes:          1,
					Status:           "status",
					Patch:            "patch",
					PreviousFilename: "previous_filename",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			commitFiles := defs.FromGitHubCommitFiles(test.githubCommitFiles)
			assert.Equal(t, test.wantCommitFiles, commitFiles)
		})
	}
}
