// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func HasFileName() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildBoolType()),
		Code:           hasFileNameCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest},
	}
}

func hasFileNameCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	fileNameStr := args[0].(*lang.StringValue)

	patch := e.GetTarget().(*target.PullRequestTarget).Patch
	for fp := range patch {
		if fp == fileNameStr.Val {
			return lang.BuildTrueValue(), nil
		}
	}

	return lang.BuildFalseValue(), nil
}
