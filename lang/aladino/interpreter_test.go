// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"fmt"
	"testing"

	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	"github.com/stretchr/testify/assert"
)

func TestBuildInternalLabelID(t *testing.T) {
	labelID := "label_id"

	wantVal := fmt.Sprintf("@label:%v", labelID)

	gotVal := aladino.BuildInternalLabelID(labelID)

	assert.Equal(t, wantVal, gotVal)
}
