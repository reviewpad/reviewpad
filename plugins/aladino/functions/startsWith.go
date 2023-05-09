// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"strings"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func StartsWith() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code:           startsWithCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func startsWithCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	str := args[0].(*aladino.StringValue).Val
	prefix := args[1].(*aladino.StringValue).Val

	return aladino.BuildBoolValue(strings.HasPrefix(str, prefix)), nil
}
