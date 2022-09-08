// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func All() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType(
			[]aladino.Type{
				aladino.BuildArrayOfType(aladino.BuildStringType()),
				aladino.BuildFunctionType(
					[]aladino.Type{aladino.BuildStringType()},
					aladino.BuildBoolType(),
				),
			},
			aladino.BuildBoolType(),
		),
		Code:           allCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func allCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	elems := args[0].(*aladino.ArrayValue).Vals
	fn := args[1].(*aladino.FunctionValue).Fn
	match := false

	for _, elem := range elems {
		match = fn([]aladino.Value{elem}).(*aladino.BoolValue).Val
		if !match {
			break
		}
	}

	return aladino.BuildBoolValue(match), nil
}
