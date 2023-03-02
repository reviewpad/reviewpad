// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package report_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v4/utils/report"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	wantError := "Error occurred! Details:\nThere is something wrong with testing"
	gotError := report.Error("There is something wrong with %v", "testing")

	assert.Equal(t, wantError, gotError)
}
