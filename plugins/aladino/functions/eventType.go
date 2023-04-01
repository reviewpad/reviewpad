// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"reflect"

	"github.com/google/go-github/v49/github"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func EventType() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildStringType()),
		Code:           eventTypeCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func eventTypePullRequest(pullRequestEventPayload *github.PullRequestEvent) (aladino.Value, error) {
	if pullRequestEventPayload == nil {
		return aladino.BuildStringValue(""), nil
	}

	if pullRequestEventPayload.Action == nil {
		return aladino.BuildStringValue(""), nil
	}

	return aladino.BuildStringValue(*pullRequestEventPayload.Action), nil
}

func eventTypeIssue(issuePayload *github.IssuesEvent) (aladino.Value, error) {
	if issuePayload == nil {
		return aladino.BuildStringValue(""), nil
	}

	if issuePayload.Action == nil {
		return aladino.BuildStringValue(""), nil
	}

	return aladino.BuildStringValue(*issuePayload.Action), nil
}

func eventTypeCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	pullRequestEvent := e.GetEventPayload()
	if pullRequestEvent == nil {
		return aladino.BuildStringValue(""), nil
	}

	switch reflect.TypeOf(pullRequestEvent).String() {
	case "*github.PullRequestEvent":
		return eventTypePullRequest(pullRequestEvent.(*github.PullRequestEvent))
	case "*github.IssuesEvent":
		return eventTypeIssue(pullRequestEvent.(*github.IssuesEvent))
	default:
		return aladino.BuildStringValue(""), nil
	}
}
