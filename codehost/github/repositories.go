// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"io"

	"github.com/google/go-github/v52/github"
	pbc "github.com/reviewpad/api/go/codehost"
)

type DownloadMethod string

const (
	DownloadMethodSHA        DownloadMethod = "withSHA"
	DownloadMethodBranchName DownloadMethod = "withBranchName"
)

type DownloadContentsOptions struct {
	Method DownloadMethod
}

func (c *GithubClient) GetRepositoryBranch(ctx context.Context, owner string, repo string, branch string, followRedirects bool) (*github.Branch, *github.Response, error) {
	return c.clientREST.Repositories.GetBranch(ctx, owner, repo, branch, followRedirects)
}

func (c *GithubClient) GetDefaultRepositoryBranch(ctx context.Context, owner string, repo string) (string, error) {
	repository, _, err := c.clientREST.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return "", err
	}

	return repository.GetDefaultBranch(), nil
}

func (c *GithubClient) DownloadContents(ctx context.Context, filePath string, branch *pbc.Branch, options *DownloadContentsOptions) ([]byte, error) {
	branchRepoOwner := branch.Repo.Owner
	branchRepoName := branch.Repo.Name

	var branchRef string
	switch options.Method {
	case DownloadMethodSHA:
		branchRef = branch.Sha
	case DownloadMethodBranchName:
		branchRef = branch.Name
	default:
		return nil, fmt.Errorf("invalid download method specified")
	}

	ioReader, _, err := c.clientREST.Repositories.DownloadContents(ctx, branchRepoOwner, branchRepoName, filePath, &github.RepositoryContentGetOptions{
		Ref: branchRef,
	})

	if err != nil {
		return nil, err
	}

	return io.ReadAll(ioReader)
}
