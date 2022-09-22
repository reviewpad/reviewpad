// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"
	"log"

	"github.com/reviewpad/reviewpad/v3/codehost"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func RemoveLabels() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType())}, nil),
		Code:           removeLabelsCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func removeLabelsCode(e aladino.Env, args []aladino.Value) error {
	t := e.GetTarget()

	labelsToBeRemoved := args[0].(*aladino.ArrayValue).Vals

	if len(labelsToBeRemoved) == 0 {
		return fmt.Errorf("removeLabels: no labels provided")
	}

	for _, labelToBeRemoved := range labelsToBeRemoved {
		labelID := labelToBeRemoved.(*aladino.StringValue).Val
		internalLabelID := aladino.BuildInternalLabelID(labelID)

		var labelName string

		if val, ok := e.GetRegisterMap()[internalLabelID]; ok {
			labelName = val.(*aladino.StringValue).Val
		} else {
			labelName = labelID
			log.Printf("[warn]: the \"%v\" label was not found in the environment", labelID)
		}

		labels, err := t.GetLabels()
		if err != nil {
			return err
		}

		if isLabelAppliedToIssue(labelName, labels) {
			err := t.RemoveLabel(labelName)
			if err != nil {
				return err
			}
		} else {
			log.Printf("[warn]: the \"%v\" label is not applied to the issue", labelName)
		}
	}

	return nil
}

func isLabelAppliedToIssue(label string, issuesLabels []*codehost.Label) bool {
	for _, issueLabel := range issuesLabels {
		if issueLabel.Name == label {
			return true
		}
	}

	return false
}
