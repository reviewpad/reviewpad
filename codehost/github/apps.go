// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github

import (
	"context"

	"github.com/google/go-github/v48/github"
)

// GetInstallations returns all installations for the github app.
// For more details see https://docs.github.com/en/rest/apps/apps#list-installations-for-the-authenticated-app
func (c *GithubAppClient) GetInstallations(ctx context.Context) ([]*github.Installation, error) {
	installations, err := PaginatedRequest(
		func() interface{} {
			return []*github.Installation{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			currentInstallations := i.([]*github.Installation)
			installations, resp, err := c.Apps.ListInstallations(ctx, &github.ListOptions{
				Page:    page,
				PerPage: maxPerPage,
			})
			if err != nil {
				return nil, nil, err
			}
			currentInstallations = append(currentInstallations, installations...)
			return currentInstallations, resp, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return installations.([]*github.Installation), nil
}

// CreateInstallationToken creates a new installation token for the github app.
// For more details see https://docs.github.com/en/rest/apps/apps#create-an-installation-access-token-for-an-app
func (c *GithubAppClient) CreateInstallationToken(ctx context.Context, id int64, opts *github.InstallationTokenOptions) (*github.InstallationToken, error) {
	token, _, err := c.Apps.CreateInstallationToken(ctx, id, opts)
	if err != nil {
		return nil, err
	}

	return token, nil
}
