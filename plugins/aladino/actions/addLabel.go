// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func AddLabel() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, nil),
		Code:           addLabelCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func addLabelCode(e aladino.Env, args []lang.Value) error {
	t := e.GetTarget()
	labelID := args[0].(*lang.StringValue).Val
	internalLabelID := aladino.BuildInternalLabelID(labelID)
	log := e.GetLogger().WithField("builtin", "addLabel")

	var labelName string

	if val, ok := e.GetRegisterMap()[internalLabelID]; ok {
		labelName = val.(*lang.StringValue).Val
	} else {
		labelName = labelID
		log.Infof("the '%v' label is not defined in the 'labels:' section", labelID)
	}

	return t.AddLabels([]string{labelName})
}
