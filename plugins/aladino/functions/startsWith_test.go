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

var startsWith = plugins_aladino.PluginBuiltIns(nil).Functions["startsWith"].Code

func TestStarts_WhenTrue(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	wantVal := aladino.BuildBoolValue(true)

	args := []aladino.Value{
		aladino.BuildStringValue("Title of a pull request"),
		aladino.BuildStringValue("Title of"),
	}

	gotVal, err := startsWith(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestStarts_WhenFalse(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	wantVal := aladino.BuildBoolValue(false)

	args := []aladino.Value{
		aladino.BuildStringValue("Pull request title"),
		aladino.BuildStringValue("Title"),
	}

	gotVal, err := startsWith(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
