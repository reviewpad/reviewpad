// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package codehost

import (
	"context"

	pbe "github.com/reviewpad/api/go/entities"
	api "github.com/reviewpad/api/go/services"
)

type CodeHostClient struct {
	HostInfo       *HostInfo
	CodehostClient api.HostsClient
	Token          string
}

type HostInfo struct {
	Host    pbe.Host
	HostUri string
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
