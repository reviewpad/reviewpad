// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"reflect"
	"strings"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
)

// workflowStatus gets GitHub workflow status. Status can be "success" or "failure".
func WorkflowStatus() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildStringType()),
		Code: workflowStatusCode,
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
	owner := utils.GetPullRequestBaseOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestBaseRepoName(e.GetPullRequest())

	checkRuns, _, err := e.GetClient().Checks.ListCheckRunsForRef(e.GetCtx(), owner, repo, headSHA, &github.ListCheckRunsOptions{})
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
