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

func TestFromGithubMilestones(t *testing.T) {
	tests := map[string]struct {
		githubMilestones []*github.Milestone
		milestones       defs.Milestones
	}{
		"when single milestone": {
			githubMilestones: []*github.Milestone{
				{
					NodeID: github.String("1"),
					Title:  github.String("title"),
				},
			},
			milestones: []defs.Milestone{
				{
					ID:    "1",
					Title: "title",
				},
			},
		},
		"when multiple milestones": {
			githubMilestones: []*github.Milestone{
				{
					NodeID: github.String("1"),
					Title:  github.String("title"),
				},
				{
					NodeID: github.String("2"),
					Title:  github.String("title"),
				},
			},
			milestones: []defs.Milestone{
				{
					ID:    "1",
					Title: "title",
				},
				{
					ID:    "2",
					Title: "title",
				},
			},
		},
		"when no milestones": {
			githubMilestones: []*github.Milestone{},
			milestones:       []defs.Milestone{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			milestones := defs.FromGithubMilestones(test.githubMilestones)
			assert.Equal(t, test.milestones, milestones)
		})
	}
}
