// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Base() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildStringType()),
		Code:           baseCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func baseCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	t := e.GetTarget().(*target.PullRequestTarget)

	return aladino.BuildStringValue(t.GetBase()), nil
}
