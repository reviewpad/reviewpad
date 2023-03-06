// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	gh "github.com/google/go-github/v49/github"
	"github.com/reviewpad/reviewpad/v4/codehost"
	"github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"golang.org/x/exp/slices"
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

	// Get all available assignees that are not authors.
	nonAuthorAvailableAssignees := filterAuthorsFromAssignees(availableAssignees, authors)

	// Fetch the total number of open pull requests that each author has already assigned to them as a reviewer.
	totalOpenPRsAsReviewerByUser, err := gitHubClient.GetOpenPullRequestsAsReviewer(env.GetCtx(), pr.GetTargetEntity().Owner, pr.GetTargetEntity().Repo, append(authors, nonAuthorAvailableAssignees...))
	if err != nil {
		return fmt.Errorf("error getting open pull requests as reviewer: %s", err)
	}

	selectedReviewers := []string{}
	for _, author := range authors {
		if isUserEligibleToReview(author, pr.PullRequest, reviewersToExclude, availableAssignees, totalOpenPRsAsReviewerByUser[author], maxAllowedAssignedReviews) {
			selectedReviewers = append(selectedReviewers, author)
		}
	}

	// If there are no authors eligible to review
	// find eligible reviewers from the available assignees.
	if len(selectedReviewers) == 0 {
		for _, assignee := range nonAuthorAvailableAssignees {
			if isUserEligibleToReview(assignee, pr.PullRequest, reviewersToExclude, availableAssignees, totalOpenPRsAsReviewerByUser[assignee], maxAllowedAssignedReviews) {
				selectedReviewers = append(selectedReviewers, assignee)
			}
		}

		selectedReviewers = getRandomUsers(selectedReviewers, totalRequiredReviewers)
	}

	// If no reviewers were found until now
	// assign a random reviewer.
	if len(selectedReviewers) == 0 {
		return assignRandomReviewerCode(env, nil)
	}

	if totalRequiredReviewers > len(selectedReviewers) {
		env.GetLogger().Warnf("number of required reviewers %d is less than available reviewers %d", totalRequiredReviewers, len(selectedReviewers))

		totalRequiredReviewers = len(selectedReviewers)
	}

	selectedReviewers = selectedReviewers[:totalRequiredReviewers]

	return pr.RequestReviewers(selectedReviewers)
}

func filterAuthorsFromAssignees(availableAssignees []*codehost.User, authors []string) []string {
	assignees := []string{}
	for _, assignee := range availableAssignees {
		if !slices.Contains(authors, assignee.Login) {
			assignees = append(assignees, assignee.Login)
		}
	}

	return assignees
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

func getRandomUsers(users []string, total int) []string {
	if total >= len(users) {
		return users
	}

	selectedUsers := make([]string, total)

	for i := 0; i < total; i++ {
		randIndex := rand.Intn(len(users))
		selectedUsers[i] = users[randIndex]

		users = append(users[:randIndex], users[randIndex+1:]...)
	}

	return selectedUsers
}
