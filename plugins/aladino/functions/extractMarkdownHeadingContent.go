// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func ExtractMarkdownHeadingContent() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildStringType(), aladino.BuildIntType()}, aladino.BuildStringType()),
		Code:           extractMarkdownHeadingContent,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func extractMarkdownHeadingContent(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	input := args[0].(*aladino.StringValue).Val
	headingTitle := args[1].(*aladino.StringValue).Val
	headingLevel := args[2].(*aladino.IntValue).Val

	if headingLevel < 1 || headingLevel > 6 {
		return nil, fmt.Errorf("invalid heading level: %d", headingLevel)
	}

	escapedTitle := regexp.QuoteMeta(headingTitle)
	pattern := fmt.Sprintf("(?s)%s %s\\s*\\n(.*?)(?:\\n%s\\s|$)", strings.Repeat("#", headingLevel), escapedTitle, strings.Repeat("#", headingLevel))

	reg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex pattern %s %w", pattern, err)
	}

	matches := extractMatches(input, reg)

	return aladino.BuildStringValue(strings.Join(matches, "\n")), nil
}
