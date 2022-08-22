// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"
	"log"

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
	t := e.GetTarget()
	totalRequiredReviewers := args[1].(*aladino.IntValue).Val
	if totalRequiredReviewers == 0 {
		return fmt.Errorf("assignReviewer: total required reviewers can't be 0")
	}

	availableReviewers := args[0].(*aladino.ArrayValue).Vals
	if len(availableReviewers) == 0 {
		return fmt.Errorf("assignReviewer: list of reviewers can't be empty")
	}

	user, err := t.GetUser()
	if err != nil {
		return err
	}

	// Remove pull request author from provided reviewers list
	for index, reviewer := range availableReviewers {
		if reviewer.(*aladino.StringValue).Val == user.Login {
			availableReviewers = append(availableReviewers[:index], availableReviewers[index+1:]...)
			break
		}
	}

	totalAvailableReviewers := len(availableReviewers)
	if totalRequiredReviewers > totalAvailableReviewers {
		log.Printf("assignReviewer: total required reviewers %v exceeds the total available reviewers %v", totalRequiredReviewers, totalAvailableReviewers)
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

	// TODO: #164 - Improve reviewer selection
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

	return t.RequestReviewers(reviewers)
}
