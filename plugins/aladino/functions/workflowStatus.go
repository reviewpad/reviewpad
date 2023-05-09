// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"reflect"
	"strings"

	"github.com/google/go-github/v52/github"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func WorkflowStatus() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildStringType()),
		Code:           workflowStatusCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func workflowStatusCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	workflowName := strings.ToLower(args[0].(*aladino.StringValue).Val)

	workflowPayload := e.GetEventPayload()
	if reflect.TypeOf(workflowPayload).String() != "*github.WorkflowRunEvent" {
		return aladino.BuildStringValue(""), nil
	}

	workflowRunPayload := workflowPayload.(*github.WorkflowRunEvent).WorkflowRun
	if workflowRunPayload == nil {
		return aladino.BuildStringValue(""), nil
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
				return aladino.BuildStringValue(*check.Conclusion), nil
			} else {
				return aladino.BuildStringValue(*check.Status), nil
			}
		}
	}

	return aladino.BuildStringValue(""), nil
}
