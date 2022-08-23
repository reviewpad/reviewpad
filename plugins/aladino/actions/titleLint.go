// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/reviewpad/go-conventionalcommits"
	"github.com/reviewpad/go-conventionalcommits/parser"
	"github.com/reviewpad/host-event-handler/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func TitleLint() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, nil),
		Code:           titleLintCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func titleLintCode(e aladino.Env, _ []aladino.Value) error {
	title := e.GetTarget().GetTitle()

	res, err := parser.NewMachine(conventionalcommits.WithTypes(conventionalcommits.TypesConventional)).Parse([]byte(title))
	if err != nil || !res.Ok() {
		body := fmt.Sprintf("**Unconventional title detected**: '%v'", title)
		if err != nil {
			body = fmt.Sprintf("**Unconventional title detected**: '%v' %v", title, err)
		}

		reportedMessages := e.GetBuiltInsReportedMessages()
		reportedMessages[aladino.SEVERITY_ERROR] = append(reportedMessages[aladino.SEVERITY_ERROR], body)
	}

	return nil
}
