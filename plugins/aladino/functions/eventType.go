// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"reflect"

	"github.com/google/go-github/v49/github"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func EventType() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildStringType()),
		Code:           eventTypeCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func eventTypeCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	pullRequestEvent := e.GetEventPayload()
	if reflect.TypeOf(pullRequestEvent).String() != "*github.PullRequestEvent" {
		return aladino.BuildStringValue(""), nil
	}

	pullRequestEventPayload := pullRequestEvent.(*github.PullRequestEvent)
	if pullRequestEventPayload == nil {
		return aladino.BuildStringValue(""), nil
	}

	if pullRequestEventPayload.Action == nil {
		return aladino.BuildStringValue(""), nil
	}

	return aladino.BuildStringValue(*pullRequestEventPayload.Action), nil
}
