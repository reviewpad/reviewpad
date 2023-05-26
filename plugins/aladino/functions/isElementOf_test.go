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

var isElementOf = plugins_aladino.PluginBuiltIns().Functions["isElementOf"].Code

func TestIsElementOf_WhenTrue(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	args := []lang.Value{
		lang.BuildStringValue("elemA"),
		lang.BuildArrayValue([]lang.Value{
			lang.BuildStringValue("elemA"),
			lang.BuildStringValue("elemB"),
		}),
	}
	gotVal, err := isElementOf(mockedEnv, args)

	wantVal := lang.BuildBoolValue(true)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestIsElementOf_WhenFalse(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	args := []lang.Value{
		lang.BuildStringValue("elemA"),
		lang.BuildArrayValue([]lang.Value{
			lang.BuildStringValue("elemB"),
			lang.BuildStringValue("elemC"),
		}),
	}
	gotVal, err := isElementOf(mockedEnv, args)

	wantVal := lang.BuildBoolValue(false)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
