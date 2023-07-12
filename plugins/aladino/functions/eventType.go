// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"reflect"

	"github.com/google/go-github/v52/github"
	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func EventType() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{}, lang.BuildStringType()),
		Code:           eventTypeCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest, event_processor.Issue},
	}
}

func eventTypePullRequest(pullRequestEventPayload *github.PullRequestEvent) (lang.Value, error) {
	if pullRequestEventPayload == nil {
		return lang.BuildStringValue(""), nil
	}

	if pullRequestEventPayload.Action == nil {
		return lang.BuildStringValue(""), nil
	}

	return lang.BuildStringValue(*pullRequestEventPayload.Action), nil
}

func eventTypeIssue(issuePayload *github.IssuesEvent) (lang.Value, error) {
	if issuePayload == nil {
		return lang.BuildStringValue(""), nil
	}

	if issuePayload.Action == nil {
		return lang.BuildStringValue(""), nil
	}

	return lang.BuildStringValue(*issuePayload.Action), nil
}

func eventTypeCode(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	pullRequestEvent := e.GetEventPayload()
	if pullRequestEvent == nil {
		return lang.BuildStringValue(""), nil
	}

	switch reflect.TypeOf(pullRequestEvent).String() {
	case "*github.PullRequestEvent":
		return eventTypePullRequest(pullRequestEvent.(*github.PullRequestEvent))
	case "*github.IssuesEvent":
		return eventTypeIssue(pullRequestEvent.(*github.IssuesEvent))
	default:
		return lang.BuildStringValue(""), nil
	}
}
