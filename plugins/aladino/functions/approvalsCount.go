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

func ApprovalsCount() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType(
			[]aladino.Type{},
			aladino.BuildIntType(),
		),
		Code:           approvalsCountCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func approvalsCountCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget)

	count, err := pullRequest.GetApprovalsCount()
	if err != nil {
		return nil, err
	}

	return lang.BuildIntValue(count), nil
}
