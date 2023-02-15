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

func AssignCodeAuthorReviewer() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{
			aladino.BuildIntType(),
			aladino.BuildArrayOfType(aladino.BuildStringType()),
			aladino.BuildIntType(),
		}, nil),
		Code:           assignCodeAuthorReviewer,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func assignCodeAuthorReviewer(env aladino.Env, args []aladino.Value) error {
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
	filteredReviewers := []github.GitBlameAuthorRank{}
	for _, reviewerRank := range reviewersRanks {
		if strings.HasSuffix(reviewerRank.Username, "[bot]") {
			continue
		}

		if pr.GetUser().GetLogin() == reviewerRank.Username {
			continue
		}

		numberOfOpenReviews, err := githubClient.GetOpenReviewsCountByUser(ctx, targetEntity.Owner, targetEntity.Repo, reviewerRank.Username)
		if err != nil {
			return fmt.Errorf("error getting number of open reviews for user %s: %w", reviewerRank.Username, err)
		}

		if maxAssignedReviews > 0 && numberOfOpenReviews > maxAssignedReviews {
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

	reviewersToRequest := []string{}
	for i := 0; i < totalRequiredReviewers; i++ {
		reviewersToRequest = append(reviewersToRequest, filteredReviewers[i].Username)
	}

	return target.RequestReviewers(reviewersToRequest)
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
