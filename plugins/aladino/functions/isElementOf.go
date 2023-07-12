// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func IsElementOf() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType(), lang.BuildArrayOfType(lang.BuildStringType())}, lang.BuildBoolType()),
		Code:           isElementOfCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest, event_processor.Issue},
	}
}

func isElementOfCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	member := args[0].(*lang.StringValue)
	group := args[1].(*lang.ArrayValue).Vals

	for _, groupMember := range group {
		if member.Equals(groupMember) {
			return lang.BuildBoolValue(true), nil
		}
	}

	return lang.BuildBoolValue(false), nil
}
