// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package github

import (
	"context"
	"errors"
	"strings"

	"github.com/shurcooL/githubv4"
)

var (
	ErrProjectHasNoStatusField    = errors.New("project has no status field")
	ErrProjectHasNoSuchField      = errors.New("project field not found")
	ErrProjectHasNoSuchFieldValue = errors.New("project field value not found")
	ErrProjectItemsNotFound       = errors.New("project items not found")
	ErrProjectNotFound            = errors.New("project not found")
	ErrProjectStatusNotFound      = errors.New("project status not found")
	ErrProjectItemNotFound        = errors.New("project item not found")
)

type ProjectV2 struct {
	ID     string
	Number uint64
	Title  string
}

type PageInfo struct {
	HasNextPage bool
	EndCursor   string
}

type AddProjectV2ItemByIdInput struct {
	ProjectID string `json:"projectId"`
	ContentID string `json:"contentId"`
	// A unique identifier for the client performing the mutation. (Optional.)
	ClientMutationID *string `json:"clientMutationId,omitempty"`
}

type FieldValue interface{}

type SingleSelectFieldValue struct {
	SingleSelectOptionId string `json:"singleSelectOptionId"`
}

type TextFieldValue struct {
	Text string `json:"text"`
}

type NumberFieldValue struct {
	Number githubv4.Float `json:"number"`
}

type UpdateProjectV2ItemFieldValueInput struct {
	ItemID    string     `json:"itemId"`
	Value     FieldValue `json:"value"`
	ProjectID string     `json:"projectId"`
	FieldID   string     `json:"fieldId"`
}

type Fields struct {
	PageInfo PageInfo
	Nodes    []FieldNode
}

type FieldNode struct {
	TypeName                 string                   `graphql:"__typename"`
	SingleSelectFieldDetails SingleSelectFieldDetails `graphql:"... on ProjectV2SingleSelectField"`
	FieldDetails             FieldDetails             `graphql:"... on ProjectV2Field"`
}

type SingleSelectFieldDetails struct {
	ID      string
	Name    string
	Options []struct {
		ID   string
		Name string
	}
}

type FieldDetails struct {
	ID       string
	Name     string
	DataType string
}

type DeleteProjectV2ItemMutation struct {
	DeleteProjectV2Item struct {
		ClientMutationID string
	} `graphql:"deleteProjectV2Item(input: $input)"`
}

type GetPullRequestProjectItemsQuery struct {
	Repository struct {
		PullRequest struct {
			ProjectItems struct {
				PageInfo PageInfo
				Nodes    []struct {
					Project struct {
						ID string
					}
					ID string
				}
			} `graphql:"projectItems(first: 100, after: $afterCursor)"`
		} `graphql:"pullRequest(number: $number)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type GetIssueProjectItemsQuery struct {
	Repository struct {
		Issue struct {
			ProjectItems struct {
				PageInfo PageInfo
				Nodes    []struct {
					Project struct {
						ID string
					}
					ID string
				}
			} `graphql:"projectItems(first: 100, after: $afterCursor)"`
		} `graphql:"issue(number: $number)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (c *GithubClient) GetProjectV2ByName(ctx context.Context, owner, repo, name string) (*ProjectV2, error) {
	// Warning: we've faced trouble before with the Github GraphQL API, we'll update this later when we find a better alternative.
	// We request the first 50 projects since GitHub can return several projects  with a partially matching name.
	// For instance, a GitHub organization with two projects "private project" and "public project",
	// a query for "private project" will return both projects. This is because both projects have the word "project".
	var getProjectV2ByNameQuery struct {
		Repository struct {
			ProjectsV2 struct {
				Nodes []ProjectV2
			} `graphql:"projectsV2(query: $name, first: 50, orderBy: {field: TITLE, direction: ASC})"`
		} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
	}

	varGetProjectV2ByNameQueryVariables := map[string]interface{}{
		"repositoryOwner": githubv4.String(owner),
		"repositoryName":  githubv4.String(repo),
		"name":            githubv4.String(name),
	}

	if err := c.clientGQL.Query(ctx, &getProjectV2ByNameQuery, varGetProjectV2ByNameQueryVariables); err != nil {
		return nil, err
	}

	if len(getProjectV2ByNameQuery.Repository.ProjectsV2.Nodes) == 0 {
		return nil, ErrProjectNotFound
	}

	var project ProjectV2
	for _, node := range getProjectV2ByNameQuery.Repository.ProjectsV2.Nodes {
		if strings.EqualFold(node.Title, name) {
			project = node
		}
	}

	return &project, nil
}

func (c *GithubClient) GetProjectFieldsByProjectNumber(ctx context.Context, owner, repo string, projectNumber uint64, retryCount int) ([]FieldNode, error) {
	fields := []FieldNode{}
	hasNextPage := true
	currentRequestRetry := 1

	var getProjectFieldsQuery struct {
		Repository struct {
			ProjectV2 *struct {
				Fields Fields `graphql:"fields(first: 50, after: $afterCursor, orderBy: {field: NAME, direction: ASC})"`
			} `graphql:"projectV2(number: $projectNumber)"`
		} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
	}

	varGQLGetProjectFieldsQuery := map[string]interface{}{
		"repositoryOwner": githubv4.String(owner),
		"repositoryName":  githubv4.String(repo),
		"projectNumber":   githubv4.Int(projectNumber),
		"afterCursor":     githubv4.String(""),
	}

	for hasNextPage {
		if err := c.clientGQL.Query(ctx, &getProjectFieldsQuery, varGQLGetProjectFieldsQuery); err != nil {
			currentRequestRetry++
			if currentRequestRetry <= retryCount {
				continue
			}
			return nil, err
		}

		project := getProjectFieldsQuery.Repository.ProjectV2
		if project == nil {
			return nil, ErrProjectNotFound
		}

		fields = append(fields, project.Fields.Nodes...)

		hasNextPage = project.Fields.PageInfo.HasNextPage

		varGQLGetProjectFieldsQuery["afterCursor"] = githubv4.String(project.Fields.PageInfo.EndCursor)
	}

	return fields, nil
}

// GetPullRequestProjectV2ItemID returns the project item ID for a pull request.
// Each item on a project board is represented by a project item.
func (c *GithubClient) GetPullRequestProjectV2ItemID(ctx context.Context, owner, repo, projectID string, number int) (string, error) {
	var getPullRequestProjectItemsQuery GetPullRequestProjectItemsQuery
	varGQLGetPullRequestProjectItemsQuery := map[string]interface{}{
		"owner":       githubv4.String(owner),
		"name":        githubv4.String(repo),
		"number":      githubv4.Int(number),
		"afterCursor": githubv4.String(""),
	}

	for {
		if err := c.clientGQL.Query(ctx, &getPullRequestProjectItemsQuery, varGQLGetPullRequestProjectItemsQuery); err != nil {
			return "", err
		}

		// Since a single pull request can be associated with multiple projects, we need to find the project item ID
		// that matches the project ID we're looking for.
		for _, node := range getPullRequestProjectItemsQuery.Repository.PullRequest.ProjectItems.Nodes {
			if node.Project.ID == projectID {
				return node.ID, nil
			}
		}

		if !getPullRequestProjectItemsQuery.Repository.PullRequest.ProjectItems.PageInfo.HasNextPage {
			break
		}

		varGQLGetPullRequestProjectItemsQuery["afterCursor"] = githubv4.String(getPullRequestProjectItemsQuery.Repository.PullRequest.ProjectItems.PageInfo.EndCursor)
	}

	return "", ErrProjectItemNotFound
}

// GetIssueProjectV2ItemID returns the project item ID for an issue.
// Each item on a project board is represented by a project item.
func (c *GithubClient) GetIssueProjectV2ItemID(ctx context.Context, owner, repo, projectID string, number int) (string, error) {
	var getIssueProjectItemsQuery GetIssueProjectItemsQuery
	varGQLGetIssueProjectItemsQuery := map[string]interface{}{
		"owner":       githubv4.String(owner),
		"name":        githubv4.String(repo),
		"number":      githubv4.Int(number),
		"afterCursor": githubv4.String(""),
	}

	for {
		if err := c.clientGQL.Query(ctx, &getIssueProjectItemsQuery, varGQLGetIssueProjectItemsQuery); err != nil {
			return "", err
		}

		// Since a single issue can be associated with multiple projects, we need to find the project item ID
		// that matches the project ID we're looking for.
		for _, node := range getIssueProjectItemsQuery.Repository.Issue.ProjectItems.Nodes {
			if node.Project.ID == projectID {
				return node.ID, nil
			}
		}

		if !getIssueProjectItemsQuery.Repository.Issue.ProjectItems.PageInfo.HasNextPage {
			break
		}

		varGQLGetIssueProjectItemsQuery["afterCursor"] = githubv4.String(getIssueProjectItemsQuery.Repository.Issue.ProjectItems.PageInfo.EndCursor)
	}

	return "", ErrProjectItemNotFound
}

func (c *GithubClient) DeleteProjectV2Item(ctx context.Context, projectID, itemID string) error {
	var deleteProjectV2ItemMutation DeleteProjectV2ItemMutation
	deleteProjectV2ItemMutationInput := githubv4.DeleteProjectV2ItemInput{
		ProjectID: projectID,
		ItemID:    itemID,
	}

	return c.clientGQL.Mutate(ctx, &deleteProjectV2ItemMutation, deleteProjectV2ItemMutationInput, nil)
}
