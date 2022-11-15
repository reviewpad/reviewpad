// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils_test

import (
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

func TestIsPullRequestReadyForReportMetrics(t *testing.T) {
	tests := map[string]struct {
		targetEntity *handler.TargetEntity
		wantVal      bool
	}{
		"when target entity is nil": {
			wantVal: false,
		},
		"when issue target entity": {
			targetEntity: &handler.TargetEntity{
				Kind: handler.Issue,
			},
			wantVal: false,
		},
		"when event name is not pull request": {
			targetEntity: &handler.TargetEntity{
				Kind:      handler.PullRequest,
				EventName: "pull_request_review",
			},
			wantVal: false,
		},
		"when event action is not closed": {
			targetEntity: &handler.TargetEntity{
				Kind:        handler.PullRequest,
				EventName:   "pull_request",
				EventAction: "opened",
			},
			wantVal: false,
		},
		"when event name is pull request and event action is closed": {
			targetEntity: &handler.TargetEntity{
				Kind:        handler.PullRequest,
				EventName:   "pull_request",
				EventAction: "closed",
			},
			wantVal: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			val := utils.IsPullRequestReadyForReportMetrics(test.targetEntity)
			assert.Equal(t, test.wantVal, val)
		})
	}
}

func TestIsReviewPadCommand(t *testing.T) {
	tests := map[string]struct {
		targetEntity *handler.TargetEntity
		wantVal      bool
	}{
		"when target entity is nil": {
			wantVal: false,
		},
		"when event name is pull request review": {
			wantVal: false,
			targetEntity: &handler.TargetEntity{
				EventName: "pull_request_review",
			},
		},
		"when comment body is nil": {
			wantVal: false,
			targetEntity: &handler.TargetEntity{
				Comment: &github.IssueComment{
					Body: nil,
				},
			},
		},
		"when comment body doesn't have /reviewpad prefix": {
			wantVal: false,
			targetEntity: &handler.TargetEntity{
				Comment: &github.IssueComment{
					Body: github.String("some comment"),
				},
			},
		},
		"when event name is issue comment and body has /reviewpad prefix": {
			wantVal: true,
			targetEntity: &handler.TargetEntity{
				EventName: "issue_comment",
				Comment: &github.IssueComment{
					Body: github.String("/reviewpad"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			val := utils.IsReviewPadCommand(test.targetEntity)
			assert.Equal(t, test.wantVal, val)
		})
	}
}
