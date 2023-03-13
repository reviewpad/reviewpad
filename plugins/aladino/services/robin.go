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
	ROBIN_SERVICE_KEY      = "robin"
	ROBIN_SERVICE_ENDPOINT = "INPUT_ROBIN_SERVICE"
)

func NewRobinServiceWithConfig(endpoint string) (services.RobinClient, *grpc.ClientConn, error) {
	return clients.NewRobinClient(endpoint)
}

func NewRobinService() (services.RobinClient, *grpc.ClientConn, error) {
	robinEndpoint := os.Getenv(ROBIN_SERVICE_ENDPOINT)

	return NewRobinServiceWithConfig(robinEndpoint)
}
