// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"

	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
)

func HasRequiredApprovals() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: lang.BuildFunctionType(
			[]lang.Type{
				lang.BuildIntType(),
				lang.BuildArrayOfType(lang.BuildStringType()),
			},
			lang.BuildBoolType(),
		),
		Code:           hasRequiredApprovalsCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest},
	}
}

func hasRequiredApprovalsCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget)
	totalRequiredApprovals := args[0].(*lang.IntValue).Val
	requiredApprovalsFrom := args[1].(*lang.ArrayValue).Vals

	if len(requiredApprovalsFrom) < totalRequiredApprovals {
		return nil, fmt.Errorf("hasRequiredApprovals: the number of required approvals exceeds the number of members from the given list of required approvals")
	}

	approvedBy, err := pullRequest.GetLatestApprovedReviews()
	if err != nil {
		return nil, err
	}

	totalApprovedReviews := 0
	for _, requiredApproval := range requiredApprovalsFrom {
		if utils.ElementOf(approvedBy, requiredApproval.(*lang.StringValue).Val) {
			totalApprovedReviews++
		}
	}

	return lang.BuildBoolValue(totalApprovedReviews >= totalRequiredApprovals), nil
}
