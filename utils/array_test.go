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

func TestElementOf_WhenElementIsPresent(t *testing.T) {
	str := "world"
	list := []string{
		"hello",
		str,
	}

	assert.True(t, utils.ElementOf(list, str), fmt.Sprintf("\"%s\" should be an element of the list %v.", str, list))
}

func TestElementOf_WhenElementIsNotPresent(t *testing.T) {
	str := "bye"
	list := []string{
		"hello",
		"world",
	}

	assert.False(t, utils.ElementOf(list, str), fmt.Sprintf("\"%s\" should not be an element of the list %v.", str, list))
}
