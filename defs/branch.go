// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package defs

import "github.com/google/go-github/v49/github"

type Branch struct {
	Label string
	Ref   string
	Sha   string
	Repo  Repository
	User  User
}

type Branches []Branch

func FromGithubBranches(githubBranches []*github.PullRequestBranch) Branches {
	branches := make(Branches, len(githubBranches))
	for i, githubBranch := range githubBranches {
		branches[i] = FromGithubBranch(githubBranch)
	}
	return branches
}

func FromGithubBranch(githubBranch *github.PullRequestBranch) Branch {
	return Branch{
		Label: githubBranch.GetLabel(),
		Ref:   githubBranch.GetRef(),
		Sha:   githubBranch.GetSHA(),
		Repo:  FromGithubRepository(githubBranch.GetRepo()),
		User:  FromGithubUser(githubBranch.GetUser()),
	}
}
