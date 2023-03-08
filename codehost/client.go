// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package codehost

import (
	"context"

	"github.com/google/uuid"
	pbe "github.com/reviewpad/api/go/entities"
	"github.com/reviewpad/api/go/services"
	api "github.com/reviewpad/api/go/services"
	"google.golang.org/grpc/metadata"
)

type CodeHostClient struct {
	HostInfo       *HostInfo
	CodehostClient services.HostsClient
	Token          string
}

type HostInfo struct {
	Host    pbe.Host
	HostUri string
}

// RequestIDKey identifies request id field in context
const RequestIDKey = "request-id"

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

	requestID := uuid.New().String()
	md := metadata.Pairs(RequestIDKey, requestID)
	reqCtx := metadata.NewOutgoingContext(ctx, md)

	_, err := c.CodehostClient.PostGeneralComment(reqCtx, req)
	return err
}
