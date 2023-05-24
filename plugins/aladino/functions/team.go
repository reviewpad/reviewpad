// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/google/go-github/v52/github"
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Team() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code:           teamCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func teamCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	teamSlug := args[0].(*lang.StringValue).Val
	orgName := e.GetTarget().GetTargetEntity().Owner

	members, _, err := e.GetGithubClient().ListTeamMembersBySlug(e.GetCtx(), orgName, teamSlug, &github.TeamListTeamMembersOptions{})
	if err != nil {
		return nil, err
	}

	membersLogin := make([]lang.Value, len(members))
	for i, member := range members {
		membersLogin[i] = lang.BuildStringValue(*member.Login)
	}

	return lang.BuildArrayValue(membersLogin), nil
}
