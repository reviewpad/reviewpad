// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/reviewpad/reviewpad/v3/codehost"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Review() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildStringType()}, nil),
		Code:           reviewCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func reviewCode(e aladino.Env, args []aladino.Value) error {
	t := e.GetTarget().(*target.PullRequestTarget)
	reviewpadBot := "reviewpad-bot"

	reviewEvent, err := parseReviewEvent(args[0].(*aladino.StringValue).Val)
	if err != nil {
		return err
	}

	reviewBody := args[1].(*aladino.StringValue).Val

	// Only a review with APPROVE does not need a review comment
	if reviewEvent != "APPROVE" && reviewBody == "" {
		return fmt.Errorf("review: comment required in %v event", reviewEvent)
	}

	reviews, err := t.GetReviews()
	if err != nil {
		return err
	}

	// If the last review of reviewpad-bot was an approval or if there were no updates to the pull request
	// since last review made by reviewpad-bot, then do nothing.
	if codehost.HasReview(reviews, reviewpadBot) {
		lastReview := codehost.LastReview(reviews, reviewpadBot)

		if lastReview.State == "APPROVED" {
			return nil
		}

		if !lastReview.SubmittedAt.Before(t.PullRequest.GetUpdatedAt()) {
			return nil
		}
	}

	return t.Review(reviewEvent, reviewBody)
}

func parseReviewEvent(reviewEvent string) (string, error) {
	switch reviewEvent {
	case "COMMENT", "REQUEST_CHANGES", "APPROVE":
		return reviewEvent, nil
	default:
		return "", fmt.Errorf("review: unsupported review event %v", reviewEvent)
	}
}
