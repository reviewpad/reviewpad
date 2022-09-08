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

	// ignoring error because filterCode will always return nil
	filteredValues, _ := filterCode(e, args)

	return aladino.BuildBoolValue(len(elems) > 0 && len(filteredValues.(*aladino.ArrayValue).Vals) == len(elems)), nil
}
