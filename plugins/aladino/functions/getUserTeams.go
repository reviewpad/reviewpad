// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func GetUserTeams() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildArrayOfType(lang.BuildStringType())),
		Code:           getUserTeams,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func getUserTeams(e aladino.Env, args []lang.Value) (lang.Value, error) {
	login := args[0].(*lang.StringValue).Val

	userTeams, err := e.GetTarget().GetUserTeams(e.GetCtx(), login)
	if err != nil {
		return nil, fmt.Errorf("error getting user teams. %v", err.Error())
	}

	teams := []lang.Value{}
	for _, team := range userTeams {
		teams = append(teams, lang.BuildStringValue(team))
	}

	return lang.BuildArrayValue(teams), nil
}
