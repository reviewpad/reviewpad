// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package defs

import "github.com/google/go-github/v49/github"

type Repository struct {
	ID       string
	Name     string
	FullName string
	Owner    User
	IsFork   bool
	URL      string
}

type Repositories []Repository

func FromGithubRepository(githubRepository *github.Repository) Repository {
	return Repository{
		ID:       githubRepository.GetNodeID(),
		Name:     githubRepository.GetName(),
		FullName: githubRepository.GetFullName(),
		IsFork:   githubRepository.GetFork(),
		Owner:    FromGithubUser(githubRepository.GetOwner()),
		URL:      githubRepository.GetURL(),
	}
}

func FromGithubRepositories(githubRepositories []*github.Repository) Repositories {
	repositories := make(Repositories, len(githubRepositories))
	for i, githubRepository := range githubRepositories {
		repositories[i] = FromGithubRepository(githubRepository)
	}
	return repositories
}
