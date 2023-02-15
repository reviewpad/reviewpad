// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"context"
	"fmt"
	"strings"

	gh "github.com/google/go-github/v49/github"
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
	totalRequiredReviewers := args[0].(*aladino.IntValue).Val
	reviewersToExclude := args[1].(*aladino.ArrayValue).Vals
	maxAllowedAssignedReviews := args[2].(*aladino.IntValue).Val

	gitHubClient := env.GetGithubClient()
	pr := env.GetTarget().(*target.PullRequestTarget)

	reviewers, err := pr.GetReviewers()
	if err != nil {
		return fmt.Errorf("error getting reviewers: %s", err)
	}

	// When there are already reviewers assigned
	// don't want to assign more.
	if len(reviewers.Users) > 0 {
		return nil
	}

	// Fetch all available assignees to which the pull request may be assigned to
	// in order to avoid assigning a reviewer that is not available.
	availableAssignees, err := pr.GetAvailableAssignees()
	if err != nil {
		return fmt.Errorf("error getting available assignees: %s", err)
	}

	// Fetch all users that have authored the changed files in the pull request.
	authors, err := getAuthorsFromGitBlame(env.GetCtx(), gitHubClient, pr)
	if err != nil {
		return fmt.Errorf("error getting authors from git blame: %s", err)
	}

	// Fetch the total number of open pull requests that each author has already assigned to them as a reviewer.
	totalOpenPRsAsReviewerByUser, err := gitHubClient.GetOpenPullRequestsAsReviewer(env.GetCtx(), pr.GetTargetEntity().Owner, pr.GetTargetEntity().Repo, authors)
	if err != nil {
		return fmt.Errorf("error getting open pull requests as reviewer: %s", err)
	}

	selectedReviewers := []string{}
	for _, author := range authors {
		if isUserEligibleToReview(author, pr.PullRequest, reviewersToExclude, availableAssignees, totalOpenPRsAsReviewerByUser[author], maxAllowedAssignedReviews) {
			selectedReviewers = append(selectedReviewers, author)
		}
	}

	if len(selectedReviewers) == 0 {
		return assignRandomReviewerCode(env, nil)
	}

	if totalRequiredReviewers > len(selectedReviewers) {
		env.GetLogger().Warnf("number of required reviewers %d is less than available code author reviewers %d", totalRequiredReviewers, len(selectedReviewers))

		totalRequiredReviewers = len(selectedReviewers)
	}

	selectedReviewers = selectedReviewers[:totalRequiredReviewers]

	return pr.RequestReviewers(selectedReviewers)
}

func getAuthorsFromGitBlame(ctx context.Context, gitHubClient *github.GithubClient, pullRequest *target.PullRequestTarget) ([]string, error) {
	authors := []string{}

	changedFilesPath := []string{}

	for _, patch := range pullRequest.Patch {
		changedFilesPath = append(changedFilesPath, patch.Repr.GetFilename())
	}

	gitBlame, err := gitHubClient.GetGitBlame(ctx, pullRequest.GetTargetEntity().Owner, pullRequest.GetTargetEntity().Repo, pullRequest.PullRequest.GetBase().GetSHA(), changedFilesPath)
	if err != nil {
		return nil, err
	}

	authorsRank := gitHubClient.ComputeGitBlameRank(gitBlame)
	for _, reviewerRank := range authorsRank {
		authors = append(authors, reviewerRank.Username)
	}

	return authors, nil
}

func isUserEligibleToReview(
	username string,
	pullRequest *gh.PullRequest,
	reviewersToExclude []aladino.Value,
	availableAssignees []*codehost.User,
	totalOpenPRsAsReviewer int,
	maxAllowedAssignedReviews int,
) bool {
	if strings.HasSuffix(username, "[bot]") {
		return false
	}

	if pullRequest.GetUser().GetLogin() == username {
		return false
	}

	if isUserExcluded(reviewersToExclude, username) {
		return false
	}

	if !isUserValidAssignee(availableAssignees, username) {
		return false
	}

	if maxAllowedAssignedReviews > 0 && totalOpenPRsAsReviewer >= maxAllowedAssignedReviews {
		return false
	}

	return true
}

func isUserExcluded(users []aladino.Value, username string) bool {
	for _, excludedReviewer := range users {
		if excludedReviewer.(*aladino.StringValue).Val == username {
			return true
		}
	}
	return false
}

func isUserValidAssignee(assignees []*codehost.User, username string) bool {
	for _, availableAssignee := range assignees {
		if availableAssignee.Login == username {
			return true
		}
	}
	return false
}
