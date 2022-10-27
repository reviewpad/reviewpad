// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/google/go-github/v48/github"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Team() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code:           teamCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func teamCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	teamSlug := args[0].(*aladino.StringValue).Val
	orgName := e.GetTarget().GetTargetEntity().Owner

	members, _, err := e.GetGithubClient().ListTeamMembersBySlug(e.GetCtx(), orgName, teamSlug, &github.TeamListTeamMembersOptions{})
	if err != nil {
		return nil, err
	}

	membersLogin := make([]aladino.Value, len(members))
	for i, member := range members {
		membersLogin[i] = aladino.BuildStringValue(*member.Login)
	}

	return aladino.BuildArrayValue(membersLogin), nil
}
