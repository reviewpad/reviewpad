// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	var err error
	singleLineInput := `$addLabel("small")`
	multiLineInput := `$addLabel("medium multiline")

`

	_, err = Parse(singleLineInput)
	assert.Nil(t, err)

	_, err = Parse(multiLineInput)
	assert.Nil(t, err)
}
