// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"strings"

	"github.com/reviewpad/go-lib/entities"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func AddToProject() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildStringType()}, aladino.BuildStringType()),
		Code:           addToProjectCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func addToProjectCode(e aladino.Env, args []lang.Value) error {
	target := e.GetTarget()
	entity := target.GetTargetEntity()
	owner := entity.Owner
	repo := entity.Repo

	projectName := args[0].(*lang.StringValue).Val
	projectStatus := strings.ToLower(args[1].(*lang.StringValue).Val)
	totalRequestTries := 2

	project, err := e.GetGithubClient().GetProjectV2ByName(e.GetCtx(), owner, repo, projectName)
	if err != nil {
		return err
	}

	fields, err := e.GetGithubClient().GetProjectFieldsByProjectNumber(e.GetCtx(), owner, repo, project.Number, totalRequestTries)
	if err != nil {
		return err
	}

	statusField := gh.SingleSelectFieldDetails{}

	for _, field := range fields {
		if strings.EqualFold(field.SingleSelectFieldDetails.Name, "status") {
			statusField = field.SingleSelectFieldDetails
			break
		}
	}

	if statusField.ID == "" {
		return gh.ErrProjectHasNoStatusField
	}

	fieldOptionID := ""

	for _, option := range statusField.Options {
		if strings.Contains(strings.ToLower(option.Name), projectStatus) {
			fieldOptionID = option.ID
			break
		}
	}

	if fieldOptionID == "" {
		return gh.ErrProjectStatusNotFound
	}

	var addProjectV2ItemByIdMutation struct {
		AddProjectV2ItemById struct {
			Item struct {
				Id string
			}
		} `graphql:"addProjectV2ItemById(input: $input)"`
	}

	input := gh.AddProjectV2ItemByIdInput{
		ProjectID: project.ID,
		ContentID: target.GetNodeID(),
	}

	// FIXME: move mutate to a separate function in the codehost.github package
	err = e.GetGithubClient().GetClientGraphQL().Mutate(e.GetCtx(), &addProjectV2ItemByIdMutation, input, nil)
	if err != nil {
		return err
	}

	var updateProjectV2ItemFieldValueMutation struct {
		UpdateProjetV2ItemFieldValue struct {
			ClientMutationID string
		} `graphql:"updateProjectV2ItemFieldValue(input: $input)"`
	}

	updateInput := gh.UpdateProjectV2ItemFieldValueInput{
		ProjectID: project.ID,
		ItemID:    addProjectV2ItemByIdMutation.AddProjectV2ItemById.Item.Id,
		Value: gh.SingleSelectFieldValue{
			SingleSelectOptionId: fieldOptionID,
		},
		FieldID: statusField.ID,
	}

	return e.GetGithubClient().GetClientGraphQL().Mutate(e.GetCtx(), &updateProjectV2ItemFieldValueMutation, updateInput, nil)
}
