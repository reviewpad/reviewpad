// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"log"
	"sort"
)

func AssignReviewerWithPolicy() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType()), aladino.BuildIntType(), aladino.BuildStringType()}, nil),
		Code:           assignReviewerWithPolicyCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func assignReviewerWithPolicyCode(e aladino.Env, args []aladino.Value) error {
	policies := map[string]bool{"random": true, "round-robin": true, "reviewpad": true}

	policy := args[2].(*aladino.StringValue).Val
	if _, ok := policies[policy]; !ok {
		return fmt.Errorf("assignReviewerWithPolicy: wrong policy (random, round-robin, reviewpad")
	}

	t := e.GetTarget().(*target.PullRequestTarget)
	totalRequiredReviewers := args[1].(*aladino.IntValue).Val
	if totalRequiredReviewers == 0 {
		return fmt.Errorf("assignReviewerWithPolicy: total required reviewers can't be 0")
	}

	availableReviewers := args[0].(*aladino.ArrayValue).Vals
	if len(availableReviewers) == 0 {
		return fmt.Errorf("assignReviewerWithPolicy: list of reviewers can't be empty")
	}

	author, err := t.GetAuthor()
	if err != nil {
		return err
	}

	// Remove pull request author from provided reviewers list
	for index, reviewer := range availableReviewers {
		if reviewer.(*aladino.StringValue).Val == author.Login {
			availableReviewers = append(availableReviewers[:index], availableReviewers[index+1:]...)
			break
		}
	}

	totalAvailableReviewers := len(availableReviewers)
	if totalRequiredReviewers > totalAvailableReviewers {
		log.Printf("assignReviewerWithPolicy: total required reviewers %v exceeds the total available reviewers %v", totalRequiredReviewers, totalAvailableReviewers)
		totalRequiredReviewers = totalAvailableReviewers
	}

	reviewers := []string{}

	reviews, err := t.GetReviews()
	if err != nil {
		return err
	}

	// Re-request current reviewers if mention on the provided reviewers list
	for _, review := range reviews {
		for index, availableReviewer := range availableReviewers {
			if availableReviewer.(*aladino.StringValue).Val == review.User.Login {
				totalRequiredReviewers--
				if review.State != "APPROVED" {
					reviewers = append(reviewers, review.User.Login)
				}
				availableReviewers = append(availableReviewers[:index], availableReviewers[index+1:]...)
				break
			}
		}
	}

	// Skip current requested reviewers if mention on the provided reviewers list
	currentRequestedReviewers, err := t.GetRequestedReviewers()
	if err != nil {
		return err
	}

	for _, requestedReviewer := range currentRequestedReviewers {
		for index, availableReviewer := range availableReviewers {
			if availableReviewer.(*aladino.StringValue).Val == requestedReviewer.Login {
				totalRequiredReviewers--
				availableReviewers = append(availableReviewers[:index], availableReviewers[index+1:]...)
				break
			}
		}
	}

	switch policy {
	case "random":
		// Select random reviewers from the list of all provided reviewers
		for i := 0; i < totalRequiredReviewers; i++ {
			selectedElementIndex := utils.GenerateRandom(len(availableReviewers))

			selectedReviewer := availableReviewers[selectedElementIndex]
			availableReviewers = append(availableReviewers[:selectedElementIndex], availableReviewers[selectedElementIndex+1:]...)

			reviewers = append(reviewers, selectedReviewer.(*aladino.StringValue).Val)
		}
	case "round-robin":
		// using the issue number and the number of required reviewers,
		// determine reviewers
		pullRequest := e.GetTarget().(*target.PullRequestTarget).PullRequest
		prNum := gh.GetPullRequestNumber(pullRequest)
		startPos := (prNum - 1) * totalRequiredReviewers
		if startPos < 0 {
			startPos = 0
		}
		for i := 0; i < totalRequiredReviewers; i++ {
			pos := (startPos + i) % len(availableReviewers)
			selectedReviewer := availableReviewers[pos]
			reviewers = append(reviewers, selectedReviewer.(*aladino.StringValue).Val)
		}
	case "reviewpad":
		// we are looking for someone with less hanging PR
		reviewersMap := map[int]string{}
		var reviewersIssues []int
		for _, reviewer := range availableReviewers {
			query := fmt.Sprintf("is:open is:pr review-requested:%s", reviewer.(*aladino.StringValue).Val)
			c := e.GetGithubClient().GetClientREST().Search
			issues, _, err := c.Issues(e.GetCtx(), query, nil)
			if err != nil {
				return err
			}

			reviewersMap[*issues.Total] = reviewer.(*aladino.StringValue).Val
			reviewersIssues = append(reviewersIssues, *issues.Total)
		}

		// sort reviewers
		sort.Ints(reviewersIssues)
		for i := 0; i < totalRequiredReviewers; i++ {
			reviewers = append(reviewers, reviewersMap[reviewersIssues[i]])
		}
	}

	if len(reviewers) == 0 {
		log.Printf("assignReviewerWithPolicy: skipping request reviewers. the pull request already has reviewers")
		return nil
	}

	return t.RequestReviewers(reviewers)
}
