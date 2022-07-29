// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var contains = plugins_aladino.PluginBuiltIns().Functions["contains"].Code

func TestContainsTrue(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil)

	wantVal := aladino.BuildBoolValue(true)

	args := []aladino.Value{
		aladino.BuildStringValue("A title of a pull request"),
		aladino.BuildStringValue("title of"),
	}

	gotVal, err := contains(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestContainsFalse(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil)

	wantVal := aladino.BuildBoolValue(false)

	args := []aladino.Value{
		aladino.BuildStringValue("Pull request title"),
		aladino.BuildStringValue("title of"),
	}

	gotVal, err := contains(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
