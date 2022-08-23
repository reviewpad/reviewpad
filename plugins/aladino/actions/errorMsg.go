// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func ErrorMsg() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code:           errorCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func errorCode(e aladino.Env, args []aladino.Value) error {
	body := args[0].(*aladino.StringValue).Val

	reportedMessages := e.GetBuiltInsReportedMessages()
	reportedMessages[aladino.SEVERITY_ERROR] = append(reportedMessages[aladino.SEVERITY_ERROR], body)

	return nil
}
