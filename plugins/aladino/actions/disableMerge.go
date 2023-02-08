// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func DisableMerge() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code:           disableMergeCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func disableMergeCode(e aladino.Env, args []aladino.Value) error {
	reason := args[0].(*aladino.StringValue).Val

	e.GetChecks()[aladino.ReviewpadMergeGateCheckName] = engine.Check{
		Name:   aladino.ReviewpadMergeGateCheckName,
		Status: engine.CheckStatusFailure,
		Reason: reason,
	}

	return nil
}
