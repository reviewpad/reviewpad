// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package cookbook

import (
	"fmt"

	"github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/cookbook/recipes"
	"github.com/sirupsen/logrus"
)

type Cookbook interface {
	GetRecipe(name string) (recipes.Recipe, error)
}

type cookbook struct {
	gitHubClient *github.GithubClient
	logger       *logrus.Entry
}

func NewCookbook(logger *logrus.Entry, gitHubClient *github.GithubClient) Cookbook {
	return &cookbook{
		gitHubClient: gitHubClient,
		logger:       logger,
	}
}

func (c *cookbook) GetRecipe(name string) (recipes.Recipe, error) {
	switch name {
	case "size":
		return recipes.NewSizeRecipe(c.logger.WithField("recipe", "size"), c.gitHubClient), nil
	default:
		return nil, fmt.Errorf("recipe %s not found", name)
	}
}
