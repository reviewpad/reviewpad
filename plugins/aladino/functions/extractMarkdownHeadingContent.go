// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func ExtractMarkdownHeadingContent() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType(), lang.BuildStringType(), lang.BuildIntType()}, lang.BuildStringType()),
		Code:           extractMarkdownHeadingContent,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest, event_processor.Issue},
	}
}

func extractMarkdownHeadingContent(e aladino.Env, args []lang.Value) (lang.Value, error) {
	input := args[0].(*lang.StringValue).Val
	headingTitle := args[1].(*lang.StringValue).Val
	headingLevel := args[2].(*lang.IntValue).Val

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

	return lang.BuildStringValue(strings.Join(matches, "\n")), nil
}
