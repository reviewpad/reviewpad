// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	doublestar "github.com/bmatcuk/doublestar/v4"
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func HasFilePattern() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code:           hasFilePatternCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func hasFilePatternCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	filePatternRegex := args[0].(*lang.StringValue)

	patch := e.GetTarget().(*target.PullRequestTarget).Patch
	for fp := range patch {
		re, err := doublestar.Match(filePatternRegex.Val, fp)
		if err != nil {
			return lang.BuildFalseValue(), err
		}
		if re {
			return lang.BuildTrueValue(), nil
		}
	}

	return lang.BuildFalseValue(), nil
}
