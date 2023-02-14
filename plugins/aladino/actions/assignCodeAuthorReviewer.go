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
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildIntType(), aladino.BuildArrayOfType(aladino.BuildStringType()), aladino.BuildIntType()}, nil),
		Code:           assignCodeAuthorReviewer,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func assignCodeAuthorReviewer(e aladino.Env, args []aladino.Value) error {
	t := e.GetTarget().(*target.PullRequestTarget)
	total := args[0].(*aladino.IntValue).Val
	excludeReviewers := args[1].(*aladino.ArrayValue).Vals
	maxReviews := args[2].(*aladino.IntValue).Val
	pr := t.PullRequest
	ctx := e.GetCtx()
	targetEntity := t.GetTargetEntity()
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	baseSHA := pr.GetBase().GetSHA()
	githubClient := e.GetGithubClient()

	reviewers, err := t.GetReviewers()
	if err != nil {
		return fmt.Errorf("error getting reviewers: %s", err.Error())
	}

	if len(reviewers.Users) > 0 {
		return nil
	}

	filePaths := []string{}
	patch := t.Patch
	for _, fp := range patch {
		filePaths = append(filePaths, fp.Repr.GetFilename())
	}

	blame, err := githubClient.GetGitBlame(ctx, targetEntity.Owner, targetEntity.Repo, baseSHA, filePaths)
	if err != nil {
		return fmt.Errorf("error getting git blame information: %s", err.Error())
	}

	availableAssignees, err := t.GetAvailableAssignees()
	if err != nil {
		return fmt.Errorf("error getting available assignees: %s", err.Error())
	}

	reviewersRanks := githubClient.ComputeGitBlameRank(blame)
	filteredReviewers := []github.GitBlameAuthorRank{}
	for _, reviewerRank := range reviewersRanks {
		if !strings.HasSuffix(reviewerRank.Username, "[bot]") && pr.GetUser().GetLogin() != reviewerRank.Username {
			numberOfOpenReviews, err := githubClient.GetOpenReviewsCountByUser(ctx, owner, repo, reviewerRank.Username)
			if err != nil {
				return fmt.Errorf("error getting number of open reviews for user %s: %w", reviewerRank.Username, err)
			}

			if maxReviews > 0 && numberOfOpenReviews > maxReviews {
				continue
			}

			if !isExcluded(excludeReviewers, reviewerRank.Username) && isAvailableAssignee(availableAssignees, reviewerRank.Username) {
				filteredReviewers = append(filteredReviewers, reviewerRank)
			}
		}
	}

	if len(filteredReviewers) == 0 {
		return assignRandomReviewerCode(e, nil)
	}

	if total > len(filteredReviewers) {
		e.GetLogger().Warnf("number of required reviewers(%d) is less than available code author reviewers(%d)", total, len(filteredReviewers))
		total = len(filteredReviewers)
	}

	reviewersToRequest := []string{}
	for i := 0; i < total; i++ {
		reviewersToRequest = append(reviewersToRequest, filteredReviewers[i].Username)
	}

	return t.RequestReviewers(reviewersToRequest)
}

func isExcluded(excludedReviewers []aladino.Value, username string) bool {
	for _, excludedReviewer := range excludedReviewers {
		if excludedReviewer.(*aladino.StringValue).Val == username {
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
