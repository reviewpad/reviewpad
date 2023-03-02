// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package recipes

import (
	"context"

	"github.com/reviewpad/reviewpad/v4/handler"
)

type Recipe interface {
	Run(context context.Context, targetEntity handler.TargetEntity) error
}
