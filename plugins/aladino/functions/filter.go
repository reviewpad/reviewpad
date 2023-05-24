// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Filter() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType(
			[]aladino.Type{
				aladino.BuildArrayOfType(aladino.BuildStringType()),
				aladino.BuildFunctionType(
					[]aladino.Type{aladino.BuildStringType()},
					aladino.BuildBoolType(),
				),
			},
			aladino.BuildArrayOfType(aladino.BuildStringType()),
		),
		Code:           filterCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func filterCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	result := make([]lang.Value, 0)
	elems := args[0].(*lang.ArrayValue).Vals
	fn := args[1].(*lang.FunctionValue).Fn

	for _, elem := range elems {
		fnResult := fn([]lang.Value{elem}).(*lang.BoolValue).Val
		if fnResult {
			result = append(result, elem)
		}
	}

	return lang.BuildArrayValue(result), nil
}
