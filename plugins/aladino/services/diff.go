// Copyright 2023 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_services

import (
	"os"

	"github.com/reviewpad/api/go/clients"
	"github.com/reviewpad/api/go/services"
	"google.golang.org/grpc"
)

const (
	DIFF_SERVICE_KEY      = "diff"
	DIFF_SERVICE_ENDPOINT = "INPUT_DIFF_SERVICE"
)

func NewDiffServiceWithConfig(endpoint string) (services.DiffClient, *grpc.ClientConn, error) {
	return clients.NewDiffClient(endpoint)
}

func NewDiffService() (services.DiffClient, *grpc.ClientConn, error) {
	diffEndpoint := os.Getenv(DIFF_SERVICE_ENDPOINT)

	return NewDiffServiceWithConfig(diffEndpoint)
}
