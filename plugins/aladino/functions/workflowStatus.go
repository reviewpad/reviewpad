// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"reflect"
	"strings"

	"github.com/google/go-github/v52/github"
	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func WorkflowStatus() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildStringType()),
		Code:           workflowStatusCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest, event_processor.Issue},
	}
}

func workflowStatusCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	workflowName := strings.ToLower(args[0].(*lang.StringValue).Val)

	workflowPayload := e.GetEventPayload()
	if reflect.TypeOf(workflowPayload).String() != "*github.WorkflowRunEvent" {
		return lang.BuildStringValue(""), nil
	}

	workflowRunPayload := workflowPayload.(*github.WorkflowRunEvent).WorkflowRun
	if workflowRunPayload == nil {
		return lang.BuildStringValue(""), nil
	}

	headSHA := workflowRunPayload.GetHeadSHA()
	entity := e.GetTarget().GetTargetEntity()
	owner := entity.Owner
	repo := entity.Repo

	checkRuns, _, err := e.GetGithubClient().ListCheckRunsForRef(e.GetCtx(), owner, repo, headSHA, &github.ListCheckRunsOptions{})
	if err != nil {
		return nil, err
	}

	for _, check := range checkRuns.CheckRuns {
		if *check.Name == workflowName {
			if *check.Status == "completed" {
				return lang.BuildStringValue(*check.Conclusion), nil
			} else {
				return lang.BuildStringValue(*check.Status), nil
			}
		}
	}

	return lang.BuildStringValue(""), nil
}
