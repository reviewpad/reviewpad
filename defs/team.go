// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package defs

import "github.com/google/go-github/v49/github"

type Team struct {
	ID   string
	Name string
	Slug string
}

type Teams []Team

func FromGithubTeam(githubTeam *github.Team) Team {
	return Team{
		ID:   githubTeam.GetNodeID(),
		Name: githubTeam.GetName(),
		Slug: githubTeam.GetSlug(),
	}
}

func FromGithubTeams(githubTeams []*github.Team) Teams {
	teams := make(Teams, len(githubTeams))
	for i, githubTeam := range githubTeams {
		teams[i] = FromGithubTeam(githubTeam)
	}
	return teams
}
