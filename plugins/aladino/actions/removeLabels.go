// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func RemoveLabels() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildArrayOfType(lang.BuildStringType())}, nil),
		Code:           removeLabelsCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func removeLabelsCode(e aladino.Env, args []lang.Value) error {
	t := e.GetTarget()
	log := e.GetLogger().WithField("builtin", "removeLabels")

	labelsToRemove := args[0].(*lang.ArrayValue).Vals
	if len(labelsToRemove) == 0 {
		return fmt.Errorf("removeLabels: no labels provided")
	}

	for _, label := range labelsToRemove {
		labelID := label.(*lang.StringValue).Val
		internalLabelID := aladino.BuildInternalLabelID(labelID)

		var labelName string

		if val, ok := e.GetRegisterMap()[internalLabelID]; ok {
			labelName = val.(*lang.StringValue).Val
		} else {
			labelName = labelID
			log.Infof("the '%v' label is not defined in the 'labels:' section", labelID)
		}

		err := t.RemoveLabel(labelName)
		if err != nil {
			return err
		}
	}

	return nil
}
