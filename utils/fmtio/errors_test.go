// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package fmtio_test

import (
	"errors"
	"testing"

	"github.com/reviewpad/reviewpad/v3/utils/fmtio"
	"github.com/stretchr/testify/assert"
)

func TestErrorf(t *testing.T) {
	wantError := errors.New("[test] testing error message")
	gotError := fmtio.Errorf("test", "testing %v", "error message")

	assert.Equal(t, wantError, gotError)
}
