// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

func TestFileExt_WhenFilePathHasNoExtension(t *testing.T) {
	fp := "test/file"

	wantFileExt := ""
	gotFileExt := utils.FileExt(fp)

	assert.Equal(t, wantFileExt, gotFileExt)
}

func TestFileExt_WhenFilePathHasExtension(t *testing.T) {
	fp := "test/file.go"

	wantFileExt := ".go"
	gotFileExt := utils.FileExt(fp)

	assert.Equal(t, wantFileExt, gotFileExt)
}
