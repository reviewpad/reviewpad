// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func DisableActions() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildArrayOfType(lang.BuildStringType())}, nil),
		Code:           disableActionsCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func disableActionsCode(e aladino.Env, args []lang.Value) error {
	actionNames := args[0].(*lang.ArrayValue).Vals

	for _, actionName := range actionNames {
		actionName := actionName.(*lang.StringValue).Val
		e.GetBuiltIns().Actions[actionName].Disabled = true
	}

	return nil
}
