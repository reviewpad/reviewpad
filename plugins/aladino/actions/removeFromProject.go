// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"errors"
	"fmt"

	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func RemoveFromProject() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, nil),
		Code:           removeFromProjectCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest, event_processor.Issue},
	}
}

func removeFromProjectCode(e aladino.Env, args []lang.Value) error {
	target := e.GetTarget()
	entity := target.GetTargetEntity()
	owner := entity.Owner
	repo := entity.Repo
	projectName := args[0].(*lang.StringValue).Val

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
