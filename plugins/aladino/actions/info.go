// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Info() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, nil),
		Code:           infoCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func infoCode(e aladino.Env, args []lang.Value) error {
	body := args[0].(*lang.StringValue).Val

	reportedMessages := e.GetBuiltInsReportedMessages()
	reportedMessages[aladino.SEVERITY_INFO] = append(reportedMessages[aladino.SEVERITY_INFO], body)

	return nil
}
