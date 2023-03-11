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

func TestFromGithubTeams(t *testing.T) {
	tests := map[string]struct {
		githubTeams []*github.Team
		teams       defs.Teams
	}{
		"when single team": {
			githubTeams: []*github.Team{
				{
					ID:     github.Int64(1),
					NodeID: github.String("1"),
					Name:   github.String("name"),
				},
			},
			teams: []defs.Team{
				{
					ID:   "1",
					Name: "name",
				},
			},
		},
		"when multiple teams": {
			githubTeams: []*github.Team{
				{
					ID:     github.Int64(1),
					NodeID: github.String("1"),
					Name:   github.String("name"),
				},
				{
					ID:     github.Int64(2),
					NodeID: github.String("2"),
					Name:   github.String("name"),
				},
			},
			teams: []defs.Team{
				{
					ID:   "1",
					Name: "name",
				},
				{
					ID:   "2",
					Name: "name",
				},
			},
		},
		"when no teams": {
			githubTeams: []*github.Team{},
			teams:       []defs.Team{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.teams, defs.FromGithubTeams(test.githubTeams))
		})
	}
}
