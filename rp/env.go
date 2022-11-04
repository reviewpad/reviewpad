// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package rp

import (
	"os"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
	"github.com/sirupsen/logrus"
)

func LoadEnVar(name string) string {
	value := os.Getenv(name)
	if value == "" {
		logrus.Fatalf("missing required environment variable %s", name)
	}
	return value
}

func MustLoadEnvVarValueFromAws(secretCache *secretcache.Cache, envVarName string) string {
	secretName := LoadEnVar(envVarName)
	if secretName == "" {
		logrus.Panicf("The env var `%s` is not set", envVarName)
	}

	secretValue := LoadAwsSecret(secretCache, secretName)
	if secretValue == "" {
		logrus.Panicf("The value defined for %s on aws secret manager is empty", secretName)
	}

	return secretValue
}
