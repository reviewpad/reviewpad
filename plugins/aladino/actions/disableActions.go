// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func DisableActions() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType())}, nil),
		Code: disableActionsCode,
	}
}

func disableActionsCode(e aladino.Env, args []aladino.Value) error {
	actionNames := args[0].(*aladino.ArrayValue).Vals

	for _, actionName := range actionNames {
		actionName := actionName.(*aladino.StringValue).Val
		e.GetBuiltIns().Actions[actionName].Disabled = true
	}

	return nil
}
