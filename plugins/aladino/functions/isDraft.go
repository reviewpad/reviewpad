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

func IsDraft() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{}, lang.BuildBoolType()),
		Code:           isDraftCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func isDraftCode(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget)

	isDraft := pullRequest.IsDraft()
	return lang.BuildBoolValue(isDraft), nil
}
