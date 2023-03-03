// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	contents := `
api-version: reviewpad.com/v3.x

mode: silent
edition: professional
metrics-on-merge: true

pipelines:
  - name: assign-reviewers
    description: Reviewer assignment for small pull requests
    trigger: '$size() <= 30'
    stages:
      - actions:
        - '$assignReviewer(["marcelosousa"])'
        until: '$reviewerStatus("marcelosousa") == "APPROVED"'
      - actions:
        - '$assignReviewer(["ferreiratiago"])'
`
	gotFile, err := parse([]byte(contents))

	metricsOnMerge := true
	wantFile := &ReviewpadFile{
		Mode:           "silent",
		MetricsOnMerge: &metricsOnMerge,
		Pipelines: []PadPipeline{
			{
				Name:        "assign-reviewers",
				Description: "Reviewer assignment for small pull requests",
				Trigger:     "$size() <= 30",
				Stages: []PadStage{
					{
						Actions: []string{
							"$assignReviewer([\"marcelosousa\"])",
						},
						Until: "$reviewerStatus(\"marcelosousa\") == \"APPROVED\"",
					},
					{
						Actions: []string{
							"$assignReviewer([\"ferreiratiago\"])",
						},
					},
				},
			},
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, wantFile, gotFile)
}
