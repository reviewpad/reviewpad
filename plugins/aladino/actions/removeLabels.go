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
		label := labelToBeRemoved.(*aladino.StringValue).Val

		labels, err := t.GetLabels()
		if err != nil {
			return err
		}

		if isLabelAppliedToIssue(label, labels) {
			err := t.RemoveLabel(label)
			if err != nil {
				return err
			}
		} else {
			log.Printf("[warn]: the \"%v\" label is not applied to the issue", label)
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
