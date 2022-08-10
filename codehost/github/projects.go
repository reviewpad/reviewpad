// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package github

import (
	"context"
	"errors"

	"github.com/shurcooL/githubv4"
)

var (
	ErrProjectNotFound = errors.New("project not found")
)

type ProjectV2 struct {
	ID     string
	Number uint64
}

type PageInfo struct {
	HasNextPage bool
	EndCursor   string
}

type Fields struct {
	PageInfo PageInfo
	Nodes    []FieldNode
}

type FieldNode struct {
	Details FieldDetails `graphql:"... on ProjectV2SingleSelectField"`
}

type FieldDetails struct {
	ID      string
	Name    string
	Options []struct {
		ID   string
		Name string
	}
}

func (c *GithubClient) GetProjectV2ByName(ctx context.Context, owner, repo, name string) (*ProjectV2, error) {
	// Warning: we've faced trouble before with the Github GraphQL API, we'll update this later when we find a better alternative.
	var getProjectV2ByNameQuery struct {
		Repository struct {
			ProjectsV2 struct {
				Nodes []ProjectV2
			} `graphql:"projectsV2(query: $name, first: 1, orderBy: {field: TITLE, direction: ASC})"`
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

	return &getProjectV2ByNameQuery.Repository.ProjectsV2.Nodes[0], nil
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
