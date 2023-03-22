// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package codehost

import (
	"context"
	"fmt"

	pbe "github.com/reviewpad/api/go/entities"
	api "github.com/reviewpad/api/go/services"
	"github.com/reviewpad/go-lib/host"
	"github.com/reviewpad/go-lib/uri"
	"github.com/reviewpad/reviewpad/v4/handler"
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
		Comment: &pbe.ReviewComment{
			Body: body,
		},
	}

	_, err := c.CodehostClient.PostGeneralComment(ctx, req)

	return err
}

func (c *CodeHostClient) GetCodeReview(ctx context.Context, targetEntity *handler.TargetEntity) (*pbe.CodeReview, error) {
	req := &api.GetCodeReviewRequest{
		AccessToken:  c.Token,
		Slug:         fmt.Sprintf("%s/%s", targetEntity.Owner, targetEntity.Repo),
		Host:         c.HostInfo.Host,
		HostUri:      c.HostInfo.HostUri,
		ReviewNumber: int32(targetEntity.Number),
	}

	resp, err := c.CodehostClient.GetCodeReview(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Review, nil
}
