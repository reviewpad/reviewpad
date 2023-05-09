// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package recipes

import (
	"context"

	"github.com/reviewpad/go-lib/entities"
)

type Recipe interface {
	Run(context context.Context, targetEntity entities.TargetEntity) error
}
