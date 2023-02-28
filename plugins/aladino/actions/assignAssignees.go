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
	assignees := args[0].(*aladino.ArrayValue).Vals
	totalRequiredAssignees := args[1].(*aladino.IntValue).Val
	log := e.GetLogger().WithField("builtin", "assignAssignees")

	if len(assignees) == 0 {
		return fmt.Errorf("assignAssignees: list of assignees can't be empty")
	}

	if totalRequiredAssignees == 0 {
		return fmt.Errorf("assignAssignees: total required assignees is invalid")
	}

	reviewpadDefaultIntValue, err := strconv.Atoi(engine.REVIEWPAD_DEFAULT_INT_VALUE)
	if err != nil {
		return err
	}

	if totalRequiredAssignees == reviewpadDefaultIntValue {
		totalRequiredAssignees = len(assignees)
	}

	totalAvailableAssignees := len(assignees)
	if totalRequiredAssignees > totalAvailableAssignees {
		log.Infof("total required assignees %d exceeds the total available assignees %d", totalRequiredAssignees, totalAvailableAssignees)
		totalRequiredAssignees = totalAvailableAssignees
	}

	if totalRequiredAssignees > 10 {
		return fmt.Errorf("assignAssignees: can only assign up to 10 assignees")
	}

	assigneesSelected := make([]string, totalRequiredAssignees)

	rand.Seed(time.Now().UnixNano())
	for totalRequiredAssignees > 0 {
		var selectedAssigneeIndex int

		// Intn returns, as an int, a non-negative pseudo-random number in the half-open interval [0,n) from the default Source. It panics if n <= 0.
		randMax := totalAvailableAssignees - 1
		if randMax <= 0 {
			selectedAssigneeIndex = 0
		} else {
			selectedAssigneeIndex = rand.Intn(randMax)
		}

		assignee := assignees[selectedAssigneeIndex].(*aladino.StringValue).Val

		if contains(assigneesSelected, assignee) {
			continue
		}

		assigneesSelected = append(assigneesSelected, assignee)
		totalRequiredAssignees--
	}

	return t.AddAssignees(assigneesSelected)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
