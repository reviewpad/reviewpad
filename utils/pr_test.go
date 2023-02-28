// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils_test

import (
	"testing"

	"github.com/google/go-github/v49/github"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

func TestIsPullRequestReadyForReportMetrics(t *testing.T) {
	tests := map[string]struct {
		eventDetails *handler.EventDetails
		wantVal      bool
	}{
		"when event data is nil": {
			wantVal: false,
		},
		"when event name is not pull request": {
			eventDetails: &handler.EventDetails{
				EventName: "pull_request_review",
			},
			wantVal: false,
		},
		"when event action is not closed": {
			eventDetails: &handler.EventDetails{
				EventName:   "pull_request",
				EventAction: "opened",
			},
			wantVal: false,
		},
		"when event name is pull request and event action is closed": {
			eventDetails: &handler.EventDetails{
				EventName:   "pull_request",
				EventAction: "closed",
			},
			wantVal: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			val := utils.IsPullRequestReadyForReportMetrics(test.eventDetails)
			assert.Equal(t, test.wantVal, val)
		})
	}
}

func TestIsReviewPadCommand(t *testing.T) {
	tests := map[string]struct {
		eventDetails *handler.EventDetails
		wantVal      bool
	}{
		"when target entity is nil": {
			wantVal: false,
		},
		"when event name is pull request review": {
			wantVal: false,
			eventDetails: &handler.EventDetails{
				EventName: "pull_request_review",
			},
		},
		"when comment body is nil": {
			wantVal: false,
			eventDetails: &handler.EventDetails{
				EventName: "issue_comment",
				Payload: &github.IssueCommentEvent{
					Comment: &github.IssueComment{
						Body: nil,
					},
				},
			},
		},
		"when comment body doesn't have /reviewpad prefix": {
			wantVal: false,
			eventDetails: &handler.EventDetails{
				EventName: "issue_comment",
				Payload: &github.IssueCommentEvent{
					Comment: &github.IssueComment{
						Body: github.String("some comment"),
					},
				},
			},
		},
		"when event action is not created": {
			wantVal: false,
			eventDetails: &handler.EventDetails{
				EventName:   "issue_comment",
				EventAction: "updated",
				Payload: &github.IssueCommentEvent{
					Comment: &github.IssueComment{
						Body: github.String("some comment"),
					},
				},
			},
		},
		"when event name is issue comment and body has /reviewpad prefix": {
			wantVal: true,
			eventDetails: &handler.EventDetails{
				EventName:   "issue_comment",
				EventAction: "created",
				Payload: &github.IssueCommentEvent{
					Comment: &github.IssueComment{
						Body: github.String("/reviewpad"),
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			val := utils.IsReviewpadCommand(test.eventDetails)
			assert.Equal(t, test.wantVal, val)
		})
	}
}

func TestIsReviewpadCommandRun(t *testing.T) {
	tests := map[string]struct {
		command string
		wantVal bool
	}{
		"when command is empty": {
			wantVal: false,
		},
		"when command is not /reviewpad run": {
			command: "/reviewpad assign-reviewers",
			wantVal: false,
		},
		"when command is /reviewpad run": {
			command: "/reviewpad run",
			wantVal: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			val := utils.IsReviewpadCommandRun(test.command)
			assert.Equal(t, test.wantVal, val)
		})
	}
}
