// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
)

// contain all the available reviewer status states
const (
	ReviewerStatusStateCommented       string = "COMMENTED"
	ReviewerStatusStateRequestChanges  string = "CHANGES_REQUESTED"
	ReviewerStatusStateRequestApproved string = "APPROVED"
)

// ReviewerStatus return the lastest review status of reviewer on a pull request
func ReviewerStatus() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildStringType()),
		Code: reviewerStatusCdoe,
	}
}

// reviewerStatusCdoe return the lastest review status (neutral, requested_changes, approved) of the reviewerLogin on a pull request
func reviewerStatusCdoe(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	reviewerLogin := args[0].(*aladino.StringValue)

	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestBaseOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestBaseRepoName(e.GetPullRequest())

	reviews, err := utils.GetPullRequestReviews(e.GetCtx(), e.GetClient(), owner, repo, prNum, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	status := "neutral"
	// iterate over the reviewersLogin and update the status for the reviewerLogin
	for _, r := range reviews {
		if *r.User.Login != *&reviewerLogin.Val {
			continue
		}

		switch *r.State {

		case ReviewerStatusStateCommented:
			status = "neutral"

		case ReviewerStatusStateRequestChanges:
			status = "requested_changes"

		case ReviewerStatusStateRequestApproved:
			status = "approved"

		}
	}

	fmt.Println("Hello", status)

	return aladino.BuildStringValue(status), nil
}
