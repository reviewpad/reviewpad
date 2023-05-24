// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func HasLinearHistory() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code:           hasLinearHistoryCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func hasLinearHistoryCode(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget)

	ghCommits, err := pullRequest.GetCommits()
	if err != nil {
		return nil, err
	}

	for _, ghCommit := range ghCommits {
		if ghCommit.ParentsCount > 1 {
			return lang.BuildBoolValue(false), nil
		}
	}

	return lang.BuildBoolValue(true), nil
}
