// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Close() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType(), lang.BuildStringType()}, nil),
		Code:           closeCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest, event_processor.Issue},
	}
}

func closeCode(e aladino.Env, args []lang.Value) error {
	comment := args[0].(*lang.StringValue).Val
	stateReason := args[1].(*lang.StringValue).Val

	return e.GetTarget().Close(comment, stateReason)
}
