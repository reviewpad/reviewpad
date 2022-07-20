// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"
	"log"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
)

func AssignReviewer() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType()), aladino.BuildIntType()}, nil),
		Code: assignReviewerCode,
	}
}

func assignReviewerCode(e aladino.Env, args []aladino.Value) error {
	totalRequiredReviewers := args[1].(*aladino.IntValue).Val
	if totalRequiredReviewers == 0 {
		return fmt.Errorf("assignReviewer: total required reviewers can't be 0")
	}

	availableReviewers := args[0].(*aladino.ArrayValue).Vals
	if len(availableReviewers) == 0 {
		return fmt.Errorf("assignReviewer: list of reviewers can't be empty")
	}

	// Remove pull request author from provided reviewers list
	for index, reviewer := range availableReviewers {
		if reviewer.(*aladino.StringValue).Val == *e.GetPullRequest().User.Login {
			availableReviewers = append(availableReviewers[:index], availableReviewers[index+1:]...)
			break
		}
	}

	totalAvailableReviewers := len(availableReviewers)
	if totalRequiredReviewers > totalAvailableReviewers {
		log.Printf("assignReviewer: total required reviewers %v exceeds the total available reviewers %v", totalRequiredReviewers, totalAvailableReviewers)
		totalRequiredReviewers = totalAvailableReviewers
	}

	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestBaseOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestBaseRepoName(e.GetPullRequest())

	reviewers := []string{}

	reviews, _, err := e.GetClient().PullRequests.ListReviews(e.GetCtx(), owner, repo, prNum, nil)
	if err != nil {
		return err
	}

	// Skip current requested reviewers if pull request already reviewed
	for _, review := range reviews {
		if review.State != nil && *review.State == "APPROVED" {
			log.Printf("assignReviewer: skipping request reviewers. the pull request already reviewed")
			return nil
		}
	}

	// Re-request current reviewers if mention on the provided reviewers list
	for _, review := range reviews {
		for index, availableReviewer := range availableReviewers {
			if availableReviewer.(*aladino.StringValue).Val == *review.User.Login {
				totalRequiredReviewers--
				reviewers = append(reviewers, *review.User.Login)
				availableReviewers = append(availableReviewers[:index], availableReviewers[index+1:]...)
				break
			}
		}
	}

	// Skip current requested reviewers if mention on the provided reviewers list
	currentRequestedReviewers := e.GetPullRequest().RequestedReviewers
	for _, requestedReviewer := range currentRequestedReviewers {
		for index, availableReviewer := range availableReviewers {
			if availableReviewer.(*aladino.StringValue).Val == *requestedReviewer.Login {
				totalRequiredReviewers--
				availableReviewers = append(availableReviewers[:index], availableReviewers[index+1:]...)
				break
			}
		}
	}

	// Select random reviewers from the list of all provided reviewers
	for i := 0; i < totalRequiredReviewers; i++ {
		selectedElementIndex := utils.GenerateRandom(len(availableReviewers))

		selectedReviewer := availableReviewers[selectedElementIndex]
		availableReviewers = append(availableReviewers[:selectedElementIndex], availableReviewers[selectedElementIndex+1:]...)

		reviewers = append(reviewers, selectedReviewer.(*aladino.StringValue).Val)
	}

	if len(reviewers) == 0 {
		log.Printf("assignReviewer: skipping request reviewers. the pull request already has reviewers")
		return nil
	}

	_, _, err = e.GetClient().PullRequests.RequestReviewers(e.GetCtx(), owner, repo, prNum, github.ReviewersRequest{
		Reviewers: reviewers,
	})

	return err
}
