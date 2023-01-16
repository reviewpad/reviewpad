// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github

import (
	"context"

	"github.com/google/go-github/v49/github"
)

func (c *GithubClient) TriggerWorkflowByFileName(ctx context.Context, owner, repo, branch, workflowFileName string) (*github.Response, error) {
	return c.clientREST.Actions.CreateWorkflowDispatchEventByFileName(ctx, owner, repo, workflowFileName,
		github.CreateWorkflowDispatchEventRequest{
			Ref: branch,
		},
	)
}
