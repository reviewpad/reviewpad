// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"
	"regexp"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func MatchString() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType(), lang.BuildStringType()}, lang.BuildBoolType()),
		Code:           matchStringCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func matchStringCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	pattern := args[0].(*lang.StringValue).Val
	str := args[1].(*lang.StringValue).Val

	reg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to parse regex %s %w", pattern, err)
	}

	return lang.BuildBoolValue(reg.MatchString(str)), nil
}
