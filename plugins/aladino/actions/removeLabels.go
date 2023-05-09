// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func RemoveLabels() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType())}, nil),
		Code:           removeLabelsCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func removeLabelsCode(e aladino.Env, args []aladino.Value) error {
	t := e.GetTarget()
	log := e.GetLogger().WithField("builtin", "removeLabels")

	labelsToRemove := args[0].(*aladino.ArrayValue).Vals
	if len(labelsToRemove) == 0 {
		return fmt.Errorf("removeLabels: no labels provided")
	}

	for _, label := range labelsToRemove {
		labelID := label.(*aladino.StringValue).Val
		internalLabelID := aladino.BuildInternalLabelID(labelID)

		var labelName string

		if val, ok := e.GetRegisterMap()[internalLabelID]; ok {
			labelName = val.(*aladino.StringValue).Val
		} else {
			labelName = labelID
			log.Warnf("the \"%v\" label was not found in the environment", labelID)
		}

		err := t.RemoveLabel(labelName)
		if err != nil {
			return err
		}
	}

	return nil
}
