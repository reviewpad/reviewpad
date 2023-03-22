// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package codehost

import (
	"context"

	pbc "github.com/reviewpad/api/go/codehost"
	pbe "github.com/reviewpad/api/go/entities"
	api "github.com/reviewpad/api/go/services"
	"github.com/reviewpad/go-lib/host"
	"github.com/reviewpad/go-lib/uri"
)

const RequestIDKey = "request-id"

type CodeHostClient struct {
	HostInfo       *HostInfo
	CodehostClient api.HostsClient
	Token          string
}

type HostInfo struct {
	Host    pbe.Host
	HostUri string
}

func GetHostInfo(url string) (*HostInfo, error) {
	uriData, err := uri.DataFrom(url)
	if err != nil {
		return nil, err
	}

	hostId, err := host.StringToHost(uriData.Host)
	if err != nil {
		return nil, err
	}

	return &HostInfo{
		Host:    hostId,
		HostUri: uriData.Prefix + uriData.Host,
	}, nil
}

func (c *CodeHostClient) PostGeneralComment(ctx context.Context, slug, repoID, reviewID string, reviewNum int32, body string) error {
	req := &api.PostGeneralCommentRequest{
		Host:             c.HostInfo.Host,
		HostUri:          c.HostInfo.HostUri,
		Slug:             slug,
		ExternalRepoId:   repoID,
		ExternalReviewId: reviewID,
		ReviewNumber:     reviewNum,
		AccessToken:      c.Token,
		Comment: &pbc.ReviewComment{
			Body: body,
		},
	}

	_, err := c.CodehostClient.PostGeneralComment(ctx, req)

	return err
}

func (c *CodeHostClient) GetPullRequest(ctx context.Context, slug string, number int64) (*pbc.PullRequest, error) {
	req := &api.GetPullRequestRequest{
		AccessToken: c.Token,
		Slug:        slug,
		Host:        c.HostInfo.Host,
		HostUri:     c.HostInfo.HostUri,
		Number:      number,
	}

	resp, err := c.CodehostClient.GetPullRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.PullRequest, nil
}

func (c *CodeHostClient) GetPullRequestFiles(ctx context.Context, slug string, number int64) ([]*pbc.File, error) {
	req := &api.GetPullRequestFilesRequest{
		AccessToken: c.Token,
		Slug:        slug,
		Host:        c.HostInfo.Host,
		HostUri:     c.HostInfo.HostUri,
		Number:      number,
	}

	resp, err := c.CodehostClient.GetPullRequestFiles(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Files, nil
}
