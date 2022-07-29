// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
)

func HasUnaddressedReviewThreads() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code: hasUnaddressedReviewThreadsCode,
	}
}

func hasUnaddressedReviewThreadsCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestBaseOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestBaseRepoName(e.GetPullRequest())
	totalRequestTries := 2

	reviewThreads, err := utils.GetReviewThreads(e.GetCtx(), e.GetClientGQL(), owner, repo, prNum, totalRequestTries)
	if err != nil {
		return nil, err
	}

	for _, reviewThread := range reviewThreads {
		if !reviewThread.IsResolved || !reviewThread.IsOutdated {
			return aladino.BuildTrueValue(), nil
		}
	}

	return aladino.BuildFalseValue(), nil
}
