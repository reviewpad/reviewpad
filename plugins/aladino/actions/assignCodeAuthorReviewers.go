// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost"
	"github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"golang.org/x/exp/slices"
)

func AssignCodeAuthorReviewers() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: lang.BuildFunctionType([]lang.Type{
			lang.BuildIntType(),
			lang.BuildArrayOfType(lang.BuildStringType()),
			lang.BuildIntType(),
		}, nil),
		Code:           assignCodeAuthorReviewersCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func assignCodeAuthorReviewersCode(env aladino.Env, args []lang.Value) error {
	totalRequiredReviewers := args[0].(*lang.IntValue).Val
	reviewersToExclude := args[1].(*lang.ArrayValue).Vals
	maxAllowedAssignedReviews := args[2].(*lang.IntValue).Val

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

	// if a pull request has no changes for example when
	// you edit a file in one commit but revert the changes in another
	// we can end up in a situation where there are no changes in the pull request
	// in those cases it's not possible to fetch the authors of the changed files
	// so we just assign a random reviewer.
	if len(pr.Patch) == 0 {
		return assignRandomReviewerCode(env, nil)
	}

	// Fetch all users that have authored the changed files in the pull request.
	authors, err := getAuthorsFromGitBlame(env.GetCtx(), gitHubClient, pr)
	if err != nil {
		return fmt.Errorf("error getting authors from git blame: %s", err)
	}

	// Fetch all reviews in order to check if there are any reviews submitted by the authors.
	reviews, err := pr.GetReviews()
	if err != nil {
		return fmt.Errorf("error getting reviews: %s", err)
	}

	// Get the reviews that were submitted by the code authors that are valid to be assigned as a reviewer.
	// A code author is valid to be assigned as a reviewer if they are not included in the list of users to exclude from review requests.
	filteredReviews := filterReviewsByNonExcludedCodeAuthors(reviews, authors, reviewersToExclude)
	if len(filteredReviews) > 0 {
		return nil
	}

	// Fetch all available assignees to which the pull request may be assigned to
	// in order to avoid assigning a reviewer that is not available.
	availableAssignees, err := pr.GetAvailableAssignees()
	if err != nil {
		return fmt.Errorf("error getting available assignees: %s", err)
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

	// we are excluding yarn.lock files
	// to investigate if they are causing an issue with the git blame.
	excludedFiles := map[string]bool{
		"yarn.lock": true,
	}

	for _, patch := range pullRequest.Patch {
		if excludedFiles[patch.Repr.GetFilename()] {
			continue
		}

		changedFilesPath = append(changedFilesPath, patch.Repr.GetFilename())
	}

	gitBlame, err := gitHubClient.GetGitBlame(ctx, pullRequest.GetTargetEntity().Owner, pullRequest.GetTargetEntity().Repo, pullRequest.PullRequest.GetBase().GetSha(), changedFilesPath)
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
	codeReview *pbc.PullRequest,
	reviewersToExclude []lang.Value,
	availableAssignees []*codehost.User,
	totalOpenPRsAsReviewer int,
	maxAllowedAssignedReviews int,
) bool {
	if strings.HasSuffix(username, "[bot]") {
		return false
	}

	if codeReview.Author.Login == username {
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

func isUserExcluded(users []lang.Value, username string) bool {
	for _, excludedReviewer := range users {
		if excludedReviewer.(*lang.StringValue).Val == username {
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

func filterReviewsByNonExcludedCodeAuthors(reviews []*codehost.Review, authors []string, excludedAuthors []lang.Value) []*codehost.Review {
	filteredReviews := make([]*codehost.Review, 0)
	for _, review := range reviews {
		if slices.Contains(authors, review.User.Login) && !isUserExcluded(excludedAuthors, review.User.Login) {
			filteredReviews = append(filteredReviews, review)
		}
	}

	return filteredReviews
}
