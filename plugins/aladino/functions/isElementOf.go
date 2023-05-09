// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func IsElementOf() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildArrayOfType(aladino.BuildStringType())}, aladino.BuildBoolType()),
		Code:           isElementOfCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func isElementOfCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	member := args[0].(*aladino.StringValue)
	group := args[1].(*aladino.ArrayValue).Vals

	for _, groupMember := range group {
		if member.Equals(groupMember) {
			return aladino.BuildBoolValue(true), nil
		}
	}

	return aladino.BuildBoolValue(false), nil
}
