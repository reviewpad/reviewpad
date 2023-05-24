// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func IsUpdatedWithBaseBranch() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code:           isUpdatedWithBaseBranchCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func isUpdatedWithBaseBranchCode(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget).PullRequest
	targetEntity := e.GetTarget().GetTargetEntity()

	pullRequestUpToDate, err := e.GetGithubClient().GetPullRequestUpToDate(
		e.GetCtx(),
		targetEntity.Owner,
		targetEntity.Repo,
		int(pullRequest.GetNumber()),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting pull request outdated information: %s", err.Error())
	}

	return lang.BuildBoolValue(pullRequestUpToDate), nil
}
