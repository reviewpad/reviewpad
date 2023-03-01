// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/reviewpad/reviewpad/v4/engine"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func AssignAssignees() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType()), aladino.BuildIntType()}, nil),
		Code:           assignAssigneesCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func assignAssigneesCode(e aladino.Env, args []aladino.Value) error {
	t := e.GetTarget()
	availableAssignees := args[0].(*aladino.ArrayValue).Vals
	totalRequiredAssignees := args[1].(*aladino.IntValue).Val
	log := e.GetLogger().WithField("builtin", "assignAssignees")

	if len(availableAssignees) == 0 {
		return fmt.Errorf("assignAssignees: list of assignees can't be empty")
	}

	if totalRequiredAssignees == 0 {
		return fmt.Errorf("assignAssignees: total required assignees is invalid. please insert a number bigger than 0.")
	}

	reviewpadDefaultIntValue, err := strconv.Atoi(engine.REVIEWPAD_DEFAULT_INT_VALUE)
	if err != nil {
		return err
	}

	if totalRequiredAssignees == reviewpadDefaultIntValue {
		totalRequiredAssignees = len(availableAssignees)
	}

	totalAvailableAssignees := len(availableAssignees)
	if totalRequiredAssignees > totalAvailableAssignees {
		log.Infof("total required assignees %d exceeds the total available assignees %d", totalRequiredAssignees, totalAvailableAssignees)
		totalRequiredAssignees = totalAvailableAssignees
	}

	if totalRequiredAssignees > 10 {
		return fmt.Errorf("assignAssignees: can only assign up to 10 assignees")
	}

	// Skip current assignees if mention on the provided assignees list
	currentAssignees, err := t.GetAssignees()
	if err != nil {
		return err
	}

	for _, assignee := range currentAssignees {
		for _, availableAssignee := range availableAssignees {
			if availableAssignee.(*aladino.StringValue).Val == assignee.Login {
				totalRequiredAssignees--
				availableAssignees = filterAssigneeFromAssignees(availableAssignees, assignee.Login)
				break
			}
		}
	}

	assignees := make([]string, totalRequiredAssignees)

	rand.Seed(time.Now().UnixNano())
	for totalRequiredAssignees > 0 {
		selectedAssigneeIndex := rand.Intn(totalAvailableAssignees)

		assignee := availableAssignees[selectedAssigneeIndex].(*aladino.StringValue).Val

		assignees = append(assignees, assignee)
		totalRequiredAssignees--
		availableAssignees = filterAssigneeFromAssignees(availableAssignees, assignee)
	}

	return t.AddAssignees(assignees)
}

func filterAssigneeFromAssignees(assignees []aladino.Value, assignee string) []aladino.Value {
	var filteredAssignees []aladino.Value
	for _, a := range assignees {
		if a.(*aladino.StringValue).Val != assignee {
			filteredAssignees = append(filteredAssignees, a)
		}
	}
	return filteredAssignees
}
