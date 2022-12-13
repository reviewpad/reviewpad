// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func IsBinary() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code:           isBinaryCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func isBinaryCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	target := e.GetTarget().(*target.PullRequestTarget)
	head := target.PullRequest.GetHead()
	fileNameStr := args[0].(*aladino.StringValue)

	isBinary, err := target.IsFileBinary(head.GetRef(), fileNameStr.Val)
	if err != nil {
		return nil, err
	}

	if isBinary {
		return aladino.BuildTrueValue(), nil
	}

	return aladino.BuildFalseValue(), nil
}
