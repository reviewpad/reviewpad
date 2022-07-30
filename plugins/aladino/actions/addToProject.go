// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"errors"
	"strings"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/shurcooL/githubv4"
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

type FieldDetails struct {
	ID      string
	Name    string
	Options []struct {
		ID   string
		Name string
	}
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

	var repositoryProjectsQuery struct {
		Repository struct {
			ProjectsV2 struct {
				Nodes []struct {
					ID     string
					Fields struct {
						Nodes []struct {
							Details FieldDetails `graphql:"... on ProjectV2SingleSelectField"`
						}
					} `graphql:"fields(first: 50, orderBy: {field: NAME, direction: ASC})"`
				}
			} `graphql:"projectsV2(query: $name, first: 1, orderBy: {field: TITLE, direction: ASC})"`
		} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
	}

	varGQLRepositoryProjectsQuery := map[string]interface{}{
		"repositoryOwner": githubv4.String(owner),
		"repositoryName":  githubv4.String(repo),
		"name":            githubv4.String(projectName),
	}

	if err := e.GetClientGQL().Query(e.GetCtx(), &repositoryProjectsQuery, varGQLRepositoryProjectsQuery); err != nil {
		return err
	}

	if len(repositoryProjectsQuery.Repository.ProjectsV2.Nodes) == 0 {
		return ErrProjectNotFound
	}

	project := repositoryProjectsQuery.Repository.ProjectsV2.Nodes[0]

	statusField := FieldDetails{}

	for _, field := range project.Fields.Nodes {
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

	err := e.GetClientGQL().Mutate(e.GetCtx(), &addProjectV2ItemByIdMutation, input, nil)

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
