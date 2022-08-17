// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_services

import (
	"github.com/reviewpad/api/go/clients"
	"github.com/reviewpad/api/go/services"
)

const SEMANTIC_SERVICE = "semantic"

func NewSemanticService(endpoint string) services.SemanticClient {
	client, _ := clients.NewSemanticClient(endpoint)

	return client
}
