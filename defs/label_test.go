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

func TestFromGithubLabels(t *testing.T) {
	tests := map[string]struct {
		githubLabels []*github.Label
		labels       defs.Labels
	}{
		"when single label": {
			githubLabels: []*github.Label{
				{
					NodeID:      github.String("1"),
					Name:        github.String("name"),
					Description: github.String("description"),
				},
			},
			labels: []defs.Label{
				{
					ID:          "1",
					Name:        "name",
					Description: "description",
				},
			},
		},
		"when multiple labels": {
			githubLabels: []*github.Label{
				{
					ID:          github.Int64(1),
					NodeID:      github.String("1"),
					Name:        github.String("name"),
					Description: github.String("description"),
				},
				{
					ID:          github.Int64(2),
					NodeID:      github.String("2"),
					Name:        github.String("name"),
					Description: github.String("description"),
				},
			},
			labels: []defs.Label{
				{
					ID:          "1",
					Name:        "name",
					Description: "description",
				},
				{
					ID:          "2",
					Name:        "name",
					Description: "description",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.labels, defs.FromGithubLabels(test.githubLabels))
		})
	}
}
