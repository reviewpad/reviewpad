// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils

import (
	"strings"

	"github.com/reviewpad/reviewpad/v3/handler"
)

func IsPullRequestReadyForReportMetrics(targetEntity *handler.TargetEntity) bool {
	return targetEntity != nil && targetEntity.Kind == handler.PullRequest && targetEntity.EventName == "pull_request" && targetEntity.EventAction == "closed"
}

func IsReviewPadCommand(targetEntity *handler.TargetEntity) bool {
	return targetEntity != nil &&
		targetEntity.EventName == "issue_comment" &&
		targetEntity.Comment.Body != nil &&
		strings.HasPrefix(*targetEntity.Comment.Body, "/reviewpad")
}
