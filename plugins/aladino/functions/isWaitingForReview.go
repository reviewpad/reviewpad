// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/google/go-github/v49/github"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func IsWaitingForReview() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code:           isWaitingForReviewCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func isWaitingForReviewCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget).CodeReview
	requestedUsers := pullRequest.RequestedReviewers.Users
	requestedTeams := pullRequest.RequestedReviewers.Teams
	log := e.GetLogger().WithField("builtin", "isWaitingForReview")

	if len(requestedUsers) > 0 || len(requestedTeams) > 0 {
		return aladino.BuildBoolValue(true), nil
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
		return aladino.BuildBoolValue(false), nil
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
			if review.GetSubmittedAt().After(lastUserReview.GetSubmittedAt()) {
				lastReviewByUser[userLogin] = review
			}
		} else {
			lastReviewByUser[userLogin] = review
		}
	}

	for _, lastUserReview := range lastReviewByUser {
		if *lastUserReview.State != "APPROVED" {
			if lastUserReview.GetSubmittedAt().Before(lastPushedDate) {
				return aladino.BuildBoolValue(true), nil
			}
		}
	}

	return aladino.BuildBoolValue(false), nil
}
