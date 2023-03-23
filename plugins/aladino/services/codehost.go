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
	CODEHOST_SERVICE_KEY      = "codehost"
	CODEHOST_SERVICE_ENDPOINT = "INPUT_CODEHOST_SERVICE"
)

func NewCodeHostServiceWithConfig(endpoint string) (services.HostClient, *grpc.ClientConn, error) {
	return clients.NewCodeHostClient(endpoint)
}

func NewCodeHostService() (services.HostClient, *grpc.ClientConn, error) {
	codehostEndpoint := os.Getenv(CODEHOST_SERVICE_ENDPOINT)

	return NewCodeHostServiceWithConfig(codehostEndpoint)
}
