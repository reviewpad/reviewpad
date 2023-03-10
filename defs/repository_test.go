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

func TestFromGithubRepositories(t *testing.T) {
	tests := map[string]struct {
		githubRepositories []*github.Repository
		repositories       defs.Repositories
	}{
		"when single repository": {
			githubRepositories: []*github.Repository{
				{
					NodeID:   github.String("1"),
					Name:     github.String("name"),
					FullName: github.String("full_name"),
					Owner: &github.User{
						NodeID: github.String("1"),
						Login:  github.String("login"),
					},
				},
			},
			repositories: []defs.Repository{
				{
					ID:       "1",
					Name:     "name",
					FullName: "full_name",
					Owner: defs.User{
						ID:    "1",
						Login: "login",
					},
				},
			},
		},
		"when multiple repositories": {
			githubRepositories: []*github.Repository{
				{
					NodeID:   github.String("1"),
					Name:     github.String("name"),
					FullName: github.String("full_name"),
				},
				{
					NodeID:   github.String("2"),
					Name:     github.String("name"),
					FullName: github.String("full_name"),
				},
			},
			repositories: []defs.Repository{
				{
					ID:       "1",
					Name:     "name",
					FullName: "full_name",
				},
				{
					ID:       "2",
					Name:     "name",
					FullName: "full_name",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repositories := defs.FromGithubRepositories(test.githubRepositories)
			assert.Equal(t, test.repositories, repositories)
		})
	}
}
