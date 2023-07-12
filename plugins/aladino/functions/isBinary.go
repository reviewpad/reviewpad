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

func IsBinary() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildBoolType()),
		Code:           isBinaryCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest},
	}
}

func isBinaryCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	target := e.GetTarget().(*target.PullRequestTarget)
	headBranch := target.PullRequest.Head.Name
	fileName := args[0].(*lang.StringValue).Val

	isBinary, err := target.IsFileBinary(headBranch, fileName)
	if err != nil {
		return nil, err
	}

	return lang.BuildBoolValue(isBinary), nil
}
