// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package defs

import (
	"github.com/google/go-github/v49/github"
)

type Milestone struct {
	ID    string
	Title string
}

type Milestones []Milestone

func FromGithubMilestone(githubMilestone *github.Milestone) *Milestone {
	if githubMilestone == nil {
		return nil
	}

	return &Milestone{
		ID:    githubMilestone.GetNodeID(),
		Title: githubMilestone.GetTitle(),
	}
}

func FromGithubMilestones(githubMilestones []*github.Milestone) Milestones {
	milestones := make(Milestones, len(githubMilestones))

	for i, githubMilestone := range githubMilestones {
		milestones[i] = *FromGithubMilestone(githubMilestone)
	}

	return milestones
}
