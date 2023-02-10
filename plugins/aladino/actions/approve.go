// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Approve() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code:           approveCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func approveCode(e aladino.Env, args []aladino.Value) error {
	t := e.GetTarget().(*target.PullRequestTarget)
	log := e.GetLogger().WithField("builtin", "approve")

	reviewBody := args[0].(*aladino.StringValue).Val

	if t.PullRequest.GetState() == "closed" {
		log.Infof("skipping approve because the pull request is closed")
		return nil
	}

	if t.PullRequest.GetDraft() {
		log.Infof("skipping approve because the pull request is in draft")
		return nil
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
		// The last push was made before the last review so a new review is not needed
		if lastPushDate.Before(*latestReview.SubmittedAt) {
			log.Infof("skipping approve because there were no updates since the last review")
			return nil
		}

		log.Infof("latest review from %v is %v with body '%v'", authenticatedUserLogin, latestReview.State, latestReview.Body)

		if latestReview.State == "APPROVED" && latestReview.Body == reviewBody {
			log.Infof("skipping approve since it's the same as the latest review")
			return nil
		}
	}
	log.Infof("creating review APPROVE with body '%v'", reviewBody)

	return t.Review("APPROVE", reviewBody)
}
