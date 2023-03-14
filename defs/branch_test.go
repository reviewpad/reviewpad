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

func TestFromGithubBranch(t *testing.T) {
	tests := map[string]struct {
		githubBranch *github.PullRequestBranch
		branch       defs.Branch
	}{
		"when single branch": {
			githubBranch: &github.PullRequestBranch{
				Label: github.String("label"),
				Ref:   github.String("ref"),
				SHA:   github.String("sha"),
				Repo: &github.Repository{
					ID:       github.Int64(1),
					NodeID:   github.String("1"),
					Name:     github.String("name"),
					FullName: github.String("login/name"),
					Owner: &github.User{
						NodeID: github.String("1"),
						Login:  github.String("login"),
					},
				},
				User: &github.User{
					NodeID: github.String("1"),
					Login:  github.String("login"),
				},
			},
			branch: defs.Branch{
				Label: "label",
				Ref:   "ref",
				SHA:   "sha",
				Repo: defs.Repository{
					ID:       "1",
					Name:     "name",
					FullName: "login/name",
					Owner: defs.User{
						ID:    "1",
						Login: "login",
					},
				},
				User: defs.User{
					ID:    "1",
					Login: "login",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			branch := defs.FromGithubBranch(test.githubBranch)
			assert.Equal(t, test.branch, branch)
		})
	}
}
