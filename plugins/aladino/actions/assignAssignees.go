// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func AssignAssignees() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildArrayOfType(lang.BuildStringType()), lang.BuildIntType()}, nil),
		Code:           assignAssigneesCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest, event_processor.Issue},
	}
}

func assignAssigneesCode(e aladino.Env, args []lang.Value) error {
	t := e.GetTarget()

	rawAvailableAssignees := args[0].(*lang.ArrayValue).Vals
	totalRequiredAssignees := args[1].(*lang.IntValue).Val

	log := e.GetLogger().WithField("builtin", "assignAssignees")

	availableAssignees := make([]string, len(rawAvailableAssignees))
	for i, v := range rawAvailableAssignees {
		availableAssignees[i] = v.(*lang.StringValue).Val
	}

	if len(availableAssignees) == 0 {
		return fmt.Errorf("assignAssignees: list of assignees can't be empty")
	}

	if totalRequiredAssignees == 0 {
		return fmt.Errorf("assignAssignees: total required assignees is invalid. please insert a number bigger than 0")
	}

	totalAvailableAssignees := len(availableAssignees)
	if totalRequiredAssignees > totalAvailableAssignees {
		log.Infof("total required assignees %d exceeds the total available assignees %d", totalRequiredAssignees, totalAvailableAssignees)
		totalRequiredAssignees = totalAvailableAssignees
	}

	if totalRequiredAssignees > 10 {
		return fmt.Errorf("assignAssignees: can only assign up to 10 assignees")
	}

	return t.AddAssignees(getRandomUsers(availableAssignees, totalRequiredAssignees))
}
