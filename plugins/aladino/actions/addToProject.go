// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"errors"
	"strings"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
)

var (
	ErrProjectNotFound         = errors.New("project not found")
	ErrProjectHasNoStatusField = errors.New("project has no status field")
	ErrProjectStatusNotFound   = errors.New("project status not found")
)

type AddProjectV2ItemByIdInput struct {
	ProjectID string `json:"projectId"`
	ContentID string `json:"contentId"`
	// A unique identifier for the client performing the mutation. (Optional.)
	ClientMutationID *string `json:"clientMutationId,omitempty"`
}

type FieldValue struct {
	SingleSelectOptionId string `json:"singleSelectOptionId"`
}

type UpdateProjectV2ItemFieldValueInput struct {
	ItemID    string     `json:"itemId"`
	Value     FieldValue `json:"value"`
	ProjectID string     `json:"projectId"`
	FieldID   string     `json:"fieldId"`
}

func AddToProject() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildStringType()}, aladino.BuildStringType()),
		Code: addToProjectCode,
	}
}

func addToProjectCode(e aladino.Env, args []aladino.Value) error {
	pr := e.GetPullRequest()
	owner := utils.GetPullRequestBaseOwnerName(pr)
	repo := utils.GetPullRequestBaseRepoName(pr)
	projectName := args[0].(*aladino.StringValue).Val
	projectStatus := strings.ToLower(args[1].(*aladino.StringValue).Val)

	project, err := utils.GetProjectV2ByName(e.GetCtx(), e.GetClientGQL(), owner, repo, projectName)
	if err != nil {
		return err
	}

	fields, err := utils.GetProjectFieldsByProjectNumber(e.GetCtx(), e.GetClientGQL(), owner, repo, "", project.Number)
	if err != nil {
		return err
	}

	statusField := utils.FieldDetails{}

	for _, field := range fields {
		if strings.EqualFold(field.Details.Name, "status") {
			statusField = field.Details
			break
		}
	}

	if statusField.ID == "" {
		return ErrProjectHasNoStatusField
	}

	fieldOptionID := ""

	for _, option := range statusField.Options {
		if strings.Contains(strings.ToLower(option.Name), projectStatus) {
			fieldOptionID = option.ID
			break
		}
	}

	if fieldOptionID == "" {
		return ErrProjectStatusNotFound
	}

	var addProjectV2ItemByIdMutation struct {
		AddProjectV2ItemById struct {
			Item struct {
				Id string
			}
		} `graphql:"addProjectV2ItemById(input: $input)"`
	}

	input := AddProjectV2ItemByIdInput{
		ProjectID: project.ID,
		ContentID: *pr.NodeID,
	}

	err = e.GetClientGQL().Mutate(e.GetCtx(), &addProjectV2ItemByIdMutation, input, nil)

	if err != nil {
		return err
	}

	var updateProjectV2ItemFieldValueMutation struct {
		UpdateProjetV2ItemFieldValue struct {
			ClientMutationID string
		} `graphql:"updateProjectV2ItemFieldValue(input: $input)"`
	}

	updateInput := UpdateProjectV2ItemFieldValueInput{
		ProjectID: project.ID,
		ItemID:    addProjectV2ItemByIdMutation.AddProjectV2ItemById.Item.Id,
		Value: FieldValue{
			SingleSelectOptionId: fieldOptionID,
		},
		FieldID: statusField.ID,
	}

	return e.GetClientGQL().Mutate(e.GetCtx(), &updateProjectV2ItemFieldValueMutation, updateInput, nil)
}
