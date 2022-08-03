// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
)

func AssignTeamReviewer() *aladino.BuiltInAction {
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
	owner := utils.GetPullRequestBaseOwnerName(pullRequest)
	repo := utils.GetPullRequestBaseRepoName(pullRequest)

	_, _, err := e.GetClient().PullRequests.RequestReviewers(e.GetCtx(), owner, repo, prNum, github.ReviewersRequest{
		TeamReviewers: teamReviewersSlugs,
	})

	return err
}
