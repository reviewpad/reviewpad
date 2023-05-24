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

func HasCodePattern() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code:           hasCodePatternCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func hasCodePatternCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	arg := args[0].(*lang.StringValue)
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
			return lang.BuildTrueValue(), nil
		}
	}

	return lang.BuildFalseValue(), nil
}
