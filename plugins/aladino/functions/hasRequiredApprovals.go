package plugins_aladino_functions

import (
	"fmt"

	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
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

	approvedBy, err := getUserLoginsFromApprovals(pullRequest)
	if err != nil {
		return nil, err
	}

	totalApprovedReviews := 0
	for _, requiredApproval := range requiredApprovalsFrom {
		if contains(approvedBy, requiredApproval.(*aladino.StringValue).Val) {
			totalApprovedReviews++
		}

		if totalApprovedReviews == totalRequiredApprovals {
			return aladino.BuildBoolValue(true), nil
		}
	}

	return aladino.BuildBoolValue(false), nil
}

func getUserLoginsFromApprovals(pullRequest *target.PullRequestTarget) ([]string, error) {
	reviews, err := pullRequest.GetReviews()
	if err != nil {
		return nil, err
	}

	approvedBy := make([]string, 0)
	for _, review := range reviews {
		user := review.User.Login

		if review.State == "APPROVED" && !contains(approvedBy, user) {
			approvedBy = append(approvedBy, user)
		}
	}

	return approvedBy, nil
}

func contains(slice []string, str string) bool {
	for _, elem := range slice {
		if elem == str {
			return true
		}
	}

	return false
}
