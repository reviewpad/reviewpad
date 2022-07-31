// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file
package utils

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

func GetProjectV2ByName(ctx context.Context, client *githubv4.Client, owner, repo, name string) (*ProjectV2, error) {
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

	if err := client.Query(ctx, &getProjectV2ByNameQuery, varGetProjectV2ByNameQueryVariables); err != nil {
		return nil, err
	}

	if len(getProjectV2ByNameQuery.Repository.ProjectsV2.Nodes) == 0 {
		return nil, ErrProjectNotFound
	}

	return &getProjectV2ByNameQuery.Repository.ProjectsV2.Nodes[0], nil
}

func GetProjectFieldsByProjectNumber(ctx context.Context, client *githubv4.Client, owner, repo, endCursor string, projectNumber uint64) ([]FieldNode, error) {
	fields := []FieldNode{}

	var repositoryProjectsQuery struct {
		Repository struct {
			ProjectV2 *struct {
				Fields Fields `graphql:"fields(first: 50, after: $afterCursor, orderBy: {field: NAME, direction: ASC})"`
			} `graphql:"projectV2(number: $projectNumber)"`
		} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
	}

	varGQLRepositoryProjectsQuery := map[string]interface{}{
		"repositoryOwner": githubv4.String(owner),
		"repositoryName":  githubv4.String(repo),
		"projectNumber":   githubv4.Int(projectNumber),
		"afterCursor":     githubv4.String(endCursor),
	}

	if err := client.Query(ctx, &repositoryProjectsQuery, varGQLRepositoryProjectsQuery); err != nil {
		return nil, err
	}

	project := repositoryProjectsQuery.Repository.ProjectV2

	if project == nil {
		return nil, ErrProjectNotFound
	}

	afterCursor := project.Fields.PageInfo.EndCursor

	fields = append(fields, project.Fields.Nodes...)

	if project.Fields.PageInfo.HasNextPage {

		fs, err := GetProjectFieldsByProjectNumber(ctx, client, owner, repo, afterCursor, projectNumber)

		if err != nil {
			return nil, err
		}

		fields = append(fields, fs...)

	}

	return fields, nil
}
