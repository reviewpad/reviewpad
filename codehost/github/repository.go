// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github

import (
	"context"
)

type Branch struct {
	Ref      string
	CloneURL string
	SSHURL   string
}

type PullRequest struct {
	Base       Branch
	Head       Branch
	Rebaseable bool
}

// TODO: Rename to GetPullRequest
func (c *GithubClient) GetPullRequest2(ctx context.Context, owner string, repo string, number int) (*PullRequest, error) {
	pr, _, err := c.clientREST.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		return nil, err
	}

	pullRequest := &PullRequest{
		Rebaseable: *pr.Rebaseable,
		Base: Branch{
			Ref:      pr.Base.GetRef(),
			CloneURL: pr.Base.GetRepo().GetCloneURL(),
			SSHURL:   pr.Base.GetRepo().GetSSHURL(),
		},
		Head: Branch{
			Ref:      pr.Head.GetRef(),
			CloneURL: pr.Head.GetRepo().GetCloneURL(),
			SSHURL:   pr.Head.GetRepo().GetSSHURL(),
		},
	}

	return pullRequest, nil
}
