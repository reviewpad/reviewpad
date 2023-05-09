// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"
	"regexp"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func MatchString() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code:           matchStringCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func matchStringCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	pattern := args[0].(*aladino.StringValue).Val
	str := args[1].(*aladino.StringValue).Val

	reg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to parse regex %s %w", pattern, err)
	}

	return aladino.BuildBoolValue(reg.MatchString(str)), nil
}
