// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package github

import (
	"context"
	"fmt"
)

type CreateCommitStatusOptions struct {
	State       string `json:"state"`
	Context     string `json:"context"`
	Description string `json:"description"`
}

type CommitStatus struct {
	ID  uint64 `json:"id"`
	URL string `json:"url"`
}

func (c *GithubClient) CreateCommitStatus(ctx context.Context, owner string, repo string, headSHA string, opt *CreateCommitStatusOptions) (*CommitStatus, error) {
	url := fmt.Sprintf("repos/%s/%s/statuses/%s", owner, repo, headSHA)
	req, err := c.clientREST.NewRequest("POST", url, opt)
	if err != nil {
		return nil, err
	}

	commitStatus := &CommitStatus{}

	_, err = c.clientREST.Do(ctx, req, commitStatus)
	if err != nil {
		return nil, err
	}

	return commitStatus, nil
}
