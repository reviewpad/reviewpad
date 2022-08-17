// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_services

import (
	"os"

	"github.com/reviewpad/api/go/clients"
	"github.com/reviewpad/api/go/services"
)

const (
	SEMANTIC_SERVICE_KEY      = "semantic"
	SEMANTIC_SERVICE_ENDPOINT = "INPUT_SEMANTIC_SERVICE"
)

func NewSemanticServiceWithConfig(endpoint string) (services.SemanticClient, error) {
	return clients.NewSemanticClient(endpoint)
}

func NewSemanticService() (services.SemanticClient, error) {
	semanticEndpoint := os.Getenv(SEMANTIC_SERVICE_ENDPOINT)

	return NewSemanticServiceWithConfig(semanticEndpoint)
}
