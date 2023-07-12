// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"strings"

	"github.com/reviewpad/go-lib/event/event_processor"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func AddToProject() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType(), lang.BuildStringType()}, lang.BuildStringType()),
		Code:           addToProjectCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest, event_processor.Issue},
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

	itemID, err := target.AddToProject(project.ID)
	if err != nil {
		return err
	}

	// If no status is provided, we don't need to update the project item field because it's already set to the default value.
	// The default value is "No status".
	if projectStatus == "" {
		return nil
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

	var updateProjectV2ItemFieldValueMutation struct {
		UpdateProjetV2ItemFieldValue struct {
			ClientMutationID string
		} `graphql:"updateProjectV2ItemFieldValue(input: $input)"`
	}

	updateInput := gh.UpdateProjectV2ItemFieldValueInput{
		ProjectID: project.ID,
		ItemID:    itemID,
		Value: gh.SingleSelectFieldValue{
			SingleSelectOptionId: fieldOptionID,
		},
		FieldID: statusField.ID,
	}

	return e.GetGithubClient().GetClientGraphQL().Mutate(e.GetCtx(), &updateProjectV2ItemFieldValueMutation, updateInput, nil)
}
