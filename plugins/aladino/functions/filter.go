// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import "github.com/reviewpad/reviewpad/v2/lang/aladino"

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
		Code: filterCode,
	}
}

func filterCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	result := make([]aladino.Value, 0)
	elems := args[0].(*aladino.ArrayValue).Vals
	fn := args[1].(*aladino.FunctionValue).Fn

	for _, elem := range elems {
		fnResult := fn([]aladino.Value{elem}).(*aladino.BoolValue).Val
		if fnResult {
			result = append(result, elem)
		}
	}

	return aladino.BuildArrayValue(result), nil
}
