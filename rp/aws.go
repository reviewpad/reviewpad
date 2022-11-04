// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package rp

import (
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

func LoadAwsSecret(secretCache *secretcache.Cache, secretName string) string {
	secret, err := secretCache.GetSecretString(secretName)
	if err != nil {
		log.WithError(err).Fatalf("failed to get secret %v from secret manager: %v", secretName, err)
	}
	return secret
}
