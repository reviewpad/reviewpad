// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"
	"sort"

	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/codehost"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
)

func AssignReviewer() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: lang.BuildFunctionType(
			[]lang.Type{
				lang.BuildArrayOfType(lang.BuildStringType()),
				lang.BuildIntType(),
				lang.BuildStringType(),
			},
			nil,
		),
		Code:           assignReviewerCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest},
	}
}

func assignReviewerCode(e aladino.Env, args []lang.Value) error {
	availableReviewers := args[0].(*lang.ArrayValue).Vals
	totalRequiredReviewers := args[1].(*lang.IntValue).Val
	policy := args[2].(*lang.StringValue).Val
	target := e.GetTarget().(*target.PullRequestTarget)
	log := e.GetLogger().WithField("builtin", "assignReviewer")

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
		log.Infof("total required reviewers %v exceeds the total available reviewers %v", totalRequiredReviewers, totalAvailableReviewers)
		totalRequiredReviewers = totalAvailableReviewers
	}

	reviewers := []string{}

	reviews, err := target.GetReviews()
	if err != nil {
		return err
	}

	lastPushDate, err := target.GetPullRequestLastPushDate()
	if err != nil {
		return err
	}

	// Re-request current reviewers only when last review status is not APPROVED
	for _, availableReviewer := range availableReviewers {
		userLogin := availableReviewer.(*lang.StringValue).Val
		if codehost.HasReview(reviews, userLogin) {
			lastReview := codehost.LastReview(reviews, userLogin)
			if lastReview.State != "APPROVED" && lastReview.SubmittedAt.Before(lastPushDate) {
				reviewers = append(reviewers, userLogin)
			} else {
				log.Infof("reviewer %v has already approved the pull request", userLogin)
			}
			totalRequiredReviewers--
			availableReviewers = filterReviewerFromReviewers(availableReviewers, userLogin)
		}
	}

	// Skip current requested reviewers if mention on the provided reviewers list
	currentRequestedReviewers := target.GetRequestedReviewers()

	for _, requestedReviewer := range currentRequestedReviewers {
		for _, availableReviewer := range availableReviewers {
			if availableReviewer.(*lang.StringValue).Val == requestedReviewer.Login {
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
		log.Infof("no reviewers were assigned")
		return nil
	}

	return target.RequestReviewers(reviewers)
}

func filterReviewerFromReviewers(reviewers []lang.Value, reviewer string) []lang.Value {
	var filteredReviewers []lang.Value
	for _, r := range reviewers {
		if r.(*lang.StringValue).Val != reviewer {
			filteredReviewers = append(filteredReviewers, r)
		}
	}
	return filteredReviewers
}

func getReviewersUsingPolicyRandom(availableReviewers []lang.Value, totalRequiredReviewers int) []string {
	reviewers := []string{}
	for i := 0; i < totalRequiredReviewers; i++ {
		selectedElementIndex := utils.GenerateRandom(len(availableReviewers))

		selectedReviewer := availableReviewers[selectedElementIndex].(*lang.StringValue).Val
		availableReviewers = filterReviewerFromReviewers(availableReviewers, selectedReviewer)

		reviewers = append(reviewers, selectedReviewer)
	}
	return reviewers
}

func getReviewersUsingPolicyRoundRobin(e aladino.Env, availableReviewers []lang.Value, totalRequiredReviewers int) []string {
	reviewers := []string{}

	// Use pull request number as a starting point to select reviewers
	pullRequest := e.GetTarget().(*target.PullRequestTarget).PullRequest
	prNum := int(pullRequest.GetNumber())
	startPos := (prNum - 1) * totalRequiredReviewers

	if startPos < 0 {
		startPos = 0
	}

	for i := 0; i < totalRequiredReviewers; i++ {
		pos := (startPos + i) % len(availableReviewers)
		selectedReviewer := availableReviewers[pos]
		reviewers = append(reviewers, selectedReviewer.(*lang.StringValue).Val)
	}

	return reviewers
}

func getReviewersUsingPolicyReviewpad(e aladino.Env, availableReviewers []lang.Value, totalRequiredReviewers int) ([]string, error) {
	reviewers := []string{}
	reviewersMap := map[int][]string{}

	for _, reviewer := range availableReviewers {
		// Look for reviewers with less hanging PRs
		query := fmt.Sprintf("is:open is:pr review-requested:%s", reviewer.(*lang.StringValue).Val)
		issues, _, err := e.GetGithubClient().GetClientREST().Search.Issues(e.GetCtx(), query, nil)
		if err != nil {
			return nil, err
		}

		r := reviewersMap[*issues.Total]
		if len(r) == 0 {
			reviewersMap[*issues.Total] = []string{reviewer.(*lang.StringValue).Val}
		} else {
			reviewersMap[*issues.Total] = append(r, reviewer.(*lang.StringValue).Val)
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
