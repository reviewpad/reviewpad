// Copyright 2023 zola
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
