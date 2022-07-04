// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils_test

import (
	"fmt"
	"testing"

	"github.com/reviewpad/reviewpad/v2/utils"
	"github.com/stretchr/testify/assert"
)

func TestElementOf_WhenTrue(t *testing.T) {
	s := []string{
		"hello",
		"world",
	}
	str := "world"

	assert.True(t, utils.ElementOf(s, str), fmt.Sprintf("\"%s\" should be an element of the list %v.", str, s))
}

func TestElementOf_WhenFalse(t *testing.T) {
	s := []string{
		"hello",
		"world",
	}
	str := "bye"

	assert.False(t, utils.ElementOf(s, str), fmt.Sprintf("\"%s\" should not be an element of the list %v.", str, s))
}
