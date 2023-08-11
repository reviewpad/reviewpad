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

func GetOrganizationTeams() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{}, lang.BuildArrayOfType(lang.BuildStringType())),
		Code:           getOrganizationTeamsCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func getOrganizationTeamsCode(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	owner := e.GetTarget().GetTargetEntity().Owner

	organizationTeams, err := e.GetTarget().GetOrganizationTeams(e.GetCtx(), owner)
	if err != nil {
		return nil, fmt.Errorf("error getting organization teams. %v", err.Error())
	}

	teams := []lang.Value{}
	for _, team := range organizationTeams {
		teams = append(teams, lang.BuildStringValue(team))
	}

	return lang.BuildArrayValue(teams), nil
}
