// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"log"

	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
)

func IsWaitingForReview() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code: isWaitingForReviewCode,
	}
}

func isWaitingForReviewCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	pullRequest := e.GetPullRequest()
	requestedUsers := pullRequest.RequestedReviewers
	requestedTeams := pullRequest.RequestedTeams

	if len(requestedUsers) > 0 || len(requestedTeams) > 0 {
		return aladino.BuildBoolValue(true), nil
	}

	prNum := utils.GetPullRequestNumber(pullRequest)
	owner := utils.GetPullRequestBaseOwnerName(pullRequest)
	repo := utils.GetPullRequestBaseRepoName(pullRequest)
	author := pullRequest.GetUser().GetLogin()

	commits, err := utils.GetPullRequestCommits(e.GetCtx(), e.GetClient(), owner, repo, prNum)
	if err != nil {
		return nil, err
	}

	if len(commits) == 0 {
		log.Printf("[WARN] No commits found for pull request %s/%s#%d.", owner, repo, prNum)
		return aladino.BuildBoolValue(false), nil
	}

	lastPushedDate, err := utils.GetPullRequestLastPushedEventDate(e.GetCtx(), e.GetClientGQL(), owner, repo, prNum)
	if err != nil {
		return nil, err
	}

	reviews, err := utils.GetPullRequestReviews(e.GetCtx(), e.GetClient(), owner, repo, prNum)
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
