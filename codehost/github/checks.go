// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github

import (
	"context"

	"github.com/google/go-github/v49/github"
)

func (c *GithubClient) ListCheckRunsForRef(ctx context.Context, owner string, repo string, ref string, opts *github.ListCheckRunsOptions) (*github.ListCheckRunsResults, *github.Response, error) {
	return c.clientREST.Checks.ListCheckRunsForRef(ctx, owner, repo, ref, opts)
}
