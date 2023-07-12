// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/google/go-github/v52/github"
	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func IsWaitingForReview() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{}, lang.BuildBoolType()),
		Code:           isWaitingForReviewCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest},
	}
}

func isWaitingForReviewCode(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget).PullRequest
	requestedUsers := pullRequest.RequestedReviewers.Users
	requestedTeams := pullRequest.RequestedReviewers.Teams
	log := e.GetLogger().WithField("builtin", "isWaitingForReview")

	if len(requestedUsers) > 0 || len(requestedTeams) > 0 {
		return lang.BuildBoolValue(true), nil
	}

	prNum := pullRequest.GetNumber()
	owner := pullRequest.GetBase().GetRepo().GetOwner()
	repo := pullRequest.GetBase().GetRepo().GetName()
	author := pullRequest.GetAuthor().GetLogin()

	commits, err := e.GetGithubClient().GetPullRequestCommits(e.GetCtx(), owner, repo, int(prNum))
	if err != nil {
		return nil, err
	}

	if len(commits) == 0 {
		log.Warnf("no commits found for pull request %s/%s#%d", owner, repo, prNum)
		return lang.BuildBoolValue(false), nil
	}

	lastPushedDate, err := e.GetGithubClient().GetPullRequestLastPushDate(e.GetCtx(), owner, repo, int(prNum))
	if err != nil {
		return nil, err
	}

	reviews, err := e.GetGithubClient().GetPullRequestReviews(e.GetCtx(), owner, repo, int(prNum))
	if err != nil {
		return nil, err
	}

	lastReviewByUser := make(map[string]*github.PullRequestReview)

	for _, review := range reviews {
		userLogin := review.User.GetLogin()
		if userLogin == "" || userLogin == author {
			continue
		}

		lastUserReview, ok := lastReviewByUser[userLogin]
		if ok {
			if review.GetSubmittedAt().After(lastUserReview.GetSubmittedAt().Time) {
				lastReviewByUser[userLogin] = review
			}
		} else {
			lastReviewByUser[userLogin] = review
		}
	}

	for _, lastUserReview := range lastReviewByUser {
		if *lastUserReview.State != "APPROVED" {
			if lastUserReview.GetSubmittedAt().Before(lastPushedDate) {
				return lang.BuildBoolValue(true), nil
			}
		}
	}

	return lang.BuildBoolValue(false), nil
}
