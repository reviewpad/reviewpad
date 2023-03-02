// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func HasCodePattern() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code:           hasCodePatternCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func hasCodePatternCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	arg := args[0].(*aladino.StringValue)
	patch := e.GetTarget().(*target.PullRequestTarget).Patch

	for _, file := range patch {
		if file == nil {
			continue
		}

		isMatch, err := file.Query(arg.Val)
		if err != nil {
			return nil, err
		}

		if isMatch {
			return aladino.BuildTrueValue(), nil
		}
	}

	return aladino.BuildFalseValue(), nil
}
