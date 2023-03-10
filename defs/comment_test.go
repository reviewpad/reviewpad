// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package defs_test

import (
	"testing"
	"time"

	"github.com/google/go-github/v49/github"
	"github.com/reviewpad/reviewpad/v4/defs"
	"github.com/stretchr/testify/assert"
)

func TestFromGithubComments(t *testing.T) {
	now := time.Now()
	tests := map[string]struct {
		githubComments []*github.IssueComment
		comment        defs.Comments
	}{
		"when single comment": {
			githubComments: []*github.IssueComment{
				{
					ID:        github.Int64(1),
					NodeID:    github.String("1"),
					Body:      github.String("body"),
					CreatedAt: &now,
					UpdatedAt: &now,
				},
			},
			comment: []defs.Comment{
				{
					ID:        "1",
					Number:    1,
					Body:      "body",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		},
		"when multiple comments": {
			githubComments: []*github.IssueComment{
				{
					ID:        github.Int64(1),
					NodeID:    github.String("1"),
					Body:      github.String("body"),
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				{
					ID:        github.Int64(2),
					NodeID:    github.String("2"),
					Body:      github.String("body"),
					CreatedAt: &now,
					UpdatedAt: &now,
				},
			},
			comment: []defs.Comment{
				{
					ID:        "1",
					Number:    1,
					Body:      "body",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "2",
					Number:    2,
					Body:      "body",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			comment := defs.FromGithubComments(test.githubComments)
			assert.Equal(t, test.comment, comment)
		})
	}

}
