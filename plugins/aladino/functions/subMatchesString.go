// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"
	"regexp"

	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func SubMatchesString() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildStringType()}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code:           subMatchesString,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func extractMatches(input string, reg *regexp.Regexp) []string {
	matches := reg.FindAllStringSubmatch(input, -1)
	if len(matches) == 0 {
		return []string{}
	}
	if len(matches[0]) == 0 {
		return []string{}
	}

	return matches[0][1:]
}

func subMatchesString(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	pattern := args[0].(*aladino.StringValue).Val
	str := args[1].(*aladino.StringValue).Val

	reg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex pattern %s %w", pattern, err)
	}

	matches := extractMatches(str, reg)
	mValues := make([]aladino.Value, len(matches))
	for i, match := range matches {
		mValues[i] = aladino.BuildStringValue(match)
	}

	return aladino.BuildArrayValue(mValues), nil
}
