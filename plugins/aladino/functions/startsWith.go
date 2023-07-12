// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"strings"

	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func StartsWith() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType(), lang.BuildStringType()}, lang.BuildBoolType()),
		Code:           startsWithCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest, event_processor.Issue},
	}
}

func startsWithCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	str := args[0].(*lang.StringValue).Val
	prefix := args[1].(*lang.StringValue).Val

	return lang.BuildBoolValue(strings.HasPrefix(str, prefix)), nil
}
