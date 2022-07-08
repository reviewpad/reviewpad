// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	"github.com/reviewpad/reviewpad/v2/utils"
)

func AssignAssignee() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType())}, nil),
		Code: assignAssigneeCode,
	}
}

func assignAssigneeCode(e aladino.Env, args []aladino.Value) error {
	assignees := args[0].(*aladino.ArrayValue).Vals
	if len(assignees) == 0 {
		return fmt.Errorf("assignAssignee: list of assignees can't be empty")
	}

	if len(assignees) > 10 {
		return fmt.Errorf("assignAssignee: can only assign up to 10 assignees")
	}

	assigneesLogin := make([]string, len(assignees))
	for i, assignee := range assignees {
		assigneesLogin[i] = assignee.(*aladino.StringValue).Val
	}

	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestRepoName(e.GetPullRequest())

	_, _, err := e.GetClient().Issues.AddAssignees(e.GetCtx(), owner, repo, prNum, assigneesLogin)

	return err
}
