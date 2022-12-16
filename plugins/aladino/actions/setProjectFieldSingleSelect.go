// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func SetProjectFieldSingleSelect() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildStringType(), aladino.BuildStringType()}, nil),
		Code:           setProjectFieldSingleSelect,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func setProjectFieldSingleSelect(e aladino.Env, args []aladino.Value) error {
	projectTitle := args[0].(*aladino.StringValue).Val
	fieldName := args[1].(*aladino.StringValue).Val
	fieldValue := args[2].(*aladino.StringValue).Val

	err := e.GetTarget().SetProjectFieldSingleSelect(projectTitle, fieldName, fieldValue)
	if err != nil {
		return err
	}

	return nil
}
