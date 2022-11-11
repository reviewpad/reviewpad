// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"
	"log"
	"sort"

	"github.com/reviewpad/reviewpad/v3/codehost"
	"github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
)

func AssignReviewer() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType(
			[]aladino.Type{
				aladino.BuildArrayOfType(aladino.BuildStringType()),
				aladino.BuildIntType(),
				aladino.BuildStringType(),
			},
			nil,
		),
		Code:           assignReviewerCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func assignReviewerCode(e aladino.Env, args []aladino.Value) error {
	availableReviewers := args[0].(*aladino.ArrayValue).Vals
	totalRequiredReviewers := args[1].(*aladino.IntValue).Val
	policy := args[2].(*aladino.StringValue).Val
	target := e.GetTarget().(*target.PullRequestTarget)
	targetEntity := target.GetTargetEntity()
	ctx := e.GetCtx()
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	prNum := targetEntity.Number

	allowedPolicies := map[string]bool{"random": true, "round-robin": true, "reviewpad": true}
	if _, ok := allowedPolicies[policy]; !ok {
		return fmt.Errorf("assignReviewer: policy %s is not supported. allowed policies %v", policy, allowedPolicies)
	}

	if totalRequiredReviewers == 0 {
		return fmt.Errorf("assignReviewer: total required reviewers can't be 0")
	}

	if len(availableReviewers) == 0 {
		return fmt.Errorf("assignReviewer: list of reviewers can't be empty")
	}

	author, err := target.GetAuthor()
	if err != nil {
		return err
	}

	availableReviewers = filterReviewerFromReviewers(availableReviewers, author.Login)

	totalAvailableReviewers := len(availableReviewers)
	if totalRequiredReviewers > totalAvailableReviewers {
		log.Printf("assignReviewer: total required reviewers %v exceeds the total available reviewers %v", totalRequiredReviewers, totalAvailableReviewers)
		totalRequiredReviewers = totalAvailableReviewers
	}

	reviewers := []string{}

	reviews, err := target.GetReviews()
	if err != nil {
		return err
	}

	lastPushDate, err := e.GetGithubClient().GetPullRequestLastPushDate(ctx, owner, repo, prNum)
	if err != nil {
		return err
	}

	// Re-request current reviewers only when last review status is not APPROVED
	for _, availableReviewer := range availableReviewers {
		userLogin := availableReviewer.(*aladino.StringValue).Val
		if codehost.HasReview(reviews, userLogin) {
			lastReview := codehost.LastReview(reviews, userLogin)
			if lastReview.State != "APPROVED" && lastReview.SubmittedAt.Before(lastPushDate) {
				reviewers = append(reviewers, userLogin)
			} else {
				log.Printf("assignReviewer: reviewer %v has already approved the pull request", userLogin)
			}
			totalRequiredReviewers--
			availableReviewers = filterReviewerFromReviewers(availableReviewers, userLogin)
		}
	}

	// Skip current requested reviewers if mention on the provided reviewers list
	currentRequestedReviewers, err := target.GetRequestedReviewers()
	if err != nil {
		return err
	}

	for _, requestedReviewer := range currentRequestedReviewers {
		for _, availableReviewer := range availableReviewers {
			if availableReviewer.(*aladino.StringValue).Val == requestedReviewer.Login {
				totalRequiredReviewers--
				availableReviewers = filterReviewerFromReviewers(availableReviewers, requestedReviewer.Login)
				break
			}
		}
	}

	switch policy {
	case "random":
		reviewers = append(reviewers, getReviewersUsingPolicyRandom(availableReviewers, totalRequiredReviewers)...)
	case "round-robin":
		reviewers = append(reviewers, getReviewersUsingPolicyRoundRobin(e, availableReviewers, totalRequiredReviewers)...)
	case "reviewpad":
		r, err := getReviewersUsingPolicyReviewpad(e, availableReviewers, totalRequiredReviewers)
		if err != nil {
			return err
		}
		reviewers = append(reviewers, r...)
	}

	if len(reviewers) == 0 {
		log.Printf("assignReviewer: no reviewers were assigned")
		return nil
	}

	return target.RequestReviewers(reviewers)
}

func filterReviewerFromReviewers(reviewers []aladino.Value, reviewer string) []aladino.Value {
	var filteredReviewers []aladino.Value
	for _, r := range reviewers {
		if r.(*aladino.StringValue).Val != reviewer {
			filteredReviewers = append(filteredReviewers, r)
		}
	}
	return filteredReviewers
}

func getReviewersUsingPolicyRandom(availableReviewers []aladino.Value, totalRequiredReviewers int) []string {
	reviewers := []string{}
	for i := 0; i < totalRequiredReviewers; i++ {
		selectedElementIndex := utils.GenerateRandom(len(availableReviewers))

		selectedReviewer := availableReviewers[selectedElementIndex].(*aladino.StringValue).Val
		availableReviewers = filterReviewerFromReviewers(availableReviewers, selectedReviewer)

		reviewers = append(reviewers, selectedReviewer)
	}
	return reviewers
}

func getReviewersUsingPolicyRoundRobin(e aladino.Env, availableReviewers []aladino.Value, totalRequiredReviewers int) []string {
	reviewers := []string{}

	// Use pull request number as a starting point to select reviewers
	pullRequest := e.GetTarget().(*target.PullRequestTarget).PullRequest
	prNum := github.GetPullRequestNumber(pullRequest)
	startPos := (prNum - 1) * totalRequiredReviewers

	if startPos < 0 {
		startPos = 0
	}

	for i := 0; i < totalRequiredReviewers; i++ {
		pos := (startPos + i) % len(availableReviewers)
		selectedReviewer := availableReviewers[pos]
		reviewers = append(reviewers, selectedReviewer.(*aladino.StringValue).Val)
	}

	return reviewers
}

func getReviewersUsingPolicyReviewpad(e aladino.Env, availableReviewers []aladino.Value, totalRequiredReviewers int) ([]string, error) {
	reviewers := []string{}
	reviewersMap := map[int][]string{}

	for _, reviewer := range availableReviewers {
		// Look for reviewers with less hanging PRs
		query := fmt.Sprintf("is:open is:pr review-requested:%s", reviewer.(*aladino.StringValue).Val)
		issues, _, err := e.GetGithubClient().GetClientREST().Search.Issues(e.GetCtx(), query, nil)
		if err != nil {
			return nil, err
		}

		r := reviewersMap[*issues.Total]
		if len(r) == 0 {
			reviewersMap[*issues.Total] = []string{reviewer.(*aladino.StringValue).Val}
		} else {
			reviewersMap[*issues.Total] = append(r, reviewer.(*aladino.StringValue).Val)
		}
	}

	reviewersMapKeys := []int{}
	for key := range reviewersMap {
		reviewersMapKeys = append(reviewersMapKeys, key)
	}
	sort.Ints(reviewersMapKeys)

	orderedReviewers := []string{}
	for _, k := range reviewersMapKeys {
		orderedReviewers = append(orderedReviewers, reviewersMap[k]...)
	}

	for i := 0; i < totalRequiredReviewers; i++ {
		reviewers = append(reviewers, orderedReviewers[i])
	}

	return reviewers, nil
}
