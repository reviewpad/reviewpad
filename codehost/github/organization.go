// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package github

import (
	"context"

	"github.com/shurcooL/graphql"
)

type OrganizationTeamsQuery struct {
	RepositoryOwner struct {
		Organization struct {
			Teams struct {
				PageInfo struct {
					HasNextPage bool   `graphql:"hasNextPage"`
					EndCursor   string `graphql:"endCursor"`
				}
				Nodes []struct {
					Slug string `graphql:"slug"`
				}
			} `graphql:"teams(first: 100, after: $cursor)"`
		} `graphql:"... on Organization"`
	} `graphql:"repositoryOwner(login: $owner)"`
}

type UserTeamsQuery struct {
	RepositoryOwner struct {
		Organization struct {
			Teams struct {
				PageInfo struct {
					HasNextPage bool   `graphql:"hasNextPage"`
					EndCursor   string `graphql:"endCursor"`
				}
				Nodes []struct {
					Slug string `graphql:"slug"`
				}
			} `graphql:"teams(first: 100, userLogins: [$login], after: $cursor)"`
		} `graphql:"... on Organization"`
	} `graphql:"repositoryOwner(login: $owner)"`
}

func (c *GithubClient) GetOrganizationTeams(ctx context.Context, owner string) ([]string, error) {
	hasNextPage := true
	teams := []string{}
	varOrganizationTeamsQuery := map[string]interface{}{
		"owner":  graphql.String(owner),
		"cursor": (*graphql.String)(nil),
	}

	for hasNextPage {
		organizationTeamsQuery := &OrganizationTeamsQuery{}

		err := c.GetClientGraphQL().Query(ctx, organizationTeamsQuery, varOrganizationTeamsQuery)
		if err != nil {
			return nil, err
		}

		varOrganizationTeamsQuery["cursor"] = graphql.String(organizationTeamsQuery.RepositoryOwner.Organization.Teams.PageInfo.EndCursor)
		hasNextPage = organizationTeamsQuery.RepositoryOwner.Organization.Teams.PageInfo.HasNextPage

		for _, team := range organizationTeamsQuery.RepositoryOwner.Organization.Teams.Nodes {
			teams = append(teams, team.Slug)
		}
	}

	return teams, nil
}

func (c *GithubClient) GetUserTeams(ctx context.Context, owner string, login string) ([]string, error) {
	hasNextPage := true
	teams := []string{}
	varUserTeamsQuery := map[string]interface{}{
		"owner":  graphql.String(owner),
		"login":  graphql.String(login),
		"cursor": (*graphql.String)(nil),
	}

	for hasNextPage {
		userTeamsQuery := &UserTeamsQuery{}

		err := c.GetClientGraphQL().Query(ctx, userTeamsQuery, varUserTeamsQuery)
		if err != nil {
			return nil, err
		}

		varUserTeamsQuery["cursor"] = graphql.String(userTeamsQuery.RepositoryOwner.Organization.Teams.PageInfo.EndCursor)
		hasNextPage = userTeamsQuery.RepositoryOwner.Organization.Teams.PageInfo.HasNextPage

		for _, team := range userTeamsQuery.RepositoryOwner.Organization.Teams.Nodes {
			teams = append(teams, team.Slug)
		}
	}

	return teams, nil
}
