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

func State() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{}, lang.BuildStringType()),
		Code:           stateCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest, event_processor.Issue},
	}
}

func stateCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	// We are lower casing the state because the enum values are in upper case
	// and we don't want to break people's code.
	return lang.BuildStringValue(strings.ToLower(e.GetTarget().GetState().String())), nil
}
