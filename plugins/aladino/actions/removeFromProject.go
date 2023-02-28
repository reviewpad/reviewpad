// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"errors"
	"fmt"

	"github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func RemoveFromProject() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code:           removeFromProjectCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func removeFromProjectCode(e aladino.Env, args []aladino.Value) error {
	target := e.GetTarget()
	entity := target.GetTargetEntity()
	owner := entity.Owner
	repo := entity.Repo
	projectName := args[0].(*aladino.StringValue).Val

	project, err := e.GetGithubClient().GetProjectV2ByName(e.GetCtx(), owner, repo, projectName)
	if err != nil {
		return fmt.Errorf("failed to get project by name: %s", err.Error())
	}

	itemID, err := target.GetProjectV2ItemID(project.ID)
	if err != nil {
		if errors.Is(err, github.ErrProjectItemNotFound) {
			return nil
		}

		return fmt.Errorf("failed to get project item id: %s", err.Error())
	}

	err = e.GetGithubClient().DeleteProjectV2Item(e.GetCtx(), project.ID, itemID)
	if err != nil {
		return fmt.Errorf("failed to remove from project: %s", err.Error())
	}

	return nil
}
