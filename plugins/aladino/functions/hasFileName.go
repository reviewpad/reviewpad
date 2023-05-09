// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func HasFileName() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code:           hasFileNameCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func hasFileNameCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	fileNameStr := args[0].(*aladino.StringValue)

	patch := e.GetTarget().(*target.PullRequestTarget).Patch
	for fp := range patch {
		if fp == fileNameStr.Val {
			return aladino.BuildTrueValue(), nil
		}
	}

	return aladino.BuildFalseValue(), nil
}
