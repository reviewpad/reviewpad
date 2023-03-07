// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"

	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
)

func HasRequiredApprovals() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType(
			[]aladino.Type{
				aladino.BuildIntType(),
				aladino.BuildArrayOfType(aladino.BuildStringType()),
			},
			aladino.BuildBoolType(),
		),
		Code:           hasRequiredApprovalsCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func hasRequiredApprovalsCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget)
	totalRequiredApprovals := args[0].(*aladino.IntValue).Val
	requiredApprovalsFrom := args[1].(*aladino.ArrayValue).Vals

	if len(requiredApprovalsFrom) < totalRequiredApprovals {
		return nil, fmt.Errorf("hasRequiredApprovals: the number of required approvals exceeds the number of members from the given list of required approvals")
	}

	approvedBy, err := pullRequest.GetApprovedReviewers()
	if err != nil {
		return nil, err
	}

	totalApprovedReviews := 0
	for _, requiredApproval := range requiredApprovalsFrom {
		if utils.Contains(approvedBy, requiredApproval.(*aladino.StringValue).Val) {
			totalApprovedReviews++
		}
	}

	return aladino.BuildBoolValue(totalApprovedReviews >= totalRequiredApprovals), nil
}
