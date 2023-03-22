// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"

	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func IsUpdatedWithBaseBranch() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code:           isUpdatedWithBaseBranchCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func isUpdatedWithBaseBranchCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget).CodeReview
	targetEntity := e.GetTarget().GetTargetEntity()

	headBehindBy, err := e.GetGithubClient().GetHeadBehindBy(
		e.GetCtx(),
		targetEntity.Owner,
		targetEntity.Repo,
		pullRequest.GetHead().GetRepo().GetOwner(),
		pullRequest.GetHead().GetName(),
		int(pullRequest.GetNumber()),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting head behind by information: %s", err.Error())
	}

	return aladino.BuildBoolValue(headBehindBy == 0), nil
}
