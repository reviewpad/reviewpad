// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Any() *aladino.BuiltInFunction {
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
		Code:           anyCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func anyCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	elems := args[0].(*lang.ArrayValue).Vals
	fn := args[1].(*lang.FunctionValue).Fn
	match := false

	for _, elem := range elems {
		match = fn([]lang.Value{elem}).(*lang.BoolValue).Val
		if match {
			break
		}
	}

	return lang.BuildBoolValue(match), nil
}
