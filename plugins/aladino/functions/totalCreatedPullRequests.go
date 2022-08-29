// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func TotalCreatedPullRequests() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildIntType()),
		Code:           totalCreatedPullRequestsCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func totalCreatedPullRequestsCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	devName := args[0].(*aladino.StringValue).Val

	entity := e.GetTarget().GetTargetEntity()
	owner := entity.Owner
	repo := entity.Repo

	issues, _, err := e.GetGithubClient().ListIssuesByRepo(e.GetCtx(), owner, repo, &github.IssueListByRepoOptions{
		Creator: devName,
		State:   "all",
	})
	if err != nil {
		return nil, err
	}

	count := 0
	for _, issue := range issues {
		if issue.IsPullRequest() {
			count++
		}
	}

	return aladino.BuildIntValue(count), nil
}
