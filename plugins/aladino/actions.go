// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	"github.com/reviewpad/reviewpad/v2/utils"
)

const ReviewpadCommentAnnotation = "<!--@annotation-reviewpad-single-comment-->"

func addLabel() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code: addLabelCode,
	}
}

func addLabelCode(e aladino.Env, args []aladino.Value) error {
	if len(args) != 1 {
		return fmt.Errorf("addLabel: expecting 1 argument, got %v", len(args))
	}

	labelVal := args[0]
	if !labelVal.HasKindOf(aladino.STRING_VALUE) {
		return fmt.Errorf("addLabel: expecting string argument, got %v", labelVal.Kind())
	}

	label := labelVal.(*aladino.StringValue).Val

	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestRepoName(e.GetPullRequest())

	_, _, err := e.GetClient().Issues.GetLabel(e.GetCtx(), owner, repo, label)
	if err != nil {
		return err
	}

	_, _, err = e.GetClient().Issues.AddLabelsToIssue(e.GetCtx(), owner, repo, prNum, []string{label})

	return err
}

func assignRandomReviewer() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, nil),
		Code: assignRandomReviewerCode,
	}
}

func assignRandomReviewerCode(e aladino.Env, _ []aladino.Value) error {
	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestRepoName(e.GetPullRequest())

	ghPr, _, err := e.GetClient().PullRequests.Get(e.GetCtx(), owner, repo, prNum)
	if err != nil {
		return err
	}

	// When there's already assigned reviewers, do nothing
	totalRequestReviewers := len(ghPr.RequestedReviewers)
	if totalRequestReviewers > 0 {
		return nil
	}

	ghUsers, _, err := e.GetClient().Repositories.ListCollaborators(e.GetCtx(), owner, repo, nil)
	if err != nil {
		return err
	}

	filteredGhUsers := []*github.User{}

	for i := range ghUsers {
		if ghUsers[i].GetLogin() != ghPr.GetUser().GetLogin() {
			filteredGhUsers = append(filteredGhUsers, ghUsers[i])
		}
	}

	if len(filteredGhUsers) == 0 {
		return fmt.Errorf("can't assign a random user because there is no users")
	}

	lucky := utils.GenerateRandom(len(filteredGhUsers))
	ghUser := filteredGhUsers[lucky]

	_, _, err = e.GetClient().PullRequests.RequestReviewers(e.GetCtx(), owner, repo, prNum, github.ReviewersRequest{
		Reviewers: []string{ghUser.GetLogin()},
	})

	return err
}

func assignReviewer() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType()), aladino.BuildIntType()}, nil),
		Code: assignReviewerCode,
	}
}

func assignReviewerCode(e aladino.Env, args []aladino.Value) error {
	if len(args) < 1 {
		return fmt.Errorf("assignReviewer: expecting at least 1 argument")
	}

	arg := args[0]
	if !arg.HasKindOf(aladino.ARRAY_VALUE) {
		return fmt.Errorf("assignReviewer: requires array argument, got %v", arg.Kind())
	}

	if !args[1].HasKindOf(aladino.INT_VALUE) {
		return fmt.Errorf("assignReviewer: the parameter total is required to be an int, instead got %v", args[1].Kind())
	}

	totalRequiredReviewers := args[1].(*aladino.IntValue).Val

	availableReviewers := arg.(*aladino.ArrayValue).Vals

	for _, reviewer := range availableReviewers {
		if !reviewer.HasKindOf(aladino.STRING_VALUE) {
			return fmt.Errorf("assignReviewer: requires array of strings, got array with value of %v", reviewer.Kind())
		}
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
	owner := utils.GetPullRequestOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestRepoName(e.GetPullRequest())

	reviewers := []string{}

	reviews, _, err := e.GetClient().PullRequests.ListReviews(e.GetCtx(), owner, repo, prNum, nil)
	if err != nil {
		return err
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

func assignTeamReviewer() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType())}, nil),
		Code: assignTeamReviewerCode,
	}
}

func assignTeamReviewerCode(e aladino.Env, args []aladino.Value) error {
	teamReviewers := args[0].(*aladino.ArrayValue).Vals

	if len(teamReviewers) < 1 {
		return fmt.Errorf("assignTeamReviewer: requires at least 1 team to request for review")
	}

	teamReviewersSlugs := make([]string, len(teamReviewers))

	for i, team := range teamReviewers {
		teamReviewersSlugs[i] = team.(*aladino.StringValue).Val
	}

	pullRequest := e.GetPullRequest()
	prNum := utils.GetPullRequestNumber(pullRequest)
	owner := utils.GetPullRequestOwnerName(pullRequest)
	repo := utils.GetPullRequestRepoName(pullRequest)

	_, _, err := e.GetClient().PullRequests.RequestReviewers(e.GetCtx(), owner, repo, prNum, github.ReviewersRequest{
		TeamReviewers: teamReviewersSlugs,
	})

	return err
}

func close() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, nil),
		Code: closeCode,
	}
}

func closeCode(e aladino.Env, args []aladino.Value) error {
	pullRequest := e.GetPullRequest()

	prNum := utils.GetPullRequestNumber(pullRequest)
	owner := utils.GetPullRequestOwnerName(pullRequest)
	repo := utils.GetPullRequestRepoName(pullRequest)

	closedState := "closed"
	pullRequest.State = &closedState
	_, _, err := e.GetClient().PullRequests.Edit(e.GetCtx(), owner, repo, prNum, pullRequest)

	return err
}

func comment() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code: commentCode,
	}
}

func commentCode(e aladino.Env, args []aladino.Value) error {
	pullRequest := e.GetPullRequest()

	prNum := utils.GetPullRequestNumber(pullRequest)
	owner := utils.GetPullRequestOwnerName(pullRequest)
	repo := utils.GetPullRequestRepoName(pullRequest)

	commentBody := args[0].(*aladino.StringValue).Val

	_, _, err := e.GetClient().Issues.CreateComment(e.GetCtx(), owner, repo, prNum, &github.IssueComment{
		Body: &commentBody,
	})

	return err
}

func commentOnce() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code: commentOnceCode,
	}
}

func commentOnceCode(e aladino.Env, args []aladino.Value) error {
	pullRequest := e.GetPullRequest()

	prNum := utils.GetPullRequestNumber(pullRequest)
	owner := utils.GetPullRequestOwnerName(pullRequest)
	repo := utils.GetPullRequestRepoName(pullRequest)

	commentBody := args[0].(*aladino.StringValue).Val
	commentBodyWithReviewpadAnnotation := fmt.Sprintf("%v%v", ReviewpadCommentAnnotation, commentBody)
	commentBodyWithReviewpadAnnotationHash := sha256.Sum256([]byte(commentBodyWithReviewpadAnnotation))

	comments, err := utils.GetPullRequestComments(e.GetCtx(), e.GetClient(), owner, repo, prNum, &github.IssueListCommentsOptions{})
	if err != nil {
		return err
	}

	for _, comment := range comments {
		commentHash := sha256.Sum256([]byte(*comment.Body))
		commentAlreadyExists := commentHash == commentBodyWithReviewpadAnnotationHash
		if commentAlreadyExists {
			return nil
		}
	}

	_, _, err = e.GetClient().Issues.CreateComment(e.GetCtx(), owner, repo, prNum, &github.IssueComment{
		Body: &commentBodyWithReviewpadAnnotation,
	})

	return err
}

func merge() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code: mergeCode,
	}
}

func mergeCode(e aladino.Env, args []aladino.Value) error {
	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestRepoName(e.GetPullRequest())

	mergeMethod, err := parseMergeMethod(args)
	if err != nil {
		return err
	}

	_, _, err = e.GetClient().PullRequests.Merge(e.GetCtx(), owner, repo, prNum, "Merged by Reviewpad", &github.PullRequestOptions{
		MergeMethod: mergeMethod,
	})
	return err
}

func parseMergeMethod(args []aladino.Value) (string, error) {
	if len(args) > 1 {
		return "", fmt.Errorf("merge: received two arguments")
	}

	if len(args) == 0 {
		return "merge", nil
	}

	arg := args[0]
	if arg.HasKindOf(aladino.STRING_VALUE) {
		mergeMethod := arg.(*aladino.StringValue).Val
		switch mergeMethod {
		case "merge", "rebase", "squash":
			return mergeMethod, nil
		default:
			return "", fmt.Errorf("merge: unexpected argument %v", mergeMethod)
		}
	} else {
		return "", fmt.Errorf("merge: expects string argument")
	}
}

func removeLabel() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code: removeLabelCode,
	}
}

func removeLabelCode(e aladino.Env, args []aladino.Value) error {
	if len(args) != 1 {
		return fmt.Errorf("removeLabel: expecting 1 argument, got %v", len(args))
	}

	labelVal := args[0]
	if !labelVal.HasKindOf(aladino.STRING_VALUE) {
		return fmt.Errorf("removeLabel: expecting string argument, got %v", labelVal.Kind())
	}

	label := labelVal.(*aladino.StringValue).Val

	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestRepoName(e.GetPullRequest())

	_, _, err := e.GetClient().Issues.GetLabel(e.GetCtx(), owner, repo, label)
	if err != nil {
		return err
	}

	var labelIsAppliedToPullRequest bool = false
	for _, ghLabel := range e.GetPullRequest().Labels {
		if ghLabel.GetName() == label {
			labelIsAppliedToPullRequest = true
			break
		}
	}

	if !labelIsAppliedToPullRequest {
		return nil
	}

	_, err = e.GetClient().Issues.RemoveLabelForIssue(e.GetCtx(), owner, repo, prNum, label)

	return err
}
