// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var contains = plugins_aladino.PluginBuiltIns().Functions["contains"].Code

func TestContainsTrue(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	wantVal := lang.BuildBoolValue(true)

	args := []lang.Value{
		lang.BuildStringValue("A title of a pull request"),
		lang.BuildStringValue("title of"),
	}

	gotVal, err := contains(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestContainsFalse(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	wantVal := lang.BuildBoolValue(false)

	args := []lang.Value{
		lang.BuildStringValue("Pull request title"),
		lang.BuildStringValue("title of"),
	}

	gotVal, err := contains(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
