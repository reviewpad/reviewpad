// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package defs

import "github.com/google/go-github/v49/github"

type Label struct {
	ID          string
	Name        string
	Description string
}

type Labels []Label

func FromGithubLabels(githubLabels []*github.Label) Labels {
	labels := make(Labels, len(githubLabels))
	for i, githubLabel := range githubLabels {
		labels[i] = FromGithubLabel(githubLabel)
	}
	return labels
}

func FromGithubLabel(githubLabel *github.Label) Label {
	return Label{
		ID:          githubLabel.GetNodeID(),
		Name:        githubLabel.GetName(),
		Description: githubLabel.GetDescription(),
	}
}
