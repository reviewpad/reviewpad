// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package fmtio_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v3/utils/fmtio"
	"github.com/stretchr/testify/assert"
)

func TestSprintf(t *testing.T) {
	wantStr := "[test] test value: TEST"
	gotStr := fmtio.Sprintf("test", "test value: %v", "TEST")

	assert.Equal(t, wantStr, gotStr)
}

func TestSprint(t *testing.T) {
	wantStr := "[test] testing"
	gotStr := fmtio.Sprint("test", "testing")

	assert.Equal(t, wantStr, gotStr)
}
