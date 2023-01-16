// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github

import (
	"context"

	"github.com/google/go-github/v49/github"
)

func (c *GithubClient) ListOrganizationMembers(ctx context.Context, org string, opts *github.ListMembersOptions) ([]*github.User, *github.Response, error) {
	return c.clientREST.Organizations.ListMembers(ctx, org, opts)
}

func (c *GithubClient) ListTeamMembersBySlug(ctx context.Context, org string, slug string, opts *github.TeamListTeamMembersOptions) ([]*github.User, *github.Response, error) {
	return c.clientREST.Teams.ListTeamMembersBySlug(ctx, org, slug, opts)
}
