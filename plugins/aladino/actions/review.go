// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
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

	log := e.GetLogger().WithField("builtin", "review")

	if t.PullRequest.IsDraft {
		log.Infof("skipping review because the pull request is in draft")
		return nil
	}

	if t.PullRequest.Status == pbc.PullRequestStatus_CLOSED {
		log.Infof("skipping review because the pull request is closed")
		return nil
	}

	reviewEvent, err := parseReviewEvent(args[0].(*aladino.StringValue).Val)
	if err != nil {
		return err
	}

	reviewBody, err := parseReviewBody(reviewEvent, args[1].(*aladino.StringValue).Val)
	if err != nil {
		return err
	}

	authenticatedUserLogin, err := e.GetGithubClient().GetAuthenticatedUserLogin()
	if err != nil {
		return err
	}

	latestReview, err := t.GetLatestReviewFromReviewer(authenticatedUserLogin)
	if err != nil {
		return err
	}

	lastPushDate, err := t.GetPullRequestLastPushDate()
	if err != nil {
		return err
	}

	if latestReview != nil {
		latestReviewEvent, err := mapReviewStateToEvent(latestReview.State)
		if err != nil {
			return err
		}

		// If the latest review is the same as the one we want to create, and the last push date is before the latest review
		// then we skip the review creation.
		if latestReviewEvent == reviewEvent && latestReview.Body == reviewBody && lastPushDate.Before(*latestReview.SubmittedAt) {
			log.Infof("skipping review because there were no updates since the last review")
			return nil
		}

		log.Infof("latest review from %v is %v with body %v", authenticatedUserLogin, latestReviewEvent, latestReview.Body)
	}

	log.Infof("creating review %v with body %v", reviewEvent, reviewBody)

	return t.Review(reviewEvent, reviewBody)
}

func parseReviewEvent(reviewEvent string) (string, error) {
	switch reviewEvent {
	case "COMMENT", "REQUEST_CHANGES", "APPROVE":
		return reviewEvent, nil
	default:
		return "", fmt.Errorf("review: unsupported review state %v", reviewEvent)
	}
}

func parseReviewBody(reviewEvent, reviewBody string) (string, error) {
	if reviewEvent != "APPROVE" && reviewBody == "" {
		return "", fmt.Errorf("review: comment required in %v state", reviewEvent)
	}

	return reviewBody, nil
}

func mapReviewStateToEvent(reviewState string) (string, error) {
	switch reviewState {
	case "COMMENTED":
		return "COMMENT", nil
	case "CHANGES_REQUESTED":
		return "REQUEST_CHANGES", nil
	case "APPROVED":
		return "APPROVE", nil
	case "DISMISSED":
		return "DISMISSED", nil
	default:
		return "", fmt.Errorf("review: unsupported review state %v", reviewState)
	}
}
