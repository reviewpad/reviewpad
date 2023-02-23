// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils

import (
	"strings"

	"github.com/google/go-github/v49/github"
	"github.com/reviewpad/reviewpad/v3/handler"
)

func IsPullRequestReadyForReportMetrics(eventDetails *handler.EventDetails) bool {
	return eventDetails != nil && eventDetails.EventName == "pull_request" && eventDetails.EventAction == "closed"
}

func IsReviewpadCommand(eventDetails *handler.EventDetails) bool {
	return eventDetails != nil &&
		eventDetails.EventName == "issue_comment" &&
		eventDetails.EventAction == "created" &&
		strings.HasPrefix(eventDetails.Payload.(*github.IssueCommentEvent).GetComment().GetBody(), "/reviewpad")
}

func IsReviewpadCommandDryRun(command string) bool {
	return strings.TrimPrefix(command, "/reviewpad ") == "dry-run"
}

func IsReviewpadCommandRun(command string) bool {
	return strings.TrimPrefix(command, "/reviewpad ") == "run"
}
