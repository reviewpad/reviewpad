// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"log"

	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func RemoveLabel() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code:           removeLabelCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func removeLabelCode(e aladino.Env, args []aladino.Value) error {
	t := e.GetTarget()
	labelID := args[0].(*aladino.StringValue).Val
	internalLabelID := aladino.BuildInternalLabelID(labelID)

	var labelName string

	if val, ok := e.GetRegisterMap()[internalLabelID]; ok {
		labelName = val.(*aladino.StringValue).Val
	} else {
		labelName = labelID
		log.Printf("[warn]: the %v label was not found in the environment", labelID)
	}

	labels, err := t.GetLabels()
	if err != nil {
		return err
	}

	labelIsAppliedToPullRequest := false
	for _, label := range labels {
		if label.Name == labelName {
			labelIsAppliedToPullRequest = true
			break
		}
	}

	if !labelIsAppliedToPullRequest {
		return nil
	}

	return t.RemoveLabel(labelName)
}
