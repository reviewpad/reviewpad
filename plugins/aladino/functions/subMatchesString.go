// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"
	"regexp"

	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func SubMatchesString() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType(), lang.BuildStringType()}, lang.BuildArrayOfType(lang.BuildStringType())),
		Code:           subMatchesString,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest, event_processor.Issue},
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

func subMatchesString(e aladino.Env, args []lang.Value) (lang.Value, error) {
	pattern := args[0].(*lang.StringValue).Val
	str := args[1].(*lang.StringValue).Val

	reg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex pattern %s %w", pattern, err)
	}

	matches := extractMatches(str, reg)
	mValues := make([]lang.Value, len(matches))
	for i, match := range matches {
		mValues[i] = lang.BuildStringValue(match)
	}

	return lang.BuildArrayValue(mValues), nil
}
