// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"
	"strings"

	"github.com/reviewpad/reviewpad/v3/codehost"
	"github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func AssignCodeAuthorReviewers() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{
			aladino.BuildIntType(),
			aladino.BuildArrayOfType(aladino.BuildStringType()),
			aladino.BuildIntType(),
		}, nil),
		Code:           assignCodeAuthorReviewersCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func assignCodeAuthorReviewersCode(env aladino.Env, args []aladino.Value) error {
	target := env.GetTarget().(*target.PullRequestTarget)
	totalRequiredReviewers := args[0].(*aladino.IntValue).Val
	reviewersToExclude := args[1].(*aladino.ArrayValue).Vals
	maxAssignedReviews := args[2].(*aladino.IntValue).Val
	pr := target.PullRequest
	ctx := env.GetCtx()
	targetEntity := target.GetTargetEntity()
	githubClient := env.GetGithubClient()

	reviewers, err := target.GetReviewers()
	if err != nil {
		return fmt.Errorf("error getting reviewers: %s", err.Error())
	}

	if len(reviewers.Users) > 0 {
		return nil
	}

	changedFilesPath := []string{}
	for _, patch := range target.Patch {
		changedFilesPath = append(changedFilesPath, patch.Repr.GetFilename())
	}

	gitBlame, err := githubClient.GetGitBlame(ctx, targetEntity.Owner, targetEntity.Repo, pr.GetBase().GetSHA(), changedFilesPath)
	if err != nil {
		return fmt.Errorf("error getting git blame information: %s", err.Error())
	}

	availableAssignees, err := target.GetAvailableAssignees()
	if err != nil {
		return fmt.Errorf("error getting available assignees: %s", err.Error())
	}

	reviewersRanks := githubClient.ComputeGitBlameRank(gitBlame)
	reviewersUsernames := getUsernamesFromAuthorRank(reviewersRanks)
	totalOpenPRsAsReviewerByUser, err := githubClient.GetOpenPullRequestsAsReviewer(ctx, targetEntity.Owner, targetEntity.Repo, reviewersUsernames)
	if err != nil {
		return fmt.Errorf("error getting open pull request as reviewer: %s", err.Error())
	}

	filteredReviewers := []github.GitBlameAuthorRank{}
	for _, reviewerRank := range reviewersRanks {
		if strings.HasSuffix(reviewerRank.Username, "[bot]") {
			continue
		}

		if pr.GetUser().GetLogin() == reviewerRank.Username {
			continue
		}

		totalOpenPRsAsReviewer := totalOpenPRsAsReviewerByUser[reviewerRank.Username]

		if maxAssignedReviews > 0 && totalOpenPRsAsReviewer > maxAssignedReviews {
			continue
		}

		if !isValueInList(reviewersToExclude, reviewerRank.Username) && isAvailableAssignee(availableAssignees, reviewerRank.Username) {
			filteredReviewers = append(filteredReviewers, reviewerRank)
		}
	}

	if len(filteredReviewers) == 0 {
		return assignRandomReviewerCode(env, nil)
	}

	if totalRequiredReviewers > len(filteredReviewers) {
		env.GetLogger().Warnf("number of required reviewers(%d) is less than available code author reviewers(%d)", totalRequiredReviewers, len(filteredReviewers))
		totalRequiredReviewers = len(filteredReviewers)
	}

	reviewersToRequest := getUsernamesFromAuthorRank(filteredReviewers[:totalRequiredReviewers])

	return target.RequestReviewers(reviewersToRequest)
}

func getUsernamesFromAuthorRank(reviewersRanks []github.GitBlameAuthorRank) []string {
	usernames := []string{}
	for _, reviewerRank := range reviewersRanks {
		usernames = append(usernames, reviewerRank.Username)
	}

	return usernames
}

func isValueInList(excludedReviewers []aladino.Value, value string) bool {
	for _, excludedReviewer := range excludedReviewers {
		if excludedReviewer.(*aladino.StringValue).Val == value {
			return true
		}
	}

	return false
}

func isAvailableAssignee(availableAssignees []*codehost.User, username string) bool {
	for _, availableAssignee := range availableAssignees {
		if availableAssignee.Login == username {
			return true
		}
	}

	return false
}
